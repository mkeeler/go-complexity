package asthelpers

import "go/ast"

type VisitFn func(node ast.Node) ast.Visitor

func (fn VisitFn) Visit(node ast.Node) ast.Visitor {
	return fn(node)
}
