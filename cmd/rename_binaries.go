package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	versionRegex = regexp.MustCompile(`\d+(\.\d+)+`)
	kernelRegex  = regexp.MustCompile(`^([a-z]+)-([^-]+)`)
	// make suffixes consistent between download URLs, even if it is technically incorrect for that OS
	suffixRegex = regexp.MustCompile(`\.dylib`)
)

func main() {
	cliName, args := os.Args[0], os.Args[1:]
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s OUT_DIR\n", filepath.Base(cliName))
		os.Exit(2)
	}

	dir := args[0]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workingDirectory, err = filepath.Abs(workingDirectory)
	if err != nil {
		panic(err)
	}
	binaryName := filepath.Base(workingDirectory)
	for _, f := range files {
		fileName := f.Name()
		newFileName := kernelRegex.ReplaceAllStringFunc(fileName, func(s string) string {
			if !strings.HasPrefix(s, binaryName+"-") {
				// only rename the output binaries
				return s
			}
			s = strings.TrimPrefix(s, binaryName+"-")
			s = strings.Title(s)
			return binaryName + "-" + s
		})
		newFileName = versionRegex.ReplaceAllLiteralString(newFileName, "")
		newFileName = strings.Replace(newFileName, "--", "-", -1)
		newFileName = strings.Replace(newFileName, "amd64", "x86_64", 1)
		newFileName = suffixRegex.ReplaceAllLiteralString(newFileName, ".so")

		if fileName == newFileName {
			continue
		}
		fmt.Printf("Moving %s to %s...\n", fileName, newFileName)
		err := os.Rename(filepath.Join(dir, fileName), filepath.Join(dir, newFileName))
		if err != nil {
			panic(err)
		}
	}
}
