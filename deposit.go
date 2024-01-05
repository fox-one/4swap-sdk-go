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

type Deposit struct {
	ID           string          `json:"id,omitempty"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	State        string          `json:"state,omitempty"`
	UserID       string          `json:"user_id,omitempty"`
	FollowID     string          `json:"follow_id,omitempty"`
	BaseAssetID  string          `json:"base_asset_id,omitempty"`
	QuoteAssetID string          `json:"quote_asset_id,omitempty"`
	BaseAmount   decimal.Decimal `json:"base_amount,omitempty"`
	QuoteAmount  decimal.Decimal `json:"quote_amount,omitempty"`
	Slippage     decimal.Decimal `json:"slippage,omitempty"`
}

func (c *Client) ReadDeposit(ctx context.Context, depositID string) (*Deposit, error) {
	const uri = "/api/deposits/{id}"
	resp, err := c.request(ctx).SetPathParams(map[string]string{
		"id": depositID,
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
