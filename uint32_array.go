package clipper

import (
	"strconv"
	"strings"
)

// -- uint32Array Value
type uint32ArrayValue struct {
	value *[]uint32
}

func newUint32ArrayValue(val []uint32, p *[]uint32) *uint32ArrayValue {
	ssv := new(uint32ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUint32ArrayValueFromCSV(val string, p *[]uint32) *uint32ArrayValue {
	ssv := new(uint32ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uint32ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint32, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseUint(s, 10, 32); err == nil {
			iv = append(iv, uint32(n))
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

func (s *uint32ArrayValue) Reset(i interface{}) {
	v := i.([]uint32)
	*s.value = v
}

func (s *uint32ArrayValue) ResetWith(val []uint32) error {
	out := make([]uint32, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uint32ArrayValue) GetUint32Array() []uint32 {
	out := make([]uint32, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uint32ArrayValue) Type() string {
	return "uint32Array"
}

func (s *uint32ArrayValue) Get() interface{} {
	return s.GetUint32Array()
}

func (s *uint32ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(uint64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddUint32Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint32Array(name, shortName string, value []uint32, p *[]uint32, help string) *Opt {
	v := newUint32ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddUint32ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint32ArrayFromCSV(name, shortName string, value string, p *[]uint32, help string) *Opt {
	v := newUint32ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
