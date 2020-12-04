package fswap

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseAssetID      = "4d8c508b-91c5-375b-92b0-ee702ed2dac5"
	quoteAssetID     = "c94ac88f-4671-3976-b60a-09064f1811e8"
	liquidityAssetID = "00608a54-c563-3e67-8312-80f4471219be"
	userID           = "318df485-02e1-3c10-8ffd-b241d10dcfd3"
	token            = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDcwNTM2MzMsImlhdCI6MTYwNzA1MDAzMywianRpIjoiZDZmN2Q5Y2MtMzAwNi00YTZmLTk3NTMtZTNlY2FiNjU3YzEyIiwic2NwIjoiRlVMTCIsInNpZCI6IjU4M2Y4ZDljLTUzY2QtNDg2Zi04NzNmLTQ2ZDU3MDYyNTA4YiIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiIzMThkZjQ4NS0wMmUxLTNjMTAtOGZmZC1iMjQxZDEwZGNmZDMifQ.hzsFOnVdgLsTJ3X8YOIHLnRpTlEx42C5aIJhXhF4OthtAp9kT7pZmZPUWNYgghFUhpSnoaNYVacJmYtEA6mONozyz4yHx4mFsn4zvaYNUY-09IkaZJqnow3ig2Lwgb9z1kN-1TV7941ohWwW7Nq-HLuAHLNXZWZ7kT5Zi7dKlbk"
)

func TestAddDeposit(t *testing.T) {
	ctx := context.Background()
	pair, err := ReadPair(ctx, baseAssetID, quoteAssetID)
	require.Nil(t, err, "read pair")

	req := AddLiquidityReq{
		UserID:      userID,
		Pair:        pair,
		BaseAmount:  pair.BaseAmount.Div(pair.BaseValue).Div(decimal.NewFromInt(100)).Truncate(8),
		QuoteAmount: pair.QuoteAmount.Div(pair.BaseValue).Div(decimal.NewFromInt(100)).Truncate(8),
		Slippage:    decimal.NewFromFloat(0.01),
	}

	deposit, err := AddLiquidity(ctx, &req)
	require.Nil(t, err, "add deposit")
	assert.Len(t, deposit.Transfers, 2, "should have 2 transfers")
}

func TestSwap(t *testing.T) {
	ctx := context.Background()
	pairs, err := ListPairs(ctx)
	require.Nil(t, err, "list pairs")

	order, err := Route(pairs, baseAssetID, quoteAssetID, decimal.NewFromFloat(0.0001))
	require.Nil(t, err, "route")

	order.UserID = userID
	action, err := Swap(ctx, order)
	require.Nil(t, err, "swap")
	assert.Len(t, action.Transfers, 1, "should have 1 transfers")
}

func TestRemoveLiquidity(t *testing.T) {
	req := RemoveLiquidityReq{
		UserID:           userID,
		LiquidityAssetID: liquidityAssetID,
		Amount:           decimal.NewFromFloat(0.0001),
	}

	ctx := context.Background()
	action, err := RemoveLiquidity(ctx, &req)
	require.Nil(t, err, "remove liquidity")
	assert.Len(t, action.Transfers, 1, "should have 1 transfers")
}
