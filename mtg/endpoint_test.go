package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwitchEndpoint(t *testing.T) {
	const endpoint = "https://swap-mtg-test-api.fox.one"
	ctx := context.Background()
	err := UseEndpoint(ctx, endpoint)
	require.Nil(t, err, "switch endpoint")
}
