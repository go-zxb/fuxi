package newapi

import (
	"fmt"
	"github.com/go-zxb/fuxi/pkg"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

type ASTRouter struct {
	Name     string
	Imports  []string
	FilePath string
}

func (a *ASTRouter) InsetCode() error {
	// 创建一个文件集
	fset := token.NewFileSet()

	// 解析文件
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	// 遍历AST
	hasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		// 查找Start函数
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == "InitRouter" {

			// 检查函数体中是否已经存在fmt.Println语句
			for _, stmt := range funcDecl.Body.List {
				if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
					if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
						if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if selExpr, ok := selExpr.X.(*ast.CallExpr); ok {
								if selExpr, ok := selExpr.Fun.(*ast.SelectorExpr); ok {
									if selExpr.Sel.Name == fmt.Sprintf("New%sRouter", pkg.InitialLetter(a.Name)) {
										hasInsert = true
									}
								}
							}
						}
					}
				}
			}

			if !hasInsert {
				// 创建一个新的表达式语句
				newStmt := &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: strings.ToLower(a.Name)},
									Sel: &ast.Ident{Name: fmt.Sprintf("New%sRouter", pkg.InitialLetter(a.Name))},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "engine",
									},
								},
							},
							Sel: ast.NewIdent(fmt.Sprintf("Init%sRouter", pkg.InitialLetter(a.Name))),
						},
					},
				}
				// 插入新的语句到Start函数体中
				funcDecl.Body.List = append(funcDecl.Body.List, newStmt)
			}
		}
		// 查找导入声明
		importDecl, ok := n.(*ast.GenDecl)
		if ok && importDecl.Tok == token.IMPORT {
			// 添加新的导入路径
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
		// 打开文件以写入修改后的内容
		file, err := os.Create(a.FilePath)
		if err != nil {
			fmt.Println("⚠️❎ Error creating file:", err)
			return err
		}
		defer file.Close()

		// 格式化并写入修改后的AST
		err = format.Node(file, fset, node)
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
