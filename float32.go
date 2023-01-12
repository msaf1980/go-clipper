package clipper

import "strconv"

// -- float32 Value
type float32Value float32

func newFloat32Value(val float32, p *float32) *float32Value {
	*p = val
	return (*float32Value)(p)
}

func (f *float32Value) Set(s string, _ bool) error {
	v, err := strconv.ParseFloat(s, 32)
	*f = float32Value(v)
	return err
}

func (f *float32Value) Reset(i interface{}) {
	v := i.(float32)
	*f = float32Value(v)
}

func (f *float32Value) Type() string {
	return "float32"
}

func (f *float32Value) Get() interface{} {
	return f.GetFloat32()
}

func (f *float32Value) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 32) }

func (f *float32Value) GetFloat32() float32 { return float32(*f) }

// AddFloat32 registers an float32 argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddFloat32(name, shortName string, value float32, p *float32) *Opt {
	v := newFloat32Value(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
