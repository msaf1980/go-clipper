package clipper

import (
	"strconv"
)

type int64Value int64

func newInt64Value(val int64, p *int64) *int64Value {
	*p = val
	return (*int64Value)(p)
}

func (i *int64Value) Set(s string, _ bool) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		*i = int64Value(v)
	}
	return err
}

func (iv *int64Value) Reset(i interface{}) {
	v := i.(int64)
	*iv = int64Value(v)
}

func (i *int64Value) Type() string {
	return "int64"
}

func (i *int64Value) Get() interface{} {
	return i.GetInt64()
}

func (i *int64Value) String() string { return strconv.FormatInt(int64(*i), 10) }

func (i *int64Value) GetInt64() int64 { return int64(*i) }

// AddInt64 registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt64(name, shortName string, value int64, p *int64, help string) *Opt {
	v := newInt64Value(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
