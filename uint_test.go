package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint_0(t *testing.T) {
	var n uint
	v := newUintValue(0, &n)

	require.Equal(t, uint(0), n)

	err := v.Set("-1", true)
	require.Error(t, err)
	require.Equal(t, uint(0), n)

	err = v.Set(strconv.FormatInt(math.MaxInt, 10), true)
	require.NoError(t, err)
	require.Equal(t, uint(math.MaxInt), n)

	err = v.Set(strconv.FormatInt(math.MinInt, 10), true)
	require.Error(t, err)
	require.Equal(t, uint(math.MaxInt), n)
}

func TestUint_12(t *testing.T) {
	var n uint
	v := newUintValue(12, &n)

	require.Equal(t, uint(12), n)

	require.Equal(t, uint(12), v.GetUint())
	i := v.Get()
	require.Equal(t, uint(12), i)
}
