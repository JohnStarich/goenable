package main

import (
	"fmt"
	"strings"

	"github.com/johnstarich/goenable/usage"
)

// Usage returns the full set of documentation for this plugin
func Usage() string {
	return strings.TrimSpace(`
'hello' is a hello world program. No arguments required.
It says hi!
`)
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
	if len(args) != 0 {
		return usage.GenericError()
	}
	fmt.Println("Hello, world!")
	return nil
}

func main() {}
