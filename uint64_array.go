package clipper

import (
	"strconv"
	"strings"
)

// -- uint64Array Value
type uint64ArrayValue struct {
	value *[]uint64
}

func newUint64ArrayValue(val []uint64, p *[]uint64) *uint64ArrayValue {
	ssv := new(uint64ArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUint64ArrayValueFromCSV(val string, p *[]uint64) *uint64ArrayValue {
	ssv := new(uint64ArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uint64ArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint64, 0, len(v))
	for _, s := range v {
		if n, err := strconv.ParseUint(s, 10, 64); err == nil {
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

func (s *uint64ArrayValue) Reset(i interface{}) {
	v := i.([]uint64)
	*s.value = v
}

func (s *uint64ArrayValue) ResetWith(val []uint64) error {
	out := make([]uint64, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uint64ArrayValue) GetUint64Array() []uint64 {
	out := make([]uint64, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uint64ArrayValue) Type() string {
	return "uint64Array"
}

func (s *uint64ArrayValue) Get() interface{} {
	return s.GetUint64Array()
}

func (s *uint64ArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(n, 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddIntArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64Array(name, shortName string, value []uint64, p *[]uint64, help string) *Opt {
	v := newUint64ArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddIntArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64ArrayFromCSV(name, shortName string, value string, p *[]uint64, help string) *Opt {
	v := newUint64ArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
