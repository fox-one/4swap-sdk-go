package curve

import (
	"math"

	"github.com/bradfitz/iter"
	"github.com/shopspring/decimal"
)

var (
	nCoins          = decimal.NewFromInt(2)
	one             = decimal.NewFromInt(1)
	two             = decimal.NewFromInt(2)
	precision int32 = 8
)

type Swap struct {
	// Amplification coefficient
	A decimal.Decimal
}

func (m Swap) Swap(x, y decimal.Decimal, dx decimal.Decimal) decimal.Decimal {
	x = x.Shift(precision)
	y = y.Shift(precision)
	dx = dx.Shift(precision)

	dy := m.exchange(x, y, dx)
	if !dy.IsPositive() {
		return decimal.Zero
	}

	dy = dy.Shift(-precision)
	return dy
}

func (m Swap) Reverse(x, y decimal.Decimal, dy decimal.Decimal) decimal.Decimal {
	if y.LessThanOrEqual(dy) {
		return decimal.Zero
	}

	x = x.Shift(precision)
	y = y.Shift(precision)
	dy = dy.Shift(precision)

	dx := m.reverseExchange(x, y, dy)
	if !dx.IsPositive() {
		return decimal.Zero
	}

	dx = dx.Shift(-precision)
	return dx
}

func (m Swap) getD(xp []decimal.Decimal) decimal.Decimal {
	sum := decimal.Sum(decimal.Zero, xp[0:]...)
	if !sum.IsPositive() {
		return decimal.Zero
	}

	dp := decimal.Zero
	d := sum
	ann := m.A.Mul(nCoins)

	for range iter.N(255) {
		_dp := d
		for _, _x := range xp {
			_dp = _dp.Mul(d).Div(_x.Mul(nCoins).Add(one))
		}

		dp = d
		d1 := ann.Sub(one).Mul(d)
		d2 := nCoins.Add(one).Mul(_dp)
		d = ann.Mul(sum).Add(_dp.Mul(nCoins)).Mul(d).Div(d1.Add(d2))

		if d.Sub(dp).Truncate(0).IsZero() {
			break
		}
	}

	return d.Truncate(0)
}

func (m Swap) getY(d decimal.Decimal, x decimal.Decimal) decimal.Decimal {
	ann := m.A.Mul(nCoins)

	c := d.Mul(d).Div(x.Mul(nCoins))
	c = c.Mul(d).Div(ann.Mul(nCoins))

	b := x.Add(d.Div(ann))

	yp := decimal.Zero
	y := d

	for range iter.N(255) {
		yp = y
		y = y.Mul(y).Add(c).Div(y.Add(y).Add(b).Sub(d))

		if y.Sub(yp).Truncate(0).IsZero() {
			break
		}
	}

	return y
}

// reverse getY
func (m Swap) getX(d decimal.Decimal, y decimal.Decimal) decimal.Decimal {
	ann := m.A.Mul(nCoins)
	k := d.Mul(d).Mul(d).Div(ann).Div(nCoins).Div(nCoins)
	j := d.Div(ann).Sub(d).Add(y).Add(y)
	n := y.Sub(j).Div(two)
	x := sqrt(k.Div(y).Add(n.Mul(n))).Add(n)
	return x
}

func (m Swap) exchange(x, y, dx decimal.Decimal) decimal.Decimal {
	xp := []decimal.Decimal{x, y}
	_x := x.Add(dx)
	d := m.getD(xp)
	_y := m.getY(d, _x)
	dy := y.Sub(_y)
	return dy
}

func (m Swap) reverseExchange(x, y, dy decimal.Decimal) decimal.Decimal {
	xp := []decimal.Decimal{x, y}
	_y := y.Sub(dy)
	d := m.getD(xp)
	_x := m.getX(d, _y)
	dx := _x.Sub(x)
	return dx
}

func sqrt(d decimal.Decimal) decimal.Decimal {
	f, _ := d.Float64()
	f = math.Sqrt(f)
	return decimal.NewFromFloat(f)
}
