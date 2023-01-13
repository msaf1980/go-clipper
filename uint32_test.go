package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint32_0(t *testing.T) {
	var n uint32
	v := newUint32Value(0, &n)

	require.Equal(t, uint32(0), n)

	err := v.Set("-1", true)
	require.Error(t, err)
	require.Equal(t, uint32(0), n)

	err = v.Set(strconv.FormatInt(math.MaxInt32, 10), true)
	require.NoError(t, err)
	require.Equal(t, uint32(math.MaxInt32), n)

	err = v.Set(strconv.FormatInt(math.MinInt, 10), true)
	require.Error(t, err)
	require.Equal(t, uint32(math.MaxInt32), n)
}

func TestUint32_12(t *testing.T) {
	var n uint32
	v := newUint32Value(12, &n)

	require.Equal(t, uint32(12), n)

	require.Equal(t, uint32(12), v.GetUint32())
	i := v.Get()
	require.Equal(t, uint32(12), i)

	v.Reset(uint32(12))
}
