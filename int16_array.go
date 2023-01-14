package clipper

import (
	"strconv"
	"strings"
)

// -- int16Array Value
type int16ArrayValue struct {
	value *[]int16
}

func newInt16ArrayValue(val []int16, p *[]int16) *int16ArrayValue {
	ssv := new(int16ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newInt16ArrayValueFromCSV(val string, p *[]int16) *int16ArrayValue {
	ssv := new(int16ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *int16ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]int16, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseInt(s, 10, 16); err == nil {
			iv = append(iv, int16(n))
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

func (s *int16ArrayValue) Reset(i interface{}) {
	v := i.([]int16)
	*s.value = v
}

func (s *int16ArrayValue) ResetWith(val []int16) error {
	out := make([]int16, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *int16ArrayValue) GetInt16Array() []int16 {
	out := make([]int16, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *int16ArrayValue) Type() string {
	return "int16Array"
}

func (s *int16ArrayValue) Get() interface{} {
	return s.GetInt16Array()
}

func (s *int16ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatInt(int64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddInt16Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) Addint16Array(name, shortName string, value []int16, p *[]int16, help string) *Opt {
	v := newInt16ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddInt16ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) Addint16ArrayFromCSV(name, shortName string, value string, p *[]int16, help string) *Opt {
	v := newInt16ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
