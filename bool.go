package clipper

import "strconv"

// -- bool Value
type boolValue bool

func newBoolValue(p *bool) *boolValue {
	return (*boolValue)(p)
}

func (b *boolValue) Set(s string, _ bool) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

func (b *boolValue) Reset(i interface{}) {
	v := i.(bool)
	*b = boolValue(v)
}

func (b *boolValue) Type() string {
	return "bool"
}

func (b *boolValue) Get() interface{} {
	return b.GetBool()
}

func (b *boolValue) String() string { return strconv.FormatBool(bool(*b)) }

func (b *boolValue) GetBool() bool { return bool(*b) }

// AddBool registers an bool configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddBool(name, shortName string, value bool, b *bool, help string) *Opt {
	*b = value
	v := newBoolValue(b)
	o := commandConfig.AddValue(name, shortName, v, help)
	return o
}
