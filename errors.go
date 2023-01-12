package clipper

import (
	"errors"
	"fmt"
)

var (
	ErrOverflow     = errors.New("length overflow")
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
