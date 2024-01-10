package fswap

import (
	"context"
	"testing"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPreOrder(t *testing.T) {
	ctx := context.Background()

	c := New()
	c.UseToken("your auth token")

	pairs, err := c.ListPairs(ctx)
	if err != nil {
		t.Fatal(err)
	}

	req := &PreOrderReq{
		PayAssetID:  "4d8c508b-91c5-375b-92b0-ee702ed2dac5",
		FillAssetID: "31d2ea9c-95eb-3355-b65b-ba096853bc18",
		PayAmount:   decimal.NewFromFloat(0.1),
	}

	preOrder, err := PreOrderWithPairs(pairs, req)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fill amount: %s", preOrder.FillAmount)

	followID := uuid.NewString()
	minAmount := preOrder.FillAmount.Mul(decimal.NewFromFloat(0.99)).Truncate(8)
	memo := BuildSwap(followID, req.FillAssetID, preOrder.Paths, minAmount)

	t.Logf("memo: %s", memo)

	group, err := c.ReadGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("target mix address: %s", group.MixAddress)

	transfer := &mixin.TransferInput{
		AssetID:    req.PayAssetID,
		OpponentID: group.MixAddress,
		Amount:     req.PayAmount,
		TraceID:    followID,
		Memo:       memo,
	}

	t.Log(mixin.URL.SafePay(transfer))

	// transfer pay asset to mix address

	// view order detail
	order, err := c.ReadOrder(ctx, followID)
	if err != nil {
		assert.True(t, IsErrorCode(err, 401))
	} else {
		t.Logf("order state: %s", order.State)
	}
}
