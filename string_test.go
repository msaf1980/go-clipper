package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	var s string
	v := newStringValue("", &s)

	require.Equal(t, "", s)

	err := v.Set("true", true)
	require.NoError(t, err)
	require.Equal(t, "true", s)

	err = v.Set("false", true)
	require.NoError(t, err)
	require.Equal(t, "false", s)

	require.Equal(t, "false", v.Get())
}

func TestStringNonEmpty(t *testing.T) {
	s := "a"
	v := newStringValue("s", &s)

	require.Equal(t, "s", s)

	err := v.Set("true", true)
	require.NoError(t, err)
	require.Equal(t, "true", s)
}
