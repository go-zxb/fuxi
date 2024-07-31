package addapi

import (
	"github.com/go-zxb/fuxi/pkg"
	"strings"
)

func (a *AddApi) Tmpl() {
	create := `
	package main
func (a *UserApi) CreateUser(ctx *gin.Context) {
	var oo model.` + pkg.InitialLetter(a.Name) + `
	if err := ctx.ShouldBindJSON(&oo); err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `(&oo)
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("成功").Success(ctx)
}
`

	updateCode := `
package main
func (a *UserApi) Tmpl(ctx *gin.Context) {
	var oo model.` + pkg.InitialLetter(a.Name) + `
	if err := ctx.ShouldBindJSON(&oo); err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `(&oo)
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("成功").Success(ctx)
}
`

	delete := `
	package main
func (a *UserApi) CreateUser(ctx *gin.Context) {
	id := strconvx.StringToUint(ctx.Param("id"))
	if id == 0 {
		a.reply.New().WithCode(400).WithMessage("id not found").Fail(ctx)
		return
	}
	err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `(id)
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("删除成功").Success(ctx)
}
`

	queryList := `
	package main
func (a *UserApi) CreateUser(ctx *gin.Context) {
	id := strconvx.StringToUint(ctx.Param("id"))
	if id == 0 {
		a.reply.New().WithCode(400).WithMessage("id not found").Fail(ctx)
		return
	}
	var oo model.` + pkg.InitialLetter(a.Name) + `Query
	oo.ID = id
	list, err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `(&oo)
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("成功").WithData(list).Success(ctx)
}
`

	queryID := `
	package main
func (a *UserApi) CreateUser(ctx *gin.Context) {
	id := strconvx.StringToUint(ctx.Param("id"))
	if id == 0 {
		a.reply.New().WithCode(400).WithMessage("id not found").Fail(ctx)
		return
	}
	
	` + a.Name + `, err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `(id)
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("成功").WithData(` + a.Name + `).Success(ctx)
}
`
	queryNoParams := `
	package main
func (a *UserApi) CreateUser(ctx *gin.Context) {
	id := strconvx.StringToUint(ctx.Param("id"))
	if id == 0 {
		a.reply.New().WithCode(400).WithMessage("id not found").Fail(ctx)
		return
	}
	
	resp,err := a.sv.` + pkg.InitialLetter(a.ApiFunc) + `()
	if err != nil {
		a.reply.New().WithCode(400).WithMessage(err.Error()).Fail(ctx)
		return
	}
	a.reply.New().WithCode(200).WithMessage("成功").WithData(resp).Success(ctx)
}
`

	switch a.Method {
	case "GET":
		if a.ISByID {
			a.code = queryID
		} else {
			if a.NoParams {
				a.code = queryNoParams
			} else {
				a.code = queryList
				if !strings.Contains(a.Api, "/:id") {
					a.Api += "/:id"
				}
			}
		}

	case "POST":
		a.code = create
	case "PUT":
		a.code = updateCode
		if !strings.Contains(a.Api, "/:id") {
			a.Api += "/:id"
		}
	case "DELETE":
		a.code = delete
	default:
		a.code = ""
	}
}
