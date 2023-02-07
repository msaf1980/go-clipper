package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCounterFlag(t *testing.T) {
	var c int
	v := newCounterFlag(0, &c)

	require.Equal(t, 0, c)

	for i := 1; i < 4; i++ {
		err := v.Set("true", true)
		require.NoError(t, err)
		require.Equal(t, i, c)
	}

	err := v.Set("false", true)
	require.NoError(t, err)
	require.Equal(t, 2, c)
}
