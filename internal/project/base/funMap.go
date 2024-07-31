package base

import (
	"github.com/go-zxb/fuxi/pkg"
	"text/template"
)

// FuncMap 模板函数，处理template模板
var FuncMap = template.FuncMap{
	"initialLetter": pkg.InitialLetter,
	"openapi":       openapi,
	"isLast":        isLast,
	"dataType":      dataType,
}

// 处理接口文档这段
func openapi(val any) any {
	switch val {
	case "string":
		return "string"
	case "uint":
		return "number"
	case "int":
		return "number"
	case "int64":
		return "number"
	case "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		return "string"
	}
}

// 处理数据类型
func dataType(val any) any {
	switch val {
	case "time.Time":
		return "int64"
	default:
		return val
	}
}

// 判断处理渲染最后一个元素用
func isLast(index int, length int) bool {
	return index == length-1
}
