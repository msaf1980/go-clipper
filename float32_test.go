package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFloat32_0(t *testing.T) {
	var n float32
	v := newFloat32Value(0, &n)

	require.Equal(t, float32(0), n)

	err := v.Set("-1.48", true)
	require.NoError(t, err)
	require.Equal(t, float32(-1.48), n)
	require.Equal(t, v.GetFloat32(), n)

	err = v.Set(strconv.FormatFloat(math.MaxFloat32, 'g', -1, 32), true)
	require.NoError(t, err)
	require.Equal(t, float32(math.MaxFloat32), n)

	err = v.Set(strconv.FormatFloat(math.MaxFloat32+1.0, 'g', -1, 32), true)
	require.NoError(t, err)
	require.Equal(t, float32(math.MaxFloat32), n)

	err = v.Set("+inf", true)
	require.NoError(t, err)
	require.Equal(t, float32(math.Inf(1)), n)

	err = v.Set("-inf", true)
	require.NoError(t, err)
	require.Equal(t, float32(math.Inf(-1)), n)

	err = v.Set("inf", true)
	require.NoError(t, err)
	require.Equal(t, float32(math.Inf(1)), n)
}

func TestFloat32_12(t *testing.T) {
	var n float32
	v := newFloat32Value(float32(12), &n)

	require.Equal(t, float32(12), n)

	require.Equal(t, float32(12), v.GetFloat32())
	i := v.Get()
	require.Equal(t, float32(12), i)
}
