package cyclomatic

import (
	"bytes"
	"encoding/json"
)



type CyclomaticComplexity struct {
	Packages map[string]*PackageCyclomaticComplexity
}

func newCyclomaticComplexity() *CyclomaticComplexity {
	return &CyclomaticComplexity{
		Packages: make(map[string]*PackageCyclomaticComplexity),
	}
}

func (c *CyclomaticComplexity) recordComplexity(pkgName, fnName string, score int) {
	pkg, ok := c.Packages[pkgName]
	if !ok {
		pkg = newPackageCyclomaticComplexity()
		c.Packages[pkgName] = pkg
	}

	pkg.recordComplexity(fnName, score)
}

func (c *CyclomaticComplexity) MarshalJSON() ([]byte, error) {
	type jsonFormatEntry struct {
		Package  string
		Function string
		Score    int
	}

	var entries []*jsonFormatEntry
	for pkg, pkgC := range c.Packages {
		for fnName, fnC := range pkgC.Functions {
			entries = append(entries, &jsonFormatEntry{Package: pkg, Function: fnName, Score: fnC.Score})
		}
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(entries)
	return buf.Bytes(), err
}

type PackageCyclomaticComplexity struct {
	Functions map[string]*FunctionCyclomaticComplexity
}

func newPackageCyclomaticComplexity() *PackageCyclomaticComplexity {
	return &PackageCyclomaticComplexity{
		Functions: make(map[string]*FunctionCyclomaticComplexity),
	}
}

func (p *PackageCyclomaticComplexity) recordComplexity(fnName string, score int) {
	fn, ok := p.Functions[fnName]
	if !ok {
		fn = newFunctionCyclomaticComplexity()
		p.Functions[fnName] = fn
	}

	fn.recordComplexity(score)
}

type FunctionCyclomaticComplexity struct {
	Score int
}

func newFunctionCyclomaticComplexity() *FunctionCyclomaticComplexity {
	return &FunctionCyclomaticComplexity{
		// The score always starts at 1
		Score: 1,
	}
}

func (f *FunctionCyclomaticComplexity) recordComplexity(score int) {
	f.Score = score
}
