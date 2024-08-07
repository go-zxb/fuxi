package addRepo

import (
	"github.com/go-zxb/fuxi/internal/ast/add/base"
	"github.com/go-zxb/fuxi/pkg"
	"go/ast"
)

func (a *AddRepo) FuncType() {
	if a.ReturnType == "objlist" {
		a.IsReturnList = true
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
	bodyQueryByIDFirst1 := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(id))
	return q.First()
}
`

	bodyQueryByIDFirst2 := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(id))
	_, err := q.First()
	if err != nil {
		return "", err
	}
	return "", nil
}
`

	bodyQueryByIDFirst3 := `
package main
func (r *` + pkg.InitialLetter(a.Name) + `Repo) ` + pkg.InitialLetter(a.ApiFunc) + `() {
	Q := r.Q.` + pkg.InitialLetter(a.Name) + `
	q := Q.Where(Q.ID.Eq(id))
	return q.Count()
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
			a.funcType = base.QueryByID(a.IsReturnList, a.Name, a.ReturnType)
			if a.IsReturnList {
				a.bodyCode = bodyQueryByIDList
			} else {
				if a.ReturnType == "obj" {
					a.bodyCode = bodyQueryByIDFirst1
				} else {
					if goReturnType(a.ReturnType) == "0" {
						a.bodyCode = bodyQueryByIDFirst3
					} else {
						a.bodyCode = bodyQueryByIDFirst2
					}
				}

			}
		} else {
			//入参Query对象 可包含更多查询字段
			if a.NoParams {
				a.funcType = base.QueryNoParams(a.IsReturnList, a.Name, a.ReturnType)
				a.bodyCode = bodyQueryNoParams
			} else {
				a.funcType = base.QueryByObj(a.IsReturnList, a.Name, a.ReturnType)
				if a.IsReturnList {
					a.bodyCode = bodyQueryList
				} else {
					a.bodyCode = bodyQueryByIDAndObjFirst
				}
			}

		}

	case "POST":
		a.funcType = base.Create(a.Name)
		a.bodyCode = bodyCreate
	case "PUT":
		a.funcType = base.Update(a.Name)
		a.bodyCode = bodyUpdate
	case "DELETE":
		a.funcType = base.Delete1()
		a.bodyCode = bodyDelete
	default:
		a.funcType = &ast.FuncType{}
		a.bodyCode = bodyQueryList
	}
}

func goReturnType(returnType string) string {
	switch returnType {
	case "int", "uint", "int64", "float64":
		return "0"
	case "string":
		return "1"
	default:
		return "1"
	}
}
