package clipper

import (
	"strconv"
)

type int16Value int16

func newint16Value(val int16, p *int16) *int16Value {
	*p = val
	return (*int16Value)(p)
}

func (i *int16Value) Set(s string, _ bool) error {
	v, err := strconv.ParseInt(s, 0, 16)
	if err == nil {
		*i = int16Value(v)
	}
	return err
}

func (iv *int16Value) Reset(i interface{}) {
	v := i.(int16)
	*iv = int16Value(v)
}

func (i *int16Value) Type() string {
	return "int16"
}

func (i *int16Value) Get() interface{} {
	return i.Getint16()
}

func (i *int16Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int16Value) Getint16() int16 { return int16(*i) }

// AddInt16 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt16(name, shortName string, value int16, p *int16, help string) *Opt {
	v := newint16Value(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
