package clipper

import (
	"strconv"
)

type uint64Value uint64

func newUint64Value(val uint64, p *uint64) *uint64Value {
	*p = val
	return (*uint64Value)(p)
}

func (u *uint64Value) Set(s string, _ bool) error {
	v, err := strconv.ParseUint(s, 0, 64)
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
	return commandConfig.AddValue(name, shortName, v, help)
}
