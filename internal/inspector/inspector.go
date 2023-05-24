package inspector

import (
	"fmt"
	"go/ast"
	"go/token"

	goinspector "golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

type Inspector struct {
	fset           *token.FileSet
	packages       []*packages.Package
	filesToPackage map[string]*packages.Package
	ins            *goinspector.Inspector
}

// NewInspector will find all .go files rooted at the
// given path and load them into the returned Inspector
func NewInspector(path string, tests bool) (*Inspector, error) {
	fset := token.NewFileSet()

	cfg := packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedModule |
			packages.NeedCompiledGoFiles,
		Fset:  fset,
		Dir:   path,
		Tests: tests,
	}

	pkgs, err := packages.Load(&cfg, "./...")
	if err != nil {
		return nil, fmt.Errorf("error loading all Go files: %w", err)
	}

	var files []*ast.File

	for _, pkg := range pkgs {
		files = append(files, pkg.Syntax...)
	}

	goinspector.New(files)

	filesToPackages := make(map[string]*packages.Package)

	for _, pkg := range pkgs {
		for _, fname := range pkg.CompiledGoFiles {
			filesToPackages[fname] = pkg
		}
	}

	return &Inspector{
		fset:           fset,
		packages:       pkgs,
		filesToPackage: filesToPackages,
		ins:            goinspector.New(files),
	}, nil
}

// Walks the inspectors ast.Nodes in depth first order executing the
// provided fn before visiting children.
func (ins *Inspector) WalkAST(types []ast.Node, fn func(n ast.Node)) {
	ins.ins.Preorder(types, fn)
}

func (ins *Inspector) WalkASTWithStack(types []ast.Node, fn func(n ast.Node, push bool, stack []ast.Node) (proceed bool)) {
	ins.ins.WithStack(types, fn)
}

// Iterates over all Go packages in the order returned by
// golang.org/x/tools/go/packages.Load and calls the provided function
// with that package as an input.
func (ins *Inspector) WalkPackages(fn func(pkg *packages.Package)) {
	for _, pkg := range ins.packages {
		fn(pkg)
	}
}

func (ins *Inspector) PackageName(pos token.Pos) string {
	f := ins.fset.File(pos)
	pkg := ins.filesToPackage[f.Name()]
	return pkg.PkgPath
}

func (ins *Inspector) FileName(pos token.Pos) string {
	f := ins.fset.File(pos)
	return f.Name()
}
