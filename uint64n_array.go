package clipper

import (
	"strconv"
	"strings"
)

// -- uint64NArray Value
type uint64NArrayValue struct {
	value *[]uint64
}

func newUint64NArrayValue(val []uint64, p *[]uint64) *uint64NArrayValue {
	ssv := new(uint64NArrayValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func newUint64NArrayValueFromCSV(val string, p *[]uint64) *uint64NArrayValue {
	ssv := new(uint64NArrayValue)
	ssv.value = p
	if err := ssv.Set(val, false); err != nil {
		panic(err)
	}
	return ssv
}

func (s *uint64NArrayValue) Set(val string, doAppend bool) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	iv := make([]uint64, 0, len(v))
	for _, s := range v {
		if n, err := uint64NFromString(s); err == nil {
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

func (s *uint64NArrayValue) Reset(i interface{}) {
	v := i.([]uint64)
	*s.value = v
}

func (s *uint64NArrayValue) ResetWith(val []uint64) error {
	out := make([]uint64, len(val))
	copy(out, val)
	*s.value = out
	return nil
}

func (s *uint64NArrayValue) GetUint64Array() []uint64 {
	out := make([]uint64, len(*s.value))
	copy(out, *s.value)
	return out
}

func (s *uint64NArrayValue) Type() string {
	return "uint64NArray"
}

func (s *uint64NArrayValue) Get() interface{} {
	return s.GetUint64Array()
}

func (s *uint64NArrayValue) String() string {
	strSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		strSlice[i] = strconv.FormatUint(n, 10)
	}

	return "[" + strings.Join(strSlice, ",") + "]"
}

// AddUint64NArray registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64NArray(name, shortName string, value []uint64, p *[]uint64, help string) *Opt {
	v := newUint64NArrayValue(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}

// AddUint64ArrayFromCSV registers an int argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddUint64NArrayFromCSV(name, shortName string, value string, p *[]uint64, help string) *Opt {
	v := newUint64NArrayValueFromCSV(value, p)
	return commandConfig.AddValue(name, shortName, v, true, help)
}
