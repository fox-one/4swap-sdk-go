package fswap

import (
	"errors"

	"github.com/shopspring/decimal"
)

var (
	SwapFee = "0.003"

	ErrInsufficientLiquiditySwapped = errors.New("insufficient liquidity swapped")
)

// Result represent Swap Result
type Result struct {
	PayAssetID  string
	PayAmount   decimal.Decimal
	FillAssetID string
	FillAmount  decimal.Decimal
	FeeAssetID  string
	FeeAmount   decimal.Decimal
	RouteID     int64
}

// swap trade in a pair
func swap(pair *Pair, payAssetID string, payAmount decimal.Decimal) (*Result, error) {
	K := pair.BaseAmount.Mul(pair.QuoteAmount)
	if !K.IsPositive() {
		return nil, ErrInsufficientLiquiditySwapped
	}

	payAmount = payAmount.Truncate(8)

	r := &Result{
		PayAssetID: payAssetID,
		PayAmount:  payAmount,
		FeeAssetID: payAssetID,
		FeeAmount:  payAmount.Mul(Decimal(SwapFee)).Truncate(8),
		RouteID:    pair.RouteID,
	}

	funds := payAmount.Sub(r.FeeAmount)
	if !funds.IsPositive() {
		return nil, errors.New("pay amount must be positive")
	}

	switch payAssetID {
	case pair.BaseAssetID:
		newBase := pair.BaseAmount.Add(funds)
		newQuote := K.Div(newBase)
		r.FillAssetID = pair.QuoteAssetID
		r.FillAmount = pair.QuoteAmount.Sub(newQuote).Truncate(8)
	case pair.QuoteAssetID:
		newQuote := pair.QuoteAmount.Add(funds)
		newBase := K.Div(newQuote)
		r.FillAssetID = pair.BaseAssetID
		r.FillAmount = pair.BaseAmount.Sub(newBase).Truncate(8)
	default:
		return nil, errors.New("invalid pay asset id")
	}

	return r, nil
}

// reverseSwap is a Reverse version of Swap
func reverseSwap(pair *Pair, fillAssetID string, fillAmount decimal.Decimal) (*Result, error) {
	K := pair.BaseAmount.Mul(pair.QuoteAmount)
	if !K.IsPositive() {
		return nil, ErrInsufficientLiquiditySwapped
	}

	fillAmount = fillAmount.Truncate(8)
	if !fillAmount.IsPositive() {
		return nil, errors.New("invalid fill amount")
	}

	r := &Result{
		FillAssetID: fillAssetID,
		FillAmount:  fillAmount,
		RouteID:     pair.RouteID,
	}

	switch fillAssetID {
	case pair.BaseAssetID:
		newBase := pair.BaseAmount.Sub(fillAmount)
		if !newBase.IsPositive() {
			return nil, ErrInsufficientLiquiditySwapped
		}

		newQuote := K.Div(newBase)
		r.PayAssetID = pair.QuoteAssetID
		r.PayAmount = newQuote.Sub(pair.QuoteAmount)
	case pair.QuoteAssetID:
		newQuote := pair.QuoteAmount.Sub(fillAmount)
		if !newQuote.IsPositive() {
			return nil, ErrInsufficientLiquiditySwapped
		}

		newBase := K.Div(newQuote)
		r.PayAssetID = pair.BaseAssetID
		r.PayAmount = newBase.Sub(pair.BaseAmount)
	default:
		return nil, errors.New("invalid fill asset id")
	}

	r.PayAmount = r.PayAmount.Div(decimal.NewFromInt(1).Sub(Decimal(SwapFee)))
	r.PayAmount = Ceil(r.PayAmount, 8)
	r.FeeAssetID = r.PayAssetID
	r.FeeAmount = r.PayAmount.Mul(Decimal(SwapFee)).Truncate(8)

	return r, nil
}
