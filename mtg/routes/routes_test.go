package routes

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoutesEncoding(t *testing.T) {
	r := Routes{18, 124, 3, 4}

	t.Run("json", func(t *testing.T) {
		b, err := r.MarshalJSON()
		require.Nil(t, err)

		t.Log(len(b), string(b))
	})

	t.Run("binary", func(t *testing.T) {
		b, err := r.MarshalBinary()
		require.Nil(t, err)

		t.Log(len(b), base64.StdEncoding.EncodeToString(b))
	})
}
