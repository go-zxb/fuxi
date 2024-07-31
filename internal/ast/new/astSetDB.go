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
)

type ASTSetDB struct {
	Name     string
	Imports  map[string]string
	FilePath string
}

func (a *ASTSetDB) InsetCode() error {
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
		if ok && funcDecl.Name.Name == "initDatabase" {

			// 检查函数体中是否已经存在fmt.Println语句
			for _, stmt := range funcDecl.Body.List {
				if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
					if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
						if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if selExpr.X.(*ast.Ident).Name == "query"+pkg.InitialLetter(a.Name) {
								hasInsert = true
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
				// 插入新的语句到Start函数体中
				//r := funcDecl.Body.List[len(funcDecl.Body.List)-2:]
				funcDecl.Body.List = append([]ast.Stmt{newStmt}, funcDecl.Body.List...)
			}
		}
		// 查找导入声明
		importDecl, ok := n.(*ast.GenDecl)
		if ok && importDecl.Tok == token.IMPORT {
			// 添加新的导入路径
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
				// 将结构｛｝添加到参数中
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

	if !hasInsert {
		// 打开文件以写入修改后的内容
		file, err := os.Create(a.FilePath)
		if err != nil {
			log.Println("⚠️❎ Error creating file:", err)
			return err
		}
		defer file.Close()

		// 格式化并写入修改后的AST
		err = format.Node(file, fset, node)
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
