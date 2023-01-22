package clipper

import (
	"strconv"
)

type int64NValue int64

func newInt64NValue(val int64, p *int64) *int64NValue {
	*p = val
	return (*int64NValue)(p)
}

func newInt64NValueFromString(val string, p *int64) *int64NValue {
	v, err := int64NFromString(val)
	if err != nil {
		panic(err)
	}
	*p = v
	return (*int64NValue)(p)
}

func int64NFromString(s string) (n int64, err error) {
	last := len(s) - 1
	if last < 1 {
		return strconv.ParseInt(s, 0, 64)
	}
	switch s[last] {
	case 'k': // SI is 1e3
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e3
	case 'm': // SI is 1e6
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e6
	case 'g': // SI is 1e9
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e9
	case 't': // SI is 1e12
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1e12
	case 'K':
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1024
	case 'M':
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1048576
	case 'G':
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
			return n, err
		}
		n *= 1073741824
	case 'T':
		if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
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
			if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1000
		case 'M': // (Mi) SI is 1e6
			if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e6
		case 'G': // (Gi) SI is 1e9
			if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e9
		case 'T': // (Ti) SI is 1e12
			if n, err = strconv.ParseInt(s[:last], 0, 64); err != nil {
				return n, err
			}
			n *= 1e12
		default:
			n, err = strconv.ParseInt(s, 0, 64)
		}
	default:
		n, err = strconv.ParseInt(s, 0, 64)
	}
	return
}

func (i *int64NValue) Set(s string, _ bool) error {
	v, err := int64NFromString(s)
	if err == nil {
		*i = int64NValue(v)
	}
	return err
}

func (iv *int64NValue) Reset(i interface{}) {
	v := i.(int64)
	*iv = int64NValue(v)
}

func (i *int64NValue) Type() string {
	return "int64N"
}

func (i *int64NValue) Get() interface{} {
	return i.GetInt64()
}

func (i *int64NValue) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int64NValue) GetInt64() int64 { return int64(*i) }

// AddInt64 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt64N(name, shortName string, value int64, p *int64, help string) *Opt {
	v := newInt64NValue(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}

// AddInt64N registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// The `value` argument represents initial value, with suffix of k (1e3), m (1e6), g (1e9), K (1024), M (1048576), G (1073741824)
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt64NFromString(name, shortName string, value string, p *int64, help string) *Opt {
	v := newInt64NValueFromString(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
