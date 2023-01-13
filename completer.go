package clipper

import (
	"sort"
	"strings"
)

func isBlanc(s string) bool {
	for _, c := range s {
		if c != ' ' {
			return false
		}
	}
	return true
}

func appendCommands(registry *Registry, c []string, last, line string) []string {
	start := len(c)
	for name := range registry.Commands {
		if name != "" && (last == "" || strings.HasPrefix(name, last)) {
			if line == "" {
				c = append(c, name)
			} else {
				c = append(c, line[:len(line)-len(last)]+name)
			}
		}
	}
	sort.Strings(c[start:])
	return c
}

func appendFlags(commandConfig *CommandConfig, c []string, last, line string) []string {
	for _, n := range commandConfig.OptsOrder {
		opt := commandConfig.Opts[n]
		name := "--" + opt.Name
		if name != last && strings.HasPrefix(name, last) {
			if line == "" {
				c = append(c, name)
			} else {
				c = append(c, line[:len(line)-len(last)]+name)
			}
		}
		// if opt.ShortName != "" {
		// 	name = "-" + opt.ShortName
		// 	if name != last && strings.HasPrefix(name, last) {
		// 		if line == "" {
		// 			c = append(c, name)
		// 		} else {
		// 			c = append(c, line[:len(line)-len(last)]+name)
		// 		}
		// 	}
		// }
	}
	if commandConfig.version != nil {
		if commandConfig.version.name != last && strings.HasPrefix(commandConfig.version.name, last) {
			if line == "" {
				c = append(c, commandConfig.version.name)
			} else {
				c = append(c, line[:len(line)-len(last)]+commandConfig.version.name)
			}
		}
		// if commandConfig.version.shortName != last &&
		// 	commandConfig.version.shortName != "" && strings.HasPrefix(commandConfig.version.shortName, last) {

		// 	if line == "" {
		// 		c = append(c, commandConfig.version.shortName)
		// 	} else {
		// 		c = append(c, line[:len(line)-len(last)]+commandConfig.version.shortName)
		// 	}
		// }
	}
	if last != "--help" && strings.HasPrefix("--help", last) {
		if line == "" {
			c = append(c, "--help")
		} else {
			c = append(c, line[:len(line)-len(last)]+"--help")
		}
	}

	return c
}

func SplitQuoted(s string) []string {
	var (
		startQuote bool
	)
	sv := make([]string, 0, 4)
	start := 0

	for i, c := range s {
		switch c {
		case ' ':
			if !startQuote {
				v := s[start:i]
				if !isBlanc(v) {
					sv = append(sv, v)
				}
				start = i + 1
			}
		case '"':
			startQuote = !startQuote
		}
	}
	v := s[start:]
	if !isBlanc(v) {
		sv = append(sv, v)
	}

	return sv
}

// CompleterAppended return slice of completer variants only without initial line
func (registry *Registry) Completer(line string) (c []string) {

	var commandName string
	sv := SplitQuoted(line)

	if len(sv) == 0 {
		c = make([]string, 0, len(registry.Commands)+4)
		commandConfig, ok := registry.Commands[commandName]
		if ok {
			c = appendFlags(commandConfig, c, "", "")
		}
		c = appendCommands(registry, c, "", "")
		return
	}
	if len(sv) > 0 {
		if !strings.HasPrefix(sv[0], "-") {
			commandName = sv[0]
		}
	}

	commandConfig, ok := registry.Commands[commandName]
	if !ok {
		c = make([]string, 0)
		if len(sv) == 1 && !strings.HasPrefix(sv[0], "-") {
			c = appendCommands(registry, c, sv[0], "")
		}
		return
	}

	if strings.HasSuffix(line, " ") {
		// return flags for command

		// new command, exclude already added
		sm := make(map[string]bool)
		for _, s := range sv {
			if isFlag(s) {
				sm[s] = true
			}
		}

		c = make([]string, 0, len(commandConfig.Opts)*2)
		c = appendFlags(commandConfig, c, "", "")

	} else {
		// complete last arg
		last := sv[len(sv)-1]
		c = make([]string, 0)
		if strings.HasPrefix(last, "-") {
			// complete flag
			c = appendFlags(commandConfig, c, last, "")
		} else if commandConfig == nil && len(sv) == 1 {
			c = appendCommands(registry, c, last, "")
		}
	}

	return
}

// CompleterAppended return slice of completer variants with prepended initial line
func (registry *Registry) CompleterAppended(line string) (c []string) {

	var commandName string
	sv := SplitQuoted(line)

	if len(sv) == 0 {
		c = make([]string, 0, len(registry.Commands)+4)
		commandConfig, ok := registry.Commands[commandName]
		if ok {
			c = appendFlags(commandConfig, c, "", line)
		}
		c = appendCommands(registry, c, "", line)
		return
	}
	if len(sv) > 0 {
		if !strings.HasPrefix(sv[0], "-") {
			commandName = sv[0]
		}
	}

	commandConfig, ok := registry.Commands[commandName]
	if !ok {
		c = make([]string, 0)
		if len(sv) == 1 && !strings.HasPrefix(sv[0], "-") {
			c = appendCommands(registry, c, commandName, line)
		}
		return
	}

	if strings.HasSuffix(line, " ") {
		// return flags for command

		// new command, exclude already added
		sm := make(map[string]bool)
		for _, s := range sv {
			if isFlag(s) {
				sm[s] = true
			}
		}

		c = make([]string, 0, len(commandConfig.Opts)*2)
		c = appendFlags(commandConfig, c, "", line)

	} else {
		// complete last arg
		last := sv[len(sv)-1]
		c = make([]string, 0)
		if strings.HasPrefix(last, "-") {
			// complete flag
			c = appendFlags(commandConfig, c, last, line)
		} else if commandConfig == nil && len(sv) == 1 {
			c = appendCommands(registry, c, last, line)
		}
	}

	return
}
