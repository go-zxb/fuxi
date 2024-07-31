package addService

import (
	"github.com/go-zxb/fuxi/pkg"
	"go/ast"
)

func (a *AddService) FuncType() {
	create := ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(a.Name)},
					Type:  &ast.StarExpr{X: ast.NewIdent("model." + pkg.InitialLetter(a.Name))},
				},
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}
	queryList := ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(a.Name)},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "model." + pkg.InitialLetter(a.Name) + "Query"},
					},
				},
			},
		},
		Results: a.GetReturnType(),
	}
	queryID := ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("id")},
					Type:  ast.NewIdent("uint"),
				},
			},
		},
		Results: a.GetReturnType(),
	}

	queryNoParams := ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{},
		},
		Results: a.GetReturnType(),
	}

	delete_ := ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("id")},
					Type:  ast.NewIdent("uint"),
				},
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}

	switch a.Method {
	case "GET":
		if a.ISByID {
			a.funcType = queryID
			a.obj = "id"
		} else {
			if a.NoParams {
				a.funcType = queryNoParams
				a.obj = ""
			} else {
				a.funcType = queryList
				a.obj = a.Name
			}
		}
	case "POST":
		a.funcType = create
		a.obj = a.Name
	case "PUT":
		a.funcType = create
		a.obj = a.Name
	case "DELETE":
		a.funcType = delete_
		a.obj = "id"
	default:
		a.funcType = ast.FuncType{}
		a.obj = a.Name
	}
}

func (a *AddService) GetReturnType() *ast.FieldList {
	if a.IsReturnList {
		return &ast.FieldList{
			List: []*ast.Field{
				{
					Type: &ast.ArrayType{
						Elt: &ast.StarExpr{
							X: &ast.Ident{Name: "model." + pkg.InitialLetter(a.Name)},
						},
					},
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		}
	} else {
		return &ast.FieldList{
			List: []*ast.Field{
				a.getReturnType(),
				{
					Type: ast.NewIdent("error"),
				},
			},
		}
	}
}

func (a *AddService) getReturnType() *ast.Field {
	switch a.ReturnType {
	case "int", "uint", "int64", "float64", "string":
		return &ast.Field{
			Type: ast.NewIdent(a.ReturnType),
		}
	default:
		return &ast.Field{
			Type: &ast.StarExpr{
				X: &ast.Ident{Name: "model." + pkg.InitialLetter(a.Name)},
			},
		}
	}
}
