package inspector

import (
	"fmt"
	"go/ast"
	"go/token"

	goinspector "golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

type Inspector struct {
	fset     *token.FileSet
	packages map[string]*packages.Package
	ins      *goinspector.Inspector
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
		fset:     fset,
		packages: filesToPackages,
		ins:      goinspector.New(files),
	}, nil
}

// Walks the inspectors ast.Nodes in depth first order executing the
// provided fn before visiting children.
func (ins *Inspector) Preorder(types []ast.Node, fn func(n ast.Node)) {
	ins.ins.Preorder(types, fn)
}

func (ins *Inspector) PackageName(pos token.Pos) string {
	f := ins.fset.File(pos)
	pkg := ins.packages[f.Name()]
	return pkg.PkgPath
}

func (ins *Inspector) FuncName(fn *ast.FuncDecl) string {
	if fn.Recv == nil {
		return fn.Name.String()
	}

	receiver := ""
	switch v := fn.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		ident, ok := v.X.(*ast.Ident)
		if ok {
			// TODO - is it important to keep around
			receiver = ident.Name
		}
	case *ast.Ident:
		receiver = v.Name
	}

	return fmt.Sprintf("%s.%s", receiver, fn.Name.String())
}
