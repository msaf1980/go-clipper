package clipper

import (
	"net"
	"strings"
)

// -- net.IPArray Value
type ipArrayValue struct {
	value *[]net.IP
}

func newIPArrayValue(val []net.IP, p *[]net.IP) *ipArrayValue {
	ssv := new(ipArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newIPArrayValueFromCSV(val string, p *[]net.IP) *ipArrayValue {
	ssv := new(ipArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *ipArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]net.IP, 0, len(v))
	for _, s := range v {
		s = strings.TrimSpace(s)
		if ip := net.ParseIP(s); ip == nil {
			return ErrorInvalidValue{s, ErrIPParse}
		} else {
			iv = append(iv, ip)
		}
	}
	if doAppend {
		*s.value = append(*s.value, iv...)
	} else {
		*s.value = iv
	}
	return nil
}

func (s *ipArrayValue) Reset(i interface{}) {
	v := i.([]net.IP)
	*s.value = v
}

func (s *ipArrayValue) ResetWith(val []net.IP) error {
	out := make([]net.IP, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *ipArrayValue) GetIPArray() []net.IP {
	out := make([]net.IP, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *ipArrayValue) Type() string {
	return "net.IPArray"
}

func (s *ipArrayValue) Get() interface{} {
	return s.GetIPArray()
}

func (s *ipArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = n.String()
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddIPArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddIPArray(name, shortName string, value []net.IP, p *[]net.IP) *Opt {
	v := newIPArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}

// Addnet.IPArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddIPArrayFromCSV(name, shortName string, value string, p *[]net.IP) *Opt {
	v := newIPArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
