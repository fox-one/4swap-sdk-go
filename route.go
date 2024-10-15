package fswap

import (
	"math"

	"github.com/fox-one/4swap-sdk-go/v2/route"
	"github.com/fox-one/4swap-sdk-go/v2/swap"
	"github.com/shopspring/decimal"
)

const (
	MaxRouteDepth = 4
)

type Graph map[string]map[string]*Pair

func (g Graph) Add(pay, fill string, pair *Pair) {
	m, ok := g[pay]
	if !ok {
		m = make(map[string]*Pair)
		g[pay] = m
	}

	m[fill] = pair
}

func (g Graph) AddPair(pair *Pair) {
	g.Add(pair.BaseAssetID, pair.QuoteAssetID, pair)
	g.Add(pair.QuoteAssetID, pair.BaseAssetID, pair)
}

type node struct {
	*Result

	p *node
	d int
}

func (n *node) Cmp(a *node) int {
	if c := n.FillAmount.Cmp(a.FillAmount); c != 0 {
		return c
	}

	return cmpDepth(n.d, a.d)
}

func (n *node) ReverseCmp(a *node) int {
	if c := a.PayAmount.Cmp(n.PayAmount); c != 0 {
		return c
	}

	return cmpDepth(n.d, a.d)
}

func (n *node) Contain(id uint16) bool {
	for iter := n; iter != nil && iter.d > 0; iter = iter.p {
		if iter.RouteID == id {
			return true
		}
	}

	return false
}

func (n *node) Results(reverse bool) []*Result {
	results := make([]*Result, n.d)
	for iter := n; iter != nil && iter.d > 0; iter = iter.p {
		idx := iter.d - 1
		if reverse {
			idx = len(results) - iter.d
		}

		results[idx] = iter.Result
	}

	return results
}

func (n *node) route(g Graph, fillAsset string, best *node) {
	if n.d >= MaxRouteDepth {
		return
	}

	for fill, pair := range g[n.FillAssetID] {
		if n.Contain(pair.RouteID) {
			continue
		}

		arrived := fill == fillAsset
		if !arrived && n.d+1 == MaxRouteDepth {
			continue
		}

		r, err := Swap(pair, n.FillAssetID, n.FillAmount)
		if err != nil {
			continue
		}

		next := &node{
			Result: r,
			p:      n,
			d:      n.d + 1,
		}

		if !arrived {
			next.route(g, fillAsset, best)
			continue
		}

		if next.Cmp(best) > 0 {
			*best = *next
		}
	}
}

func (n *node) reverseRoute(g Graph, payAsset string, best *node) {
	if n.d >= MaxRouteDepth {
		return
	}

	for pay, pair := range g[n.PayAssetID] {
		if n.Contain(pair.RouteID) {
			continue
		}

		arrived := pay == payAsset
		if !arrived && n.d+1 == MaxRouteDepth {
			continue
		}

		r, err := ReverseSwap(pair, n.PayAssetID, n.PayAmount)
		if err != nil {
			continue
		}

		next := &node{
			Result: r,
			p:      n,
			d:      n.d + 1,
		}

		if !arrived {
			next.reverseRoute(g, payAsset, best)
			continue
		}

		if next.ReverseCmp(best) > 0 {
			*best = *next
		}
	}
}

func Route(pairs []*Pair, payAssetID, fillAssetID string, payAmount decimal.Decimal) (*Order, error) {
	g := make(Graph)
	for _, pair := range pairs {
		g.AddPair(pair)
	}

	best := &node{Result: &Result{}}
	root := &node{
		Result: &Result{
			FillAssetID: payAssetID,
			FillAmount:  payAmount,
		},
	}

	root.route(g, fillAssetID, best)
	if best.d == 0 {
		return nil, ErrInsufficientLiquiditySwapped
	}

	var (
		ids         []uint16
		one         = decimal.NewFromInt(1)
		priceImpact = one
	)

	for _, r := range best.Results(false) {
		ids = append(ids, r.RouteID)

		p := g[r.PayAssetID][r.FillAssetID]
		x := p.BaseAmount.Div(p.QuoteAmount)
		updatePairWithResult(p, r)
		y := p.BaseAmount.Div(p.QuoteAmount)

		if p.SwapMethod != swap.MethodCurve {
			z := y.Sub(x).Abs().Div(x).Add(one)
			priceImpact = priceImpact.Mul(z)
		}
	}

	order := &Order{
		PayAssetID:  payAssetID,
		PayAmount:   payAmount,
		FillAssetID: fillAssetID,
		FillAmount:  best.FillAmount,
		Paths:       route.Single(ids...),
		PriceImpact: priceImpact.Sub(one),
	}

	return order, nil
}

func ReverseRoute(pairs []*Pair, payAssetID, fillAssetID string, fillAmount decimal.Decimal) (*Order, error) {
	g := make(Graph)
	for _, pair := range pairs {
		g.AddPair(pair)
	}

	best := &node{Result: &Result{PayAmount: decimal.NewFromInt(math.MaxInt64)}}
	root := &node{
		Result: &Result{
			PayAssetID: fillAssetID,
			PayAmount:  fillAmount,
		},
	}

	root.reverseRoute(g, payAssetID, best)
	if best.d == 0 {
		return nil, ErrInsufficientLiquiditySwapped
	}

	var (
		ids         []uint16
		one         = decimal.NewFromInt(1)
		priceImpact = one
	)

	for _, r := range best.Results(true) {
		ids = append(ids, r.RouteID)

		p := g[r.PayAssetID][r.FillAssetID]
		x := p.BaseAmount.Div(p.QuoteAmount)
		updatePairWithResult(p, r)
		y := p.BaseAmount.Div(p.QuoteAmount)

		if p.SwapMethod != swap.MethodCurve {
			z := y.Sub(x).Abs().Div(x).Add(one)
			priceImpact = priceImpact.Mul(z)
		}
	}

	order := &Order{
		PayAssetID:  payAssetID,
		PayAmount:   best.PayAmount,
		FillAssetID: fillAssetID,
		FillAmount:  fillAmount,
		Paths:       route.Single(ids...),
		PriceImpact: priceImpact.Sub(one),
	}

	return order, nil
}

func MergeOrders(orders []*Order) *Order {
	var m Order

	for _, order := range orders {
		m.PayAssetID = order.PayAssetID
		m.PayAmount = m.PayAmount.Add(order.PayAmount)
		m.FillAssetID = order.FillAssetID
		m.FillAmount = m.FillAmount.Add(order.FillAmount)
	}

	for _, order := range orders {
		w := order.PayAmount.Div(m.PayAmount).Shift(2)

		for _, path := range order.Paths {
			path.Weight = uint8(path.Share().Mul(w).IntPart())
			m.Paths = mergePaths(m.Paths, path)
		}
	}

	return &m
}

func cmpDepth(a, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	} else {
		return 0
	}
}

func updatePairWithResult(p *Pair, r *Result) {
	if r.PayAssetID == p.BaseAssetID {
		p.BaseAmount = p.BaseAmount.Add(r.PayAmount).Sub(r.ProfitAmount)
		p.QuoteAmount = p.QuoteAmount.Sub(r.FillAmount)
	} else {
		p.QuoteAmount = p.QuoteAmount.Add(r.PayAmount).Sub(r.ProfitAmount)
		p.BaseAmount = p.BaseAmount.Sub(r.FillAmount)
	}
}

func equalPairs(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}

	for idx, id := range a {
		if id != b[idx] {
			return false
		}
	}

	return true
}

func mergePaths(paths route.Paths, p route.Path) route.Paths {
	for idx, path := range paths {
		if equalPairs(path.Pairs, p.Pairs) {
			path.Weight += p.Weight
			paths[idx] = path
			return paths
		}
	}

	return append(paths, p)
}
