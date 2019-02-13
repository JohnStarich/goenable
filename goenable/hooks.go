package goenable

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/johnstarich/bash-go-loader/stringutil"
	"github.com/johnstarich/bash-go-loader/usage"
	"github.com/pkg/errors"
)

const (
	functionPrefixSeparator = "-"
	outputEnvVar            = "__GOENABLE_OUTPUT"
)

var (
	importCache = make(map[string]*plugin.Plugin)
	importNames = make(map[string]string)
)

// LoadFunc is a handler for loading Go plugins
type LoadFunc = func() error

// RunFunc is a handler for running Go plugins
type RunFunc = func(args []string) error

// Name is the string users invoke to execute this loadable
func Name() string {
	return "goenable"
}

// UsageShort returns a short summary of usage information, usually indicating the arguments that should be provided to the loadable
func UsageShort() string {
	return "goenable load|run GO_PLUGIN [args ...]"
}

// Usage returns the full set of documentation for this loadable
func Usage() string {
	return strings.TrimSpace(`
'goenable' is a utility to load Go plugins as scripts

goenable load GO_PLUGIN
	Load the given plugin. Be sure to 'eval' the output of this call to begin using the new plugin.

goenable run GO_PLUGIN [args ...]
	Run the given plugin, optionally with arguments. This is typically reserved for internal use.
`)
}

// Run executes this loadable with the given arguments
func Run(args []string) int {
	if err := run(args); err != nil {
		switch err.(type) {
		case usage.Error:
			if message := err.Error(); message != "" {
				fmt.Fprintln(os.Stderr, message)
			}
			fmt.Fprintln(os.Stderr, "Usage: "+UsageShort())
			return 2
		default:
			fmt.Fprintln(os.Stderr, "Error:", err)
			return 1
		}
	}
	return 0
}

// Load runs any set up required by this loadable
func Load(name string) int {
	return 1
}

// Unload runs any tear down required by this loadable
func Unload() {
}

func run(args []string) error {
	if len(args) < 2 {
		return usage.Errorf(stringutil.Dedent(`
			Provide a command and a plugin.
				Available commands: load, run
		`))
	}
	command, fileName, args := args[0], args[1], args[2:]
	name := filepath.Base(fileName)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	switch command {
	case "run":
		pluginPath, loaded := importNames[name]
		if !loaded {
			return usage.Errorf("Plugin not loaded yet: " + name)
		}
		p := importCache[pluginPath]
		runSym, err := p.Lookup("Run")
		if err != nil {
			return err
		}
		runFunc, ok := runSym.(RunFunc)
		if !ok {
			return fmt.Errorf("Run() for plugin is not a goenable.RunFunc: Run = %T", runSym)
		}
		return runFunc(args)
	case "load":
		absPath, err := filepath.Abs(fileName)
		if err != nil {
			return err
		}
		if _, ok := importCache[absPath]; ok {
			// plugin was loaded already, skipping...
			return nil
		}
		p, err := plugin.Open(fileName)
		if err != nil {
			return err
		}
		if err := tryLoad(p); err != nil {
			return errors.Wrap(err, "Failed to run load for plugin "+name)
		}
		os.Setenv(outputEnvVar, makeModule(name))
		importCache[absPath] = p
		importNames[name] = absPath
	default:
		return usage.Errorf("Invalid command: " + command)
	}
	return nil
}

func tryLoad(p *plugin.Plugin) error {
	loadSym, err := p.Lookup("Load")
	if err != nil {
		return err
	}
	load, ok := loadSym.(LoadFunc)
	if !ok {
		return fmt.Errorf("Load() for plugin is not a goenable.LoadFunc: Load = %T", loadSym)
	}
	return load()
}

func makeModule(name string) string {
	return stringutil.Dedent(`
		` + name + `() {
			goenable run ` + name + ` "$@"
		}
	`)
}
