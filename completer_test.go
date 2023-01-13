package clipper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitQuoted(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{
			s:    `--version 1.0.1 --dir "/c/P F" --dir "/opt" -v`,
			want: []string{"--version", "1.0.1", "--dir", `"/c/P F"`, "--dir", `"/opt"`, "-v"},
		},
		{
			s:    ` --version 1.0.1 --dir "/c/P F" --dir "/opt" -v `,
			want: []string{"--version", "1.0.1", "--dir", `"/c/P F"`, "--dir", `"/opt"`, "-v"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %s", i, tt.s), func(t *testing.T) {
			if got := splitQuoted(tt.s); !reflect.DeepEqual(got, tt.want) {
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

		typesNum     []int
		typesVerbose []bool
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

	typesCommand, _ := registry.Register("types", "")                              // sub-command
	typesCommand.AddIntArray("num", "n", []int{1, 24, -2}, &typesNum, "int array") // --num, -n | default value: []int
	typesCommand.AddMultiFlag("verbose", "v", &typesVerbose, "")                   // --verbose, -v | default value: []

	tests := map[string][]string{
		"": {
			//root command flags
			"--force", "-f",
			"--verbose", "-v",
			"--dir", "-d",
			"--version", "-V",
			// all commands except root ("")
			"ghost",
			"info",
			"list",
			"types",
		},
		"g": {
			// all commands from "g"
			"ghost",
		},
		"i -v": {
			// no command "i", so can't complete
		},
		"-": {
			// all flags for root cmd
			"--force", "-f",
			"--verbose", "-v",
			"--dir", "-d",
			"--version", "-V",
		},
		"--v": {
			// all flags for root cmd from "--v"
			"--verbose",
			"--version",
		},
		"-v": {
			// not flags for root cmd starts with "-v" except "-v"
		},
		"info": {
			// command completed, add space for flags compete
		},
		"info ": {
			// all flags for info cmd
			"--verbose", "-v",
			"--version", "-V",
			"--output", "-o",
			"--clean", "-N",
		},
		"info -": {
			// all flags for info cmd
			"--verbose", "-v",
			"--version", "-V",
			"--output", "-o",
			"--clean", "-N",
		},
		"info --v": {
			// all flags for info cmd from "--v"
			"--verbose",
			"--version",
		},
		"info -v": {
			// not flags for info cmd starts with "-v" except "-v"
		},
	}
	for line, want := range tests {
		t.Run(line, func(t *testing.T) {
			got := registry.Completer(line)
			assert.Equal(t, want, got)
		})
	}
}
