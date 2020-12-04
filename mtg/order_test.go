package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadOrder(t *testing.T) {
	const (
		followID = "2212be57-b730-4c3b-a50a-aea1a27b2758"
	)

	ctx := context.Background()
	ctx = WithToken(ctx, token)
	order, err := ReadOrder(ctx, followID)
	require.Nil(t, err, "request should be ok")

	t.Log(order.PayAssetID, order.PayAmount)
	t.Log(order.FillAssetID, order.FillAmount)
}
