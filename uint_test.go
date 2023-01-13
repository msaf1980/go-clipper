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

	err = v.Set("a", true)
	require.Error(t, err)
	require.Equal(t, uint(0), n)

	err = v.Set(strconv.FormatUint(math.MaxUint32, 10), true)
	require.NoError(t, err)
	require.Equal(t, uint(math.MaxUint32), n)
}

func TestUint_12(t *testing.T) {
	var n uint
	v := newUintValue(12, &n)

	require.Equal(t, uint(12), n)

	require.Equal(t, uint(12), v.GetUint())
	i := v.Get()
	require.Equal(t, uint(12), i)

	v.Reset(uint(12))
}
