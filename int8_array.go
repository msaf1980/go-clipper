package clipper

import (
	"strconv"
	"strings"
)

// -- int8Array Value
type int8ArrayValue struct {
	value *[]int8
}

func newInt8ArrayValue(val []int8, p *[]int8) *int8ArrayValue {
	ssv := new(int8ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newInt8ArrayValueFromCSV(val string, p *[]int8) *int8ArrayValue {
	ssv := new(int8ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *int8ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]int8, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseInt(s, 10, 8); err == nil {
			iv = append(iv, int8(n))
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

func (s *int8ArrayValue) Reset(i interface{}) {
	v := i.([]int8)
	*s.value = v
}

func (s *int8ArrayValue) ResetWith(val []int8) error {
	out := make([]int8, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *int8ArrayValue) GetInt8Array() []int8 {
	out := make([]int8, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *int8ArrayValue) Type() string {
	return "int8Array"
}

func (s *int8ArrayValue) Get() interface{} {
	return s.GetInt8Array()
}

func (s *int8ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatInt(int64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddInt8Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt8Array(name, shortName string, value []int8, p *[]int8) *Opt {
	v := newInt8ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}

// AddInt8ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt8ArrayFromCSV(name, shortName string, value string, p *[]int8) *Opt {
	v := newInt8ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
