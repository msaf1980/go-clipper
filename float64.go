package clipper

import (
	"strconv"
)

// -- float64 Value
type float64Value float64

func newFloat64Value(val float64, p *float64) *float64Value {
	*p = val
	return (*float64Value)(p)
}

func (f *float64Value) Set(s string, _ bool) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = float64Value(v)
	return err
}

func (f *float64Value) Reset(i interface{}) {
	v := i.(float64)
	*f = float64Value(v)
}

func (f *float64Value) Type() string {
	return "float64"
}

func (f *float64Value) Get() interface{} {
	return f.GetFloat64()
}

func (f *float64Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

func (f *float64Value) GetFloat64() float64 { return float64(*f) }

// AddFloat64 registers an float64 argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddFloat64(name, shortName string, value float64, p *float64, help string) *Opt {
	v := newFloat64Value(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
