package addapi

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

// 添加apiHandle
func (a *AddApi) InsertApiHandle() error {
	// 预定义的模板代码
	a.Tmpl()

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, a.FilePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	templateFile, err := parser.ParseFile(fset, "", a.code, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var templateBlock *ast.BlockStmt
	for _, decl := range templateFile.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			templateBlock = funcDecl.Body
			break
		}
	}

	if templateBlock == nil {
		log.Fatal("template block not found")
	}

	// 创建注释组
	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text: "// This is a comment for the xx function",
			},
		},
	}

	hasInsert := false
	ast.Inspect(node, func(n ast.Node) bool {
		if file, ok := n.(*ast.File); ok {

			// 检查是否已经存在 xx 方法
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Name.Name == pkg.InitialLetter(a.ApiFunc) {
					fmt.Println(pkg.InitialLetter(a.ApiFunc), "ApiHandle 方法已经存在")
					hasInsert = true
					return false // 已经存在 xx 方法
				}
			}

			if !hasInsert {
				// 插入 UpdateUser 方法
				updateUserMethod := &ast.FuncDecl{
					Recv: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{ast.NewIdent("a")},
								Type:  &ast.StarExpr{X: ast.NewIdent(pkg.InitialLetter(a.Name) + "Api")},
							},
						},
					},
					Name: ast.NewIdent(pkg.InitialLetter(a.ApiFunc)),
					Type: &ast.FuncType{
						Params: &ast.FieldList{
							List: []*ast.Field{
								{
									Names: []*ast.Ident{ast.NewIdent("ctx")},
									Type:  &ast.StarExpr{X: ast.NewIdent("gin.Context")},
								},
							},
						},
						Results: &ast.FieldList{},
					},
					Body: templateBlock,
					Doc:  commentGroup,
				}
				// 在文件末尾插入 xx 方法
				i := len(file.Decls)
				file.Decls = append(file.Decls[:i], append([]ast.Decl{updateUserMethod}, file.Decls[i:]...)...)
			}
			return false
		}
		return true
	})

	if !hasInsert {
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
			log.Println("✅ ApiHandle 生成代码写入文件时出错:", err)
			return err
		}
		log.Println("✅ ApiHandle 生成代码成功。")
	}
	return nil
}
