package main

import (
	"fmt"
	"os"

	"github.com/msaf1980/go-clipper"
)

func main() {
	var (
		rootForce, rootVerbose bool
		rootVersion, rootDir   string
		root                   []string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		list, listDir []string
	)

	// create a new registry
	registry := clipper.NewRegistry()

	// register the root command
	if _, ok := os.LookupEnv("NO_ROOT"); !ok {
		rootCommand, _ := registry.Register("") // root command
		// rootCommand.AddArg("output", "")                    //
		rootCommand.AddFlag("force", "f", &rootForce)             // --force, -f | default value: "false"
		rootCommand.AddFlag("verbose", "v", &rootVerbose)         // --verbose, -v | default value: "false"
		rootCommand.AddString("version", "V", "", &rootVersion)   // --version, -V | default value: ""
		rootCommand.AddString("dir", "d", "/var/users", &rootDir) // --dir <value> | default value: "/var/users"
		rootCommand.AddStringArgs(-1, &root)
	}

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info")              // sub-command
	infoCommand.AddFlag("verbose", "v", &infoVerbose)        // --verbose, -v | default value: "false"
	infoCommand.AddString("version", "V", "", &infoVersion). // --version, -V | default value: "false"
									SetValidValues([]string{"", "1.0.1", "2.0.0"})
	infoCommand.AddString("output", "o", "./", &infoOutput) // --output, -o <value> | default value: "./"
	infoCommand.AddFlag("no-clean", "N", &infoNoClean)      // --no-clean | default value: "true"

	listCommand, _ := registry.Register("list")                     // sub-command
	listCommand.AddStringArray("dir", "d", []string{"a"}, &listDir) // --output, -o <value> | default value: "./"
	listCommand.AddStringArgs(-1, &list)
	// listCommand.Args.SetMinLen(1) // set minimal length (at parse step) | default value: 0

	// register the `ghost` sub-command
	registry.Register("ghost")

	/*----------------*/

	// parse command-line arguments
	command, err := registry.Parse(os.Args[1:])

	/*----------------*/

	// check for error
	if err != nil {
		fmt.Printf("error => %#v\n", err)
		return
	}

	// get executed sub-command name
	fmt.Printf("sub-command => %#v\n  Dump variables\n", command)
	c := registry[command]
	for _, name := range c.OptsOrder {
		opt := c.Opts[name]
		fmt.Printf("    %s=%q\n", name, opt.Value.String())
	}
	// get unnamed args
	if args := c.Args.String(); args != "" {
		fmt.Printf("    args=%s\n", args)
	}
}
