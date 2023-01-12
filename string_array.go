package clipper

// -- stringArray Value
type stringArrayValue struct {
	value *[]string
}

func newStringArrayValue(val []string, p *[]string) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newStringArrayValueFromCSV(val string, p *[]string) *stringArrayValue {
	ssv := new(stringArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *stringArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if doAppend {
		*s.value = append(*s.value, v...)
	} else {
		*s.value = v
	}
	return nil
}

func (s *stringArrayValue) Reset(i interface{}) {
	v := i.([]string)
	*s.value = v
}

func (s *stringArrayValue) ResetWith(val []string) error {
	out := make([]string, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *stringArrayValue) GetStringArray() []string {
	out := make([]string, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *stringArrayValue) Type() string {
	return "stringArray"
}

func (s *stringArrayValue) Get() interface{} {
	return s.GetStringArray()
}

func (s *stringArrayValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}

// -- stringArray Value (limited length)
type stringArrayLValue struct {
	stringArrayValue

	min int
	max int
}

func newStringArrayLValue(val []string, p *[]string, max int) *stringArrayLValue {
	ssv := new(stringArrayLValue)
	ssv.value = p
	*ssv.value = val
	ssv.max = max
	return ssv
}

func (s *stringArrayLValue) SetMaxLen(max int) Arg {
	s.max = max
	return s
}

func (s *stringArrayLValue) MaxLen() int {
	return s.max
}

func (s *stringArrayLValue) SetMinLen(min int) Arg {
	s.min = min
	return s
}

func (s *stringArrayLValue) MinLen() int {
	return s.min
}

func (s *stringArrayLValue) CheckLen() error {
	return CheckLen("length", len(*s.value), s.min, s.max)
}

func (s *stringArrayLValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if s.max >= 0 && len(*s.value)+len(v) > s.max {
		return ErrorLengthOverflow{"length at argument=" + val, ">", s.max}
	}
	if doAppend {
		*s.value = append(*s.value, v...)
	} else {
		*s.value = v
	}
	return nil
}

func (s *stringArrayLValue) ResetWith(val []string) error {
	if s.max >= 0 && len(val) > s.max {
		return ErrorLengthOverflow{"length", ">", s.max}
	}
	out := make([]string, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

// AddStringArray registers an string argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddStringArray(name, shortName string, value []string, p *[]string, help string) *Opt {
	v := newStringArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}

// AddStringArrayFromCSV registers an string argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddStringArrayFromCSV(name, shortName string, value string, p *[]string, help string) *Opt {
	v := newStringArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
