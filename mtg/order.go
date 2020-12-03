package fswap

import (
	"context"
	"time"

	"github.com/fox-one/4swap-sdk-go/mtg/routes"
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
	UserID     string    `json:"user_id,omitempty"`
	State      string    `json:"state,omitempty"`
	PayAssetID string    `json:"pay_asset_id,omitempty"`
	// pay amount
	PayAmount   decimal.Decimal `json:"pay_amount,omitempty"`
	FillAssetID string          `json:"fill_asset_id,omitempty"`
	// fill amount
	FillAmount decimal.Decimal `json:"fill_amount,omitempty"`
	// 最少购买量
	MinAmount decimal.Decimal `json:"min_amount,omitempty"`
	// route id
	Routes      routes.Routes   `json:"routes,omitempty"`
	RouteAssets []string        `json:"route_assets,omitempty"`
	RoutePrice  decimal.Decimal `json:"route_price,omitempty"`
	PriceImpact decimal.Decimal `json:"price_impact,omitempty"`
	FollowID    string          `json:"follow_id,omitempty"`
}

// ReadOrder return order detail by id
// WithToken required
func ReadOrder(ctx context.Context, followID string) (*Order, error) {
	const uri = "/api/orders/{follow_id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"follow_id": followID,
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
