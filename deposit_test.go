package fswap

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddDeposit(t *testing.T) {
	const (
		baseAssetID  = "4d8c508b-91c5-375b-92b0-ee702ed2dac5"
		quoteAssetID = "815b0b1a-2764-3736-8faa-42d694fa620a"
		token        = `eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDAwNzg2NDksImlhdCI6MTYwMDA3ODA0OSwianRpIjoiZDFjOGE5NmUtNWFlNy00MzVlLWFiNTgtOWJkYjJjZGI4YWM3Iiwic2NwIjoiRlVMTCIsInNpZCI6ImRiMmYzMmJiLWYyYTUtNDJiMS1iOTQ2LTYzYTRlMTI5YjAyYyIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1YzRmMzBhNi0xZjQ5LTQzYzMtYjM3Yi1jMDFhYWU1MTkxYWYifQ.HQ_JCrpymNqcVPcZrvGE8yW8ekkt-_h7B2gIrGr-_FP1rk3E9Eh1x9EzItFAOrh5gPN32gtvf_2w4-mlf_6X9eCR9xT_rCdBmhnopCDWeMJaOYn9TIJpSQPVaCix_5feKJROXnzs8MpwFKNdIgUsACrWn7t7nv2nmtMmImrXWxY`
	)

	req := AddDepositReq{
		BaseAssetID:  baseAssetID,
		BaseAmount:   decimal.NewFromInt(100),
		QuoteAssetID: quoteAssetID,
		QuoteAmount:  decimal.NewFromInt(100),
		Slippage:     decimal.NewFromFloat(0.01),
	}

	ctx := context.Background()
	ctx = WithToken(ctx, token)
	deposit, err := AddDeposit(ctx, &req)
	require.Nil(t, err, "add deposit")
	assert.Len(t, deposit.Transfers, 2, "should have 2 transfers")
}
