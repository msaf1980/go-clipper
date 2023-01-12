package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUIntArrayReset(t *testing.T) {
	var sa []uint
	v := newUintArrayValue([]uint{2, 10}, &sa)

	require.Equal(t, sa, []uint{2, 10})

	v.Reset([]uint{1})
	require.Equal(t, sa, []uint{1})

	err := v.ResetWith([]uint{1, 5, 11, 0})
	require.NoError(t, err)
	require.Equal(t, sa, []uint{1, 5, 11, 0})

	require.Equal(t, []uint{1, 5, 11, 0}, v.Get())
}

func TestUintArray(t *testing.T) {
	var sa []uint
	v := newUintArrayValue([]uint{}, &sa)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	err := v.Set("1", true)
	require.NoError(t, err)
	require.Equal(t, sa, []uint{1})

	err = v.Set(`1,5,11`, true) // read from CSV
	require.NoError(t, err)
	require.Equal(t, sa, []uint{1, 1, 5, 11})

	v.Reset([]uint{2})
	require.NoError(t, err)
	require.Equal(t, sa, []uint{2})

	err = v.Set("-1", true)
	require.Error(t, err)
	require.Equal(t, sa, []uint{2})

	err = v.Set("a", true)
	require.Error(t, err)
	require.Equal(t, sa, []uint{2})
}
