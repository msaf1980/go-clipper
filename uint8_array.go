package clipper

import (
	"strconv"
	"strings"
)

// -- uint8Array Value
type uint8ArrayValue struct {
	value *[]uint8
}

func newUint8ArrayValue(val []uint8, p *[]uint8) *uint8ArrayValue {
	ssv := new(uint8ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUint8ArrayValueFromCSV(val string, p *[]uint8) *uint8ArrayValue {
	ssv := new(uint8ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uint8ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint8, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseUint(s, 10, 16); err == nil {
			iv = append(iv, uint8(n))
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

func (s *uint8ArrayValue) Reset(i interface{}) {
	v := i.([]uint8)
	*s.value = v
}

func (s *uint8ArrayValue) ResetWith(val []uint8) error {
	out := make([]uint8, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uint8ArrayValue) GetUint8Array() []uint8 {
	out := make([]uint8, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uint8ArrayValue) Type() string {
	return "uint8Array"
}

func (s *uint8ArrayValue) Get() interface{} {
	return s.GetUint8Array()
}

func (s *uint8ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(uint64(n), 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddUint8Array registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint8Array(name, shortName string, value []uint8, p *[]uint8, help string) *Opt {
	v := newUint8ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddUint8ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint8ArrayFromCSV(name, shortName string, value string, p *[]uint8, help string) *Opt {
	v := newUint8ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
