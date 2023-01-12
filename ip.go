package clipper

import (
	"net"
	"strings"
)

// -- net.IP value
type ipValue net.IP

func newIPValue(val net.IP, p *net.IP) *ipValue {
	*p = val
	return (*ipValue)(p)
}

func newIPValueFromString(val string, p *net.IP) *ipValue {
	val = strings.TrimSpace(val)
	ip := net.ParseIP(val)
	if ip == nil {
		panic(ErrorInvalidValue{val, ErrIPParse})
	}
	*p = ip
	return (*ipValue)(p)
}

func (i *ipValue) String() string { return net.IP(*i).String() }

func (i *ipValue) Set(s string, _ bool) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return ErrorInvalidValue{s, ErrIPParse}
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return ErrorInvalidValue{s, ErrIPParse}
	}
	*i = ipValue(ip)
	return nil
}

func (i *ipValue) Reset(p interface{}) {
	v := p.(net.IP)
	*i = ipValue(v)
}

func (i *ipValue) Get() interface{} {
	return i.GetIP()
}

func (i *ipValue) GetIP() net.IP {
	ip := net.IP(*i)
	out := make(net.IP, len(ip))
	copy(out, ip)
	return ip
}

func (i *ipValue) Type() string {
	return "ip"
}

// AddIP registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddIP(name, shortName string, value net.IP, p *net.IP) *Opt {
	v := newIPValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
