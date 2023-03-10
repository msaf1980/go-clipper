package clipper

import (
	"strconv"
	"strings"
)

// -- int64NArray Value
type int64NArrayValue struct {
	value *[]int64
}

func newInt64NArrayValue(val []int64, p *[]int64) *int64NArrayValue {
	ssv := new(int64NArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newInt64NArrayValueFromCSV(val string, p *[]int64) *int64NArrayValue {
	ssv := new(int64NArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *int64NArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]int64, 0, len(v))
	for _, s := range v {
		if n, err := int64NFromString(s); err == nil {
			iv = append(iv, int64(n))
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

func (s *int64NArrayValue) Reset(i interface{}) {
	v := i.([]int64)
	*s.value = v
}

func (s *int64NArrayValue) ResetWith(val []int64) error {
	out := make([]int64, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *int64NArrayValue) GetInt64Array() []int64 {
	out := make([]int64, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *int64NArrayValue) Type() string {
	return "int64NArray"
}

func (s *int64NArrayValue) Get() interface{} {
	return s.GetInt64Array()
}

func (s *int64NArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatInt(n, 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddInt64Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddInt64NArray(name, shortName string, value []int64, p *[]int64, help string) *Opt {
	v := newInt64NArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddInt64ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) Addint64NArrayFromCSV(name, shortName string, value string, p *[]int64, help string) *Opt {
	v := newInt64NArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
