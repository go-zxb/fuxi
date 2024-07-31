package addService

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

type AddService struct {
	Name         string
	FilePath     string
	Method       string
	Api          string //接口
	ApiFunc      string //接口方法
	Results      []*ast.Field
	funcType     ast.FuncType
	ISByID       bool
	IsReturnList bool
	NoParams     bool
	obj          string
	ReturnType   string
}

func (a *AddService) InsertService() error {

	a.FuncType()

	block := &ast.BlockStmt{
		List: []ast.Stmt{
			// 创建一个赋值语句
			//&ast.AssignStmt{
			//	Lhs: []ast.Expr{
			//		ast.NewIdent("_"),
			//		ast.NewIdent("err"),
			//	},
			//	Tok: token.DEFINE,
			//	Rhs: []ast.Expr{
			//		&ast.CallExpr{
			//			Fun: &ast.SelectorExpr{
			//				X:   ast.NewIdent("s.repo"),
			//				Sel: ast.NewIdent(pkg.InitialLetter(a.ApiFunc)),
			//			},
			//			Args: []ast.Expr{
			//				ast.NewIdent(a.Name),
			//			},
			//		},
			//	},
			//},
			// 创建一个返回语句
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s.repo"),
							Sel: ast.NewIdent(pkg.InitialLetter(a.ApiFunc)),
						},
						Args: []ast.Expr{
							ast.NewIdent(a.obj),
						},
					},
				},
			},
		},
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
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
					Type: &a.funcType,
					Body: block,
					Doc: &ast.CommentGroup{
						List: []*ast.Comment{
							{Text: "// " + pkg.InitialLetter(a.ApiFunc) + a.Method + " " + a.Name},
						},
					},
				}
				// 在文件末尾插入 UpdateUser 方法
				i := len(file.Decls)
				file.Decls = append(file.Decls[:i], append([]ast.Decl{updateUserMethod}, file.Decls[i:]...)...)
			}
			return false

		}
		return true
	})

	//添加repo依赖倒置接口
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
							// 在接口定义中插入新的方法
							newMethod := &ast.Field{
								Names: []*ast.Ident{ast.NewIdent(pkg.InitialLetter(a.ApiFunc))},
								Type:  &a.funcType,
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
	if !hasInsert || !apiHasInsert {
		// 打开文件以写入修改后的内容
		f, err := os.Create(a.FilePath)
		if err != nil {
			log.Println("❎ Error creating node:", err)
			return err
		}
		defer f.Close()

		// 格式化并写入修改后的AST
		err = format.Node(f, fset, node)
		if err != nil {
			log.Println("✅ AddService 生成代码写入文件时出错:", err)
			return err
		}
		log.Println("✅ AddService 生成代码成功。")
	}
	return nil
}
