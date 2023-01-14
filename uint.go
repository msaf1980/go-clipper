package clipper

import (
	"strconv"
)

type uintValue uint

func newUintValue(val uint, p *uint) *uintValue {
	*p = val
	return (*uintValue)(p)
}

func (u *uintValue) Set(s string, _ bool) error {
	v, err := strconv.ParseUint(s, 0, 32)
	if err == nil {
		*u = uintValue(v)
	}
	return err
}
func (u *uintValue) Reset(i interface{}) {
	v := i.(uint)
	*u = uintValue(v)
}

func (u *uintValue) Type() string {
	return "uint"
}

func (u *uintValue) Get() interface{} {
	return u.GetUint()
}

func (u *uintValue) String() string { return strconv.FormatUint(uint64(*u), 10) }

func (u *uintValue) GetUint() uint { return uint(*u) }

// AddUint registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint(name, shortName string, value uint, p *uint, help string) *Opt {
	v := newUintValue(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
