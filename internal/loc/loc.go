package loc

import (
	"go/ast"
	"go/token"

	asthelpers "github.com/mkeeler/go-complexity/internal/ast-helpers"
	"github.com/mkeeler/go-complexity/internal/inspector"
	"golang.org/x/tools/go/packages"
)

type LinesOfCode struct {
	Packages map[string]*PackageLinesOfCode
}

type PackageLinesOfCode struct {
	Lines     LineCounts
	Files     map[string]LineCounts
	Functions map[string]LineCounts
}

type LineCounts struct {
	Code     int
	Comments int
}

// CacluclateLinesOfCode will calculate the number of non-comment lines of code
// within each package and return them within the corresponding map type
func CalculateLinesOfCode(ins *inspector.Inspector) *LinesOfCode {
	loc := &LinesOfCode{
		Packages: make(map[string]*PackageLinesOfCode),
	}
	ins.WalkPackages(func(pkg *packages.Package) {
		loc.Packages[pkg.PkgPath] = calculatePackageLOC(pkg)
	})

	return loc
}

// calculatePackageLOC will calculate the number of non-comment lines of code
// in the package. This requires both the Fset and the Syntax fields of the
// Package to be properly configured. If using the Inspector from
// github.com/mkeeler/go-complexity/internal/inspector then this will already
// be the case.
func calculatePackageLOC(pkg *packages.Package) *PackageLinesOfCode {
	loc := &PackageLinesOfCode{
		Files:     make(map[string]LineCounts),
		Functions: make(map[string]LineCounts),
	}

	for _, file := range pkg.Syntax {
		lineCounts := calculateFileLOC(pkg.Fset, file)
		loc.Lines.Code += lineCounts.Code
		loc.Lines.Comments += lineCounts.Comments
		loc.Files[file.Name.String()] = lineCounts

		var v ast.Visitor
		v = asthelpers.VisitFn(func(node ast.Node) ast.Visitor {
			switch n := node.(type) {
			case *ast.FuncDecl:
				loc.Functions[asthelpers.FuncName(n)] = calculateFunctionLOC(pkg.Fset, n)
				return nil
			default:
				return v
			}
		})

		ast.Walk(v, file)
	}

	return loc
}

// calculateFileLOC will calculate the number of non-comment lines of code
// within a single ast.File. The corresponding token.FileSet must have been
// used during parsing of the file and should contain position information
// for the file.
func calculateFileLOC(fset *token.FileSet, file *ast.File) LineCounts {
	var counts LineCounts
	f := fset.File(file.Pos())
	counts.Code = f.Line(file.End())

	for _, cg := range file.Comments {
		startLine := f.Line(cg.Pos())
		endLine := f.Line(cg.End())

		commentLoC := endLine - startLine + 1
		counts.Comments += commentLoC
		counts.Code -= commentLoC
	}

	return counts
}

// calculateFunctionLOC will calculate the number of non-comment lines of
// code within a single function declaration. The corresponding token.FileSet
// must have been used during parsing of the file and should contain position
// information for the file.
func calculateFunctionLOC(fset *token.FileSet, fn *ast.FuncDecl) LineCounts {
	f := fset.File(fn.Pos())

	start := f.Line(fn.Pos())
	end := f.Line(fn.End())

	v := funcBodyVisitor{
		f: f,
	}

	ast.Walk(&v, fn)

	return LineCounts{Code: end - start + 1 - v.commentLines, Comments: v.commentLines}
}

type funcBodyVisitor struct {
	f            *token.File
	commentLines int
}

func (v *funcBodyVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CommentGroup:
		start := v.f.Line(n.Pos())
		end := v.f.Line(n.End())

		v.commentLines += end - start + 1
	}

	return v
}
