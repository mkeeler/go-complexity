package cyclomatic

import (
	"go/ast"
	"go/token"

	"github.com/mkeeler/go-complexity/internal/inspector"
)

func CalculateComplexity(ins *inspector.Inspector) (*CyclomaticComplexity, error) {
	c := newCyclomaticComplexity()

	ins.Preorder([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			score := calculateFunctionCyclomaticComplexity(fn)
			c.recordComplexity(ins.PackageName(fn.Pos()), ins.FuncName(fn), score)
		}
	})

	return c, nil
}

type visitor struct {
	decisions int
	exits     int
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CommClause, *ast.RangeStmt, *ast.ForStmt, *ast.IfStmt, *ast.CaseClause:
		v.decisions += 1
	case *ast.BinaryExpr:
		switch n.Op {
		case token.LAND, token.LOR:
			v.decisions += 1
		}
	case *ast.ReturnStmt:
		v.exits += 1
	}
	return v
}

func calculateFunctionCyclomaticComplexity(fn *ast.FuncDecl) int {
	v := visitor{}
	ast.Walk(&v, fn)
	return v.decisions - v.exits + 2
}
