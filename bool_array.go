package clipper

import (
	"strconv"
	"strings"
)

// -- boolArray Value
type boolArrayValue struct {
	value *[]bool
}

func newBoolArrayValue(val []bool, p *[]bool) *boolArrayValue {
	ssv := new(boolArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newboolArrayValueFromCSV(val string, p *[]bool) *boolArrayValue {
	ssv := new(boolArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *boolArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]bool, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseBool(s); err == nil {
			iv = append(iv, n)
		} else {
			return err
		}
	}
	if doAppend {
		*s.value = append(*s.value, iv...)
	} else {
		*s.value = iv
	}
	return nil
}

func (s *boolArrayValue) Reset(i interface{}) {
	v := i.([]bool)
	*s.value = v
}

func (s *boolArrayValue) ResetWith(val []bool) error {
	out := make([]bool, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *boolArrayValue) GetSlice() []bool {
	out := make([]bool, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *boolArrayValue) Type() string {
	return "boolArray"
}

func (s *boolArrayValue) Get() interface{} {
	return s.GetSlice()
}

func (s *boolArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatBool(n)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddBoolArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddBoolArray(name, shortName string, value []bool, p *[]bool, help string) *Opt {
	v := newBoolArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddBoolArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddBoolArrayFromCSV(name, shortName string, value string, p *[]bool, help string) *Opt {
	v := newboolArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
