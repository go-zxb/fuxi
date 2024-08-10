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

type ASTSetDB struct {
	Name     string
	Imports  map[string]string
	FilePath string
	FuXiAst  base.FuXiAst
}

func (a *ASTSetDB) InsetCode() error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	hasInsert := false
	dbInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "initDatabase" {

			// 检测是否插入DB注册语句
			dbInsert = a.FuXiAst.HasFunctionCall(node, "initDatabase", "query"+pkg.InitialLetter(a.Name))

			if !dbInsert {
				// 创建一个新的表达式语句
				newStmt := &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "query" + pkg.InitialLetter(a.Name)},
							Sel: &ast.Ident{Name: "SetDefault"},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "db",
							},
						},
					},
				}
				funcDecl.Body.List = append([]ast.Stmt{newStmt}, funcDecl.Body.List...)
			}
		}
		// 添加包路径
		importDecl, ok := n.(*ast.GenDecl)
		if ok && importDecl.Tok == token.IMPORT {
			importSpec := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "query" + pkg.InitialLetter(a.Name) + fmt.Sprintf("%q", a.Imports["query"]),
				},
			}
			importDecl.Specs = append(importDecl.Specs, importSpec)
			importSpec = &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: a.Name + "Model" + fmt.Sprintf("%q", a.Imports["model"]),
				},
			}
			importDecl.Specs = append(importDecl.Specs, importSpec)
		}
		return true
	})

	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if selectorExpr.Sel.Name == "AutoMigrate" {
			hasInsert = false
			for _, arg := range callExpr.Args {
				compositeLit, ok := arg.(*ast.CompositeLit)
				if ok {
					if compositeLit.Type.(*ast.SelectorExpr).Sel.Name == pkg.InitialLetter(a.Name) {
						hasInsert = true
					}
				}
			}

			if !hasInsert {
				// 将xx结构体｛｝添加到参数中
				appType := &ast.SelectorExpr{
					X:   &ast.Ident{Name: fmt.Sprintf("\n%sModel", a.Name)},
					Sel: &ast.Ident{Name: pkg.InitialLetter(a.Name)},
				}
				userLit := &ast.CompositeLit{Type: appType}
				callExpr.Args = append(callExpr.Args, userLit)
			}
		}

		return true
	})

	if !hasInsert || !dbInsert {
		err = a.FuXiAst.SaveNode(node, fset, a.FilePath)
		if err != nil {
			log.Println("⚠️❎ 数据库相关代码生成写入文件时出错:", err)
			return err
		}
		log.Println("✅ 数据库相关代码生成成功。")
	} else {
		log.Println("⚠️❎ 数据库相关代码已存在，无需重复生成！")
	}
	return nil
}
