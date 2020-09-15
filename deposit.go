package fswap

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

const (
	DepositStatePending = "Pending"
	DepositStateCancelled = "Cancelled"
	DepositStateDone = "Done"
)

type Deposit struct {
	ID string `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	State string `json:"state,omitempty"`
	UserID string `json:"user_id,omitempty"`
	Transfers []*TransferReq `json:"transfers,omitempty"`
}

type AddDepositReq struct {
	BaseAssetID string `json:"base_asset_id,omitempty"`
	BaseAmount decimal.Decimal `json:"base_amount,omitempty"`
	QuoteAssetID string `json:"quote_asset_id,omitempty"`
	QuoteAmount decimal.Decimal `json:"quote_amount,omitempty"`
	// optional, default 0.01
	Slippage decimal.Decimal `json:"slippage,omitempty"`
}

func AddDeposit(ctx context.Context,req *AddDepositReq) (*Deposit,error) {
	const uri = "/api/pairs/{base_asset_id}/{quote_asset_id}/deposit"
	resp,err := Request(ctx).SetPathParams(map[string]string{
		"base_asset_id": req.BaseAssetID,
		"quote_asset_id": req.QuoteAssetID,
	}).SetBody(req).Post(uri)
	if err != nil {
		return nil, err
	}

	var deposit Deposit
	if err := UnmarshalResponse(resp,&deposit);err != nil {
		return nil,err
	}

	return &deposit,err
}

func ReadDeposit(ctx context.Context,depositID string) (*Deposit,error) {
	const uri = "/api/deposits/{id}"
	resp,err := Request(ctx).SetPathParams(map[string]string{
		"id": depositID,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var deposit Deposit
	if err := UnmarshalResponse(resp,&deposit);err != nil {
		return nil,err
	}

	return &deposit,err
}
