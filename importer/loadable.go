package importer

import "fmt"

// Name is the string users invoke to execute this loadable
func Name() string {
	return "import"
}

// UsageShort returns a short summary of usage information, usually indicating the arguments that should be provided to the loadable
func UsageShort() string {
	return "usage synopsis goes here"
}

// Usage returns the full set of documentation for this loadable
func Usage() string {
	return `
long documentation strings go here
`
}

// Run executes this loadable with the given arguments
func Run(args []string) int {
	fmt.Println("Args", args)
	return 0
}

// Load runs any set up required by this loadable
func Load(name string) int {
	fmt.Println("Loaded!", name)
	return 1
}

// Unload runs any tear down required by this loadable
func Unload() {
	fmt.Println("Unloading!")
}
