package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mkeeler/go-complexity/internal/cyclomatic"
	"github.com/mkeeler/go-complexity/internal/inspector"
)

func outputPretty(c *cyclomatic.CyclomaticComplexity) {
	for pkg, pkgC := range c.Packages {
		fmt.Printf("Package: %s\n", pkg)
		for fnName, fnC := range pkgC.Functions {
			fmt.Printf("\tFunction: %s, Complexity: %d\n", fnName, fnC.Score)
		}
	}
}

type jsonFormatEntry struct {
	Package  string
	Function string
	Score    int
}

func outputJSON(c *cyclomatic.CyclomaticComplexity) {
	var entries []*jsonFormatEntry
	for pkg, pkgC := range c.Packages {
		for fnName, fnC := range pkgC.Functions {
			entries = append(entries, &jsonFormatEntry{Package: pkg, Function: fnName, Score: fnC.Score})
		}
	}
	json.NewEncoder(os.Stdout).Encode(entries)
}

func main() {
	tests := false
	jsonOutput := false
	flag.BoolVar(&tests, "tests", false, "Enable analyzing test files in addition to the main code")
	flag.BoolVar(&jsonOutput, "json", false, "Output in JSON format instead of prettified output")
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: complexity <directory>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ins, _ := inspector.NewInspector(args[0], tests)

	c, err := cyclomatic.CalculateComplexity(ins)

	if err != nil {
		fmt.Println(err)
		return
	}

	if jsonOutput {
		outputJSON(c)
	} else {
		outputPretty(c)
	}
}
