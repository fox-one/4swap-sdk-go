package fswap

import (
	"context"
	"time"

	"github.com/fox-one/4swap-sdk-go/v2/route"
	"github.com/shopspring/decimal"
)

const (
	OrderStateTrading  = "Trading"
	OrderStateRejected = "Rejected"
	OrderStateDone     = "Done"
)

type Order struct {
	ID          string          `json:"id,omitempty"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UserID      string          `json:"user_id,omitempty"`
	State       string          `json:"state,omitempty"`
	PayAssetID  string          `json:"pay_asset_id,omitempty"`
	PayAmount   decimal.Decimal `json:"pay_amount,omitempty"`
	FillAssetID string          `json:"fill_asset_id,omitempty"`
	FillAmount  decimal.Decimal `json:"fill_amount,omitempty"`
	MinAmount   decimal.Decimal `json:"min_amount,omitempty"`
	Paths       route.Paths     `json:"paths,omitempty"`

	Transfers []*Transfer `json:"transfers,omitempty"`
}

type PreOrderReq struct {
	PayAssetID  string `json:"pay_asset_id,omitempty"`
	FillAssetID string `json:"fill_asset_id,omitempty"`
	// pay amount 和 fill amount 二选一
	PayAmount  decimal.Decimal `json:"pay_amount,omitempty"`
	FillAmount decimal.Decimal `json:"fill_amount,omitempty"`
}

// PreOrder 预下单
//
// 如果要同时对多个交易对预下单，不建议使用这个方法；而是先调用 ListPairs
// 然后重复使用 Pairs 去 Route 或者 ReverseRoute，这样只需要调用一次 /pairs 接口
// 不会那么容易触发 Rate Limit
func (c *Client) PreOrder(ctx context.Context, req *PreOrderReq) (*Order, error) {
	pairs, err := c.ListPairs(ctx)
	if err != nil {
		return nil, err
	}

	return PreOrderWithPairs(pairs, req)
}

func PreOrderWithPairs(pairs []*Pair, req *PreOrderReq) (*Order, error) {
	var (
		order *Order
		err   error
	)

	if req.PayAmount.IsPositive() {
		order, err = Route(pairs, req.PayAssetID, req.FillAssetID, req.PayAmount)
	} else {
		order, err = ReverseRoute(pairs, req.PayAssetID, req.FillAssetID, req.FillAmount)
	}

	if err != nil {
		return nil, err
	}

	return order, nil
}

// ReadOrder return order detail by id
// WithToken required
func (c *Client) ReadOrder(ctx context.Context, id string) (*Order, error) {
	const uri = "/api/orders/{id}"
	resp, err := c.request(ctx).SetPathParams(map[string]string{
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
