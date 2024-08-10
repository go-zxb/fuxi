package base

import (
	"go/ast"

	"github.com/go-zxb/fuxi/pkg"
)

func Create(name string) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(name)},
					Type:  &ast.StarExpr{X: ast.NewIdent("model." + pkg.InitialLetter(name))},
				},
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}
}
