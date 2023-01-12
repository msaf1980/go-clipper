package clipper

import (
	"strconv"
	"strings"
)

// -- uint16Array Value
type uint16ArrayValue struct {
	value *[]uint16
}

func newUint16ArrayValue(val []uint16, p *[]uint16) *uint16ArrayValue {
	ssv := new(uint16ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUint16ArrayValueFromCSV(val string, p *[]uint16) *uint16ArrayValue {
	ssv := new(uint16ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uint16ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint16, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseUint(s, 10, 16); err == nil {
			iv = append(iv, uint16(n))
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

func (s *uint16ArrayValue) Reset(i interface{}) {
	v := i.([]uint16)
	*s.value = v
}

func (s *uint16ArrayValue) ResetWith(val []uint16) error {
	out := make([]uint16, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uint16ArrayValue) GetUint16Array() []uint16 {
	out := make([]uint16, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uint16ArrayValue) Type() string {
	return "uint16Array"
}

func (s *uint16ArrayValue) Get() interface{} {
	return s.GetUint16Array()
}

func (s *uint16ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(uint64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddUint16Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint16Array(name, shortName string, value []uint16, p *[]uint16, help string) *Opt {
	v := newUint16ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}

// AddUint16ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint16ArrayFromCSV(name, shortName string, value string, p *[]uint16, help string) *Opt {
	v := newUint16ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, help)
}
