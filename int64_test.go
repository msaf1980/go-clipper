package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64_0(t *testing.T) {
	var n int64
	v := newInt64Value(0, &n)

	require.Equal(t, int64(0), n)

	err := v.Set("-1", true)
	require.NoError(t, err)
	require.Equal(t, int64(-1), n)

	err = v.Set(strconv.FormatInt(math.MaxInt, 10), true)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt), n)

	err = v.Set(strconv.FormatInt(math.MinInt, 10), true)
	require.NoError(t, err)
	require.Equal(t, int64(math.MinInt), n)

	err = v.Set("-1Ki", true)
	require.Error(t, err)
	require.Equal(t, int64(math.MinInt), n)
}

func TestInt64_12(t *testing.T) {
	var n int64
	v := newInt64Value(12, &n)

	require.Equal(t, int64(12), n)

	require.Equal(t, int64(12), v.GetInt64())
	i := v.Get()
	require.Equal(t, int64(12), i)

	v.Reset(int64(12))
}
