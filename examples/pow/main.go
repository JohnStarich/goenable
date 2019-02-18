package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
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
func Run(args []string) (int, error) {
	if len(args) != 2 {
		// Print usage and return exit code 2
		return 2, errors.New(Usage())
	}
	// If we encounter errors and return them, then the exit code is 1
	x, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return 1, err
	}
	y, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return 1, err
	}
	// Print pow result and return success with exit code 0
	fmt.Println(math.Pow(x, y))
	return 0, nil
}

func main() {}
