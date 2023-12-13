package swap

import (
	"log"

	"github.com/fox-one/4swap-sdk-go/v2/swap/curve"
	"github.com/fox-one/4swap-sdk-go/v2/swap/uni"
	"github.com/shopspring/decimal"
)

const (
	MethodUni   = "uni"
	MethodCurve = "curve"
)

var registered map[string]Interface

func init() {
	registered = map[string]Interface{
		MethodUni:   &uni.Swap{},
		MethodCurve: &curve.Swap{A: decimal.NewFromInt(200)},
	}
}

type Interface interface {
	Swap(x, y decimal.Decimal, dx decimal.Decimal) decimal.Decimal
	Reverse(x, y decimal.Decimal, dy decimal.Decimal) decimal.Decimal
}

func Imp(method string) Interface {
	if method == "" {
		method = MethodUni
	}

	imp, ok := registered[method]
	if !ok {
		log.Fatalf("unknown swap method %q", method)
	}

	return imp
}
