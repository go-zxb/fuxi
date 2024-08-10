package base

import (
	"go/ast"

	"github.com/go-zxb/fuxi/pkg"
)

func Update(name string) *ast.FuncType {
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
					Type: ast.NewIdent("gen.ResultInfo"),
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}
}

func UpdateSvc(name string) *ast.FuncType {
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
