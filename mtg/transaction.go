package fswap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

//go:generate stringer -type TransactionType -trimprefix TransactionType
type TransactionType int

const (
	_                     TransactionType = iota
	TransactionTypeAdd                    // 1 加注
	TransactionTypeRemove                 // 2 提取
	TransactionTypeSwap                   // 3 兑换
)

type (
	Transaction struct {
		ID            string          `json:"id,omitempty"`
		FollowID      string          `json:"follow_id,omitempty"`
		CreatedAt     time.Time       `json:"created_at,omitempty"`
		UserID        string          `json:"user_id,omitempty"`
		Type          string          `json:"type,omitempty"`
		BaseAssetID   string          `json:"base_asset_id,omitempty"`
		QuoteAssetID  string          `json:"quote_asset_id,omitempty"`
		BaseAmount    decimal.Decimal `json:"base_amount,omitempty"`
		QuoteAmount   decimal.Decimal `json:"quote_amount,omitempty"`
		FeeAssetID    string          `json:"fee_asset_id,omitempty"`
		FeeAmount     decimal.Decimal `json:"fee_amount,omitempty"`
		PayAssetID    string          `json:"pay_asset_id,omitempty"`
		FilledAssetID string          `json:"filled_asset_id,omitempty"`
		Funds         decimal.Decimal `json:"funds,omitempty"`
		Amount        decimal.Decimal `json:"amount,omitempty"`
		Value         decimal.Decimal `json:"value,omitempty"`
		FeeValue      decimal.Decimal `json:"fee_value,omitempty"`
	}

	ListTransactionsReq struct {
		UserID       string          `json:"user_id,omitempty"`
		BaseAssetID  string          `json:"base_asset_id,omitempty"`
		QuoteAssetID string          `json:"quote_asset_id,omitempty"`
		Type         TransactionType `json:"type,omitempty"`
		Cursor       string          `json:"cursor,omitempty"`
		Limit        int             `json:"limit,omitempty"`
	}

	ListTransactionsResp struct {
		Transactions []*Transaction `json:"transactions,omitempty"`
		Pagination   *Pagination    `json:"pagination,omitempty"`
		Summary      struct {
			TotalAddBaseAmount     decimal.Decimal `json:"total_add_base_amount,omitempty"`
			TotalAddQuoteAmount    decimal.Decimal `json:"total_add_quote_amount,omitempty"`
			TotalRemoveBaseAmount  decimal.Decimal `json:"total_remove_base_amount,omitempty"`
			TotalRemoveQuoteAmount decimal.Decimal `json:"total_remove_quote_amount,omitempty"`
		} `json:"summary"`
	}
)

func ParseTransactionType(t string) TransactionType {
	for idx := 0; idx < len(_TransactionType_index)-1; idx++ {
		l, r := _TransactionType_index[idx], _TransactionType_index[idx+1]
		if typ := _TransactionType_name[l:r]; strings.EqualFold(typ, t) {
			return TransactionType(idx + 1)
		}
	}

	return 0
}

func ReadTransaction(ctx context.Context, baseAssetID, quoteAssetID, followID string) (*Transaction, error) {
	const uri = "/api/transactions/{base_asset_id}/{quote_asset_id}/mine/{follow_id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"base_asset_id":  baseAssetID,
		"quote_asset_id": quoteAssetID,
		"follow_id":      followID,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var tx Transaction
	if err := UnmarshalResponse(resp, &tx); err != nil {
		return nil, err
	}

	return &tx, err
}

func ListTransactions(ctx context.Context, req *ListTransactionsReq) (*ListTransactionsResp, error) {
	const uri = "/api/transactions/{base_asset_id}/{quote_asset_id}"
	resp, err := Request(ctx).
		SetPathParams(map[string]string{
			"base_asset_id":  req.BaseAssetID,
			"quote_asset_id": req.QuoteAssetID,
		}).
		SetQueryParams(map[string]string{
			"type":    req.Type.String(),
			"cursor":  req.Cursor,
			"user_id": req.UserID,
			"limit":   fmt.Sprint(req.Limit),
		}).Get(uri)
	if err != nil {
		return nil, err
	}

	var result ListTransactionsResp
	if err := UnmarshalResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, err
}
