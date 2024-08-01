package base

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

type FuXiAst struct {
}

// ContainsMethod 检查是否包含xx方法签名
func (FuXiAst) ContainsMethod(file *ast.File, methodName string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if ts, ok := n.(*ast.FuncDecl); ok {
			if ts.Name.Name == methodName {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// ContainsInterface 检测是否包含接口
func (FuXiAst) ContainsInterface(file *ast.File, typeName string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if ifs, ok := n.(*ast.InterfaceType); ok {
			for _, field := range ifs.Methods.List {
				if len(field.Names) > 0 && field.Names[0].Name == typeName {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

// AddMethodComment 添加方法注释
func (FuXiAst) AddMethodComment(file *ast.File, methodName, comment string) {
	ast.Inspect(file, func(n ast.Node) bool {
		if fd, ok := n.(*ast.FuncDecl); ok {
			if fd.Name.Name == methodName {
				newComment := &ast.Comment{
					Slash: fd.Pos() - 1,
					Text:  comment,
				}
				if fd.Doc != nil {
					fd.Doc.List = append([]*ast.Comment{newComment}, fd.Doc.List...)
				} else {
					fd.Doc = &ast.CommentGroup{
						List: []*ast.Comment{newComment},
					}
				}
				return false
			}
		}
		return true
	})
}

// PrintAST 格式化并输出 AST
func (FuXiAst) PrintAST(fset *token.FileSet, file *ast.File) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		fmt.Println("Error formatting new code:", err)
		return
	}
	fmt.Println(buf.String())
}

// HasFunctionCall 检查某方法内是否调用了指定函数
func (FuXiAst) HasFunctionCall(file *ast.File, methodName, targetFunc string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == methodName {
			for _, stmt := range funcDecl.Body.List {
				if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
					if selExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
						if selector1, ok := selExpr.Fun.(*ast.SelectorExpr); ok {
							if x, ok := selector1.X.(*ast.Ident); ok && x.Name == targetFunc {
								found = true
								return false
							}
							if selector2, ok := selector1.X.(*ast.CallExpr); ok {
								if selector3, ok := selector2.Fun.(*ast.SelectorExpr); ok {
									if selector3.Sel.Name == targetFunc {
										found = true
										return false
									}
								}
							} else {
								if selector1.Sel.Name == targetFunc {
									found = true
									return false
								}
							}
						} else {
							ident, ok := selExpr.Fun.(*ast.Ident)
							if ok {
								if ident.Name == targetFunc {
									found = true
									return false
								}
							}
						}
					}
				}
			}
		}
		return true
	})
	return found
}

// HasImport 检查是否包含导入
func (FuXiAst) HasImport(file *ast.File, importName string) bool {
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if im, ok := n.(*ast.GenDecl); ok {
			for _, field := range im.Specs {
				if val, ok := field.(*ast.ImportSpec); ok && val.Path.Value == fmt.Sprintf("%q", importName) {
					found = true
					return false
				}
			}
		}
		return true
	})
	return found
}

func (FuXiAst) SaveNode(node *ast.File, fset *token.FileSet, filePath string) error {
	// 打开文件以写入修改后的内容
	fd, err := os.Create(filePath)
	if err != nil {
		return errors.New("open file error:" + err.Error())
	}
	defer fd.Close()

	// 格式化并写入修改后的AST
	return format.Node(fd, fset, node)

}
