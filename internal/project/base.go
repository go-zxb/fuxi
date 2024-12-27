package project

import (
	"fmt"
	"log"

	newapi "github.com/go-zxb/fuxi/internal/ast/new"
	"github.com/go-zxb/fuxi/internal/model"
)

// InsertInitRouterCode 插入初始化路由代码
func InsertInitRouterCode(moduleName, middle, filename string) {
	// 自动插入初始化路由代码
	a := newapi.ASTRouter{
		Name:     filename,
		Imports:  []string{fmt.Sprintf("%s/internal/router/%s%s", moduleName, middle, filename)},
		FilePath: "core/core.go",
	}
	// 自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		log.Println(err)
	}
	e := newapi.Enforcer{
		AST: model.AST{
			Name:        "casbin",
			ImportsName: []string{"casbinService"},
			Imports:     []string{fmt.Sprintf("%s/internal/service/system/%s", moduleName, "casbin")},
			FilePath:    "core/core.go",
		},
	}
	e.Enforcer()
}

// InsertGormGenCode 插入gormGen代码
func InsertGormGenCode(moduleName, middle, filename string) {
	// 自动插入初始化路由代码
	a := newapi.ASTGormGen{
		Name:     filename,
		Imports:  []string{fmt.Sprintf("%s/%s/%s%s", moduleName, gormGenPath, middle, filename)},
		FilePath: gormGenPath + "/main.go",
	}
	// 自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		fmt.Println(err)
		return
	}
}

func InsertSetDB(moduleName, middle, filename string) {
	// 自动插入初始化路由代码
	a := newapi.ASTSetDB{
		Name: filename,
		Imports: map[string]string{
			"query": fmt.Sprintf("%s/%s/%s%s/query", moduleName, modelPath, middle, filename),
			"model": fmt.Sprintf("%s/%s/%s%s", moduleName, modelPath, middle, filename)},
		FilePath: "core/data/mysql.go",
	}
	// 自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		fmt.Println(err)
		return
	}
}
