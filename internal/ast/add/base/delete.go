package base

import "go/ast"

func Delete1() *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("id")},
					Type:  ast.NewIdent("uint"),
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

func Delete2() *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("id")},
					Type:  ast.NewIdent("uint"),
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
