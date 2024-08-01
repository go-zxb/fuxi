package findfield

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

type FiledInfo struct {
	Name      string
	Type      string
	Tag       string
	Comment   string
	Child     []*FiledInfo
	ChildType string
}

// FindStruct 查找方法名是否存在并获取结构体字段
func FindStruct(funcName string) (files []*FiledInfo, structName string, err error) {
	err = filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			comments := FunctionComments(path, funcName)
			if comments == "" {
				return nil
			}
			//处理注释
			//index := strings.Index(comments, "@model")
			cmmSlice := strings.Split(comments, " ")
			slice := strings.Split(cmmSlice[len(cmmSlice)-1:][0], ".")
			//获取结构体字段
			fields, _ := GetStruct(".", slice...)
			if len(fields) == 0 {
				log.Println("❌", "没有找到结构体数据,请注意注释格式")
			}
			log.Println("✅", "获取结构体数据成功", funcName, slice[1])
			structName = slice[1]
			files = fields
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		log.Fatalln("❌", err.Error())
	}
	return
}

// FunctionComments 获取注释
func FunctionComments(path string, funcName string) (val string) {

	// 创建一个新的文件集
	fset := token.NewFileSet()

	// 解析源代码
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历 AST 并获取方法的注释
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Doc != nil && x.Name.Name == funcName {
				//fmt.Printf("Method: %s\n", x.Name.Name)
				for _, comment := range x.Doc.List {
					if strings.Contains(comment.Text, "@model") {
						val = comment.Text
						return false
					}

				}
			}
		}
		return true
	})
	return
}

// GetStruct 根据包名和结构体名称,查找获取结构体字段 最后返回结构体字段
func GetStruct(path string, pack ...string) ([]*FiledInfo, string) {
	var fields = make([]*FiledInfo, 0)
	var modelPath string
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			//log.Println(path)
			// 创建一个新的文件集
			fset := token.NewFileSet()
			// 解析源代码
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			ast.Inspect(node, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.TypeSpec:
					if structType, ok := x.Type.(*ast.StructType); ok {
						for _, field := range structType.Fields.List {
							if len(field.Names) > 0 && node.Name.Name == strings.TrimSpace(pack[0]) && x.Name.Name == strings.TrimSpace(pack[1]) {
								modelPath = path
								f := HandleField(field, node.Name.Name, x.Name.Name)
								if f != nil {
									fields = append(fields, f)
								}
							}
						}
					}
					if len(fields) > 0 {
						return false
					}
				case *ast.CommentGroup:

				}
				return true
			})
			if modelPath != "" {
				return filepath.SkipAll
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalln("❌", err.Error())
	}
	return fields, modelPath
}

func HandleField(field *ast.Field, pack, end string) *FiledInfo {
	var field_ *FiledInfo
	switch reflect.TypeOf(field.Type).String() {
	case "*ast.Ident":
		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}
		field_ = &FiledInfo{
			Name:      field.Names[0].Name,
			Type:      fmt.Sprintf("%s", field.Type),
			Tag:       tag,
			Comment:   fmt.Sprintf("%v", strings.TrimSpace(field.Comment.Text())),
			ChildType: "ast.Ident",
		}
		//fmt.Println("*ast.Ident", field.Names[0].Name, field.Type.(*ast.Ident).Name)
	case "*ast.StarExpr":
		_, ok := field.Type.(*ast.StarExpr)
		if ok {
			//fmt.Println("继承结构体名称:", array.X.(*ast.SelectorExpr).Sel.Name)
			//fmt.Println("引用包名:", array.X.(*ast.SelectorExpr).X.(*ast.Ident).Name)
			name_x := field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
			name_sel := field.Type.(*ast.SelectorExpr).Sel.Name
			var fd = make([]*FiledInfo, 0)
			if end != name_sel {
				fd, _ = GetStruct(".", name_x, name_sel)
			}

			t := fmt.Sprintf("%s.%s", name_x, name_sel)
			field_ = handleF(field, t, name_sel, end, "ast.StarExpr", fd)

		}
	case "*ast.SelectorExpr":
		//fmt.Println("类型:SelectorExpr", "字段名:", field.Names[0].Name, "继承结构体名称:",
		//	field.Type.(*ast.SelectorExpr).Sel.Name, "引用包名:", field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name)
		name_x := field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
		name_sel := field.Type.(*ast.SelectorExpr).Sel.Name
		var fd = make([]*FiledInfo, 0)
		if end != name_sel {
			fd, _ = GetStruct(".", name_x, name_sel)
		}

		t := fmt.Sprintf("%s.%s", name_x, name_sel)
		field_ = handleF(field, t, name_sel, end, "ArrayType.ast.SelectorExpr", fd)

	case "*ast.ArrayType":
		//fmt.Println("切片类型:", reflect.TypeOf(field.Type.(*ast.ArrayType).Elt).String())
		switch reflect.TypeOf(field.Type.(*ast.ArrayType).Elt).String() {
		case "*ast.Ident":
			//fmt.Println("ArrayType ast.Ident", field.Names[0].Name, field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name)
			name_x := pack
			name_sel := field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name
			var fd = make([]*FiledInfo, 0)
			if end != name_sel {
				fd, _ = GetStruct(".", name_x, name_sel)
			}

			t := fmt.Sprintf("%s.%s", name_x, name_sel)
			field_ = handleF(field, t, name_sel, end, "ArrayType.ast.Ident", fd)
		case "*ast.SelectorExpr":
			array, ok := field.Type.(*ast.ArrayType)
			if ok {
				//fmt.Println("类型:", "ArrayType,SelectorExpr")
				//fmt.Println("字段名:", field.Names[0].Name)
				//fmt.Println("字段结构体SelectorExpr:", array.Elt)
				//fmt.Println("继承结构体名称:", array.Elt.(*ast.SelectorExpr).Sel.Name)
				//fmt.Println("引用包名:", array.Elt.(*ast.SelectorExpr).X.(*ast.Ident).Name)
				name_x := array.Elt.(*ast.SelectorExpr).X.(*ast.Ident).Name
				name_sel := array.Elt.(*ast.SelectorExpr).Sel.Name
				var fd = make([]*FiledInfo, 0)
				if end != name_sel {
					fd, _ = GetStruct(".", name_x, name_sel)
				}

				t := fmt.Sprintf("%s.%s", name_x, name_sel)
				field_ = handleF(field, t, name_sel, end, "ArrayType.ast.SelectorExpr", fd)

			}
		case "*ast.StarExpr":
			array, ok := field.Type.(*ast.ArrayType).Elt.(*ast.StarExpr)
			if ok {
				//fmt.Println("类型:", "ArrayType,StarExpr")
				//fmt.Println("字段名:", field.Names[0].Name)
				//fmt.Println("字段结构体SelectorExpr:", array.X.(*ast.SelectorExpr))
				//fmt.Println("继承结构体名称:", array.X.(*ast.SelectorExpr).Sel.Name)
				//fmt.Println("引用包名:", array.X.(*ast.SelectorExpr).X.(*ast.Ident).Name)

				name_x := array.X.(*ast.SelectorExpr).X.(*ast.Ident).Name
				name_sel := array.X.(*ast.SelectorExpr).Sel.Name
				var fd = make([]*FiledInfo, 0)
				if end != name_sel {
					fd, _ = GetStruct(".", name_x, name_sel)
				}

				t := fmt.Sprintf("%s.%s", name_x, name_sel)
				field_ = handleF(field, t, name_sel, end, "ArrayType.ast.StarExpr", fd)

			}
		}
	case "*ast.MapType":
		key := field.Type.(*ast.MapType).Key
		ty := field.Type.(*ast.MapType).Value
		//fmt.Println("map type", key,ty)
		tag := ""
		if field.Tag != nil {
			tag = field.Tag.Value
		}
		field_ = &FiledInfo{
			Name:      field.Names[0].Name,
			Type:      fmt.Sprintf("map[%s][%s]", key, ty),
			Tag:       tag,
			Comment:   fmt.Sprintf("%v", strings.TrimSpace(field.Comment.Text())),
			ChildType: "ast.MapType",
		}
	default:
		fmt.Println("未处理类型:", reflect.TypeOf(field.Type).String())
	}
	return field_
}

func handleF(field *ast.Field, t, nameSel, end, ChildType string, fd []*FiledInfo) *FiledInfo {
	tag := ""
	if field.Tag != nil {
		tag = field.Tag.Value
	}

	field_ := &FiledInfo{
		Name:      field.Names[0].Name,
		Type:      t,
		Tag:       tag,
		Comment:   fmt.Sprintf("%v", strings.TrimSpace(field.Comment.Text())),
		Child:     fd,
		ChildType: ChildType,
	}

	if end == nameSel {
		field_.Child = append(field_.Child, &FiledInfo{
			Name:      field_.Name,
			Type:      field_.Type,
			Tag:       field_.Tag,
			Comment:   field_.Comment,
			Child:     nil,
			ChildType: ChildType,
		})
	}
	return field_
}
