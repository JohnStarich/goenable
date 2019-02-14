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
	x, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return err
	}
	// Print pow result and return success with exit code 0
	fmt.Println(math.Pow(x, y))
	return nil
}

func main() {}
