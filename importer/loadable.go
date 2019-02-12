package importer

import (
	"fmt"
	"os"
	"strings"
)

// Name is the string users invoke to execute this loadable
func Name() string {
	return "import"
}

// UsageShort returns a short summary of usage information, usually indicating the arguments that should be provided to the loadable
func UsageShort() string {
	return "import SCRIPT"
}

// Usage returns the full set of documentation for this loadable
func Usage() string {
	return strings.TrimSpace(`
'import' is a utility to source scripts with special shell extensions.
`)
}

// Run executes this loadable with the given arguments
func Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: "+UsageShort())
		return 2
	}
	fmt.Println("echo", args[0])
	return 0
}

// Load runs any set up required by this loadable
func Load(name string) int {
	return 1
}

// Unload runs any tear down required by this loadable
func Unload() {
}
