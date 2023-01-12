package clipper

import (
	"strconv"
	"strings"
)

// -- int32Array Value
type int32ArrayValue struct {
	value *[]int32
}

func newInt32ArrayValue(val []int32, p *[]int32) *int32ArrayValue {
	ssv := new(int32ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newInt32ArrayValueFromCSV(val string, p *[]int32) *int32ArrayValue {
	ssv := new(int32ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *int32ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]int32, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			iv = append(iv, int32(n))
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

func (s *int32ArrayValue) Reset(i interface{}) {
	v := i.([]int32)
	*s.value = v
}

func (s *int32ArrayValue) ResetWith(val []int32) error {
	out := make([]int32, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *int32ArrayValue) GetInt32Array() []int32 {
	out := make([]int32, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *int32ArrayValue) Type() string {
	return "int32Array"
}

func (s *int32ArrayValue) Get() interface{} {
	return s.GetInt32Array()
}

func (s *int32ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatInt(int64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddInt32Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt32Array(name, shortName string, value []int32, p *[]int32, help string) *Opt {
	v := newInt32ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}

// AddInt32ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt32ArrayFromCSV(name, shortName string, value string, p *[]int32, help string) *Opt {
	v := newInt32ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
