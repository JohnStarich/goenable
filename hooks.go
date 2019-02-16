package main

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/johnstarich/goenable/stringutil"
	"github.com/johnstarich/goenable/usage"
	"github.com/pkg/errors"
)

const (
	functionPrefixSeparator = "-"
	exitSuccess             = 0
	exitGeneralError        = 1
	exitUsageError          = 2
)

var (
	importCache = make(map[string]*plugin.Plugin)
	importNames = make(map[string]string)
)

// UnloadFunc is a handler for loading Go plugins
type UnloadFunc = func() error

// LoadFunc is a handler for loading Go plugins
type LoadFunc = func() error

// RunFunc is a handler for running Go plugins
type RunFunc = func(args []string) error

// UsageFunc is a handler for returning usage for Go plugins
type UsageFunc = func() string

// Name is the string users invoke to execute this loadable
func Name() string {
	return "goenable"
}

// UsageShort returns a short summary of usage information, usually indicating the arguments that should be provided to the loadable
func UsageShort() string {
	return "a utility to run Go plugins as shell functions"
}

// Usage returns the full set of documentation for this loadable
func Usage() string {
	return strings.TrimSpace(`
goenable load GO_PLUGIN ENV_VAR
	Load the given plugin. Stores 'eval'-able output into ENV_VAR.
	Be sure to run 'eval "${ENV_VAR}"' to begin using the new plugin.
	Local variables are strongly recommended.

goenable run GO_PLUGIN [args ...]
	Run the given plugin, optionally with arguments. This is typically reserved for internal use.
	i.e. Using eval on the output from 'goenable load' will add your plugin as a command, which calls 'goenable run' underneath.
`)
}

// Run executes this loadable with the given arguments
func Run(args []string) int {
	if p, err := run(args); err != nil {
		switch err.(type) {
		case usage.Error:
			if message := err.Error(); message != "" {
				logError(message)
				return exitUsageError
			}
			if p == nil {
				logError("Error loading plugin")
				return exitGeneralError
			}
			usageSym, err := p.Lookup("Usage")
			if err != nil {
				logError("Error loading plugin usage: no Usage() available")
				return exitGeneralError
			}
			usage, ok := usageSym.(UsageFunc)
			if !ok {
				logError("Usage() for plugin is not a goenable.UsageFunc: %T", usageSym)
				return exitGeneralError
			}
			logError("Usage:\n%s", usage())
			return exitUsageError
		default:
			logError("Error:", err)
			return exitGeneralError
		}
	}
	return exitSuccess
}

// Load runs any set up required by this loadable
func Load(name string) bool {
	return true
}

// Unload runs any tear down required by this loadable
func Unload() {
	for name, p := range importCache {
		unloadSym, err := p.Lookup("Unload")
		if err != nil {
			logError("Unload failed: Plugin '%s' does not support Unload()", name)
			continue
		}
		unload, ok := unloadSym.(UnloadFunc)
		if !ok {
			logError("Unload() for plugin '%s' is not a goenable.UnloadFunc: %T", name, unloadSym)
			continue
		}
		err = unload()
		if err != nil {
			logError("Unload failed for '%s': %s", name, err.Error())
		}
	}
}

func logError(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message+"\n", args...)
}

func run(args []string) (*plugin.Plugin, error) {
	if len(args) == 0 {
		return nil, usage.Errorf("Usage: goenable load|run [args ...]")
	}
	if len(args) < 2 {
		return nil, usage.Errorf("Usage:\n%s", Usage())
	}
	var p *plugin.Plugin
	command, fileName, args := args[0], args[1], args[2:]
	name := filepath.Base(fileName)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	switch command {
	case "run":
		pluginPath, loaded := importNames[name]
		if !loaded {
			return nil, usage.Errorf("Plugin not loaded yet: " + name)
		}
		p = importCache[pluginPath]
		runSym, err := p.Lookup("Run")
		if err != nil {
			return p, err
		}
		runFunc, ok := runSym.(RunFunc)
		if !ok {
			return p, fmt.Errorf("Run() for plugin is not a goenable.RunFunc: Run = %T", runSym)
		}
		return p, runFunc(args)
	case "load":
		if len(args) != 1 {
			return nil, usage.Errorf("Usage:\n%s", Usage())
		}
		outputEnvVar := args[0]
		absPath, err := filepath.Abs(fileName)
		if err != nil {
			return nil, err
		}
		var loaded bool
		if p, loaded = importCache[absPath]; loaded {
			// plugin was loaded already, skipping...
			return p, nil
		}
		if p, err = plugin.Open(fileName); err != nil {
			return nil, err
		}
		if err := tryLoad(p); err != nil {
			return p, errors.Wrap(err, "Failed to run load for plugin "+name)
		}
		os.Setenv(outputEnvVar, makeModule(name))
		importCache[absPath] = p
		importNames[name] = absPath
	default:
		return nil, usage.Errorf("Invalid command: " + command)
	}
	return p, nil
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
