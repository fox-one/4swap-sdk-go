package fswap

import (
	"math"

	"github.com/shopspring/decimal"
	"github.com/speps/go-hashids"
)

const (
	MaxRouteDepth = 4
)

func EncodeRoutes(ids []int64) string {
	hd := hashids.NewData()
	hd.Salt = "uniswap routes"
	h, _ := hashids.NewWithData(hd)
	id, _ := h.EncodeInt64(ids)
	return id
}

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

func (n *node) Contain(id int64) bool {
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

	order := &Order{}
	ids := make([]int64, 0, best.d)
	for idx, r := range best.Results(false) {
		if idx == 0 {
			order.PayAssetID = r.PayAssetID
			order.PayAmount = r.PayAmount
			order.RouteAssets = append(order.RouteAssets, order.PayAssetID)
		}

		order.FillAssetID = r.FillAssetID
		order.FillAmount = r.FillAmount
		order.RouteAssets = append(order.RouteAssets, order.FillAssetID)
		ids = append(ids, r.RouteID)
	}

	order.Routes = EncodeRoutes(ids)
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

	order := &Order{}
	ids := make([]int64, 0, best.d)
	for idx, r := range best.Results(true) {
		if idx == 0 {
			order.PayAssetID = r.PayAssetID
			order.PayAmount = r.PayAmount
			order.RouteAssets = append(order.RouteAssets, order.PayAssetID)
		}

		order.FillAssetID = r.FillAssetID
		order.FillAmount = r.FillAmount
		order.RouteAssets = append(order.RouteAssets, order.FillAssetID)
		ids = append(ids, r.RouteID)
	}

	order.Routes = EncodeRoutes(ids)
	return order, nil
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
