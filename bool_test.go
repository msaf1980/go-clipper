package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoolFalse(t *testing.T) {
	var b bool
	v := newBoolValue(false, &b)

	require.False(t, b)

	err := v.Set("true", true)
	require.NoError(t, err)
	require.True(t, b)

	err = v.Set("false", true)
	require.NoError(t, err)
	require.False(t, b)
}

func TestBoolTrue(t *testing.T) {
	var b bool
	v := newBoolValue(true, &b)

	require.True(t, b)

	err := v.Set("false", true)
	require.NoError(t, err)
	require.False(t, b)
}

func TestBoolReset(t *testing.T) {
	var b bool
	v := newBoolValue(true, &b)

	require.True(t, b)

	var i interface{} = false

	v.Reset(i)
	require.False(t, b)

	require.False(t, v.GetBool())
	i = v.Get()
	require.Equal(t, false, i)
}
