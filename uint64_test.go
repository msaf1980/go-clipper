package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint64_0(t *testing.T) {
	var n uint64
	v := newUint64Value(0, &n)

	require.Equal(t, uint64(0), n)

	err := v.Set("-1", true)
	require.Error(t, err)
	require.Equal(t, uint64(0), n)

	err = v.Set(strconv.FormatInt(math.MaxInt64, 10), true)
	require.NoError(t, err)
	require.Equal(t, uint64(math.MaxInt64), n)

	err = v.Set(strconv.FormatInt(math.MinInt64, 10), true)
	require.Error(t, err)
	require.Equal(t, uint64(math.MaxInt64), n)
}

func TestUint64_12(t *testing.T) {
	var n uint64
	v := newUint64Value(12, &n)

	require.Equal(t, uint64(12), n)

	require.Equal(t, uint64(12), v.GetUint64())
	i := v.Get()
	require.Equal(t, uint64(12), i)
}
