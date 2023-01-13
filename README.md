# clipper
A simple Go package to parse command-line arguments [_getopt(3) style_](http://man7.org/linux/man-pages/man3/getopt.3.html). Designed especially for making CLI based libraries with ease. It has built-in support for sub-commands, long and short flag name combination (_for example `--version` <==> `-v`_), `--flag=<value>` syntax, inverted flag (_for example `--no-clean`_), variadic arguments for long-style flags(_for example `--dir... /data1 /data2`), etc.

Main advantage - state might be reset to initial state (to default values and set unchanged).
So can be simple reused in embedded interactive CLI.

Based on [clapper](https://github.com/thatisuday/clapper), but typed (not string for all) and has more features.

Can be simple extended for additional types (see [Value](https://github.com/msaf1980/go-clipper/blob/master/value.go#L11) interface, base on extended [pflag](https://github.com/spf13/pflag)).

## Documentation
[**pkg.go.dev**](https://pkg.go.dev/github.com/msaf1980/go-clipper?tab=doc)

## Installation
```
$ go get "github.com/msaf1980/go-clipper"
```

## Usage

```go
// cmd.go
package main

import (
	"fmt"
	"os"

	"github.com/msaf1980/go-clipper"
)

var (
	VERSION = "0.0.1"
)

func main() {

	var (
		rootForce, rootVerbose bool
		rootDir   string
		root                   []string

		infoVerbose, infoNoClean bool
		infoVersion, infoOutput  string

		list, listDir []string
		listVerbose []bool
	)

	// create a new registry
	registry := clipper.NewRegistry("programm description message")

	// register the root command
	if _, ok := os.LookupEnv("NO_ROOT"); !ok {
		rootCommand, _ := registry.Register("root help") // root command
		// rootCommand.AddArg("output", "")                    //
		rootCommand.AddFlag("force", "f", &rootForce, "flag help")             // --force, -f | default value: "false"
		rootCommand.AddFlag("verbose", "v", &rootVerbose, "flag help")         // --verbose, -v | default value: "false"
		rootCommand.AddString("dir", "d", "/var/users", &rootDir, "flag help") // --dir <value> | default value: "/var/users"
		rootCommand.AddStringArgs(-1, &root, "args help") // root unnamed args
        rootCommand.AddVersionHelper("version", "V", registry.Description, VERSION)
	}

	// register the `info` sub-command
	infoCommand, _ := registry.Register("info", "info help")              // sub-command
	infoCommand.AddFlag("verbose", "v", &infoVerbose, "flag help")        // --verbose, -v | default value: "false"
	infoCommand.AddString("version", "V", "", &infoVersion, "flag help"). // --version, -V | default value: "false"
									SetValidValues([]string{"", "1.0.1", "2.0.0"}). // valid versions
									SetRequired(true)                               // version are required
	infoCommand.AddString("output", "o", "./", &infoOutput, "flag help") // --output, -o <value> | default value: "./"
	infoCommand.AddFlag("no-clean", "N", &infoNoClean, "flag help")      // --no-clean | default value: "true"

	listCommand, _ := registry.Register("list", "list help")                     // sub-command
	listCommand.AddStringArray("dir", "d", []string{"a"}, &listDir, "flag help") // --output, -o <value> | default value: "./"
	listCommand.AddStringArgs(-1, &list, "args help")
	listCommand.AddMultiFlag("verbose", "v", &listVerbose, "multi-flag verbose") // --verbose, -v | default value: []

	// register the `ghost` sub-command
	ghostCommand, _ := registry.Register("ghost", "ghost help")
	ghostCommand.AddVersionHelper("version", "V", registry.Description, VERSION)

	/*----------------*/

	// parse command-line arguments
	command, err := registry.Parse(os.Args[1:], true)

   	// For interactive use (don't exit after help print, check helpRequested and break command execution if set)
   	// command, helpRequested, err := registry.ParseInteract(os.Args[1:], false)
   	// if !helpRequested {
   	// // execute command
   	//     ..
   	// }

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
	}}
}
```

In the above example, we have registred a **root** command and an `info` command. The `registry` can parse arguments passed to the command that executed this program.

#### Example 1
When the **root command** is executed with no command-line arguments.

```
$ go run demo/cmd.go

sub-command => ""
  Dump variables
    force="false"
    verbose="false"
    version=""
    dir="/var/users"
    args=[]
```

#### Example 2
When the **root command** is executed but not registered.

```
$ NO_ROOT=TRUE go run demo/cmd.go

error => clipper.ErrorUnknownCommand{Name:""}
```

#### Example 3
When the **root command** is executed with short/long flag names as well as by changing the positions of the arguments.

```
$ go run demo/cmd.go userinfo -V 1.0.1 -v --force --dir ./sub/dir
$ go run demo/cmd.go -V 1.0.1 --verbose --force userinfo --dir ./sub/dir
$ go run demo/cmd.go -V 1.0.1 -v --force --dir ./sub/dir userinfo
$ go run demo/cmd.go --version 1.0.1 --verbose --force --dir ./sub/dir userinfo

sub-command => ""
  Dump variables
    force="true"
    verbose="true"
    version="1.0.1"
    dir="./sub/di"
    args=[userinfo]
```

#### Example 4
When an **unregistered flag** is provided in the command-line arguments.

```

$ go run demo/cmd.go userinfo -V 1.0.1 -v --force --d ./sub/dir
error => clipper.ErrorUnknownFlag{Name:"--d"}

$ go run demo/cmd.go userinfo -V 1.0.1 -v --force --directory ./sub/dir
error => clipper.ErrorUnknownFlag{Name:"--directory"}

```


#### Example 5
When `information` was intended to be a sub-command but not registered and the root command accepts arguments.

```
$ go run demo/cmd.go information --force

sub-command => ""
  Dump variables
    force="true"
    verbose="false"
    version=""
    dir="/var/users"
    args=[information]
```

#### Example 6
When an **unnamed args (not allowed)** is provided in the command-line arguments.

```
$ go run demo/cmd.go info student -V -v --output ./opt/dir

error => clipper.ErrorUnsupportedFlag{Name:"student"}
```

#### Example 7
When a command is executed with an **inverted** flag (flag that starts with `--no-` prefix).

```
$ go run demo/cmd.go info -V -v --output ./opt/dir --no-clean

sub-command => "info"
  Dump variables
    verbose="true"
    version=""
    output="./opt/dir"
    clean="false
```

#### Example 8
When the position of argument values are changed and variadic arguments are provided.

```
$ go run demo/cmd.go list student --dir... /data1 /data2

sub-command => "list"
  Dump variables
    dir="[/data1,/data2]"
    args=[student]
```

#### Example 9
When a **sub-command** is registered without any flags.

```
$ go run demo/cmd.go ghost -v thatisuday -V 2.0.0 teachers

error => clipper.ErrorUnknownFlag{Name:"-v"}
```

#### Example 10
When a **sub-command** is registered without any arguments.

```
$ go run demo/cmd.go ghost
$ go run demo/cmd.go ghost thatisuday extra

sub-command => "ghost
```

#### Example 11
When the **root command** is not registered or the **root command** is registered with no arguments.

```
$ NO_ROOT=TRUE go run demo/cmd.go information
error => clipper.ErrorUnknownCommand{Name:"information"}

$ go run cmd.go ghost
sub-command => "ghost"
```

#### Example 12
When unsupported flag format is provided.

```
$ go run demo/cmd.go ---version 
error => clipper.ErrorUnsupportedFlag{Name:"---version"}

$ go run demo/cmd.go ---v=1.0.0 
error => clipper.ErrorUnsupportedFlag{Name:"---v"}

$ go run demo/cmd.go -version 
error => clipper.ErrorUnsupportedFlag{Name:"-version"}

$ go run demo/cmd.go list student -d... /data1 /data2
error => clipper.ErrorUnsupportedFlag{Name:"-d..."}
```

## Contribution
A lot of improvements can be made to this library, one of which is the support for combined short flags, like `-abc`. If you are willing to contribute, create a pull request and mention your bug fixes or enhancements in the comment.
