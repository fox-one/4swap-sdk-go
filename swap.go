package fswap

import (
	"errors"

	"github.com/fox-one/4swap-sdk-go/v2/swap"
	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientLiquiditySwapped = errors.New("insufficient liquidity swapped")
)

// Result represent Swap Result
type Result struct {
	PayAssetID   string
	PayAmount    decimal.Decimal
	FillAssetID  string
	FillAmount   decimal.Decimal
	FeeAssetID   string
	FeeAmount    decimal.Decimal
	ProfitAmount decimal.Decimal
	RouteID      uint16
}

// Swap trade in a pair
func Swap(pair *Pair, payAssetID string, payAmount decimal.Decimal) (*Result, error) {
	m := swap.Imp(pair.SwapMethod)

	payAmount = payAmount.Truncate(8)

	r := &Result{
		PayAssetID:   payAssetID,
		PayAmount:    payAmount,
		FeeAssetID:   payAssetID,
		FeeAmount:    payAmount.Mul(pair.FeePercent).Truncate(8),
		ProfitAmount: payAmount.Mul(pair.ProfitRate).Truncate(8),
		RouteID:      pair.RouteID,
	}

	funds := payAmount.Sub(r.FeeAmount)
	if !funds.IsPositive() {
		return nil, errors.New("pay amount must be positive")
	}

	switch payAssetID {
	case pair.BaseAssetID:
		r.FillAssetID = pair.QuoteAssetID
		r.FillAmount = m.Swap(pair.BaseAmount, pair.QuoteAmount, funds).Truncate(8)
	case pair.QuoteAssetID:
		r.FillAssetID = pair.BaseAssetID
		r.FillAmount = m.Swap(pair.QuoteAmount, pair.BaseAmount, funds).Truncate(8)
	default:
		return nil, errors.New("invalid pay asset id")
	}

	return r, nil
}

// MultiHopSwap trade in a given path
func MultiHopSwap(pairs []*Pair, payAssetID string, payAmount decimal.Decimal) (result *Result, err error) {
	for _, pair := range pairs {
		result, err = Swap(pair, payAssetID, payAmount)
		if err != nil {
			return
		}
		payAssetID = result.FillAssetID
		payAmount = result.FillAmount
	}
	return
}

// ReverseSwap is a Reverse version of Swap
func ReverseSwap(pair *Pair, fillAssetID string, fillAmount decimal.Decimal) (*Result, error) {
	m := swap.Imp(pair.SwapMethod)

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
		r.PayAssetID = pair.QuoteAssetID
		r.PayAmount = m.Reverse(pair.QuoteAmount, pair.BaseAmount, fillAmount)
	case pair.QuoteAssetID:
		r.PayAssetID = pair.BaseAssetID
		r.PayAmount = m.Reverse(pair.BaseAmount, pair.QuoteAmount, fillAmount)
	default:
		return nil, errors.New("invalid fill asset id")
	}

	if !r.PayAmount.IsPositive() {
		return nil, ErrInsufficientLiquiditySwapped
	}

	r.PayAmount = r.PayAmount.Div(decimal.NewFromInt(1).Sub(pair.FeePercent))
	r.PayAmount = Ceil(r.PayAmount, 8)
	r.FeeAssetID = r.PayAssetID
	r.FeeAmount = r.PayAmount.Mul(pair.FeePercent).Truncate(8)
	r.ProfitAmount = r.PayAmount.Mul(pair.ProfitRate).Truncate(8)

	return r, nil
}

// ReverseMultiHopSwap is a Reverse version of MultiHopSwap
func ReverseMultiHopSwap(pairs []*Pair, fillAssetID string, fillAmount decimal.Decimal) (result *Result, err error) {
	size := len(pairs)
	for i := size - 1; i >= 0; i-- {
		result, err = ReverseSwap(pairs[i], fillAssetID, fillAmount)
		if err != nil {
			return
		}
		fillAssetID = result.PayAssetID
		fillAmount = result.PayAmount
	}
	return
}
