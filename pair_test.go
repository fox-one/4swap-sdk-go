package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadPair(t *testing.T) {
	const (
		baseAssetID  = "4d8c508b-91c5-375b-92b0-ee702ed2dac5"
		quoteAssetID = "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	)

	ctx := context.Background()
	c := New()
	c.Resty().Debug = true

	pair, err := c.ReadPair(ctx, baseAssetID, quoteAssetID)
	require.Nil(t, err, "request should be ok")

	t.Log(pair.BaseAssetID, pair.BaseAmount)
	t.Log(pair.QuoteAssetID, pair.QuoteAmount)
	t.Log(pair.FeePercent, pair.Liquidity)
}

func TestListPairs(t *testing.T) {
	ctx := context.Background()
	c := New()

	pairs, err := c.ListPairs(ctx)
	require.Nil(t, err, "request should be ok")
	assert.NotEmpty(t, pairs, "pairs should not be empty")
}
