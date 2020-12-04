package fswap

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"

	"github.com/fox-one/4swap-sdk-go/mtg/encoder"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

type (
	AddLiquidityReq struct {
		UserID      string          `json:"user_id,omitempty"`
		Pair        *Pair           `json:"pair,omitempty"`
		BaseAmount  decimal.Decimal `json:"base_amount,omitempty"`
		QuoteAmount decimal.Decimal `json:"quote_amount,omitempty"`
		// optional, default 600
		Expire int64 `json:"expire,omitempty"`
		// optional, default 0.01
		Slippage decimal.Decimal `json:"slippage,omitempty"`
	}

	RemoveLiquidityReq struct {
		UserID           string          `json:"user_id,omitempty"`
		LiquidityAssetID string          `json:"liquidity_asset_id,omitempty"`
		Amount           decimal.Decimal `json:"amount,omitempty"`
	}

	TransferReq struct {
		AssetID   string          `json:"asset_id,omitempty"`
		Receivers []string        `json:"receivers,omitempty"`
		Threshold uint            `json:"threshold,omitempty"`
		TraceID   string          `json:"trace_id,omitempty"`
		Amount    decimal.Decimal `json:"amount,omitempty"`
		Memo      string          `json:"memo,omitempty"`
	}

	ActionResult struct {
		FollowID  string         `json:"follow_id,omitempty"`
		Transfers []*TransferReq `json:"transfers,omitempty"`
	}
)

func AddLiquidity(ctx context.Context, req *AddLiquidityReq) (*ActionResult, error) {
	userID, err := uuid.FromString(req.UserID)
	if err != nil {
		return nil, err
	}

	baseAssetID, err := uuid.FromString(req.Pair.BaseAssetID)
	if err != nil {
		return nil, err
	}

	quoteAssetID, err := uuid.FromString(req.Pair.QuoteAssetID)
	if err != nil {
		return nil, err
	}

	if !req.Slippage.IsPositive() {
		req.Slippage = decimal.NewFromFloat(0.01)
	}

	if req.Expire < 1 {
		req.Expire = 10 * 60
	}

	followID := uuid.Must(uuid.NewV4())
	baseMemo, err := encodeMemo(
		int(TransactionTypeAdd),
		userID,
		followID,
		quoteAssetID,
		req.Slippage,
		req.Expire)

	if err != nil {
		return nil, err
	}

	quoteMemo, err := encodeMemo(
		int(TransactionTypeAdd),
		userID,
		followID,
		baseAssetID,
		req.Slippage,
		req.Expire)

	if err != nil {
		return nil, err
	}

	return &ActionResult{
		FollowID: followID.String(),
		Transfers: []*TransferReq{
			{
				AssetID:   req.Pair.BaseAssetID,
				Receivers: group.Members,
				Threshold: group.Threshold,
				TraceID:   uuid.Must(uuid.NewV4()).String(),
				Amount:    req.BaseAmount,
				Memo:      baseMemo,
			},
			{
				AssetID:   req.Pair.QuoteAssetID,
				Receivers: group.Members,
				Threshold: group.Threshold,
				TraceID:   uuid.Must(uuid.NewV4()).String(),
				Amount:    req.QuoteAmount,
				Memo:      quoteMemo,
			},
		},
	}, nil
}

func Swap(ctx context.Context, order *Order) (*ActionResult, error) {
	userID, err := uuid.FromString(order.UserID)
	if err != nil {
		return nil, err
	}

	fillAssetId, err := uuid.FromString(order.FillAssetID)
	if err != nil {
		return nil, err
	}

	followID := uuid.Must(uuid.NewV4())
	memo, err := encodeMemo(
		int(TransactionTypeSwap),
		userID,
		followID,
		fillAssetId,
		order.Routes,
		order.MinAmount)

	if err != nil {
		return nil, err
	}

	return &ActionResult{
		FollowID: followID.String(),
		Transfers: []*TransferReq{
			{
				AssetID:   order.PayAssetID,
				Receivers: group.Members,
				Threshold: group.Threshold,
				TraceID:   followID.String(),
				Amount:    order.PayAmount,
				Memo:      memo,
			},
		},
	}, nil
}

func RemoveLiquidity(ctx context.Context, req *RemoveLiquidityReq) (*ActionResult, error) {
	userID, err := uuid.FromString(req.UserID)
	if err != nil {
		return nil, err
	}

	followID := uuid.Must(uuid.NewV4())
	memo, err := encodeMemo(
		int(TransactionTypeRemove),
		userID,
		followID)
	if err != nil {
		return nil, err
	}

	return &ActionResult{
		FollowID: followID.String(),
		Transfers: []*TransferReq{
			{
				AssetID:   req.LiquidityAssetID,
				Receivers: group.Members,
				Threshold: group.Threshold,
				TraceID:   followID.String(),
				Amount:    req.Amount,
				Memo:      memo,
			},
		},
	}, nil
}

func encodeMemo(values ...interface{}) (string, error) {
	action, err := encoder.Encode(values...)
	if err != nil {
		return "", err
	}

	key := mixin.GenerateEd25519Key()
	var pub ed25519.PublicKey
	copy(pub[:], group.PublicKey[:])
	action, err = encoder.Encrypt(action, key, pub)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(action), nil
}
