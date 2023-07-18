package fswap

import (
	"math"

	"github.com/fox-one/4swap-sdk-go/routes"
	"github.com/shopspring/decimal"
)

const (
	MaxRouteDepth = 4
)

// EncodeRoutes encode route ids to string with hashids
// deprecated use routes.Routes.HashString() instead
func EncodeRoutes(ids []int64) string {
	return routes.Routes(ids).HashString()
}

// DecodeRoutes decode route ids from string with hashids
// deprecated use routes.ParseHashedRoutes() instead
func DecodeRoutes(id string) []int64 {
	return routes.ParseHashedRoutes(id)
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

	order := &Order{
		RouteResults: best.Results(false),
	}
	ids := make([]int64, 0, best.d)
	for idx, r := range order.RouteResults {
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

	order := &Order{
		RouteResults: best.Results(true),
	}

	ids := make([]int64, 0, best.d)
	for idx, r := range order.RouteResults {
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

// UpdatePairsWithRouteResults update base & quote amount for pairs after a route
func UpdatePairsWithRouteResults(pairs []*Pair, results []*Result) {
	m := make(map[int64]*Result)
	for _, r := range results {
		m[r.RouteID] = r
	}

	for _, p := range pairs {
		r, ok := m[p.RouteID]
		if !ok {
			continue
		}

		profit := r.PayAmount.Mul(p.ProfitRate).Truncate(8)
		pay := r.PayAmount.Sub(profit)
		fill := r.FillAmount

		if p.BaseAssetID == r.PayAssetID {
			p.BaseAmount = p.BaseAmount.Add(pay)
		} else {
			p.QuoteAmount = p.QuoteAmount.Add(pay)
		}

		if p.BaseAssetID == r.FillAssetID {
			p.BaseAmount = p.BaseAmount.Sub(fill)
		} else {
			p.QuoteAmount = p.QuoteAmount.Sub(fill)
		}
	}
}

// GroupOrderRoutes generate a group routes from orders
func GroupOrderRoutes(orders []*Order) routes.Group {
	if len(orders) == 0 {
		panic("empty orders")
	}

	var (
		g    routes.Group
		m    = map[string]int{}
		pay  = orders[0].PayAssetID
		fill = orders[0].FillAssetID
	)

	for _, order := range orders {
		if order.PayAssetID != pay || order.FillAssetID != fill {
			panic("pay or fill asset id not match")
		}

		if idx, ok := m[order.Routes]; ok {
			path := g[idx]
			path.Amount = path.Amount.Add(order.PayAmount)
			g[idx] = path
		} else {
			path := routes.Path{
				Amount: order.PayAmount,
				Routes: routes.ParseHashedRoutes(order.Routes),
			}

			m[order.Routes] = len(g)
			g = append(g, path)
		}
	}

	return g
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
