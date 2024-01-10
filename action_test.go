package fswap

import (
	"testing"

	"github.com/fox-one/4swap-sdk-go/v2/route"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBuildSwapMemo(t *testing.T) {
	fillAssetID := "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	memo := BuildSwap(uuid.NewString(), fillAssetID, route.Single(1, 2, 3), Decimal("1.2"))
	assert.NotEmpty(t, memo)
}
