package clipper

import (
	"strconv"
)

type uint8Value uint8

func newUint8Value(val uint8, p *uint8) *uint8Value {
	*p = val
	return (*uint8Value)(p)
}

func (u *uint8Value) Set(s string, _ bool) error {
	v, err := strconv.ParseUint(s, 0, 8)
	if err == nil {
		*u = uint8Value(v)
	}
	return err
}

func (u *uint8Value) Reset(i interface{}) {
	v := i.(uint8)
	*u = uint8Value(v)
}

func (u *uint8Value) Type() string {
	return "uint8"
}

func (u *uint8Value) Get() interface{} {
	return u.GetUint8()
}

func (u *uint8Value) String() string { return strconv.FormatUint(uint64(*u), 10) }

func (u *uint8Value) GetUint8() uint8 { return uint8(*u) }

// AddUint8 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint8(name, shortName string, value uint8, p *uint8) *Opt {
	v := newUint8Value(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
