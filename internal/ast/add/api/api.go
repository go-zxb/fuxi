package addapi

import (
	"errors"
	"fmt"
	"github.com/go-zxb/fuxi/internal/ast/base"
	"github.com/go-zxb/fuxi/pkg"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

type AddApi struct {
	Name         string
	FilePath     string
	Method       string
	Api          string //接口
	ApiFunc      string //接口方法
	Results      []*ast.Field
	ISByID       bool
	IsReturnList bool
	code         string
	NoParams     bool
	FuXiAst      base.FuXiAst
}

func (a *AddApi) InsertRouter() error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return err
	}

	hasInsert := false
	apiHasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		//注册路由
		if block, ok := n.(*ast.BlockStmt); ok {
			for _, stmt := range block.List {
				if blockStmt, ok := stmt.(*ast.BlockStmt); ok {
					// 检查是否已注册
					for _, stmt := range blockStmt.List {
						if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
							if call, ok := exprStmt.X.(*ast.CallExpr); ok {
								if len(call.Args) > 1 {
									if lit, ok := call.Args[1].(*ast.SelectorExpr); ok && lit.Sel.Name == pkg.InitialLetter(a.ApiFunc) {
										fmt.Println(a.Method, pkg.InitialLetter(a.ApiFunc), "apiHandle interface接口已存在")
										hasInsert = true
									}
								}
							}
						}
					}

					if !hasInsert {
						// 在 {} 块的结尾插入 POST 接口
						postCall := &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "groupV1"},
								Sel: &ast.Ident{Name: a.Method},
							},
							Args: []ast.Expr{
								&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"/%s"`, a.Api)},
								&ast.Ident{Name: "r.api." + pkg.InitialLetter(a.ApiFunc)},
							},
						}
						postExprStmt := &ast.ExprStmt{X: postCall}
						blockStmt.List = append(blockStmt.List, postExprStmt)
					}
					return false
				}
			}
		}

		//插入interface接口
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
						// 检查是否已存在
						for _, method := range interfaceType.Methods.List {
							ident := method.Names[0].Name
							if ident == pkg.InitialLetter(a.ApiFunc) {
								fmt.Println(pkg.InitialLetter(a.ApiFunc), "路由接口已经存在")
								apiHasInsert = true
								return false
							}
						}

						if !apiHasInsert {
							// 插入接口方法签名
							newMethod := &ast.Field{
								Names: []*ast.Ident{ast.NewIdent(pkg.InitialLetter(a.ApiFunc))},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Names: []*ast.Ident{ast.NewIdent("ctx")},
												Type:  &ast.StarExpr{X: ast.NewIdent("gin.Context")},
											},
										},
									},
								},
							}
							interfaceType.Methods.List = append(interfaceType.Methods.List, newMethod)
						}
						return false
					}
				}
			}
		}

		return true
	})

	if !hasInsert || !apiHasInsert {
		err = a.FuXiAst.SaveNode(node, fset, a.FilePath)
		if err != nil {
			log.Println("✅ 生成路由代码写入文件时出错:", err)
			return err
		}
		log.Println("✅  生成路由代码成功。")
	} else {
		return errors.New("路由接口已经存在")
	}

	return nil
}
