package clipper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	var d time.Duration
	v := newDurationValue(0, &d)

	require.Equal(t, time.Duration(0), d)

	err := v.Set("1s", true)
	require.NoError(t, err)
	require.Equal(t, time.Second, d)
}

func TestDuration_1s(t *testing.T) {
	var d time.Duration
	v := newDurationValue(time.Second, &d)

	require.Equal(t, time.Second, d)

	require.Equal(t, time.Second, v.GetDuration())
	i := v.Get()
	require.Equal(t, time.Second, i)
}
