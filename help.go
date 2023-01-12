package clipper

import (
	"fmt"
	"os"
	"sort"
)

func PrintHelp(registry *Registry, commandName string, commandConfig *CommandConfig) {

	if commandConfig == nil {
		fmt.Fprintf(os.Stderr, "\n  command %q not registered\n", commandName)
		return
	}

	if commandName == "" {
		if registry.Description != "" {
			fmt.Printf("%s\n\n\n", registry.Description)
		}
		fmt.Printf("Usage:\n  %s [command] [options]\n\n", os.Args[0])

		commands := make([]string, 0, len(registry.Commands))
		for command := range registry.Commands {
			// exclude root command
			if command != "" {
				commands = append(commands, command)
			}
		}
		sort.Strings(commands)
		fmt.Println("Available Commands:")
		for _, command := range commands {
			fmt.Printf("  %s\t%s\n", command, registry.Commands[command].Help)
		}

		fmt.Printf("\n%s\n\n", commandConfig.Help)
	} else {
		fmt.Printf("%s\n\n", commandConfig.Help)
	}

	fmt.Println("Flags:")
	for _, flagName := range commandConfig.OptsOrder {
		opt := commandConfig.Opts[flagName]
		nameAndArgs := "--" + opt.Name
		shortName := ""
		if opt.ShortName != "" {
			shortName += "-" + opt.ShortName + " | "
		}
		if !opt.IsBool {
			nameAndArgs += " " + opt.Value.Type()
		}
		fmt.Printf("  %5s%-25s\t%s", shortName, nameAndArgs, opt.Help)
		if opt.IsRequired {
			fmt.Printf(" (required)\n")
			// } else if v, ok := opt.defaultValue.(bool); ok {
			// 	fmt.Printf(" (default %v)\n", v)
		} else {
			fmt.Printf(" (default: %q)\n", opt.Value.String())
		}
	}
	fmt.Printf("  -h | %-25s\thelp", "--help")
	if commandName == "" {
		fmt.Println("")
	} else {
		fmt.Printf(" for %s\n", commandName)
	}

	if _, ok := commandConfig.Args.(None); !ok {
		fmt.Printf("\nUnnamed args (%s)", commandConfig.ArgsHelp)
		min := commandConfig.Args.MinLen()
		max := commandConfig.Args.MaxLen()
		if min > 0 || max >= 0 {
			fmt.Println(":")
			if min > 0 {
				fmt.Printf(" min length %d", min)
			}
			if max >= 0 {
				fmt.Printf(" max length %d", max)
			}
		}
		fmt.Println("")
	}
}