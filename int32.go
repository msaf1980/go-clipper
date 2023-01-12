package clipper

import (
	"strconv"
)

type int32Value int32

func newint32Value(val int32, p *int32) *int32Value {
	*p = val
	return (*int32Value)(p)
}

func (i *int32Value) Set(s string, _ bool) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		*i = int32Value(v)
	}
	return err
}

func (iv *int32Value) Reset(i interface{}) {
	v := i.(uint32)
	*iv = int32Value(v)
}

func (i *int32Value) Type() string {
	return "int32"
}

func (i *int32Value) Get() interface{} {
	return i.Getint32()
}

func (i *int32Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int32Value) Getint32() int32 { return int32(*i) }

// AddInt32 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt32(name, shortName string, value int32, p *int32, help string) *Opt {
	v := newint32Value(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
