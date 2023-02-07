package clipper

import (
	"strconv"
)

type counterFlag int

func newCounterFlag(val int, p *int) *counterFlag {
	*p = val
	return (*counterFlag)(p)
}

func (c *counterFlag) Set(s string, doAppend bool) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	if doAppend {
		if v {
			*c++
		} else {
			*c--
		}
	} else if v {
		*c = 1
	} else {
		*c = 0
	}
	return nil
}

func (c *counterFlag) Reset(i interface{}) {
	*c = i.(counterFlag)
}

func (c *counterFlag) Type() string {
	return "counter"
}

func (c *counterFlag) Get() interface{} {
	return c.GetInt()
}

func (c *counterFlag) String() string { return strconv.Itoa(int(*c)) }

func (c *counterFlag) GetInt() int { return int(*c) }

// AddCounterFlag registers an counter (direct) multi-flag with the command (for cases like -vvv).
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddCounterFlag(name, shortName string, c *int, help string) *Opt {
	var val bool
	val, name = isInvertedFlag(name)
	if val {
		panic(name + "can't be inverted")
	}
	v := newCounterFlag(0, c)
	o := commandConfig.AddValue(name, shortName, v, false, help)
	o.IsFlag = true
	o.IsMultiValue = true
	return o
}
