package base

import (
	"go/ast"

	"github.com/go-zxb/fuxi/pkg"
)

func QueryByObj(isList bool, name, returnType string) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(name)},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "model." + pkg.InitialLetter(name) + "Query"},
					},
				},
			},
		},
		Results: getReturnType(isList, name, returnType),
	}
}

func QueryByID(isList bool, name, returnType string) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("id")},
					Type:  ast.NewIdent("uint"),
				},
			},
		},
		Results: getReturnType(isList, name, returnType),
	}
}

func QueryNoParams(isList bool, name, returnType string) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{},
		},
		Results: getReturnType(isList, name, returnType),
	}
}
