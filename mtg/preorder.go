package fswap

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	PreOrderReq struct {
		PayAssetID  string `json:"pay_asset_id,omitempty"`
		FillAssetID string `json:"fill_asset_id,omitempty"`
		// funds 和 amount 二选一
		Funds  decimal.Decimal `json:"funds,omitempty"`
		Amount decimal.Decimal `json:"amount,omitempty"`

		// deprecated
		MinAmount decimal.Decimal `json:"min_amount,omitempty"`
	}
)

// PreOrder 预下单
//
// 如果要同时对多个交易对预下单，不建议使用这个方法；而是先调用 ListPairs
// 然后重复使用 Pairs 去 Route 或者 ReverseRoute，这样只需要调用一次 /pairs 接口
// 不会那么容易触发 Rate Limit
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
