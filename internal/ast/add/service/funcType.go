package addService

import (
	"go/ast"

	"github.com/go-zxb/fuxi/internal/ast/add/base"
)

func (a *AddService) FuncType() {
	if a.ReturnType == "objlist" {
		a.IsReturnList = true
	}

	switch a.Method {
	case "GET":
		if a.ISByID {
			a.Body = base.BodySimple1("id", a.ApiFunc)
			a.funcType = base.QueryByID(a.IsReturnList, a.Name, a.ReturnType)
			a.funcTypeSv = base.QueryByID(a.IsReturnList, a.Name, a.ReturnType)
		} else {
			if a.NoParams {
				// if a.ReturnType == "" {
				// 	a.Body = &ast.BlockStmt{
				// 		List: []ast.Stmt{
				// 			&ast.ReturnStmt{
				// 				Results: []ast.Expr{
				// 					&ast.Ident{Name: "nil"},
				// 				},
				// 			},
				// 		},
				// 	}
				// } else {
				// 	a.Body = base.BodySimple1("", a.ApiFunc)
				// }
				a.Body = base.BodySimple1("", a.ApiFunc)
				a.funcType = base.QueryNoParams(a.IsReturnList, a.Name, a.ReturnType)
				a.funcTypeSv = base.QueryNoParams(a.IsReturnList, a.Name, a.ReturnType)
			} else {
				a.Body = base.BodySimple1(a.Name, a.ApiFunc)
				a.funcType = base.QueryByObj(a.IsReturnList, a.Name, a.ReturnType)
				a.funcTypeSv = base.QueryByObj(a.IsReturnList, a.Name, a.ReturnType)
			}
		}

	case "POST":
		a.funcType = base.Create(a.Name)
		a.funcTypeSv = base.Create(a.Name)
		a.Body = base.BodySimple1(a.Name, a.ApiFunc)
	case "PUT":
		a.funcType = base.Update(a.Name)
		a.funcTypeSv = base.UpdateSvc(a.Name)
		a.Body = base.BodySimple2(a.Name, a.ApiFunc)
	case "DELETE":
		a.funcType = base.Delete1()
		a.funcTypeSv = base.Delete2()
		a.Body = base.BodySimple2("id", a.ApiFunc)
	default:
		a.Body = base.BodySimple1("", a.ApiFunc)
		a.funcType = &ast.FuncType{}
		a.funcTypeSv = &ast.FuncType{}
	}
}
