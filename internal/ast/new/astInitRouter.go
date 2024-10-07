package newapi

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/go-zxb/fuxi/internal/ast/base"
	"github.com/go-zxb/fuxi/pkg"
)

type ASTRouter struct {
	Name     string
	Imports  []string
	FilePath string
	FuXiAst  base.FuXiAst
}

func (a *ASTRouter) InsetCode() error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	hasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "InitRouter" {

			// 检测是否注册了该路由
			hasInsert = a.FuXiAst.HasFunctionCall(node, "InitRouter", fmt.Sprintf("New%sRouter", pkg.InitialLetter(a.Name)))

			if !hasInsert {
				// 创建一个新的表达式语句
				newStmt := &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: a.Name},
									Sel: &ast.Ident{Name: fmt.Sprintf("New%sRouter", pkg.InitialLetter(a.Name))},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "engine",
									},
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "enforcer",
									},
								},
							},
							Sel: ast.NewIdent(fmt.Sprintf("Init%sRouter", pkg.InitialLetter(a.Name))),
						},
					},
				}
				funcDecl.Body.List = append(funcDecl.Body.List, newStmt)
			}
		}
		// 添加包路径
		importDecl, ok := n.(*ast.GenDecl)
		if ok && importDecl.Tok == token.IMPORT {
			for _, s := range a.Imports {
				importSpec := &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: a.Name + fmt.Sprintf("%q", s),
					},
				}
				importDecl.Specs = append(importDecl.Specs, importSpec)
			}
		}
		return true
	})

	if !hasInsert {
		err = a.FuXiAst.SaveNode(node, fset, a.FilePath)
		if err != nil {
			log.Println("⚠️❎ Router 生成代码写入文件时出错:", err)
			return err
		}
		log.Println("✅ Router 生成代码已成功更新。")
	} else {
		log.Println("⚠️❎ Router 代码已存在 无须重复生成")
	}
	return nil
}
