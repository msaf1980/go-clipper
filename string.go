package clipper

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string, _ bool) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Reset(i interface{}) {
	v := i.(string)
	*s = stringValue(v)
}

func (s *stringValue) Type() string {
	return "string"
}

func (s *stringValue) Get() interface{} {
	return s.String()
}

func (s *stringValue) GetString() string { return string(*s) }

func (s *stringValue) String() string { return s.GetString() }

// func stringConv(sval string) (interface{}, error) {
// 	return sval, nil
// }

// AddString registers an string argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddString(name, shortName string, value string, p *string) *Opt {
	v := newStringValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
