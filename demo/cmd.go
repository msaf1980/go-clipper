package main

import (
	"fmt"
	"os"
	"time"

	"github.com/msaf1980/go-clipper"
)

var (
	VERSION = "0.0.1"
)

func main() {
	var (
		rootForce, rootVerbose bool
		rootDir                string
		root                   []string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		list, listDir []string

		typesNum        []int
		typesVerbose    []bool
		typesTime       time.Time
		typesSTime      time.Time
		timeLayout      = "2006-01-02 15:04:05"
		typesLayoutTime time.Time
	)

	// create a new registry
	registry := clipper.NewRegistry("clipper demo")

	// register the root command
	if _, ok := os.LookupEnv("NO_ROOT"); !ok {
		rootCommand, _ := registry.Register("", "root command help") // root command
		// rootCommand.AddArg("output", "")                    //
		rootCommand.AddFlag("force", "f", &rootForce, "force")           // --force, -f | default value: "false"
		rootCommand.AddFlag("verbose", "v", &rootVerbose, "verbose")     // --verbose, -v | default value: "false"
		rootCommand.AddString("dir", "d", "/var/users", &rootDir, "dir") // --dir <value> | default value: "/var/users"
		rootCommand.AddStringArgs(-1, &root, "root args")
		rootCommand.AddVersionHelper("version", "V", registry.Description, VERSION)
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
	ghostCommand, _ := registry.Register("ghost", "ghost help")
	ghostCommand.AddVersionHelper("version", "V", registry.Description, VERSION)

	tm, _ := time.Parse(time.RFC3339Nano, "2023-01-15T21:41:29.98589753+05:00")

	typesCommand, _ := registry.Register("types", "")                               // sub-command
	typesCommand.AddIntArray("num", "n", []int{1, 24, -2}, &typesNum, "int array"). // --num, -n | default value: []int
											AttachEnv("TYPES_NUM") // try to read value from os env var
	typesCommand.AddMultiFlag("verbose", "v", &typesVerbose, "") // --verbose, -v | default value: []
	typesCommand.AddTime("time", "", tm, &typesTime, time.RFC3339Nano, "time with time.RFC3339Nano layout")
	typesCommand.AddTimeFromString("stime", "", tm.Format(time.RFC3339Nano), &typesSTime, time.RFC3339Nano, "time from string with time.RFC3339Nano layout")
	typesCommand.AddTimeFromString("ltime", "", tm.Format(timeLayout), &typesLayoutTime, timeLayout,
		"time from string with "+timeLayout+" layout").
		SetCompeterValue(tm.Format(timeLayout)) // compeleter value to default value

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
