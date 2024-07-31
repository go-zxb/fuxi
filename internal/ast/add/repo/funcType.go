package addRepo

import (
	"github.com/go-zxb/fuxi/pkg"
	"go/ast"
)

func (a *AddRepo) FuncType() {
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

	update := ast.FuncType{
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
					Type: ast.NewIdent("gen.ResultInfo"),
				},
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
					Type: ast.NewIdent("gen.ResultInfo"),
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}

	bodyCreate := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	return r.Q.` + pkg.InitialLetter(a.Name) + `.Create(` + a.Name + `)
}
`

	bodyUpdate := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	return Q.Where(Q.ID.Eq(user.ID)).Updates(` + a.Name + `)
}
`

	bodyQueryList := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(` + a.Name + `.ID))
	return q.Find()
}
`

	bodyQueryByIDList := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(id))
	return q.Find()
}
`
	bodyQueryByIDFirst := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(id))
	return q.First()
}
`

	bodyQueryByIDAndObjFirst := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(` + a.Name + `.ID))
	return q.First()
}
`

	bodyQueryNoParams := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	//Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	//q := Q.Where()
	return "",nil
}
`

	bodyDelete := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	return Q.Where(Q.ID.Eq(id)).Delete()
}
`

	switch a.Method {
	case "GET":
		if a.ISByID {
			a.funcType = queryID
			if a.IsReturnList {
				a.bodyCode = bodyQueryByIDList
			} else {
				a.bodyCode = bodyQueryByIDFirst
			}
		} else {
			//入参Query对象 可包含更多查询字段
			if a.NoParams {
				a.funcType = queryNoParams
				a.bodyCode = bodyQueryNoParams
			} else {
				a.funcType = queryList
				if a.IsReturnList {
					a.bodyCode = bodyQueryList
				} else {
					a.bodyCode = bodyQueryByIDAndObjFirst
				}
			}

		}

	case "POST":
		a.funcType = create
		a.bodyCode = bodyCreate
	case "PUT":
		a.funcType = update
		a.bodyCode = bodyUpdate
	case "DELETE":
		a.funcType = delete_
		a.bodyCode = bodyDelete
	default:
		a.funcType = ast.FuncType{}
		a.bodyCode = bodyQueryList
	}
}

func (a *AddRepo) GetReturnType() *ast.FieldList {
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

func (a *AddRepo) getReturnType() *ast.Field {
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
