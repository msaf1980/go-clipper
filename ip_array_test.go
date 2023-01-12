package clipper

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPArray(t *testing.T) {
	var sa []net.IP
	v := newIPArrayValue([]net.IP{}, &sa)

	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}

	err := v.Set("::1", true)
	require.NoError(t, err)
	require.Equal(t, sa, []net.IP{net.IPv6loopback})

	err = v.Set("192.168.0.1,::1", true)
	require.NoError(t, err)
	require.Equal(t, sa, []net.IP{
		net.IPv6loopback,
		net.IPv4(192, 168, 0, 1),
		net.IPv6loopback,
	})

	err = v.Set("a", true)
	require.Error(t, err)
	require.Equal(t, sa, []net.IP{
		net.IPv6loopback,
		net.IPv4(192, 168, 0, 1),
		net.IPv6loopback,
	})

	v.Reset([]net.IP{})
	if len(sa) != 0 {
		t.Fatalf("%#v must be empty", sa)
	}
}
