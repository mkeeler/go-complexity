package goroutines

import (
	"go/ast"

	asthelpers "github.com/mkeeler/go-complexity/internal/ast-helpers"
	"github.com/mkeeler/go-complexity/internal/inspector"
	"golang.org/x/tools/go/packages"
)

type GoRoutines struct {
	Packages map[string]*PackageGoRoutines
}

type PackageGoRoutines struct {
	Count     int
	Files     map[string]int
	Functions map[string]int
}

func CountGoRoutineInvocations(ins *inspector.Inspector) *GoRoutines {
	routines := &GoRoutines{
		Packages: make(map[string]*PackageGoRoutines),
	}

	ins.WalkPackages(func(pkg *packages.Package) {
		routines.Packages[pkg.PkgPath] = &PackageGoRoutines{
			Files:     make(map[string]int),
			Functions: make(map[string]int),
		}

		for _, file := range pkg.CompiledGoFiles {
			routines.Packages[pkg.PkgPath].Files[file] = 0
		}
	})

	ins.WalkASTWithStack([]ast.Node{(*ast.GoStmt)(nil), (*ast.FuncDecl)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if !push {
				break
			}
			packageName := ins.PackageName(n.Pos())
			functionName := asthelpers.FuncName(n)

			routines.Packages[packageName].Functions[functionName] = 0
		case *ast.GoStmt:
			packageName := ins.PackageName(n.Pos())
			fileName := ins.FileName(n.Pos())

			var functionName string
			// iterating from the element before this GoStmt to the
			// second element in the list (first is the ast.File)
			for i := len(stack) - 2; i > 0; i-- {
				if fn, ok := stack[i].(*ast.FuncDecl); ok {
					functionName = asthelpers.FuncName(fn)
					break
				}
			}

			packageRoutines := routines.Packages[packageName]
			packageRoutines.Count += 1
			packageRoutines.Files[fileName] += 1
			packageRoutines.Functions[functionName] += 1
		}

		return true
	})

	return routines
}
