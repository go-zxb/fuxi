package base

import (
	"go/ast"
	"go/token"

	"github.com/go-zxb/fuxi/pkg"
)

func BodySimple1(obj, apiFunc string) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			// 创建一个返回语句
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s.repo"),
							Sel: ast.NewIdent(pkg.InitialLetter(apiFunc)),
						},
						Args: []ast.Expr{
							ast.NewIdent(obj),
						},
					},
				},
			},
		},
	}
}

func BodySimple2(obj, apiFunc string) *ast.BlockStmt {
	return &ast.BlockStmt{
		List: []ast.Stmt{
			// 创建一个赋值语句
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("_"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("s.repo"),
							Sel: ast.NewIdent(pkg.InitialLetter(apiFunc)),
						},
						Args: []ast.Expr{
							ast.NewIdent(obj),
						},
					},
				},
			},
			&ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("err"),
				},
			},
		},
	}
}
