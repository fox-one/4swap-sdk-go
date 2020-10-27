package fswap

import (
	"context"
	"testing"
)

func TestRoute(t *testing.T) {
	ctx := context.Background()
	payAssetID := "c94ac88f-4671-3976-b60a-09064f1811e8"
	fillAssetID := "c6d0c728-2624-429b-8e0d-d9d19b6592fa"

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

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})

	t.Run("0.1 xin", func(t *testing.T) {
		order, err := Route(pairs, payAssetID, fillAssetID, Decimal("0.1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})

	t.Run("1 xin", func(t *testing.T) {
		order, err := Route(pairs, payAssetID, fillAssetID, Decimal("1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})
}

func TestReverseRoute(t *testing.T) {
	ctx := context.Background()
	payAssetID := "c94ac88f-4671-3976-b60a-09064f1811e8"
	fillAssetID := "c6d0c728-2624-429b-8e0d-d9d19b6592fa"

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

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})

	t.Run("0.01 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("0.01"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})

	t.Run("0.1 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("0.1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})

	t.Run("1 btc", func(t *testing.T) {
		order, err := ReverseRoute(pairs, payAssetID, fillAssetID, Decimal("1"))
		if err != nil {
			t.Error("Route", err)
			return
		}

		t.Log(order.Funds, order.Amount, order.RouteAssets, order.Routes)
	})
}
