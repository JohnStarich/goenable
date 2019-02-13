package namespace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/syntax"
)

const (
	functionPrefixSeparator = "-"
)

var (
	importCache = make(map[string]bool)
)

// Name is the string users invoke to execute this loadable
func Name() string {
	return "namespace"
}

// UsageShort returns a short summary of usage information, usually indicating the arguments that should be provided to the loadable
func UsageShort() string {
	return "namespace SCRIPT"
}

// Usage returns the full set of documentation for this loadable
func Usage() string {
	return strings.TrimSpace(`
'namespace' is a utility to load scripts and make them namespace-friendly.
Namespaces make it easier to create reusable modules and don't conflict in a global bash context.
`)
}

// Run executes this loadable with the given arguments
func Run(args []string) int {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: "+UsageShort())
		return 2
	}
	if err := run(args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return 1
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
	fileName := args[0]
	parser := syntax.NewParser()
	reader, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer reader.Close()
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}
	if importCache[absPath] {
		return nil
	}
	name := filepath.Base(fileName)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	f, err := parser.Parse(reader, name)
	if err != nil {
		return err
	}

	extraScript := mutate(f, name)
	fmt.Fprintf(os.Stdout, extraScript)
	printer := syntax.NewPrinter()
	printer.Print(os.Stdout, f)
	importCache[absPath] = true
	return nil
}

func mutate(f *syntax.File, name string) string {
	extraScript := ""
	functionNames := make(map[string]bool)
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.FuncDecl:
			functionNames[x.Name.Value] = true
			x.Name.Value = name + functionPrefixSeparator + x.Name.Value
		}
		return true
	})
	prefix := name + functionPrefixSeparator
	allFunctionNames := ""
	for name := range functionNames {
		allFunctionNames += " " + singleQuote(name)
	}

	if !functionNames["usage"] {
		functionNames["usage"] = true
		extraScript += dedent(`
			` + prefix + `usage() {
				echo 'Usage: ` + name + ` COMMAND' >&2
				echo 'Available commands: '` + allFunctionNames + ` >&2
			}
		`)
	}
	if !functionNames[name] {
		extraScript += dedent(`
			` + name + `() {
				local subCommand=$1
				shift
				if [[ -z "$subCommand" ]]; then
					` + prefix + `usage
					return 2
				fi
				if ! command -v "` + prefix + `${subCommand}" >/dev/null; then
					echo "Invalid subcommand: ${subCommand}" >&2
					` + prefix + `usage
					return 2
				fi
				"` + prefix + `${subCommand}" "$@"
			}
		`)
	}
	if !functionNames["complete"] {
		extraScript += dedent(`
			` + prefix + `complete() {
				local options=(` + allFunctionNames + `)
				local prev=${COMP_WORDS[COMP_CWORD - 1]}
				if [[ "$prev" != ` + name + ` ]]; then
					return
				fi
				COMPREPLY+=( $(compgen -W "${options[*]}" -- "${COMP_WORDS[COMP_CWORD]}") )
			}
		`)
	}
	extraScript += fmt.Sprintf("complete -F %s %s\n\n", prefix+"complete", name)
	syntax.Walk(f, func(node syntax.Node) bool {
		switch x := node.(type) {
		case *syntax.CallExpr:
			if len(x.Args) > 0 && len(x.Args[0].Parts) == 1 {
				switch funcName := x.Args[0].Parts[0].(type) {
				case *syntax.Lit:
					if functionNames[funcName.Value] {
						funcName.Value = name + functionPrefixSeparator + funcName.Value
					}
				}
			}
		}
		return true
	})
	return extraScript
}
