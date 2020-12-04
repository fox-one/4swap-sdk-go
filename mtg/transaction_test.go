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

func TestReadTransaction(t *testing.T) {
	const (
		followID = "45e9fa8f-24fc-4e50-b522-9aa1a4dd4a43"
	)

	ctx := context.Background()
	ctx = WithToken(ctx, token)
	tx, err := ReadTransaction(ctx, baseAssetID, quoteAssetID, followID)
	require.Nil(t, err, "request should be ok")

	t.Log(tx.PayAssetID, tx.Funds)
	t.Log(tx.FilledAssetID, tx.Amount)
	t.Log(tx.FeeAssetID, tx.FeeAmount)
}
