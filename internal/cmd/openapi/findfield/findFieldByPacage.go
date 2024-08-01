package findfield

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"log"
)

// GetModelFieldByPackage 根据包路径获取结构体字段
func GetModelFieldByPackage(pack, target string) (fields []*FiledInfo) {
	// 加载包
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, pack)
	if err != nil {
		log.Fatal(err)
	}

	if len(pkgs) == 0 {
		log.Fatal("No packages found")
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					structType, ok := x.Type.(*ast.StructType)
					if ok {
						if x.Name.Name == target {
							fmt.Printf("找到结构体: %s.%s\n", pkg.Name, x.Name.Name)
							for _, field := range structType.Fields.List {
								if len(field.Names) > 0 {
									f := HandleField(field, pkg.Name, x.Name.Name)
									if f != nil {
										fields = append(fields, f)
									}
								}
							}
						}

					}
				}
				return true
			})
		}
	}
	return
}
