package fswap

import (
	"github.com/shopspring/decimal"
)

func Decimal(v string) decimal.Decimal {
	d, _ := decimal.NewFromString(v)
	return d
}

func Ceil(d decimal.Decimal, precision int32) decimal.Decimal {
	return d.Shift(precision).Ceil().Shift(-precision)
}
