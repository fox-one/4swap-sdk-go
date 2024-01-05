package fswap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGroup(t *testing.T) {
	ctx := context.Background()
	c := New()

	group, err := c.ReadGroup(ctx)
	require.Nil(t, err, "read group")
	assert.NotEmpty(t, group.Members, "receivers should not be empty")
}
