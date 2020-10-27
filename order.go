package fswap

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	OrderStateTrading  = "Trading"
	OrderStateRejected = "Rejected"
	OrderStateDone     = "Done"
)

type Order struct {
	ID         string    `json:"id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	State      string    `json:"state,omitempty"`
	PayAssetID string    `json:"pay_asset_id,omitempty"`
	// pay amount
	Funds       decimal.Decimal `json:"funds,omitempty"`
	FillAssetID string          `json:"fill_asset_id,omitempty"`
	// fill amount
	Amount decimal.Decimal `json:"amount,omitempty"`
	// 最少购买量
	MinAmount   decimal.Decimal `json:"min_amount,omitempty"`
	PriceImpact decimal.Decimal `json:"price_impact,omitempty"`
	RouteAssets []string        `json:"route_assets,omitempty"`
	// route id
	Routes string `json:"routes,omitempty"`
}

type PreOrderReq struct {
	PayAssetID  string `json:"pay_asset_id,omitempty"`
	FillAssetID string `json:"fill_asset_id,omitempty"`
	// funds 和 amount 二选一
	Funds  decimal.Decimal `json:"funds,omitempty"`
	Amount decimal.Decimal `json:"amount,omitempty"`

	// deprecated
	MinAmount decimal.Decimal `json:"min_amount,omitempty"`
}

// PreOrder 预下单
func PreOrder(ctx context.Context, req *PreOrderReq) (*Order, error) {
	pairs, err := ListPairs(ctx)
	if err != nil {
		return nil, err
	}

	if req.Funds.IsPositive() {
		return Route(pairs, req.PayAssetID, req.FillAssetID, req.Funds)
	}

	return ReverseRoute(pairs, req.PayAssetID, req.FillAssetID, req.Amount)
}

// ReadOrder return order detail by id
// WithToken required
func ReadOrder(ctx context.Context, id string) (*Order, error) {
	const uri = "/api/order/{id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"id": id,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := UnmarshalResponse(resp, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
