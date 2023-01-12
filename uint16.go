package clipper

import (
	"strconv"
)

type uint16Value uint16

func newUint16Value(val uint16, p *uint16) *uint16Value {
	*p = val
	return (*uint16Value)(p)
}

func (u *uint16Value) Set(s string, _ bool) error {
	v, err := strconv.ParseUint(s, 0, 16)
	if err == nil {
		*u = uint16Value(v)
	}
	return err
}

func (u *uint16Value) Reset(i interface{}) {
	v := i.(uint16)
	*u = uint16Value(v)
}

func (u *uint16Value) Type() string {
	return "uint16"
}

func (u *uint16Value) Get() interface{} {
	return u.GetUint16()
}

func (u *uint16Value) String() string { return strconv.FormatUint(uint64(*u), 10) }

func (u *uint16Value) GetUint16() uint16 { return uint16(*u) }

// AddUint16 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint16(name, shortName string, value uint16, p *uint16, help string) *Opt {
	v := newUint16Value(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
