package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntArrayReset(t *testing.T) {
	var sa []int
	v := newIntArrayValue([]int{2, 10}, &sa)

	require.Equal(t, sa, []int{2, 10})

	v.Reset([]int{1})
	require.Equal(t, sa, []int{1})

	err := v.ResetWith([]int{1, 5, 11, 0, -1})
	require.NoError(t, err)
	require.Equal(t, sa, []int{1, 5, 11, 0, -1})

	require.Equal(t, []int{1, 5, 11, 0, -1}, v.Get())
}

func TestIntArray(t *testing.T) {
	var sa []int
	v := newIntArrayValue([]int{}, &sa)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	err := v.Set("1", true)
	require.NoError(t, err)
	require.Equal(t, sa, []int{1})

	err = v.Set(`1,5,11`, true) // read from CSV
	require.NoError(t, err)
	require.Equal(t, sa, []int{1, 1, 5, 11})

	v.Reset([]int{2})
	require.NoError(t, err)
	require.Equal(t, sa, []int{2})

	err = v.Set("a", true)
	require.Error(t, err)
	require.Equal(t, sa, []int{2})
}
