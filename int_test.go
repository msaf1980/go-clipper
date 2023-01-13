package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt_0(t *testing.T) {
	var n int
	v := newIntValue(0, &n)

	require.Equal(t, 0, n)

	err := v.Set("-1", true)
	require.NoError(t, err)
	require.Equal(t, -1, n)

	err = v.Set(strconv.FormatInt(math.MaxInt32, 10), true)
	require.NoError(t, err)
	require.Equal(t, int(math.MaxInt32), n)

	err = v.Set(strconv.FormatInt(math.MinInt32, 10), true)
	require.NoError(t, err)
	require.Equal(t, int(math.MinInt32), n)
}

func TestInt_12(t *testing.T) {
	var n int
	v := newIntValue(12, &n)

	require.Equal(t, 12, n)

	require.Equal(t, int(12), v.GetInt())
	i := v.Get()
	require.Equal(t, int(12), i)

	v.Reset(int(12))
}
