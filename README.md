# goenable

Write Bash builtins in Go.

`goenable` makes it easy to extend Bash with Go, without executing in a separate process. Using Bash's `enable` builtin, `goenable` is loaded into the Bash runtime and provides helpers to load [custom plugins](#write-a-plugin).

## Quick start

```bash
make  # Create the goenable.so binary and build the example plugins

enable -f ./out/goenable.so goenable  # Load goenable
help goenable                         # Print usage

goenable load ./out/hello output      # Load the hello plugin
eval "$output"                        # Prepare functions from the hello plugin

hello
# Hello, world!
```

## Write a plugin

`goenable` uses [Go plugins](https://golang.org/pkg/plugin/) to load and run custom Go code inside the Bash process.

Start with this `pow` plugin as your `main.go`:

```go
package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/johnstarich/goenable/usage"
)

// Usage returns the full set of documentation for this plugin
func Usage() string {
	return "pow X Y\nPrints X to the power of Y"
}

// Load runs any set up required by this plugin
func Load() error {
	return nil
}

// Unload runs any tear down required by this plugin
func Unload() {
}

// Run executes this plugin with the given arguments
func Run(args []string) error {
	if len(args) != 2 {
		// Print usage and return exit code 2
		return usage.GenericError()
	}
	// If we encounter errors and return them, then the exit code is 1
	x, _ := strconv.ParseFloat(args[0], 64)
	y, _ := strconv.ParseFloat(args[1], 64)
	// Print pow result and return success with exit code 0
	fmt.Println(math.Pow(x, y))
	return nil
}

func main() {}
```

Build it with `go build -o pow -buildmode=plugin main.go`

Make sure you have a Bash terminal open and the `goenable.so` binary, then run:

```bash
# Start goenable with the 'enable' Bash builtin
enable -f ./path/to/your/goenable.so goenable

goenable load ./pow output  # Load the pow plugin
eval "$output"              # Prepare the pow plugin function

pow 1 10
# 1
pow 2 2.2
# 4.59479341998814
```

Check out other plugins like [hello](examples/hello/main.go) for more working [examples](examples).

## Contributing

All contributions are welcome!

If you have suggestions for new features or run into a problem, please [submit an issue](https://github.com/JohnStarich/goenable/issues/new).
