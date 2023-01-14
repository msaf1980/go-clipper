package clipper

import (
	"strconv"
)

type uint32Value uint32

func newUint32Value(val uint32, p *uint32) *uint32Value {
	*p = val
	return (*uint32Value)(p)
}

func (u *uint32Value) Set(s string, _ bool) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err == nil {
		*u = uint32Value(v)
	}
	return err
}

func (u *uint32Value) Reset(i interface{}) {
	v := i.(uint32)
	*u = uint32Value(v)
}

func (u *uint32Value) Type() string {
	return "uint32"
}

func (u *uint32Value) Get() interface{} {
	return u.GetUint32()
}

func (u *uint32Value) String() string { return strconv.FormatUint(uint64(*u), 10) }

func (u *uint32Value) GetUint32() uint32 { return uint32(*u) }

// AddUint32 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint32(name, shortName string, value uint32, p *uint32, help string) *Opt {
	v := newUint32Value(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
