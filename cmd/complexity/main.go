package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	cyclomatic "github.com/mkeeler/go-complexity/internal/cyclomatic-complexity"
	goroutines "github.com/mkeeler/go-complexity/internal/go-routines"
	"github.com/mkeeler/go-complexity/internal/inspector"
	"github.com/mkeeler/go-complexity/internal/loc"
)

func contains[T comparable](s []T, check T) bool {
	for _, v := range s {
		if v == check {
			return true
		}
	}
	return false
}

type allMetrics struct {
	Cyclomatic *cyclomatic.CyclomaticComplexity `json:",omitempty"`
	LOC        *loc.LinesOfCode                 `json:",omitempty"`
	GoRoutines *goroutines.GoRoutines           `json:",omitempty"`
}

func main() {
	tests := false
	jsonOutput := false
	var analysisTypes StringSliceValue

	flag.BoolVar(&tests, "tests", false, "Enable analyzing test files in addition to the main code")
	flag.BoolVar(&jsonOutput, "json", false, "Output in JSON format instead of prettified output")
	flag.Var(&analysisTypes, "types", "Types of complexity analysis to perform. If unspecified all will be executed. If specified multiple times then each one specified will be execute. Possible values are: cyclomatic, loc, go-routine")
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: complexity <directory>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(analysisTypes) < 1 {
		analysisTypes = StringSliceValue{"cyclomatic", "loc"}
	}

	ins, err := inspector.NewInspector(args[0], tests)
	if err != nil {
		fmt.Printf("Error building package inspector: %v\n", err)
		return
	}

	var m allMetrics
	if contains(analysisTypes, "cyclomatic") {
		c, err := cyclomatic.CalculateComplexity(ins)
		if err != nil {
			fmt.Printf("Error computing cyclomatic complexity: %v\n", err)
			return
		}
		m.Cyclomatic = c
	}

	if contains(analysisTypes, "loc") {
		m.LOC = loc.CalculateLinesOfCode(ins)
	}

	if contains(analysisTypes, "go-routine") {
		m.GoRoutines = goroutines.CountGoRoutineInvocations(ins)
	}

	json.NewEncoder(os.Stdout).Encode(m)
}
