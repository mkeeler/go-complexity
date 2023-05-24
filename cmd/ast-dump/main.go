package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/mkeeler/go-complexity/internal/inspector"
	"golang.org/x/tools/go/packages"
)

func main() {
	tests := false

	flag.BoolVar(&tests, "tests", false, "Enable analyzing test files in addition to the main code")
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: complexity <directory>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Loading files from: %s", args[0])
	ins, err := inspector.NewInspector(args[0], tests)
	if err != nil {
		fmt.Printf("Error building package inspector: %v\n", err)
		return
	}

	ins.WalkPackages(func(pkg *packages.Package) {
		for i, fname := range pkg.CompiledGoFiles {
			fmt.Printf("File: %s\n\n", fname)
			spew.Dump(pkg.Syntax[i])
		}
	})

}
