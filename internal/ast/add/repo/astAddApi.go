package addRepo

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

type AddRepo struct {
	Name         string
	FilePath     string
	Method       string
	Api          string //接口
	ApiFunc      string //接口方法
	Results      []*ast.Field
	funcType     *ast.FuncType
	bodyCode     string
	ISByID       bool
	IsReturnList bool
	NoParams     bool
	ReturnType   string
	FuXiAst      base.FuXiAst
}

func (a *AddRepo) InsertRepo() error {
	// 预定义的模板代码

	updateBlock := &ast.BlockStmt{
		List: []ast.Stmt{
			// 创建一个赋值语句
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("Q"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent("Q"),
						},
						Sel: ast.NewIdent("User"),
					},
				},
			},
			// 创建一个返回语句
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("Q"),
									Sel: ast.NewIdent("Where"),
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("Q"),
												Sel: ast.NewIdent("ID"),
											},
											Sel: ast.NewIdent("Eq"),
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("user"),
												Sel: ast.NewIdent("ID"),
											},
										},
									},
								},
							},
							Sel: ast.NewIdent("Updates"),
						},
						Args: []ast.Expr{
							ast.NewIdent("user"),
						},
					},
				},
			},
		},
	}

	if updateBlock != nil {
	}

	a.FuncType()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		return errors.New(err.Error())
	}

	templateFile, err := parser.ParseFile(fset, "", a.bodyCode, parser.ParseComments)
	if err != nil {
		return errors.New(err.Error())
	}

	var templateBlock *ast.BlockStmt
	for _, decl := range templateFile.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			templateBlock = funcDecl.Body
			break
		}
	}

	if templateBlock == nil {
		return errors.New("template block not found")
	}

	hasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		if file, ok := n.(*ast.File); ok {

			// 检查是否已经存在 xxx 方法
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == pkg.InitialLetter(a.ApiFunc) {
					fmt.Println(pkg.InitialLetter(a.ApiFunc), "Repo方法已经存在")
					hasInsert = true
					return false // 已经存在 xxx 方法
				}
			}

			if !hasInsert {
				// 插入 XXX 方法
				updateUserMethod := &ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("r")},
								Type:  &ast.StarExpr{X: ast.NewIdent(pkg.InitialLetter(a.Name) + "Repo")},
							},
						},
					},
					Name: ast.NewIdent(pkg.InitialLetter(a.ApiFunc)),
					Type: a.funcType,
					Body: templateBlock,
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

	name, _ := pkg.GetModuleName("go.mod")
	if !a.FuXiAst.HasImport(node, name+"/internal/model/"+a.Name) {
		a.FuXiAst.AddImport(node, name+"/internal/model/"+a.Name)
		hasInsert = false
	}

	if !hasInsert {
		err = a.FuXiAst.SaveNode(node, fset, a.FilePath)
		if err != nil {
			log.Println("✅ AddRepo 生成代码写入文件时出错:", err)
			return err
		}
		log.Println("✅ AddRepo 生成代码成功。")
	} else {
		fmt.Println(pkg.InitialLetter(a.ApiFunc) + "Repo方法已经存在")
		return errors.New(pkg.InitialLetter(a.ApiFunc) + "Repo方法已经存在")
	}
	return nil
}
