package clipper

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*----------------*/

func TestRegistry_Parse_RootArgs(t *testing.T) {
	var (
		rootForce, rootVerbose bool
		rootVersion, rootDir   string
		root                   []string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		listVerbose   bool
		listDir, list []string
	)
	store := make(map[string]map[string]interface{})
	storeValues := make(map[string]map[string]*Opt)

	store[""] = make(map[string]interface{})
	storeValues[""] = make(map[string]*Opt)

	store["info"] = make(map[string]interface{})
	storeValues["info"] = make(map[string]*Opt)

	store["ghost"] = make(map[string]interface{})
	storeValues["ghost"] = make(map[string]*Opt)

	store["list"] = make(map[string]interface{})
	storeValues["list"] = make(map[string]*Opt)

	// create a new registry
	registry := NewRegistry("")

	// register the root command
	rootCommand, _ := registry.Register("", "") // root command
	// rootCommand.AddArg("output", "")                    //
	// rootForceDefault := false
	store[""]["force"] = &rootForce
	storeValues[""]["force"] = rootCommand.AddFlag("force", "f", &rootForce, "") // --force, -f | default value: "false"
	// rootVerboseDefault := false
	store[""]["verbose"] = &rootVerbose
	storeValues[""]["verbose"] = rootCommand.AddFlag("verbose", "v", &rootVerbose, "") // --verbose, -v | default value: "false"
	rootVersionDefault := ""
	store[""]["version"] = &rootVersion
	storeValues[""]["version"] = rootCommand.AddString("version", "V", rootVersionDefault, &rootVersion, "") // --version, -V | default value: ""
	rootDirDefault := "/var/users"
	store[""]["dir"] = &rootDir
	storeValues[""]["dir"] = rootCommand.AddString("dir", "d", rootDirDefault, &rootDir, "") // --dir <value> | default value: "/var/users"
	store[""][""] = &root
	rootCommand.AddStringArgs(-1, &root, "")

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "") // sub-command
	// infoVerboseDefault := false
	store["info"]["verbose"] = &infoVerbose
	storeValues["info"]["verbose"] = infoCommand.AddFlag("verbose", "v", &infoVerbose, "") // --verbose, -v | default value: "false"                // --verbose, -v | default value: "false"
	infoVersionDefault := ""
	store["info"]["version"] = &infoVersion
	storeValues["info"]["version"] = infoCommand.AddString("version", "V", infoVersionDefault, &infoVersion, ""). // --version, -V | default value: ""
															SetValidValues([]string{"", "1.0.1", "2.0.0"})
	infoOutputDefault := "./"
	store["info"]["output"] = &infoOutput
	storeValues["info"]["output"] = infoCommand.AddString("output", "o", infoOutputDefault, &infoOutput, "") // --output, -o <value> | default value: "./"
	// infoNoCleanDefault := true
	store["info"]["clean"] = &infoNoClean
	storeValues["info"]["clean"] = infoCommand.AddFlag("no-clean", "N", &infoNoClean, "") // --no-clean | default value: "true"

	// register the `list` sub-command
	listCommand, _ := registry.Register("list", "") // sub-command
	store["list"]["verbose"] = &listVerbose
	storeValues["list"]["verbose"] = listCommand.AddFlag("verbose", "v", &listVerbose, "") // --verbose, -v | default value: "false"
	listDirDefault := []string{"a"}
	store["list"]["dir"] = &listDir
	storeValues["list"]["dir"] = listCommand.AddStringArray("dir", "d", listDirDefault, &listDir, "") // --output, -o <value> | default value: "./"
	store["list"][""] = &list
	listCommand.AddStringArgs(-1, &list, "")

	// register the `ghost` sub-command
	registry.Register("ghost", "")

	tests := []struct {
		values   []string
		want     string
		wantVars map[string]interface{}
		wantArgs []string
		wantErr  string
	}{
		{
			// unknown must detect ad uunamed arg, not a unregistered command
			values: []string{"unknown"},
			want:   "",
			wantVars: map[string]interface{}{
				"version": "",
				"verbose": false,
				"force":   false,
				"dir":     "/var/users",
			},
			wantArgs: []string{"unknown"},
		},
		{
			// rootarg must detect ad uunamed arg, not a unregistered command
			values: []string{"rootarg", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
			want:   "",
			wantVars: map[string]interface{}{
				"version": "1.0.1",
				"verbose": true,
				"force":   true,
				"dir":     "./sub/dir",
			},
			wantArgs: []string{"rootarg"},
		},
		{
			values: []string{"-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
			want:   "",
			wantVars: map[string]interface{}{
				"version": "1.0.1",
				"verbose": true,
				"force":   true,
				"dir":     "./sub/dir",
			},
		},
		{
			values: []string{"-V=1.0.1", "-v", "--dir=./sub/dir"},
			want:   "",
			wantVars: map[string]interface{}{
				"version": "1.0.1",
				"verbose": true,
				"force":   false, //default
				"dir":     "./sub/dir",
			},
		},
		{
			values: []string{"info", "-V=1.0.1", "-v", "--no-clean", "--output=/sub/dir"},
			want:   "info",
			wantVars: map[string]interface{}{
				"version": "1.0.1",
				"verbose": true,
				"clean":   false,
				"output":  "/sub/dir",
			},
		},
		{
			values: []string{"info", "--output=/sub/dir"},
			want:   "info",
			wantVars: map[string]interface{}{
				"version": "",    //default
				"verbose": false, //default
				"clean":   true,  //default
				"output":  "/sub/dir",
			},
		},
		{
			values:  []string{"info", "-V=1.0.2", "-v", "--no-clean", "--output=/sub/dir"},
			want:    "info",
			wantErr: "unsupported value version",
		},
		{
			values: []string{"list", "-d", "f", "--dir", `e,"a b"`},
			want:   "list",
			wantVars: map[string]interface{}{
				"dir":     []string{"f", "e", "a b"},
				"verbose": false,
			},
		},
		{
			values: []string{"list", "student", "-d", "f", "--dir", `e,"a b"`, "manager"},
			want:   "list",
			wantVars: map[string]interface{}{
				"dir":     []string{"f", "e", "a b"},
				"verbose": false,
			},
			wantArgs: []string{"student", "manager"},
		},
		{
			// variadic
			values: []string{"list", "--dir...", "e", "a b", "-v"},
			want:   "list",
			wantVars: map[string]interface{}{
				"dir":     []string{"e", "a b"},
				"verbose": true,
			},
			wantArgs: []string{"student", "manager"},
		},
		{
			// no variadic for short-style flag
			values:  []string{"list", "-d...", "/data1", "/data2"},
			want:    "list",
			wantErr: "unsupported flag \"-d...\"",
		},
		{
			values:   []string{"ghost"},
			want:     "ghost",
			wantVars: map[string]interface{}{},
		},
		{
			values:  []string{"ghost", "-v"},
			want:    "ghost",
			wantErr: `unknown flag "-v" found`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %v", i, tt.values), func(t *testing.T) {
			registry.Reset()
			got, _, err := registry.Parse(tt.values, true)
			if err == nil && tt.wantErr != "" {
				t.Errorf("Registry.Parse() wantErr %q", tt.wantErr)
				return
			} else if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Registry.Parse() error = %q, wantErr %q", err.Error(), tt.wantErr)
					return
				}
			} else if got == tt.want {
				gotV := registry.Commands[got]
				for k, o := range gotV.Opts {
					wantVar := tt.wantVars[k]
					assert.Equal(t, wantVar, o.Value.Get(), k)
					assert.Equal(t, o.Value.Get(), storeValues[got][k].Value.Get(), k)
					switch v := store[got][k].(type) {
					case *string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *bool:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]string:
						assert.Equal(t, o.Value.Get(), *v, k)
					default:
						t.Fatalf("%T is unhandled", v)
					}
				}
				a := store[got][""]
				if len(tt.wantArgs) > 0 && a != nil && len(*(a.(*[]string))) > 0 {
					assert.Equal(t, &tt.wantArgs, a, "args")
				}
			} else {
				t.Errorf("Registry.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*----------------*/

func TestRegistry_Parse_RootNoArgs(t *testing.T) {
	var (
		rootForce, rootVerbose bool
		rootVersion, rootDir   string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string
	)
	store := make(map[string]map[string]interface{})
	storeValues := make(map[string]map[string]*Opt)

	store[""] = make(map[string]interface{})
	storeValues[""] = make(map[string]*Opt)

	store["info"] = make(map[string]interface{})
	storeValues["info"] = make(map[string]*Opt)

	store["ghost"] = make(map[string]interface{})
	storeValues["ghost"] = make(map[string]*Opt)

	// create a new registry
	registry := NewRegistry("")

	// register the root command
	rootCommand, _ := registry.Register("", "") // root command
	// rootCommand.AddArg("output", "")                    //
	// rootForceDefault := false
	store[""]["force"] = &rootForce
	storeValues[""]["force"] = rootCommand.AddFlag("force", "f", &rootForce, "") // --force, -f | default value: "false"
	// rootVerboseDefault := false
	store[""]["verbose"] = &rootVerbose
	storeValues[""]["verbose"] = rootCommand.AddFlag("verbose", "v", &rootVerbose, "") // --verbose, -v | default value: "false"
	rootVersionDefault := ""
	store[""]["version"] = &rootVersion
	storeValues[""]["version"] = rootCommand.AddString("version", "V", rootVersionDefault, &rootVersion, "") // --version, -V | default value: ""
	rootDirDefault := "/var/users"
	store[""]["dir"] = &rootDir
	storeValues[""]["dir"] = rootCommand.AddString("dir", "d", rootDirDefault, &rootDir, "") // --dir <value> | default value: "/var/users"

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "") // sub-command
	// infoVerboseDefault := false
	store["info"]["verbose"] = &infoVerbose
	storeValues["info"]["verbose"] = infoCommand.AddFlag("verbose", "v", &infoVerbose, "") // --verbose, -v | default value: "false"                // --verbose, -v | default value: "false"
	infoVersionDefault := ""
	store["info"]["version"] = &infoVersion
	storeValues["info"]["version"] = infoCommand.AddString("version", "V", infoVersionDefault, &infoVersion, ""). // --version, -V | default value: ""
															SetValidValues([]string{"", "1.0.1", "2.0.0"})
	infoOutputDefault := "./"
	store["info"]["output"] = &infoOutput
	storeValues["info"]["output"] = infoCommand.AddString("output", "o", infoOutputDefault, &infoOutput, "") // --output, -o <value> | default value: "./"
	// infoNoCleanDefault := true
	store["info"]["clean"] = &infoNoClean
	storeValues["info"]["clean"] = infoCommand.AddFlag("no-clean", "N", &infoNoClean, "") // --no-clean | default value: "true"

	// register the `ghost` sub-command
	registry.Register("ghost", "")

	tests := []struct {
		values   []string
		want     string
		wantVars map[string]interface{}
		wantArgs []string
		wantErr  string
	}{
		{
			values:  []string{"unknown"},
			want:    "",
			wantErr: `unsupported flag "unknown" found`,
		},
		{
			// rootarg must detect ad uunamed arg, not a unregistered command
			values:  []string{"rootarg", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
			want:    "",
			wantErr: `unsupported flag "rootarg" found`,
		},
		{
			values:   []string{"ghost"},
			want:     "ghost",
			wantVars: map[string]interface{}{},
		},
		{
			values:  []string{"ghost", "-v"},
			want:    "ghost",
			wantErr: `unknown flag "-v" found`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %v", i, tt.values), func(t *testing.T) {
			registry.Reset()
			got, _, err := registry.Parse(tt.values, true)
			if err == nil && tt.wantErr != "" {
				t.Errorf("Registry.Parse() wantErr %q", tt.wantErr)
				return
			} else if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Registry.Parse() error = %q, wantErr %q", err.Error(), tt.wantErr)
					return
				}
			} else if got == tt.want {
				gotV := registry.Commands[got]
				for k, o := range gotV.Opts {
					wantVar := tt.wantVars[k]
					assert.Equal(t, wantVar, o.Value.Get(), k)
					assert.Equal(t, o.Value.Get(), storeValues[got][k].Value.Get(), k)
					switch v := store[got][k].(type) {
					case *string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *bool:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]string:
						assert.Equal(t, o.Value.Get(), *v, k)
					default:
						t.Fatalf("%T is unhandled", v)
					}
				}
				a := store[got][""]
				if len(tt.wantArgs) > 0 && a != nil && len(*(a.(*[]string))) > 0 {
					assert.Equal(t, &tt.wantArgs, a, "args")
				}
			} else {
				t.Errorf("Registry.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*----------------*/

func TestRegistry_Parse_NoRoot(t *testing.T) {
	var (
		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string
	)
	store := make(map[string]map[string]interface{})
	storeValues := make(map[string]map[string]*Opt)

	store["info"] = make(map[string]interface{})
	storeValues["info"] = make(map[string]*Opt)

	store["ghost"] = make(map[string]interface{})
	storeValues["ghost"] = make(map[string]*Opt)

	// create a new registry
	registry := NewRegistry("")

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "") // sub-command
	// infoVerboseDefault := false
	store["info"]["verbose"] = &infoVerbose
	storeValues["info"]["verbose"] = infoCommand.AddFlag("verbose", "v", &infoVerbose, "") // --verbose, -v | default value: "false"                // --verbose, -v | default value: "false"
	infoVersionDefault := ""
	store["info"]["version"] = &infoVersion
	storeValues["info"]["version"] = infoCommand.AddString("version", "V", infoVersionDefault, &infoVersion, ""). // --version, -V | default value: ""
															SetValidValues([]string{"", "1.0.1", "2.0.0"}).SetRequired(true)
	infoOutputDefault := "./"
	store["info"]["output"] = &infoOutput
	storeValues["info"]["output"] = infoCommand.AddString("output", "o", infoOutputDefault, &infoOutput, "") // --output, -o <value> | default value: "./"
	// infoNoCleanDefault := true
	store["info"]["clean"] = &infoNoClean
	storeValues["info"]["clean"] = infoCommand.AddFlag("no-clean", "N", &infoNoClean, "") // --no-clean | default value: "true"

	// register the `ghost` sub-command
	registry.Register("ghost", "")

	tests := []struct {
		values   []string
		want     string
		wantVars map[string]interface{}
		wantArgs []string
		wantErr  string
	}{
		{
			values:  []string{"unknown"},
			want:    "",
			wantErr: `unknown command "unknown" found`,
		},
		{
			// rootarg must detect ad uunamed arg, not a unregistered command
			values:  []string{"rootarg", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
			want:    "",
			wantErr: `unknown command "rootarg" found`,
		},
		{
			// version are required
			values: []string{"info", "-V", "1.0.1", "-v"},
			want:   "info",
			wantVars: map[string]interface{}{
				"clean":   true,
				"verbose": true,
				"version": "1.0.1",
				"output":  "./",
			},
		},
		{
			// version are required
			values:  []string{"info"},
			want:    "info",
			wantErr: `required flag "version" not found in the arguments`,
		},
		{
			values:   []string{"ghost"},
			want:     "ghost",
			wantVars: map[string]interface{}{},
		},
		{
			values:  []string{"ghost", "-v"},
			want:    "ghost",
			wantErr: `unknown flag "-v" found`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %v", i, tt.values), func(t *testing.T) {
			registry.Reset()
			got, _, err := registry.Parse(tt.values, true)
			if err == nil && tt.wantErr != "" {
				t.Errorf("Registry.Parse() wantErr %q", tt.wantErr)
				return
			} else if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Registry.Parse() error = %q, wantErr %q", err.Error(), tt.wantErr)
					return
				}
			} else if got == tt.want {
				gotV := registry.Commands[got]
				for k, o := range gotV.Opts {
					wantVar := tt.wantVars[k]
					assert.Equal(t, wantVar, o.Value.Get(), k)
					assert.Equal(t, o.Value.Get(), storeValues[got][k].Value.Get(), k)
					switch v := store[got][k].(type) {
					case *string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *bool:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]string:
						assert.Equal(t, o.Value.Get(), *v, k)
					default:
						t.Fatalf("%T is unhandled", v)
					}
				}
				a := store[got][""]
				if len(tt.wantArgs) > 0 && a != nil && len(*(a.(*[]string))) > 0 {
					assert.Equal(t, &tt.wantArgs, a, "args")
				}
			} else {
				t.Errorf("Registry.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*----------------*/

func TestRegistry_Parse_Limits(t *testing.T) {
	var (
		rootForce, rootVerbose bool
		rootVersion, rootDir   string
		root                   []string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		ghost []string
	)
	store := make(map[string]map[string]interface{})
	storeValues := make(map[string]map[string]*Opt)

	store[""] = make(map[string]interface{})
	storeValues[""] = make(map[string]*Opt)

	store["info"] = make(map[string]interface{})
	storeValues["info"] = make(map[string]*Opt)

	store["ghost"] = make(map[string]interface{})
	storeValues["ghost"] = make(map[string]*Opt)

	// create a new registry
	registry := NewRegistry("")

	// register the root command
	rootCommand, _ := registry.Register("", "") // root command
	// rootCommand.AddArg("output", "")                    //
	// rootForceDefault := false
	store[""]["force"] = &rootForce
	storeValues[""]["force"] = rootCommand.AddFlag("force", "f", &rootForce, "") // --force, -f | default value: "false"
	// rootVerboseDefault := false
	store[""]["verbose"] = &rootVerbose
	storeValues[""]["verbose"] = rootCommand.AddFlag("verbose", "v", &rootVerbose, "") // --verbose, -v | default value: "false"
	rootVersionDefault := ""
	store[""]["version"] = &rootVersion
	storeValues[""]["version"] = rootCommand.AddString("version", "V", rootVersionDefault, &rootVersion, "") // --version, -V | default value: ""
	rootDirDefault := "/var/users"
	store[""]["dir"] = &rootDir
	storeValues[""]["dir"] = rootCommand.AddString("dir", "d", rootDirDefault, &rootDir, "") // --dir <value> | default value: "/var/users"
	store[""][""] = &root
	rootCommand.AddStringArgs(-1, &root, "")

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "") // sub-command
	// infoVerboseDefault := false
	store["info"]["verbose"] = &infoVerbose
	storeValues["info"]["verbose"] = infoCommand.AddFlag("verbose", "v", &infoVerbose, "") // --verbose, -v | default value: "false"                // --verbose, -v | default value: "false"
	infoVersionDefault := ""
	store["info"]["version"] = &infoVersion
	storeValues["info"]["version"] = infoCommand.AddString("version", "V", infoVersionDefault, &infoVersion, ""). // --version, -V | default value: ""
															SetValidValues([]string{"", "1.0.1", "2.0.0"})
	infoOutputDefault := "./"
	store["info"]["output"] = &infoOutput
	storeValues["info"]["output"] = infoCommand.AddString("output", "o", infoOutputDefault, &infoOutput, "") // --output, -o <value> | default value: "./"
	// infoNoCleanDefault := true
	store["info"]["clean"] = &infoNoClean
	storeValues["info"]["clean"] = infoCommand.AddFlag("no-clean", "N", &infoNoClean, "") // --no-clean | default value: "true"

	// register the `ghost` sub-command
	registry.Register("ghost", "")
	store["ghost"][""] = &ghost
	infoCommand.AddStringArgs(-1, &ghost, "")

	tests := []struct {
		values   []string
		minArgs  int
		maxArgs  int
		want     string
		wantVars map[string]interface{}
		wantArgs []string
		wantErr  string
	}{
		{
			values:  []string{"-V", "1.0.1"},
			maxArgs: 1,
			want:    "",
			wantVars: map[string]interface{}{
				"force":   false,
				"verbose": false,
				"version": "1.0.1",
				"dir":     "/var/users",
			},
		},
		{
			values:  []string{"rootarg", "-V", "1.0.1"},
			minArgs: 1,
			maxArgs: 1,
			want:    "",
			wantVars: map[string]interface{}{
				"force":   false,
				"verbose": false,
				"version": "1.0.1",
				"dir":     "/var/users",
			},
			wantArgs: []string{"rootarg"},
		},
		{
			values:  []string{"rootarg", "-V", "1.0.1", "overflow"},
			maxArgs: 1,
			want:    "",
			wantErr: `"" unnamed args length at argument=overflow > 1`,
		},
		{
			values:  []string{"rootarg", "-V", "1.0.1"},
			maxArgs: 2,
			minArgs: 2,
			want:    "",
			wantErr: `"" unnamed args length < 2`,
		},
		{
			values:  []string{"info", "no_args"},
			want:    "info",
			wantErr: `"info" unnamed args length at argument=no_args > 0`,
		},
		{
			// no args are alloewd for None
			values:  []string{"ghost", "test"},
			want:    "ghost",
			wantErr: `unsupported flag "test" found in the arguments`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %v", i, tt.values), func(t *testing.T) {
			registry.Reset()
			registry.Commands[tt.want].Args.SetMinLen(tt.minArgs)
			registry.Commands[tt.want].Args.SetMaxLen(tt.maxArgs)
			got, _, err := registry.Parse(tt.values, true)
			if err == nil && tt.wantErr != "" {
				t.Errorf("Registry.Parse() wantErr %q", tt.wantErr)
				return
			} else if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Registry.Parse() error = %q, wantErr %q", err.Error(), tt.wantErr)
					return
				}
			} else if got == tt.want {
				gotV := registry.Commands[got]
				for k, o := range gotV.Opts {
					wantVar := tt.wantVars[k]
					assert.Equal(t, wantVar, o.Value.Get(), k)
					assert.Equal(t, o.Value.Get(), storeValues[got][k].Value.Get(), k)
					switch v := store[got][k].(type) {
					case *string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *bool:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]string:
						assert.Equal(t, o.Value.Get(), *v, k)
					default:
						t.Fatalf("%T is unhandled", v)
					}
				}
				a := store[got][""]
				if len(tt.wantArgs) > 0 && a != nil && len(*(a.(*[]string))) > 0 {
					assert.Equal(t, &tt.wantArgs, a, "args")
				}
			} else {
				t.Errorf("Registry.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*----------------*/

func TestRegistry_Parse_CustomTypes_Error(t *testing.T) {
	var (
		infoNum []int
	)
	store := make(map[string]map[string]interface{})
	storeValues := make(map[string]map[string]*Opt)

	store["info"] = make(map[string]interface{})
	storeValues["info"] = make(map[string]*Opt)

	// create a new registry
	registry := NewRegistry("")

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "") // sub-command
	store["info"]["num"] = &infoNum
	storeValues["info"]["num"] = infoCommand.AddIntArray("num", "n", []int{}, &infoNum, "") // --num, -n | default value: []int

	tests := []struct {
		values   []string
		want     string
		wantVars map[string]interface{}
		wantArgs []string
		wantErr  string
	}{
		{
			values: []string{"info", "-n", "1,0,2"},
			want:   "info",
			wantVars: map[string]interface{}{
				"num": []int{1, 0, 2},
			},
		},
		{
			values:  []string{"info", "-n", "1,0,a"},
			want:    "info",
			wantErr: `parsing "a": invalid syntax`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("[%d] %v", i, tt.values), func(t *testing.T) {
			registry.Reset()
			got, _, err := registry.Parse(tt.values, true)
			if err == nil && tt.wantErr != "" {
				t.Errorf("Registry.Parse() wantErr %q", tt.wantErr)
				return
			} else if err != nil {
				if tt.wantErr == "" || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Registry.Parse() error = %q, wantErr %q", err.Error(), tt.wantErr)
					return
				}
			} else if got == tt.want {
				gotV := registry.Commands[got]
				for k, o := range gotV.Opts {
					wantVar := tt.wantVars[k]
					assert.Equal(t, wantVar, o.Value.Get(), k)
					assert.Equal(t, o.Value.Get(), storeValues[got][k].Value.Get(), k)
					switch v := store[got][k].(type) {
					case *string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]string:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *bool:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *int:
						assert.Equal(t, o.Value.Get(), *v, k)
					case *[]int:
						assert.Equal(t, o.Value.Get(), *v, k)
					default:
						t.Fatalf("%T is unhandled", v)
					}
				}
				a := store[got][""]
				if len(tt.wantArgs) > 0 && a != nil && len(*(a.(*[]string))) > 0 {
					assert.Equal(t, &tt.wantArgs, a, "args")
				}
			} else {
				t.Errorf("Registry.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*----------------*/

// test unsupported flag
func TestUnsupportedAssignment(t *testing.T) {

	// options
	options := map[string][]string{
		"---version": {"---version"},
		"---v":       {"---v=1.0.0"},
		"-version":   {"-version"},
	}

	for flag, options := range options {
		t.Run(flag, func(t *testing.T) {
			// command
			r := NewRegistry("")
			rootCmd, _ := r.Register("", "")

			var version string
			rootCmd.AddString("version", "v", "", &version, "")

			_, _, err := r.Parse(options, true)
			require.Error(t, err)
		})
	}
}

// test empty root command
func TestEmptyRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		lines := []string{
			`sub-command => ""`,
			`  Dump variables`,
			`    force="false"`,
			`    verbose="false"`,
			`    version=""`,
			`    dir="/var/users"`,
		}

		out := string(output)
		for _, line := range lines {
			if !strings.Contains(out, line) {
				t.Fatalf("got\n%q\nwant line\n%q", output, line)
			}
		}
	}
}

// test root command when not registered
func TestUnregisteredRootCommand(t *testing.T) {
	// command
	cmd := exec.Command("go", "run", "demo/cmd.go")
	cmd.Env = append(os.Environ(), "NO_ROOT=TRUE")

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		lines := []string{
			`error => clipper.ErrorUnknownCommand{Name:""}`,
		}

		out := string(output)
		for _, line := range lines {
			if !strings.Contains(out, line) {
				t.Fatalf("got\n%q\nwant line\n%q", output, line)
			}
		}
	}
}

// test an unregistered flag
func TestUnregisteredFlag(t *testing.T) {

	// flags
	flags := map[string][]string{
		"--forc":      {"-V", "1.0.1", "-v", "--forc", "-d", "./sub/dir"},
		"--m":         {"-V", "1.0.1", "-v", "--force", "--m", "./sub/dir"},
		"--directory": {"-V", "1.0.1", "-v", "--force", "--directory", "./sub/dir"},
	}

	for flag, options := range flags {
		t.Run(flag, func(t *testing.T) {
			// command
			cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

			// get output
			if output, err := cmd.Output(); err != nil {
				t.Fatalf("Error: %v, out: %q", err, string(output))
			} else {
				out := string(output)
				errStr := fmt.Sprintf(`error => clipper.ErrorUnknownFlag{Name:"%s"}`, flag)
				if !strings.Contains(out, errStr) {
					t.Fatalf("got\n%s\nwant\n%s", out, errStr)
				}
			}
		})
	}
}

func TestUnsupportedFlag(t *testing.T) {

	// flags
	flags := map[string][]string{
		"-force":  {"-V", "1.0.1", "-v", "-force", "-d", "./sub/dir"},
		"student": {"info", "student", "-V", "-v", "--output", "./opt/dir", "--no-clean"},
	}

	for flag, options := range flags {
		t.Run(flag, func(t *testing.T) {
			// command
			cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

			// get output
			if output, err := cmd.Output(); err != nil {
				t.Fatalf("Error: %v, out: %q", err, string(output))
			} else {
				out := string(output)
				errStr := fmt.Sprintf(`error => clipper.ErrorUnsupportedFlag{Name:"%s"}`, flag)
				if !strings.Contains(out, errStr) {
					t.Fatalf("got\n%s\nwant\n%s", out, errStr)
				}
			}
		})
	}
}

// test for valid inverted flag values
func TestValidInvertFlagValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"info", "-V", "1.0.1", "-v", "--output", "./opt/dir", "--no-clean"},
		{"info", "--version=1.0.1", "--no-clean", "--output", "./opt/dir", "--verbose"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "info"`,
				`  Dump variables`,
				`    verbose="true"`,
				`    version="1.0.1"`,
				`    output="./opt/dir"`,
				`    clean="false"`,
			}

			out := string(output)
			for _, line := range lines {
				if !strings.Contains(out, line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test for invalid flag error when an inverted flag is used without `--no-` prefix
func TestErrorUnknownFlagForInvertFlags(t *testing.T) {

	// options list
	optionsList := map[string][]string{
		"--no-dump": {"info", "--version", "--no-dump", "--output", "./opt/dir", "--verbose"},
	}

	for flag, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			out := string(output)
			errStr := fmt.Sprintf(`error => clipper.ErrorUnknownFlag{Name:"%s"}`, flag)
			if !strings.Contains(out, errStr) {
				t.Fatalf("got\n%s\nwant\n%s", out, errStr)
			}
		}
	}
}

// test `--flag=value` syntax
func TestFlagAssignmentSyntax(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"list", "student", "--dir", "/opt", "thatisuday"},
		{"list", "student", "thatisuday", "--dir", "/opt"},
	}

	for i, options := range optionsList {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// command
			cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

			// get output
			if output, err := cmd.Output(); err != nil {
				t.Fatalf("Error: %v, out: %q", err, string(output))
			} else {
				lines := []string{
					`sub-command => "list"`,
					`  Dump variables`,
					`    dir="[/opt]"`,
					`    args=[student,thatisuday]`,
				}

				out := string(output)
				for _, line := range lines {
					if !strings.Contains(out, line) {
						t.Fatalf("got\n%q\nwant line\n%q", output, line)
					}
				}
			}
		})
	}
}

// test for valid variadic argument values
func TestValidVariadicArgumentValues(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"list", "student", "--dir...", "./opt/dir", "/data"},
		{"list", "student", "--dir...", "./opt/dir", "--dir...", "/data"},
	}

	for i, options := range optionsList {
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			// command
			cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

			// get output
			if output, err := cmd.Output(); err != nil {
				t.Fatalf("Error: %v, out: %q", err, string(output))
			} else {
				lines := []string{
					`sub-command => "list"`,
					`  Dump variables`,
					`    dir="[./opt/dir,/data]"`,
					`    args=[student]`,
				}

				out := string(output)
				for _, line := range lines {
					if !strings.Contains(out, line) {
						t.Fatalf("got\n%q\nwant line\n%q", output, line)
					}
				}
			}
		})
	}
}

/*-------------------*/

// test root command with options
func TestRootCommandWithOptions(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"userinfo", "-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir"},
		{"-V", "1.0.1", "--verbose", "--force", "userinfo", "--dir", "./sub/dir"},
		{"-V", "1.0.1", "-v", "--force", "--dir", "./sub/dir", "userinfo"},
		{"--version", "1.0.1", "--verbose", "--force", "--dir", "./sub/dir", "userinfo"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => ""`,
				`  Dump variables`,
				`    force="true"`,
				`    verbose="true"`,
				`    version="1.0.1"`,
				`    dir="./sub/dir"`,
				`    args=[userinfo]`,
			}

			out := string(output)
			for _, line := range lines {
				if !strings.Contains(out, line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test sub-command with options
func TestSubCommandWithOptions(t *testing.T) {

	// options list
	optionsList := [][]string{
		{"list", "student", "manager", "--dir...", "./opt/dir1", "./opt/dir2"},
		{"list", "student", "--dir", "./opt/dir1", "-d", "./opt/dir2", "manager"},
	}

	for _, options := range optionsList {
		// command
		cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

		// get output
		if output, err := cmd.Output(); err != nil {
			t.Fatalf("Error: %v, out: %q", err, string(output))
		} else {
			lines := []string{
				`sub-command => "list"`,
				`  Dump variables`,
				`    dir="[./opt/dir1,./opt/dir2]"`,
				`    args=[student,manager]`,
			}

			out := string(output)
			for _, line := range lines {
				if !strings.Contains(out, line) {
					t.Fatalf("got\n%q\nwant line\n%q", output, line)
				}
			}
		}
	}
}

// test validate arg
func TestInvalidArg(t *testing.T) {
	// options
	options := []string{"info", "worker", "-V", "-v", "2.0.0"}

	// command
	cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		out := string(output)
		want := "error => clipper.ErrorUnsupportedFlag{Name:\"worker\"}\n"
		if out != want {
			t.Fatalf("got\n%q\nwant\n%q", out, want)
		}
	}
}

// test validate flag
func TestDiabledArgs(t *testing.T) {
	// options
	options := []string{"info", "student", "-V", "2.0.1", "-v"}

	// command
	cmd := exec.Command("go", append([]string{"run", "demo/cmd.go"}, options...)...)

	// get output
	if output, err := cmd.Output(); err != nil {
		t.Fatalf("Error: %v, out: %q", err, string(output))
	} else {
		out := string(output)
		want := "error => clipper.ErrorUnsupportedFlag{Name:\"student\"}\n"
		if out != want {
			t.Fatalf("got\n%q\nwant\n%q", out, want)
		}
	}
}
