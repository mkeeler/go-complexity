package asthelpers

import (
	"fmt"
	"go/ast"
)

func FuncName(fn *ast.FuncDecl) string {
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
