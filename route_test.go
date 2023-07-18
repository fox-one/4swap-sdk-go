package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	ctx := context.Background()
	payAssetID := "c94ac88f-4671-3976-b60a-09064f1811e8"
	fillAssetID := "c6d0c728-2624-429b-8e0d-d9d19b6592fa"

	UseEndpoint(MtgEndpoint)
	pairs, err := ListPairs(ctx)
	if err != nil {
		t.Error("ListPairs", err)
		return
	}

	t.Run("0.01 xin", func(t *testing.T) {
		order, err := Route(pairs, payAssetID, fillAssetID, Decimal("0.01"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})

	t.Run("0.1 xin", func(t *testing.T) {
		order, err := Route(pairs, payAssetID, fillAssetID, Decimal("0.1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})

	t.Run("1 xin", func(t *testing.T) {
		order, err := Route(pairs, payAssetID, fillAssetID, Decimal("1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})
}

func TestReverseRoute(t *testing.T) {
	ctx := context.Background()
	payAssetID := "c94ac88f-4671-3976-b60a-09064f1811e8"
	fillAssetID := "c6d0c728-2624-429b-8e0d-d9d19b6592fa"

	UseEndpoint(MtgEndpoint)
	pairs, err := ListPairs(ctx)
	if err != nil {
		t.Error("ListPairs", err)
		return
	}

	t.Run("0.00001 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("0.00001"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})

	t.Run("0.01 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("0.01"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})

	t.Run("0.1 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("0.1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})

	t.Run("1 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.PayAmount, order.FillAmount, order.RouteAssets, order.Routes)
	})
}

func TestDecodeRoutes(t *testing.T) {
	routes := []int64{173, 171, 1, 2, 10}
	routeId := EncodeRoutes(routes)
	decodeRoutes := DecodeRoutes(routeId)
	assert.Equal(t, routes, decodeRoutes)
	t.Log(routes, routeId, decodeRoutes)
}

func TestUpdatePairsWithRouteResults(t *testing.T) {
	pairs := []*Pair{
		{
			RouteID:      1,
			BaseAssetID:  "btc",
			QuoteAssetID: "usdt",
			BaseAmount:   Decimal("100"),
			QuoteAmount:  Decimal("10000"),
			FeePercent:   Decimal("0.003"),
			ProfitRate:   Decimal("0.001"),
		},
		{
			RouteID:      2,
			BaseAssetID:  "btc",
			QuoteAssetID: "xin",
			BaseAmount:   Decimal("100"),
			QuoteAmount:  Decimal("10000"),
			FeePercent:   Decimal("0.003"),
			ProfitRate:   Decimal("0.001"),
		},
	}

	results := []*Result{
		{
			PayAssetID:  "usdt",
			PayAmount:   Decimal("100"),
			FillAssetID: "btc",
			FillAmount:  Decimal("0.01"),
			FeeAssetID:  "usdt",
			FeeAmount:   Decimal("0.3"),
			RouteID:     1,
		},
		{
			PayAssetID:  "btc",
			PayAmount:   Decimal("0.01"),
			FillAssetID: "xin",
			FillAmount:  Decimal("8.6"),
			FeeAssetID:  "btc",
			FeeAmount:   Decimal("0.00003"),
			RouteID:     2,
		},
	}

	UpdatePairsWithRouteResults(pairs, results)

	{
		p := pairs[0]
		assert.Equal(t, Decimal("99.99").String(), p.BaseAmount.String())
		assert.Equal(t, Decimal("10099.9").String(), p.QuoteAmount.String())
	}

	{
		p := pairs[1]
		assert.Equal(t, Decimal("100.00999").String(), p.BaseAmount.String())
		assert.Equal(t, Decimal("9991.4").String(), p.QuoteAmount.String())
	}
}

func TestGroupOrderRoutes(t *testing.T) {
	orders := []*Order{
		{
			PayAssetID:  "usdt",
			FillAssetID: "btc",
			PayAmount:   Decimal("101"),
			Routes:      EncodeRoutes([]int64{2, 5, 8, 1}),
		},
		{
			PayAssetID:  "usdt",
			FillAssetID: "btc",
			PayAmount:   Decimal("102"),
			Routes:      EncodeRoutes([]int64{1, 2, 3}),
		},
		{
			PayAssetID:  "usdt",
			FillAssetID: "btc",
			PayAmount:   Decimal("103"),
			Routes:      EncodeRoutes([]int64{1, 2}),
		},
		{
			PayAssetID:  "usdt",
			FillAssetID: "btc",
			PayAmount:   Decimal("104"),
			Routes:      EncodeRoutes([]int64{1, 2, 3}),
		},
	}

	g := GroupOrderRoutes(orders)
	assert.Len(t, g, 3)
	assert.Equal(t, "101:2,5,8,1|206:1,2,3|103:1,2", g.String())
	assert.Equal(t, "410", g.Sum().String())
}
