package clipper

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat64_0(t *testing.T) {
	var n float64
	v := newFloat64Value(0, &n)

	require.Equal(t, float64(0), n)

	err := v.Set("-1.48", true)
	require.NoError(t, err)
	require.Equal(t, float64(-1.48), n)
	require.Equal(t, v.GetFloat64(), n)

	err = v.Set("+inf", true)
	require.NoError(t, err)
	require.Equal(t, math.Inf(1), n)

	err = v.Set("-inf", true)
	require.NoError(t, err)
	require.Equal(t, math.Inf(-1), n)

	err = v.Set("inf", true)
	require.NoError(t, err)
	require.Equal(t, math.Inf(1), n)
}

func TestFloat64_12(t *testing.T) {
	var n float64
	v := newFloat64Value(float64(12), &n)

	require.Equal(t, float64(12), n)

	require.Equal(t, float64(12), v.GetFloat64())
	i := v.Get()
	require.Equal(t, float64(12), i)
}
