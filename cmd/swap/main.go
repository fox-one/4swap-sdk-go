package main

import (
	"fmt"
	"os"

	"github.com/fox-one/4swap-sdk-go/swap"
	"github.com/shopspring/decimal"
)

func main() {
	args := os.Args[1:]
	method := args[0]
	imp := swap.Imp(method)

	x, _ := decimal.NewFromString(args[1])
	y, _ := decimal.NewFromString(args[2])

	dx, _ := decimal.NewFromString(args[3])
	dy := imp.Swap(x, y, dx)
	_, _ = fmt.Fprint(os.Stdout, dy.String())
}
