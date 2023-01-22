package clipper

import (
	"strconv"
)

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func newUint64NValue(val string, p *uint64) *uint64Value {
	v, err := uint64NFromString(val)
	if err != nil {
		panic(err)
	}
	*p = v
	return (*uint64Value)(p)
}

func uint64NFromString(s string) (n uint64, err error) {
	last := len(s) - 1
	if last < 1 {
		return strconv.ParseUint(s, 0, 64)
	}
	switch s[last] {
	case 'k': // SI is 1e3
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e3
	case 'm': // SI is 1e6
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e6
	case 'g': // SI is 1e9
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e9
	case 't': // SI is 1e12
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e12
	case 'K':
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1024
	case 'M':
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1048576
	case 'G':
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1073741824
	case 'T':
		if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1099511627776
	case 'i':
		if len(s) < 3 {
			strconv.ParseInt(s, 0, 64)
		}
		last--
		switch s[last] {
		case 'K': // (Ki) SI is 1e3
			if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1000
		case 'M': // (Mi) SI is 1e6
			if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e6
		case 'G': // (Gi) SI is 1e9
			if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e9
		case 'T': // (Ti) SI is 1e12
			if n, err = strconv.ParseUint(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e12
		default:
			n, err = strconv.ParseUint(s, 0, 64)
		}
	default:
		n, err = strconv.ParseUint(s, 0, 64)
	}
	return
}

func (u *uint64Value) Set(s string, _ bool) error {
	v, err := uint64NFromString(s)
	if err == nil {
		*u = uint64Value(v)
	}
	return err
}

func (u *uint64Value) Reset(i interface{}) {
	v := i.(uint64)
	*u = uint64Value(v)
}

func (u *uint64Value) Type() string {
	return "uint64"
}

func (u *uint64Value) Get() interface{} {
	return u.GetUint64()
}

func (u *uint64Value) String() string { return strconv.FormatUint(uint64(*u), 10) }

func (u *uint64Value) GetUint64() uint64 { return uint64(*u) }

// AddUint64 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64(name, shortName string, value uint64, p *uint64, help string) *Opt {
	v := newUint64Value(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}

// AddUint64N registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// The `value` argument represents initial value, with suffix of k (1e3), m (1e6), g (1e9), K (1024), M (1048576), G (1073741824)
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64N(name, shortName string, value string, p *uint64, help string) *Opt {
	v := newUint64NValue(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
