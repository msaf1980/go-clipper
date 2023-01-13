package clipper

import "time"

// -- time.Time Value
type timeValue struct {
	layout string

	value *time.Time
}

func newTimeValue(val time.Time, p *time.Time, layout string) *timeValue {
	tv := new(timeValue)
	tv.layout = layout
	tv.value = p
	*tv.value = val
	return tv
}

func newTimeValueFromString(val string, p *time.Time, layout string) *timeValue {
	v, err := time.Parse(layout, val)
	if err == nil {
		panic(err)
	}
	tv := new(timeValue)
	tv.layout = layout
	tv.value = p
	*tv.value = v
	return tv
}

func (t *timeValue) Set(s string, _ bool) error {
	v, err := time.Parse(t.layout, s)
	if err == nil {
		*t.value = v
	}
	return err
}

func (t *timeValue) Reset(i interface{}) {
	v := i.(time.Time)
	*t.value = v
}

func (t *timeValue) Type() string {
	return "time"
}

func (t *timeValue) Get() interface{} {
	return *t.value
}

func (t *timeValue) String() string { return t.value.Format(t.layout) }

func (t *timeValue) GetTime() time.Time { return *t.value }

// AddTime registers an duration argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddTime(name, shortName string, value time.Time, p *time.Time, layout string, help string) *Opt {
	v := newTimeValue(value, p, layout)
	return commandConfig.AddValue(name, shortName, v, help)
}

// AddTimeFromString registers an duration argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddTimeFromString(name, shortName string, value string, p *time.Time, layout string, help string) *Opt {
	v := newTimeValueFromString(value, p, layout)
	return commandConfig.AddValue(name, shortName, v, help)
}
