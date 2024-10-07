package newapi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/go-zxb/fuxi/internal/model"
)

type Enforcer struct {
	model.AST
}

func (e Enforcer) Enforcer() {
	// 解析文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, e.FilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}

	// 遍历文件中的所有函数声明
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// 检查函数名是否为 InitRouter
			if x.Name.Name == "InitRouter" {
				// 在函数体中插入代码
				insertCasbinEnforcer(x)
			}
		}
		// 添加包路径
		importDecl, ok := n.(*ast.GenDecl)
		if ok && importDecl.Tok == token.IMPORT {
			for i, s := range e.Imports {
				importSpec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: e.ImportsName[i] + fmt.Sprintf("%q", s),
					},
				}
				importDecl.Specs = append(importDecl.Specs, importSpec)
			}
		}
		return true
	})

	// 保存修改后的数据
	err = e.FuXiAst.SaveNode(node, fset, e.FilePath)
	if err != nil {
		fmt.Println("Error formatting modified AST:", err)
	}
}

func insertCasbinEnforcer(funcDecl *ast.FuncDecl) {
	// 检查是否已经存在 enforcer := service.NewCasbin(data.GetDB())
	for _, stmt := range funcDecl.Body.List {
		if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			if len(assignStmt.Lhs) == 1 && len(assignStmt.Rhs) == 1 {
				if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok && ident.Name == "enforcer" {
					return
				}
			}
		}
	}
	// 创建要插入的代码
	casbinEnforcerStmt := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "enforcer"}},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "casbinService"},
				Sel: &ast.Ident{Name: "NewCasbin"},
			},
			Args: []ast.Expr{&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "data"},
					Sel: &ast.Ident{Name: "GetDB"},
				},
			}},
		}},
	}

	// 将代码插入到函数体中
	funcDecl.Body.List = append([]ast.Stmt{casbinEnforcerStmt}, funcDecl.Body.List...)
}
