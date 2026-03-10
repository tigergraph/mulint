package mulint

import (
	"fmt"
	"go/ast"
	"go/types"
)

// getCallInfo returns the package/type and method/function name for a call.
// It mirrors the call resolution behavior we previously relied on from gosec.
func getCallInfo(n ast.Node, pkg *types.Package, info *types.Info) (string, string, error) {
	switch node := n.(type) {
	case *ast.CallExpr:
		switch fn := node.Fun.(type) {
		case *ast.SelectorExpr:
			switch expr := fn.X.(type) {
			case *ast.Ident:
				if expr.Obj != nil && expr.Obj.Kind == ast.Var {
					t := info.TypeOf(expr)
					if t != nil {
						return t.String(), fn.Sel.Name, nil
					}
					return "undefined", fn.Sel.Name, fmt.Errorf("missing type info")
				}
				return expr.Name, fn.Sel.Name, nil
			case *ast.SelectorExpr:
				if expr.Sel != nil {
					t := info.TypeOf(expr.Sel)
					if t != nil {
						return t.String(), fn.Sel.Name, nil
					}
					return "undefined", fn.Sel.Name, fmt.Errorf("missing type info")
				}
			case *ast.CallExpr:
				switch call := expr.Fun.(type) {
				case *ast.Ident:
					if call.Name == "new" && len(expr.Args) > 0 {
						t := info.TypeOf(expr.Args[0])
						if t != nil {
							return t.String(), fn.Sel.Name, nil
						}
						return "undefined", fn.Sel.Name, fmt.Errorf("missing type info")
					}
					if call.Obj != nil {
						switch decl := call.Obj.Decl.(type) {
						case *ast.FuncDecl:
							ret := decl.Type.Results
							if ret != nil && len(ret.List) > 0 {
								ret1 := ret.List[0]
								if ret1 != nil {
									t := info.TypeOf(ret1.Type)
									if t != nil {
										return t.String(), fn.Sel.Name, nil
									}
									return "undefined", fn.Sel.Name, fmt.Errorf("missing type info")
								}
							}
						}
					}
				}
			}
		case *ast.Ident:
			return pkg.Name(), fn.Name, nil
		}
	}

	return "", "", fmt.Errorf("unable to determine call info")
}
