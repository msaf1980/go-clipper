package clipper

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt64N_0(t *testing.T) {
	var n int64
	v := newInt64NValue(0, &n)

	require.Equal(t, int64(0), n)

	err := v.Set("-1", true)
	require.NoError(t, err)
	require.Equal(t, int64(-1), n)

	err = v.Set(strconv.FormatInt(math.MaxInt, 10), true)
	require.NoError(t, err)
	require.Equal(t, int64(math.MaxInt), n)

	err = v.Set(strconv.FormatInt(math.MinInt, 10), true)
	require.NoError(t, err)
	require.Equal(t, int64(math.MinInt), n)

	err = v.Set("-1Ki", true)
	require.NoError(t, err)
	require.Equal(t, int64(-1000), n)
}

func TestInt64N_12(t *testing.T) {
	var n int64
	v := newInt64NValue(12, &n)

	require.Equal(t, int64(12), n)

	require.Equal(t, int64(12), v.GetInt64())
	i := v.Get()
	require.Equal(t, int64(12), i)

	v.Reset(int64(12))
}

func Test_int64NFromString(t *testing.T) {
	tests := []struct {
		val     string
		wantN   int64
		wantErr bool
	}{
		{val: "", wantN: 0, wantErr: true},
		{val: "k", wantN: 0, wantErr: true},
		{val: "-k", wantN: 0, wantErr: true},
		{val: "1m1k", wantN: 0, wantErr: true},
		{val: "-1k", wantN: -1e3},
		{val: "1k", wantN: 1e3},
		{val: "1Ki", wantN: 1e3},
		{val: "1K", wantN: 1024},
		{val: "1m", wantN: 1e6},
		{val: "1Mi", wantN: 1e6},
		{val: "1M", wantN: 1048576},
		{val: "1g", wantN: 1e9},
		{val: "1Gi", wantN: 1e9},
		{val: "1G", wantN: 1073741824},
		{val: "1t", wantN: 1e12},
		{val: "1Ti", wantN: 1e12},
		{val: "12Ti", wantN: 12 * 1e12},
		{val: "10T", wantN: 10995116277760},
	}
	for _, tt := range tests {
		t.Run(tt.val, func(t *testing.T) {
			gotN, err := int64NFromString(tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("int64NFromString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotN != tt.wantN {
				t.Errorf("int64NFromString() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
