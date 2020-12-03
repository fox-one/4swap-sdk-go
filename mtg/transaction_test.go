package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListTransactions(t *testing.T) {
	ctx := context.Background()
	req := ListTransactionsReq{
		BaseAssetID:  baseAssetID,
		QuoteAssetID: quoteAssetID,
		Type:         TransactionTypeSwap,
	}
	_, err := ListTransactions(ctx, &req)
	require.Nil(t, err, "request should be ok")
}
