package model

import "github.com/go-zxb/fuxi/internal/ast/base"

type TmplData struct {
	ModuleName  string
	StructName  string
	ProjectName string
	GormGenPath string
	Import      map[string]string
	Data        []*GenCodeStruct
	Table       string
	Desc        string
}

type Import struct {
	Name string `json:"name"`
}

type GenCodeStruct struct {
	Name        string `json:"name"`
	StructType  string `json:"structType"`
	SqlType     string `json:"sqlType"`
	StructTag   string `json:"structTag"`
	SqlTag      string `json:"SqlTag"`
	Description string `json:"description"`
	Value       any    `json:"value"`
}

type CodeModel struct {
	Fields     []*GenCodeStruct `json:"fields"`
	StructName string           `json:"structName"`
	Table      string           `json:"table"`
	Desc       string           `json:"desc"`
}

type AST struct {
	Name        string
	ImportsName []string
	Imports     []string
	FilePath    string
	FuXiAst     base.FuXiAst
}
