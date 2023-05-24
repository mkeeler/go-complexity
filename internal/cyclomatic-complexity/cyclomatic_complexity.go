package cyclomatic

import (
	"go/ast"
	"go/token"

	asthelpers "github.com/mkeeler/go-complexity/internal/ast-helpers"
	"github.com/mkeeler/go-complexity/internal/inspector"
)

func CalculateComplexity(ins *inspector.Inspector) (*CyclomaticComplexity, error) {
	c := newCyclomaticComplexity()

	ins.WalkAST([]ast.Node{(*ast.FuncDecl)(nil)}, func(n ast.Node) {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			score := calculateFunctionCyclomaticComplexity(fn)
			c.recordComplexity(ins.PackageName(fn.Pos()), asthelpers.FuncName(fn), score)
		}
	})

	return c, nil
}

type visitor struct {
	decisions     int
	exits         int
	embeddedFuncs []*visitor
}

func (v *visitor) complexity() int {
	exits := v.exits
	if exits == 0 {
		exits = 1
	}

	score := v.decisions - v.exits + 2

	// include embedded function complexity
	for _, w := range v.embeddedFuncs {
		score += w.complexity()
	}

	return score
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
	case *ast.FuncLit:
		w := &visitor{}
		v.embeddedFuncs = append(v.embeddedFuncs, w)
		return w
	}

	return v
}

func calculateFunctionCyclomaticComplexity(fn *ast.FuncDecl) int {
	v := visitor{}
	ast.Walk(&v, fn)
	return v.complexity()
}
