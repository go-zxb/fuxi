package base

import (
	"go/ast"

	"github.com/go-zxb/fuxi/pkg"
)

func getReturnType(isList bool, name, returnType string) *ast.FieldList {
	if isList {
		return &ast.FieldList{
			List: []*ast.Field{
				{
					Type: &ast.ArrayType{
						Elt: &ast.StarExpr{
							X: &ast.Ident{Name: "model." + pkg.InitialLetter(name)},
						},
					},
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		}
	} else {
		// if returnType == "" {
		// 	return &ast.FieldList{
		// 		List: []*ast.Field{
		// 			{
		// 				Type: ast.NewIdent("error"),
		// 			},
		// 		},
		// 	}
		// }
		return &ast.FieldList{
			List: []*ast.Field{
				goReturnType(returnType, name),
				{
					Type: ast.NewIdent("error"),
				},
			},
		}
	}
}

func goReturnType(returnType, name string) *ast.Field {
	switch returnType {
	case "int", "uint", "int64", "float64", "string":
		return &ast.Field{
			Type: ast.NewIdent(returnType),
		}
	default:
		return &ast.Field{
			Type: &ast.StarExpr{
				X: &ast.Ident{Name: "model." + pkg.InitialLetter(name)},
			},
		}
	}
}
