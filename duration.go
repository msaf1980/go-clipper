package clipper

import "time"

// -- time.Duration Value
type durationValue time.Duration

func newDurationValue(val time.Duration, p *time.Duration) *durationValue {
	*p = val
	return (*durationValue)(p)
}

func (d *durationValue) Set(s string, _ bool) error {
	v, err := time.ParseDuration(s)
	if err == nil {
		*d = durationValue(v)
	}
	return err
}

func (d *durationValue) Reset(i interface{}) {
	v := i.(time.Duration)
	*d = durationValue(v)
}

func (d *durationValue) Type() string {
	return "duration"
}

func (d *durationValue) Get() interface{} {
	return d.GetDuration()
}

func (d *durationValue) String() string { return (*time.Duration)(d).String() }

func (d *durationValue) GetDuration() time.Duration { return time.Duration(*d) }

// AddDuration registers an duration argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddDuration(name, shortName string, value time.Duration, p *time.Duration, help string) *Opt {
	v := newDurationValue(value, p)
	return commandConfig.AddValue(name, shortName, v, false, help)
}
