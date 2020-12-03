package fswap

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	DepositStatePending   = "Pending"
	DepositStateCancelled = "Cancelled"
	DepositStateDone      = "Done"
)

type (
	Deposit struct {
		ID           string          `json:"id,omitempty"`
		CreatedAt    time.Time       `json:"created_at,omitempty"`
		State        string          `json:"state,omitempty"`
		UserID       string          `json:"user_id,omitempty"`
		FollowID     string          `json:"follow_id,omitempty"`
		BaseAssetID  string          `json:"base_asset_id,omitempty"`
		BaseAmount   decimal.Decimal `json:"base_amount,omitempty"`
		QuoteAssetID string          `json:"quote_asset_id,omitempty"`
		QuoteAmount  decimal.Decimal `json:"quote_amount,omitempty"`
		Slippage     decimal.Decimal `json:"slippage,omitempty"`
		Transfers    []*TransferReq  `json:"transfers,omitempty"`
	}
)

func ReadDeposit(ctx context.Context, followID string) (*Deposit, error) {
	const uri = "/api/deposits/{follow_id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"follow_id": followID,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var deposit Deposit
	if err := UnmarshalResponse(resp, &deposit); err != nil {
		return nil, err
	}

	return &deposit, err
}
