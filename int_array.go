package clipper

import (
	"strconv"
	"strings"
)

// -- intArray Value
type intArrayValue struct {
	value *[]int
}

func newIntArrayValue(val []int, p *[]int) *intArrayValue {
	ssv := new(intArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newIntArrayValueFromCSV(val string, p *[]int) *intArrayValue {
	ssv := new(intArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *intArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]int, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			iv = append(iv, int(n))
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

func (s *intArrayValue) Reset(i interface{}) {
	v := i.([]int)
	*s.value = v
}

func (s *intArrayValue) ResetWith(val []int) error {
	out := make([]int, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *intArrayValue) GetIntArray() []int {
	out := make([]int, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *intArrayValue) Type() string {
	return "intArray"
}

func (s *intArrayValue) Get() interface{} {
	return s.GetIntArray()
}

func (s *intArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.Itoa(n)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddIntArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddIntArray(name, shortName string, value []int, p *[]int) *Opt {
	v := newIntArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}

// AddIntArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddIntArrayFromCSV(name, shortName string, value string, p *[]int) *Opt {
	v := newIntArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
