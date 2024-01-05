package fswap

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAsset(t *testing.T) {
	const (
		assetID = "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	)

	ctx := context.Background()
	c := New()

	asset, err := c.ReadAsset(ctx, assetID)
	require.Nil(t, err, "request should be ok")
	require.True(t, asset.Price.GreaterThan(decimal.Zero), "price should be greater than zero")
	require.True(t, len(asset.Chain.Symbol) > 0, "require chain symbol exists")

	t.Log(asset.Symbol, asset.Price)
}

func TestListAssets(t *testing.T) {
	ctx := context.Background()
	c := New()

	assets, err := c.ListAssets(ctx)
	require.Nil(t, err, "request should be ok")
	assert.NotEmpty(t, assets, "assets should not be empty")
}
