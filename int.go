package clipper

import (
	"strconv"
)

type intValue int

func newIntValue(val int, p *int) *intValue {
	*p = val
	return (*intValue)(p)
}

func (i *intValue) Set(s string, _ bool) error {
	v, err := strconv.ParseInt(s, 0, 32)
	if err == nil {
		*i = intValue(v)
	}
	return err
}

func (iv *intValue) Reset(i interface{}) {
	v := i.(uint)
	*iv = intValue(v)
}

func (i *intValue) Type() string {
	return "int"
}

func (i *intValue) Get() interface{} {
	return i.GetInt()
}

func (i *intValue) String() string { return strconv.Itoa(int(*i)) }

func (i *intValue) GetInt() int { return int(*i) }

// AddInt registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt(name, shortName string, value int, p *int, help string) *Opt {
	v := newIntValue(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
