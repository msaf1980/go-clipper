package clipper

import (
	"strconv"
	"strings"
)

// -- uintArray Value
type uintArrayValue struct {
	value *[]uint
}

func newUintArrayValue(val []uint, p *[]uint) *uintArrayValue {
	ssv := new(uintArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUintArrayValueFromCSV(val string, p *[]uint) *uintArrayValue {
	ssv := new(uintArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uintArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseUint(s, 10, 32); err == nil {
			iv = append(iv, uint(n))
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

func (s *uintArrayValue) Reset(i interface{}) {
	v := i.([]uint)
	*s.value = v
}

func (s *uintArrayValue) ResetWith(val []uint) error {
	out := make([]uint, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uintArrayValue) GetUintArray() []uint {
	out := make([]uint, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uintArrayValue) Type() string {
	return "uintArray"
}

func (s *uintArrayValue) Get() interface{} {
	return s.GetUintArray()
}

func (s *uintArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(uint64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddintArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUintArray(name, shortName string, value []uint, p *[]uint) *Opt {
	v := newUintArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}

// AddintArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUintArrayFromCSV(name, shortName string, value []uint, p *[]uint) *Opt {
	v := newUintArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v)
}
