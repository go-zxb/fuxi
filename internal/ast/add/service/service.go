package addService

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/go-zxb/fuxi/internal/ast/base"
	"github.com/go-zxb/fuxi/pkg"
)

type AddService struct {
	Name         string
	FilePath     string
	Method       string
	Api          string // 接口
	ApiFunc      string // 接口方法
	Results      []*ast.Field
	funcType     *ast.FuncType
	funcTypeSv   *ast.FuncType
	Body         *ast.BlockStmt
	ISByID       bool
	IsReturnList bool
	NoParams     bool
	ReturnType   string
	FuXiAst      base.FuXiAst
	AddRepo      bool
}

func (a *AddService) InsertService() error {

	a.FuncType()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		return errors.New(err.Error())
	}

	hasInsert := false
	apiHasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		if file, ok := n.(*ast.File); ok {
			// 检查是否已经存在 xxx 方法
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == pkg.InitialLetter(a.ApiFunc) {
					fmt.Println(a.ApiFunc, " service方法已经存在")
					hasInsert = true
					return false // 已经存在 xxx 方法
				}
			}

			if !hasInsert {
				// 插入 xxx 方法
				updateUserMethod := &ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("s")},
								Type:  &ast.StarExpr{X: ast.NewIdent(pkg.InitialLetter(a.Name) + "Service")},
							},
						},
					},
					Name: ast.NewIdent(pkg.InitialLetter(a.ApiFunc)),
					Type: a.funcTypeSv,
					Body: a.Body,
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{
							{Text: "// " + pkg.InitialLetter(a.ApiFunc) + a.Method + " " + a.Name},
						},
					},
				}
				i := len(file.Decls)
				file.Decls = append(file.Decls[:i], append([]ast.Decl{updateUserMethod}, file.Decls[i:]...)...)
			}
			return false

		}
		return true
	})

	if a.AddRepo {
		// 添加repo依赖倒置接口
		ast.Inspect(node, func(n ast.Node) bool {
			if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if interfaceType, ok := typeSpec.Type.(*ast.InterfaceType); ok {
							// 检查是否已经存在 CreateUser 方法
							for _, method := range interfaceType.Methods.List {
								ident := method.Names[0].Name
								if ident == pkg.InitialLetter(a.ApiFunc) {
									fmt.Println(pkg.InitialLetter(a.ApiFunc), "repo interface接口已存在")
									apiHasInsert = true
									return false
								}
							}

							if !apiHasInsert {
								newMethod := &ast.Field{
									Names: []*ast.Ident{ast.NewIdent(pkg.InitialLetter(a.ApiFunc))},
									Type:  a.funcType,
								}
								interfaceType.Methods.List = append(interfaceType.Methods.List, newMethod)
							}
							return false

						}
					}
				}
			}
			return true
		},
		)
	}

	name, _ := pkg.GetModuleName("go.mod")
	if !a.FuXiAst.HasImport(node, name+"/internal/model/"+a.Name) {
		a.FuXiAst.AddImport(node, name+"/internal/model/"+a.Name)
		hasInsert = false
	}

	if !hasInsert || !apiHasInsert {
		err = a.FuXiAst.SaveNode(node, fset, a.FilePath)
		if err != nil {
			log.Println("✅ AddService 生成代码写入文件时出错:", err)
			return err
		}
		log.Println("✅ AddService 生成代码成功。")
	} else {
		fmt.Println(pkg.InitialLetter(a.ApiFunc) + "Service 方法已经存在")
		return errors.New(pkg.InitialLetter(a.ApiFunc) + "Service 方法已经存在")
	}
	return nil
}
