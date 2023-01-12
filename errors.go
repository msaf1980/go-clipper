package clipper

import (
	"errors"
	"fmt"
)

var (
	ErrTypeMismatch = errors.New("type mismatch")
)

// ErrorUnknownCommand represents an error when command-line arguments contain an unregistered command.
type ErrorUnknownCommand struct {
	Name string
}

func (e ErrorUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command %q found in the arguments", e.Name)
}

// ErrorUnknownFlag represents an error when command-line arguments contain an unregistered flag.
type ErrorUnknownFlag struct {
	Name string
}

func (e ErrorUnknownFlag) Error() string {
	return fmt.Sprintf("unknown flag %q found in the arguments", e.Name)
}

// ErrorRequiredFlag represents an error when command-line arguments not contain an required flag.
type ErrorRequiredFlag struct {
	Name string
}

func (e ErrorRequiredFlag) Error() string {
	return fmt.Sprintf("required flag %q not found in the arguments", e.Name)
}

// ErrorUnsupportedFlag represents an error when command-line arguments contain an unsupported flag.
type ErrorUnsupportedFlag struct {
	Name string
}

func (e ErrorUnsupportedFlag) Error() string {
	return fmt.Sprintf("unsupported flag %q found in the arguments", e.Name)
}

// ErrorUnsupportedValue represents an error when command-line arguments contain an unsupported value.
type ErrorUnsupportedValue struct {
	Name  string
	Value string
}

func (e ErrorUnsupportedValue) Error() string {
	return fmt.Sprintf("unsupported value %s=%s found in the arguments", e.Name, e.Value)
}

// ErrorUnsupportedValue represents an error when command-line arguments contain an unsupported value.
type ErrorLengthOverflow struct {
	Name  string
	Cmp   string
	Value int
}

func (e ErrorLengthOverflow) Error() string {
	return fmt.Sprintf("%s %s %d", e.Name, e.Cmp, e.Value)
}

// ErrorWrapped represents an wrapped error
type ErrorWrapped struct {
	Prefix string
	err    error
}

func (e ErrorWrapped) Error() string {
	return e.Prefix + " " + e.err.Error()
}

func WrapLengthOverflow(prefix string, err error) error {
	switch v := err.(type) {
	case ErrorLengthOverflow:
		v.Name = prefix + " " + v.Name
		return v
	default:
		return err
	}
}
