package fswap

import (
	"container/list"
	"sort"

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

type Path struct {
	list.List

	subPaths []*Path
}

func (path *Path) Depth() int {
	return path.Len()
}

func (path *Path) Add(result *Result) {
	path.PushBack(result)
}

func (path *Path) Last() (*Result, bool) {
	if last := path.Back(); last != nil {
		return last.Value.(*Result), true
	}

	return nil, false
}

func (path *Path) Results(reverse bool) []*Result {
	results := make([]*Result, 0, path.Len())
	for e := path.Front(); e != nil; e = e.Next() {
		results = append(results, e.Value.(*Result))
	}

	if reverse {
		for i, j := 0, len(results)-1; i < j; i, j = i+1, j-1 {
			results[i], results[j] = results[j], results[i]
		}
	}

	return results
}

func (path *Path) ContainPair(pair *Pair) bool {
	for _, r := range path.Results(false) {
		if r.RouteID == pair.RouteID {
			return true
		}
	}

	return false
}

func (path *Path) Sub() *Path {
	sub := &Path{}
	for _, r := range path.Results(false) {
		sub.Add(r)
	}

	path.subPaths = append(path.subPaths, sub)
	return sub
}

func (path *Path) ListAllPaths(filter func(p *Path) bool) []*Path {
	if filter(path) {
		return []*Path{path}
	}

	var paths []*Path
	for _, sub := range path.subPaths {
		paths = append(paths, sub.ListAllPaths(filter)...)
	}

	return paths
}

func (g Graph) Route(p *Path, payAsset, fillAsset string, payAmount decimal.Decimal) {
	if payAsset == fillAsset {
		return
	}

	if p.Depth() >= MaxRouteDepth {
		return
	}

	for _, pair := range g[payAsset] {
		if p.ContainPair(pair) {
			continue
		}

		r, err := Swap(pair, payAsset, payAmount)
		if err != nil {
			continue
		}

		sub := p.Sub()
		sub.Add(r)

		g.Route(sub, r.FillAssetID, fillAsset, r.FillAmount)
	}
}

func (g Graph) ReverseRoute(p *Path, payAsset, fillAsset string, fillAmount decimal.Decimal) {
	if payAsset == fillAsset {
		return
	}

	if p.Depth() >= MaxRouteDepth {
		return
	}

	for _, pair := range g[fillAsset] {
		if p.ContainPair(pair) {
			continue
		}

		r, err := ReverseSwap(pair, fillAsset, fillAmount)
		if err != nil {
			continue
		}

		sub := p.Sub()
		sub.Add(r)

		g.ReverseRoute(sub, payAsset, r.PayAssetID, r.PayAmount)
	}
}

func Route(pairs []*Pair, payAssetID, fillAssetID string, payAmount decimal.Decimal) (*Order, error) {
	g := make(Graph)
	for _, pair := range pairs {
		g.AddPair(pair)
	}

	root := &Path{}
	g.Route(root, payAssetID, fillAssetID, payAmount)

	paths := root.ListAllPaths(func(p *Path) bool {
		last, ok := p.Last()
		return ok && last.FillAssetID == fillAssetID
	})

	if len(paths) == 0 {
		return nil, ErrInsufficientLiquiditySwapped
	}

	sort.Slice(paths, func(i, j int) bool {
		li, _ := paths[i].Last()
		lj, _ := paths[j].Last()
		return li.FillAmount.GreaterThan(lj.FillAmount)
	})

	best := paths[0]

	order := &Order{}
	ids := make([]int64, 0, best.Depth())
	for idx, r := range best.Results(false) {
		if idx == 0 {
			order.PayAssetID = r.PayAssetID
			order.Funds = r.PayAmount
			order.RouteAssets = append(order.RouteAssets, order.PayAssetID)
		}

		order.FillAssetID = r.FillAssetID
		order.Amount = r.FillAmount
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

	root := &Path{}
	g.ReverseRoute(root, payAssetID, fillAssetID, fillAmount)

	paths := root.ListAllPaths(func(p *Path) bool {
		last, ok := p.Last()
		return ok && last.PayAssetID == payAssetID
	})

	if len(paths) == 0 {
		return nil, ErrInsufficientLiquiditySwapped
	}

	sort.Slice(paths, func(i, j int) bool {
		li, _ := paths[i].Last()
		lj, _ := paths[j].Last()
		return li.PayAmount.LessThan(lj.PayAmount)
	})

	best := paths[0]

	order := &Order{}
	ids := make([]int64, 0, best.Depth())
	for idx, r := range best.Results(true) {
		if idx == 0 {
			order.PayAssetID = r.PayAssetID
			order.Funds = r.PayAmount
			order.RouteAssets = append(order.RouteAssets, order.PayAssetID)
		}

		order.FillAssetID = r.FillAssetID
		order.Amount = r.FillAmount
		order.RouteAssets = append(order.RouteAssets, order.FillAssetID)
		ids = append(ids, r.RouteID)
	}

	order.Routes = EncodeRoutes(ids)
	return order, nil
}
