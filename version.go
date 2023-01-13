package clipper

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// -- string Value
type versionHelper struct {
	name        string
	shortName   string
	description string
	commandName string
	value       string
}

func newVersionHelper(name, shortName, description, commandName, version string) *versionHelper {
	if version == "" {
		version = "<unset>"
	}
	if !strings.HasPrefix(name, "--") {
		name = "--" + name
	}
	if shortName != "" && !strings.HasPrefix(shortName, "-") {
		shortName = "-" + shortName[:1]
	}
	if description == "" {
		description = path.Base(os.Args[0])
	}
	return &versionHelper{
		name:        name,
		shortName:   shortName,
		value:       version,
		description: description,
		commandName: commandName,
	}
}

func (v *versionHelper) print() {
	if v.commandName == "" {
		fmt.Printf("%s (%s)\n", v.value, v.description)
	} else {
		fmt.Printf("%s (%s) command %q\n", v.value, v.description, v.commandName)
	}
}

// AddVersion method registers a version callback.
func (commandConfig *CommandConfig) AddVersionHelper(name, shortName, description, version string) {
	commandConfig.version = newVersionHelper(name, shortName, description, commandConfig.Name, version)
}
