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

		typesNum []int
	)

	// create a new registry
	registry := clipper.NewRegistry("clipper demo")

	// register the root command
	if _, ok := os.LookupEnv("NO_ROOT"); !ok {
		rootCommand, _ := registry.Register("", "root command help") // root command
		// rootCommand.AddArg("output", "")                    //
		rootCommand.AddFlag("force", "f", &rootForce, "force")                 // --force, -f | default value: "false"
		rootCommand.AddFlag("verbose", "v", &rootVerbose, "verbose")           // --verbose, -v | default value: "false"
		rootCommand.AddString("version", "V", "", &rootVersion, "set version") // --version, -V | default value: ""
		rootCommand.AddString("dir", "d", "/var/users", &rootDir, "dir")       // --dir <value> | default value: "/var/users"
		rootCommand.AddStringArgs(-1, &root, "root args")
	}

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "info help")                // sub-command
	infoCommand.AddFlag("verbose", "v", &infoVerbose, "verbose")            // --verbose, -v | default value: "false"
	infoCommand.AddString("version", "V", "", &infoVersion, "set version"). // --version, -V | default value: ""
										SetValidValues([]string{"", "1.0.1", "2.0.0"}). // valid versions
										SetRequired(true)                               // version are required
	infoCommand.AddString("output", "o", "./", &infoOutput, "output dir") // --output, -o <value> | default value: "./"
	infoCommand.AddFlag("no-clean", "N", &infoNoClean, "disable clean")   // --no-clean | default value: "true"

	listCommand, _ := registry.Register("list", "list help")               // sub-command
	listCommand.AddStringArray("dir", "d", []string{"a"}, &listDir, "dir") // --output, -o <value> | default value: "./"
	listCommand.AddStringArgs(-1, &list, "list args")
	// listCommand.Args.SetMinLen(1) // set minimal length (at parse step) | default value: 0

	// register the `ghost` sub-command
	registry.Register("ghost", "ghost help")

	typesCommand, _ := registry.Register("types", "")                              // sub-command
	typesCommand.AddIntArray("num", "n", []int{1, 24, -2}, &typesNum, "int array") // --num, -n | default value: []int

	/*----------------*/

	// parse command-line arguments
	command, _, err := registry.Parse(os.Args[1:], true)

	/*----------------*/

	// check for error
	if err != nil {
		fmt.Printf("error => %#v\n", err)
		return
	}

	// get executed sub-command name
	fmt.Printf("sub-command => %#v\n  Dump variables\n", command)
	c := registry.Commands[command]
	for _, name := range c.OptsOrder {
		opt := c.Opts[name]
		fmt.Printf("    %s=%q\n", name, opt.Value.String())
	}
	// get unnamed args
	if args := c.Args.String(); args != "" {
		fmt.Printf("    args=%s\n", args)
	}
}
