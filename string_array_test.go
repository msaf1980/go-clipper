package clipper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArrayReset(t *testing.T) {
	var sa []string
	v := newStringArrayValue([]string{"a", "c"}, &sa)

	require.Equal(t, sa, []string{"a", "c"})

	v.Reset([]string{"b"})
	require.Equal(t, sa, []string{"b"})

	err := v.ResetWith([]string{"a", "b", "cd", "e, j", " f"})
	require.NoError(t, err)
	require.Equal(t, sa, []string{"a", "b", "cd", "e, j", " f"})

	require.Equal(t, []string{"a", "b", "cd", "e, j", " f"}, v.Get())
}

func TestStringArray(t *testing.T) {
	var sa []string
	v := newStringArrayValue([]string{}, &sa)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	err := v.Set("a", true)
	require.NoError(t, err)
	require.Equal(t, sa, []string{"a"})

	err = v.Set(`b,cd,"e, j", f`, true) // read from CSV
	require.NoError(t, err)
	require.Equal(t, sa, []string{"a", "b", "cd", "e, j", " f"})

	v.Reset([]string{"z"})
	require.NoError(t, err)
	require.Equal(t, sa, []string{"z"})
}

func TestStringArrayCSV(t *testing.T) {
	var sa []string
	newStringArrayValueFromCSV("", &sa)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	newStringArrayValueFromCSV(`b,cd,"e, j", f`, &sa)
	require.Equal(t, sa, []string{"b", "cd", "e, j", " f"})
}

func TestStringArrayWithLimit(t *testing.T) {
	var sa []string
	v := newStringArrayLValue([]string{}, &sa, 1)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	err := v.Set("a", true)
	require.NoError(t, err)
	require.Equal(t, sa, []string{"a"})

	err = v.Set(`b,cd,"e, j", f`, true)
	errIs := ErrorLengthOverflow{`length at argument=b,cd,"e, j", f`, ">", 1}
	require.ErrorIs(t, err, errIs)
	require.Equal(t, sa, []string{"a"})

	err = v.Replace([]string{"z"})
	require.NoError(t, err)
	require.Equal(t, sa, []string{"z"})

	err = v.Replace([]string{"b", "cd"})
	errIs = ErrorLengthOverflow{"length", ">", 1}
	require.ErrorIs(t, errIs, err)
	require.Equal(t, sa, []string{"z"})

	// check minimum
	v.SetMinLen(1)

	err = v.Replace([]string{})
	require.NoError(t, err)
	require.Equal(t, sa, []string{})
	err = v.CheckLen()
	errIs = ErrorLengthOverflow{"length", "<", 1}
	require.ErrorIs(t, err, errIs)
}
