package uni

import (
	"github.com/shopspring/decimal"
)

type Swap struct{}

func (_ Swap) Swap(x, y decimal.Decimal, dx decimal.Decimal) decimal.Decimal {
	k := x.Mul(y)
	if !k.IsPositive() {
		return decimal.Zero
	}

	_x := x.Add(dx)
	_y := k.Div(_x)
	dy := y.Sub(_y)

	return dy
}

func (_ Swap) Reverse(x, y decimal.Decimal, dy decimal.Decimal) decimal.Decimal {
	k := x.Mul(y)
	if !k.IsPositive() {
		return decimal.Zero
	}

	_y := y.Sub(dy)
	if !_y.IsPositive() {
		return decimal.Zero
	}

	_x := k.Div(_y)
	dx := _x.Sub(x)

	return dx
}
