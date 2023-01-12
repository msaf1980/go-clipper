// MIT License

// Copyright (c) 2020 Uday Hiwarale

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package clipper processes the command-line arguments of getopt(3) syntax.
// This package provides the ability to process the root command, sub commands,
// command-line arguments and command-line flags.
package clipper

import (
	"strconv"
	"strings"
)

/***********************************************
        PRIVATE FUNCTIONS AND VARIABLES
***********************************************/

// format command-line argument values
func formatCommandValues(values []string) (formatted []string) {

	formatted = make([]string, 0)

	// split a value by `=`
	for _, value := range values {
		if isFlag(value) {
			parts := strings.Split(value, "=")

			for _, part := range parts {
				if strings.Trim(part, " ") != "" {
					formatted = append(formatted, part)
				}
			}
		} else {
			formatted = append(formatted, value)
		}
	}

	return
}

// check if value is a short flag
func isShortFlag(value string) (bool, string) {
	if strings.HasPrefix(value, "-") {
		if len(value) == 2 && value != "--" {
			return true, value[1:2]
		}
	}
	return false, value
}

// check if value is a long flag
func isLongFlag(value string) (bool, string) {
	if len(value) > 2 && strings.HasPrefix(value, "--") {
		return true, value[2:]
	}
	return false, value
}

// check if value starts with `no-` prefix
func isInvertedFlag(value string) (bool, string) {
	if strings.HasPrefix(value, "no-") {
		return true, value[3:]
	}
	return false, value
}

// check if flag is unsupported
func isUnsupportedFlag(value string) bool {

	// a flag should be at least two characters log
	if len(value) >= 2 {

		// if short flag, it should start with `-` but not with `--`
		if len(value) == 2 {
			return !strings.HasPrefix(value, "-") || strings.HasPrefix(value, "--")
		}

		// if long flag, it should start with `--` and not with `---`
		return !strings.HasPrefix(value, "--") || strings.HasPrefix(value, "---")
	}

	return false
}

// check if flag
func isFlag(value string) bool {
	return len(value) >= 2 && strings.HasPrefix(value, "-")
}

// check if value ends with `...` sufix
func isVariadicArgument(value string) (bool, string) {
	if strings.HasSuffix(value, "...") {
		return true, value[:len(value)-3] // trim `...` suffix
	}

	return false, value
}

// check if values corresponds to the root command and return command
func getCommand(values []string, registry Registry) (commandName string, valuesToProcess []string,
	commandConfig *CommandConfig, err error) {

	// TRUE: if all `values` are empty or the first `value` is a flag
	var ok bool
	if len(values) == 0 || isFlag(values[0]) {
		// root coomand is empty
		commandConfig, ok = registry[commandName]
		if !ok {
			err = ErrorUnknownCommand{commandName}
		}
		valuesToProcess = values
	} else {
		// get `CommandConfig` object from the registry
		// if command is not registered, return `ErrorUnknownCommand` error
		commandName, valuesToProcess = nextValue(values)
		commandConfig, ok = registry[commandName]
		if !ok {
			commandName = ""
			valuesToProcess = values
			commandConfig, ok = registry[commandName]
			if !ok {
				if len(values) == 0 || isFlag(values[0]) {
					err = ErrorUnknownCommand{commandName}
				} else {
					err = ErrorUnknownCommand{values[0]}
				}
			}
		}
	}

	return
}

// return next value and remaining values of a slice of strings
func nextValue(slice []string) (v string, newSlice []string) {

	if len(slice) == 0 {
		v, newSlice = "", slice
		return
	}

	v = slice[0]

	if len(slice) > 1 {
		newSlice = slice[1:]
	} else {
		newSlice = slice[:0]
	}

	return
}

// remove whitespaces from a value
func removeWhitespaces(value string) string {
	return strings.ReplaceAll(value, " ", "_")
}

/***********************************************/

// Registry holds the configuration of the registered commands.
type Registry map[string]*CommandConfig

// Register method registers a command.
// The "name" argument should be a simple string.
// If "name" is an empty string, it is considered as a root command.
// If a command is already registered, the registered `*CommandConfig` object is returned.
// If the command is already registered, second return value will be `true`.
func (registry Registry) Register(name string) (*CommandConfig, bool) {

	// remove all whitespaces
	commandName := removeWhitespaces(name)

	// check if command is already registered, if found, return existing entry
	if _commandConfig, ok := registry[commandName]; ok {
		return _commandConfig, true
	}

	// construct new `CommandConfig` object
	commandConfig := &CommandConfig{
		Name:  commandName,
		Opts:  make(map[string]*Opt),
		short: make(map[string]string),
		Args:  None{}, // by default disable unnamed args
	}

	// add entry to the registry
	registry[commandName] = commandConfig

	return commandConfig, false
}

// Reset method reset values to it's default values.
func (registry Registry) Reset() {
	for _, cmd := range registry {
		for _, opt := range cmd.Opts {
			opt.Reset()
		}
		cmd.Args.Reset([]string{})
	}
}

// Parse method parses command-line arguments and returns an appropriate "*CommandConfig" object registered in the registry.
// If command is not registered, it return `ErrorUnknownCommand` error.
// If there is an error parsing a flag, it can return an `ErrorUnknownFlag` or `ErrorUnsupportedFlag` error.
func (registry Registry) Parse(values []string) (string, error) {

	commandName, valuesToProcess, commandConfig, err := getCommand(values, registry)
	if err != nil {
		return commandName, err
	}

	// format command-line argument values to process
	valuesToProcess = formatCommandValues(valuesToProcess)

	// check for invalid flag structure
	for _, val := range valuesToProcess {
		if isFlag(val) && isUnsupportedFlag(val) {
			return commandName, ErrorUnsupportedFlag{val}
		}
	}

	// process all command-line arguments (except command name)
	for {

		// get current command-line argument value
		var value string
		value, valuesToProcess = nextValue(valuesToProcess)

		// if `value` is empty, break the loop
		if len(value) == 0 {
			break
		}

		// check if `value` is a `flag` or an `argument`
		if isFlag(value) {

			var (
				opt        *Opt
				isInverted bool
			)

			isVariadic, value := isVariadicArgument(value)

			if ok, name := isShortFlag(value); ok {

				// get long flag name
				flagName, ok := commandConfig.short[name]
				if !ok {
					return commandName, ErrorUnknownFlag{value}
				}

				if opt, ok = commandConfig.Opts[flagName]; !ok {
					return commandName, ErrorUnknownFlag{value}
				}
			} else if ok, name := isLongFlag(value); ok {
				isInverted, name = isInvertedFlag(name)
				if opt, ok = commandConfig.Opts[name]; !ok {
					return commandName, ErrorUnknownFlag{value}
				}
			} else {
				return commandName, ErrorUnknownFlag{value}
			}

			// set flag value
			if opt.IsBool {
				if isInverted {
					opt.Set("false") // if flag is an inverted flag, its value will be `false`
				} else {
					opt.Set("true")
				}
			} else {
				for {
					if nextValue, nextValuesToProcess := nextValue(valuesToProcess); len(nextValue) != 0 && !isFlag(nextValue) {
						if !opt.Validate(nextValue) {
							return commandName, ErrorUnsupportedValue{opt.Name, nextValue}
						}
						if err = opt.Set(nextValue); err != nil {
							return commandName, WrapInvalidValue(strconv.Quote(commandName), err)
						}
						valuesToProcess = nextValuesToProcess
					} else {
						break
					}
					if !isVariadic {
						break
					}
				}
			}
		} else {
			if err := commandConfig.Args.Set(value, true); err != nil {
				return commandName, WrapInvalidValue(strconv.Quote(commandName)+" unnamed args", err)
			}
		}
	}

	if err = commandConfig.Args.CheckLen(); err != nil {
		return commandName, WrapInvalidValue(strconv.Quote(commandName)+" unnamed args", err)
	}

	for _, opt := range commandConfig.Opts {
		if opt.IsRequired && !opt.Changed {
			return commandName, ErrorRequiredFlag{opt.Name}
		}
	}

	return commandName, nil
}

// NewRegistry returns new instance of the "Registry"
func NewRegistry() Registry {
	return make(Registry)
}

/*---------------------*/

// CommandConfig type holds the structure and values of the command-line arguments of command.
type CommandConfig struct {

	// name of the sub-command ("" for the root command)
	Name string

	// named command-line options order (for display help)
	OptsOrder []string

	// named command-line options or boolean flags
	Opts map[string]*Opt

	// mapping of the short  options/flag names with long  options/flag names
	short map[string]string

	// Unnamed args
	Args Arg
}

// AddStringArgs set unnamed argument configuration with the command.
// The `max` argument represents maximum length of unnamed args (-1 - unlimited).
// `Arg` object returned.
func (commandConfig *CommandConfig) AddStringArgs(max int, p *[]string) Arg {
	commandConfig.Args = newStringArrayLValue([]string{}, p, max)
	return commandConfig.Args
}

// DisableArgs disable unnamed argument configuration with the command.
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) DisableArgs() Arg {
	commandConfig.Args = None{}
	return commandConfig.Args
}

// AddValue registers an argument configuration with the command.
// The `name` argument represents the name of the argument.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddValue(name, shortName string, v Value) *Opt {

	// clean argument values
	name = removeWhitespaces(name)
	if name == "" {
		panic("name can not be empty")
	}
	if strings.HasPrefix(name, "-") {
		panic(ErrorUnsupportedFlag{name}) // check for - symbol, must be striped
	}

	if shortName == "-" || len(shortName) > 1 {
		panic(ErrorUnsupportedFlag{shortName}) // check for - symbol, must be striped
	}

	// check if argument is variadic
	// isVariadic := false
	// if _, ok := v.(VariadicValue); ok {
	// 	isVariadic = true
	// }

	// return if argument is already registered
	if _, ok := commandConfig.Opts[name]; ok {
		panic(name + " already registered")
	}

	// create `Arg` object
	opt := &Opt{
		Name:      name,
		ShortName: shortName,
		Value:     v,
	}

	// register argument with the command-config
	commandConfig.Opts[name] = opt

	if shortName != "" {
		if _, exist := commandConfig.short[shortName]; exist {
			panic(name + " alias " + shortName + " already registered")
		}
		commandConfig.short[shortName] = name
	}

	// store init value as default for restore on reset
	opt.defaultValue = opt.Value.Get()
	commandConfig.OptsOrder = append(commandConfig.OptsOrder, name)

	return opt
}

// AddFlag registers an bool (direct/inverted) flag with the command.
// The `name` argument represents the name of the argument.
// If value of the `name` argument starts with `no-` prefix, then it is a inverted flag.
// The `shortName` argument represents the short alias of the argument.
// If an argument with given `name` is already registered, then panic
// registered `*Opt` object returned.
func (commandConfig *CommandConfig) AddFlag(name, shortName string, b *bool) *Opt {
	*b, name = isInvertedFlag(name)
	v := newBoolValue(b)
	o := commandConfig.AddValue(name, shortName, v)
	o.IsBool = true
	o.IsInverted = *b
	return o
}

/*---------------------*/

// Opt type holds the structured information about a flag.
type Opt struct {
	Name       string // long name of the flag
	ShortName  string // short name of the flag
	Usage      string // help message
	IsBool     bool   // boolean flag
	IsInverted bool   // inverted boolean flag
	// IsVariadic   bool   // true if can take multiple values
	IsRequired  bool            // required value
	ValidValues map[string]bool // valid values

	Changed      bool        // if the user set the value (or if left to default)
	Value        Value       // value as set
	defaultValue interface{} // store init value, used for reset
}

// Set set
// Use with valid backend Value (may be slice) or values can be lost/corrupted
// `*Opt` object returned.
func (o *Opt) Set(s string) error {
	err := o.Value.Set(s, o.Changed)
	if err == nil {
		o.Changed = true
	}
	return err
}

// Reset reset opt (changed flag is cleared)
func (o *Opt) Reset() {
	o.Value.Reset(o.defaultValue)
	o.Changed = false
}

// SetVariadic enable/disable variadic
// Use with valid backend Value (may be slice) or values can be lost/corrupted
// `*Opt` object returned.
// func (o *Opt) SetVariadic(variadic bool) *Opt {
// 	o.IsVariadic = variadic
// 	return o
// }

// SetRequired enable/disable required
// `*Opt` object returned.
func (o *Opt) SetRequired(required bool) *Opt {
	o.IsRequired = required
	return o
}

// SetDefault set default value (as string)
// `*Opt` object returned.
// func (o *Opt) SetDefault(defaultValue string) *Opt {
// 	o.DefaultValue = trimWhitespaces(defaultValue)
// 	return o
// }

// SetValidValues set values for validate
// `*Opt` object returned.
func (o *Opt) SetValidValues(validValues []string) *Opt {
	o.ValidValues = map[string]bool{}
	for _, v := range validValues {
		o.ValidValues[v] = true
	}
	return o
}

func (o *Opt) Validate(s string) (isValid bool) {
	if len(o.ValidValues) > 0 {
		if !validateByValues(s, o.ValidValues) {
			return false
		}
	}
	return true
}

// SetUsage enable/disable required
// `*Opt` object returned.
func (o *Opt) SetUsage(usage string) *Opt {
	o.Usage = usage
	return o
}

/*---------------------*/
