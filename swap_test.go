package fswap

import (
	"testing"

	"github.com/fox-one/4swap-sdk-go/v2/swap"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSwap(t *testing.T) {
	pair := &Pair{
		BaseAssetID:  "base-asset-id",
		QuoteAssetID: "quote-asset-id",
		BaseAmount:   decimal.NewFromFloat(1000),
		QuoteAmount:  decimal.NewFromFloat(2000),
		FeePercent:   decimal.NewFromFloat(0.01),
		ProfitRate:   decimal.NewFromFloat(0.02),
		RouteID:      1,
		SwapMethod:   swap.MethodCurve,
	}
	payAmount := decimal.NewFromFloat(100)
	result, err := Swap(pair, "base-asset-id", payAmount)
	assert.NoError(t, err)
	assert.Equal(t, result.FillAssetID, "quote-asset-id")
	assert.Equal(t, result.FillAmount.String(), "99.35838442")
}

func TestReverseSwap(t *testing.T) {
	pair := &Pair{
		BaseAssetID:  "base-asset-id",
		QuoteAssetID: "quote-asset-id",
		BaseAmount:   decimal.NewFromFloat(1000),
		QuoteAmount:  decimal.NewFromFloat(2000),
		FeePercent:   decimal.NewFromFloat(0.01),
		ProfitRate:   decimal.NewFromFloat(0.02),
		RouteID:      1,
		SwapMethod:   swap.MethodCurve,
	}

	fillAmount := decimal.NewFromFloat(50)
	result, err := ReverseSwap(pair, "quote-asset-id", fillAmount)
	assert.NoError(t, err)
	assert.Equal(t, result.PayAssetID, "base-asset-id")
	assert.Equal(t, result.PayAmount.String(), "50.30905306")
}

func TestMultiHopSwap(t *testing.T) {
	pairs := []*Pair{
		{
			BaseAssetID:  "asset1",
			QuoteAssetID: "asset2",
			BaseAmount:   decimal.NewFromFloat(1000),
			QuoteAmount:  decimal.NewFromFloat(2000),
			FeePercent:   decimal.NewFromFloat(0.01),
			ProfitRate:   decimal.NewFromFloat(0.02),
			RouteID:      1,
			SwapMethod:   swap.MethodCurve,
		},
		{
			BaseAssetID:  "asset2",
			QuoteAssetID: "asset3",
			BaseAmount:   decimal.NewFromFloat(2000),
			QuoteAmount:  decimal.NewFromFloat(3000),
			FeePercent:   decimal.NewFromFloat(0.01),
			ProfitRate:   decimal.NewFromFloat(0.02),
			RouteID:      2,
			SwapMethod:   swap.MethodCurve,
		},
	}

	payAmount := decimal.NewFromFloat(100)
	result, err := MultiHopSwap(pairs, "asset1", payAmount)
	assert.NoError(t, err)
	assert.Equal(t, result.FillAssetID, "asset3")
	assert.Equal(t, result.FillAmount.String(), "98.55357136")
}

func TestReverseMultiHopSwap(t *testing.T) {
	pairs := []*Pair{
		{
			BaseAssetID:  "asset1",
			QuoteAssetID: "asset2",
			BaseAmount:   decimal.NewFromFloat(1000),
			QuoteAmount:  decimal.NewFromFloat(2000),
			FeePercent:   decimal.NewFromFloat(0.01),
			ProfitRate:   decimal.NewFromFloat(0.02),
			RouteID:      1,
			SwapMethod:   swap.MethodCurve,
		},
		{
			BaseAssetID:  "asset2",
			QuoteAssetID: "asset3",
			BaseAmount:   decimal.NewFromFloat(2000),
			QuoteAmount:  decimal.NewFromFloat(3000),
			FeePercent:   decimal.NewFromFloat(0.01),
			ProfitRate:   decimal.NewFromFloat(0.02),
			RouteID:      2,
			SwapMethod:   swap.MethodCurve,
		},
	}

	fillAmount := decimal.NewFromFloat(50)
	result, err := ReverseMultiHopSwap(pairs, "asset3", fillAmount)
	assert.NoError(t, err)
	assert.Equal(t, result.PayAssetID, "asset1")
	assert.Equal(t, result.PayAmount.String(), "50.71408239")
}
