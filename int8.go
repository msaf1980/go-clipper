package clipper

import (
	"strconv"
)

type int8Value int8

func newint8Value(val int8, p *int8) *int8Value {
	*p = val
	return (*int8Value)(p)
}

func (i *int8Value) Set(s string, _ bool) error {
	v, err := strconv.ParseInt(s, 0, 8)
	if err == nil {
		*i = int8Value(v)
	}
	return err
}

func (iv *int8Value) Reset(i interface{}) {
	v := i.(uint8)
	*iv = int8Value(v)
}

func (i *int8Value) Type() string {
	return "int8"
}

func (i *int8Value) Get() interface{} {
	return i.Getint8()
}

func (i *int8Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int8Value) Getint8() int8 { return int8(*i) }

// AddInt8 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt8(name, shortName string, value int8, p *int8, help string) *Opt {
	v := newint8Value(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
