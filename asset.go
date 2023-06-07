package fswap

import (
	"context"

	"github.com/shopspring/decimal"
)

type (
	Asset struct {
		ID            string          `json:"id,omitempty"`
		Name          string          `json:"name,omitempty"`
		Symbol        string          `json:"symbol,omitempty"`
		DisplaySymbol string          `json:"display_symbol,omitempty"`
		Logo          string          `json:"logo,omitempty"`
		ChainID       string          `json:"chain_id,omitempty"`
		Price         decimal.Decimal `json:"price,omitempty"`
		Chain         struct {
			ID      string          `json:"id,omitempty"`
			Name    string          `json:"name,omitempty"`
			Symbol  string          `json:"symbol,omitempty"`
			Logo    string          `json:"logo,omitempty"`
			ChainID string          `json:"chain_id,omitempty"`
			Price   decimal.Decimal `json:"price,omitempty"`
			Tag     string          `json:"tag,omitempty"`
		} `json:"chain,omitempty"`
	}
)

// ReadAsset read asset
func ReadAsset(ctx context.Context, assetID string) (*Asset, error) {
	const uri = "/api/assets/{id}"
	resp, err := Request(ctx).SetPathParams(map[string]string{
		"id": assetID,
	}).Get(uri)
	if err != nil {
		return nil, err
	}

	var asset Asset
	if err := UnmarshalResponse(resp, &asset); err != nil {
		return nil, err
	}

	return &asset, nil
}

// ListAssets list all assets
func ListAssets(ctx context.Context) ([]*Asset, error) {
	const uri = "/api/assets"
	resp, err := Request(ctx).Get(uri)
	if err != nil {
		return nil, err
	}

	var body struct {
		Assets []*Asset `json:"assets,omitempty"`
	}

	if err := UnmarshalResponse(resp, &body); err != nil {
		return nil, err
	}

	return body.Assets, nil
}
