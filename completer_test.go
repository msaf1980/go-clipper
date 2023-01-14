package clipper

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_splitQuoted(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{
			s:    `--version 1.0.1 --dir "/c/P F" --dir "/opt" -v`,
			want: []string{"--version", "1.0.1", "--dir", `/c/P F`, "--dir", `/opt`, "-v"},
		},
		{
			s:    ` --version 1.0.1 --dir "/c/P F" --dir "/opt" -v `,
			want: []string{"--version", "1.0.1", "--dir", `/c/P F`, "--dir", `/opt`, "-v"},
		},
		{
			s:    ` --version 1.0.1 -s "" -S " "  -v `,
			want: []string{"--version", "1.0.1", "-s", "", "-S", " ", "-v"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tt.s), func(t *testing.T) {
			if got := SplitQuoted(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestCompleter(t *testing.T) {
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
	registry := NewRegistry("clipper demo")

	// register the root command
	rootCommand, _ := registry.Register("", "root command help") // root command
	// rootCommand.AddArg("output", "")                    //
	rootCommand.AddFlag("force", "f", &rootForce, "force")           // --force, -f | default value: "false"
	rootCommand.AddFlag("verbose", "v", &rootVerbose, "verbose")     // --verbose, -v | default value: "false"
	rootCommand.AddString("dir", "d", "/var/users", &rootDir, "dir") // --dir <value> | default value: "/var/users"
	rootCommand.AddStringArgs(-1, &root, "root args")
	rootCommand.AddVersionHelper("version", "V", registry.Description, "0.0.1")

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

	timeVal, _ := time.Parse(time.RFC3339Nano, "2023-01-14T15:11:56.606523617+05:00")

	typesCommand, _ := registry.Register("types", "")                              // sub-command
	typesCommand.AddIntArray("num", "n", []int{1, 24, -2}, &typesNum, "int array") // --num, -n | default value: []int
	typesCommand.AddMultiFlag("verbose", "v", &typesVerbose, "")                   // --verbose, -v | default value: []
	typesCommand.AddTime("time", "", timeVal, &typesTime, time.RFC3339Nano, "time with time.RFC3339Nano layout")
	typesCommand.AddTimeFromString("stime", "", timeVal.Format(time.RFC3339Nano), &typesSTime, time.RFC3339Nano,
		"time from string with time.RFC3339Nano layout").
		SetCompeterValue(timeVal.Format(time.RFC3339Nano)) // compeleter value to default value
	typesCommand.AddTimeFromString("ltime", "", timeVal.Format(timeLayout), &typesLayoutTime, timeLayout,
		"time from string with "+timeLayout+" layout").
		SetCompeterValue(timeVal.Format(timeLayout)) // compeleter value to default value

	tests := []struct {
		line string
		want []string
	}{
		{
			line: "",
			want: []string{
				//root command flags
				"--force",
				// "-f",
				"--verbose",
				// "-v",
				"--dir",
				// "-d",
				"--version",
				// "-V",
				"--help",
				// all commands except root ("")
				"ghost",
				"info",
				"list",
				"types",
			},
		},
		{
			line: "g",
			want: []string{
				// all commands from "g"
				"ghost",
			},
		},
		{
			// no command "i", so can't complete
			line: "i -v",
			want: []string{},
		},
		{
			line: "-",
			want: []string{
				// all flags for root cmd
				"--force",
				// "-f",
				"--verbose",
				// "-v",
				"--dir",
				// "-d",
				"--version",
				// "-V",
				"--help",
			},
		},
		{
			line: "--v",
			want: []string{
				// all flags for root cmd from "--v"
				"--verbose",
				"--version",
			},
		},
		{
			// not flags for root cmd starts with "-v" except "-v"
			line: "-v",
			want: []string{},
		},
		{
			// command completed, add space for flags compete
			line: "info",
			want: []string{},
		},
		{
			line: "info ",
			want: []string{
				// all flags for info cmd
				"info --verbose",
				// "-v",
				"info --version",
				// "-V",
				"info --output",
				// "-o",
				"info --clean",
				// "-N",
				"info --help",
			},
		},
		{
			line: "info -",
			want: []string{
				// all flags for info cmd
				"info --verbose",
				// "-v",
				"info --version",
				// "-V",
				"info --output",
				// "-o",
				"info --clean",
				// "-N",
				"info --help",
			},
		},
		{
			line: "info --v",
			want: []string{
				// all flags for info cmd from "--v"
				"info --verbose",
				"info --version",
			},
		},
		{
			line: "info -v",
			// not flags for info cmd starts with "-v" except "-v"
			want: []string{},
		},
		{
			line: "info -v ",
			// not flags for info cmd starts with "-v" except "-v"
			want: []string{
				"info -v --version",
				"info -v --output",
				"info -v --clean",
				"info -v --help",
			},
		},
		{
			line: "info -V ",
			// version value helper
			want: []string{
				"info -V string",
			},
		},
		{
			line: "info --version ",
			// version value helper
			want: []string{
				"info --version string",
			},
		},
		{
			line: "types --n 1 --time ",
			// version value helper
			want: []string{
				"types --n 1 --time 2006-01-02T15:04:05.999999999Z07:00",
			},
		},
		{
			line: "types --n 1 --stime ",
			// version value helper
			want: []string{
				"types --n 1 --stime 2023-01-14T15:11:56.606523617+05:00",
			},
		},
		{
			line: "types --n 1 --ltime ",
			// version value helper
			want: []string{
				`types --n 1 --ltime "2023-01-14 15:11:56"`,
			},
		},
		// TODO: check for variadic and inverted flag
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := registry.Completer(tt.line)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCompleter_NoRoot(t *testing.T) {
	var (
		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		list, listDir []string

		typesNum     []int
		typesVerbose []bool
	)

	// create a new registry
	registry := NewRegistry("clipper demo")

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
	typesCommand.AddMultiFlag("verbose", "v", &typesVerbose, "")                   // --verbose, -v | default value: []

	tests := []struct {
		line string
		want []string
	}{
		{
			line: "",
			want: []string{
				// all commands
				"ghost",
				"info",
				"list",
				"types",
			},
		},
		{
			line: "g",
			want: []string{
				// all commands from "g"
				"ghost",
			},
		},
		{
			// no command "i", so can't complete
			line: "i -v",
			want: []string{},
		},
		{
			line: "-",
			// all flags for root cmd
			want: []string{},
		},
		{
			// all flags for root cmd from "--v"
			line: "--v",
			want: []string{},
		},
		{
			line: "-v",
			// not flags for root cmd starts with "-v" except "-v"
			want: []string{},
		},
		{
			line: "info",
			// command completed, add space for flags compete
			want: []string{},
		},
		{
			line: "info ",
			want: []string{
				// all flags for info cmd
				"info --verbose",
				// "-v",
				"info --version",
				// "-V",
				"info --output",
				// "-o",
				"info --clean",
				// "-N",
				"info --help",
			},
		},
		{
			line: "info -",
			want: []string{
				// all flags for info cmd
				"info --verbose",
				// "-v",
				"info --version",
				// "-V",
				"info --output",
				// "-o",
				"info --clean",
				// "-N",
				"info --help",
			},
		},
		{
			line: "info --v",
			want: []string{
				// all flags for info cmd from "--v"
				"info --verbose",
				"info --version",
			},
		},
		{
			// not flags for info cmd starts with "-v" except "-v"
			line: "info -v",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := registry.Completer(tt.line)
			assert.Equal(t, tt.want, got)
		})
	}
}
