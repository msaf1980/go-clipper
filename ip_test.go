package clipper

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIP(t *testing.T) {
	ip := net.IPv4(0, 0, 0, 0)
	v := newIPValue(ip, &ip)

	require.Equal(t, "0.0.0.0", ip.String())

	// IPv4
	err := v.Set("192.168.0.255", true)
	require.NoError(t, err)
	require.Equal(t, net.IPv4(192, 168, 0, 255), ip)
	require.Equal(t, net.IPv4(192, 168, 0, 255), v.Get())
	require.Equal(t, "192.168.0.255", v.String())

	err = v.Set("192.168.0.25a", true)
	require.Error(t, err)
	require.Equal(t, "192.168.0.255", v.String())

	// IPv6
	err = v.Set("::1", true)
	require.NoError(t, err)
	require.Equal(t, net.IPv6loopback, ip)
	require.Equal(t, net.IPv6loopback, v.Get())
	require.Equal(t, "::1", v.String())
}
