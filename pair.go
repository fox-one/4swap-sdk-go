package fswap

import (
	"context"

	"github.com/shopspring/decimal"
)

type Pair struct {
	BaseAssetID  string          `json:"base_asset_id,omitempty"`
	BaseAmount   decimal.Decimal `json:"base_amount,omitempty"`
	QuoteAssetID string          `json:"quote_asset_id,omitempty"`
	QuoteAmount  decimal.Decimal `json:"quote_amount,omitempty"`
	FeePercent   decimal.Decimal `json:"fee_percent,omitempty"`
	RouteID      int64           `json:"route_id,omitempty"`
	// 池子总的流动性份额
	LiquidityAssetID string          `json:"liquidity_asset_id,omitempty"`
	Liquidity        decimal.Decimal `json:"liquidity,omitempty"`
	// 我的流动性份额，需要 WithToken
	Share      decimal.Decimal `json:"share,omitempty"`
	SwapMethod string          `json:"swap_method,omitempty"`
	Version    int64           `json:"version,omitempty"`
	// volume
	Volume24h      decimal.Decimal `json:"volume_24h,omitempty"`
	BaseVolume24h  decimal.Decimal `json:"base_volume_24h,omitempty"`
	QuoteVolume24h decimal.Decimal `json:"quote_volume_24h,omitempty"`
}

func (pair *Pair) reverse() {
	pair.BaseAssetID, pair.QuoteAssetID = pair.QuoteAssetID, pair.BaseAssetID
	pair.BaseAmount, pair.QuoteAmount = pair.QuoteAmount, pair.BaseAmount
	pair.BaseVolume24h, pair.QuoteVolume24h = pair.QuoteVolume24h, pair.BaseVolume24h
}

// ReadPair return pair detail by base asset id & quote asset id
func ReadPair(ctx context.Context, base, quote string) (*Pair, error) {
	const uri = "/api/pairs/{base_asset_id}/{quote_asset_id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"base_asset_id":  base,
		"quote_asset_id": quote,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var pair Pair
	if err := UnmarshalResponse(resp, &pair); err != nil {
		return nil, err
	}

	if pair.QuoteAssetID == base {
		pair.reverse()
	}

	return &pair, err
}

// ReadPairs list all pairs
func ListPairs(ctx context.Context) ([]*Pair, error) {
	const uri = "/api/pairs"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var body struct {
		Pairs []*Pair `json:"pairs,omitempty"`
	}

	if err := UnmarshalResponse(resp, &body); err != nil {
		return nil, err
	}

	return body.Pairs, nil
}
