package openapi

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/cmd/openapi/findfield"
	"github.com/spf13/cobra"
)

type Tag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OpenAPIDoc struct {
	OpenAPI string `json:"openapi"`
	Info    struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Version     string `json:"version"`
	} `json:"info"`
	Tags       []Tag                           `json:"tags"`
	Paths      map[string]map[string]Operation `json:"paths"`
	Components struct {
		Schemas         map[string]Schema         `json:"schemas"`
		SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
	} `json:"components"`
	Servers []Server `json:"servers"`
}

type Operation struct {
	Summary     string                `json:"summary"`
	Description string                `json:"description"`
	Tags        []string              `json:"tags"`
	Parameters  []Parameter           `json:"parameters"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Schema      Schema      `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

type RequestBody struct {
	Content map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content"`
}

type Schema struct {
	Type        string            `json:"type,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Description string            `json:"description"`
	Example     any               `json:"example,omitempty"`
}

type SecurityScheme struct {
	Type   string `json:"type"`
	Scheme string `json:"scheme,omitempty"`
	Name   string `json:"name,omitempty"`
	In     string `json:"in"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Route struct {
	Method    string
	Path      string
	Handler   string
	Group     string
	GroupName string
}

var routes []*Route
var group map[string]string

var output string

func init() {
	OpenapiCmd.Flags().StringVarP(&output, "output", "o", "docs/openapi/openapi.json", "output file path")
}

func genDoc() {
	group = make(map[string]string)
	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, "router") {
			if !info.IsDir() {
				if strings.HasSuffix(path, ".go") {
					gp := getGroup(path)
					if gp == "" {
						gp = "没有分组的Api"
					}
					ParseRoutes(path, gp)
				}
			}
		}
		return nil
	})

	// 生成 OpenAPI 文档
	openapiDoc := generateOpenAPIDoc()
	_ = os.MkdirAll(path.Dir(output), 0644)
	err := os.WriteFile(output, []byte(openapiDoc), 0644)
	if err != nil {
		log.Fatalln("❌", err.Error())
	}
	log.Println("✅", "接口文档生成成功", "👌")
}

func getGroup(path string) (group string) {
	// 解析 Go 文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if strings.Contains(x.Name.Name, "Router") {
				_, ok := x.Type.(*ast.StructType)
				if !ok {
					return true
				}

				group = strings.ReplaceAll(x.Name.Name, "Router", "")
			}
		}

		return true
	})
	return group
}

// 递归检查是否存在Group调用
func findGroupPath(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "Group" {
				if len(e.Args) != 1 {
					return "", false
				}
				lit, ok := e.Args[0].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					return "", false
				}
				return strings.Trim(lit.Value, `"`), true
			}
			// 递归检查X部分是否存在Group调用
			if path, ok := findGroupPath(sel.X); ok {
				return path, true
			}
		}
		return findGroupPath(e.Fun)
	case *ast.SelectorExpr:
		return findGroupPath(e.X)
	default:
		return "", false
	}
}

func ParseRoutes(filePath, getGroup string) []*Route {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	vars := make(map[string]string) // 变量名到Group路径的映射

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			// 处理变量定义，例如 groupV1 := ...
			if x.Tok != token.DEFINE {
				return true
			}
			if len(x.Lhs) != 1 || len(x.Rhs) != 1 {
				return true
			}
			lhsIdent, ok := x.Lhs[0].(*ast.Ident)
			if !ok {
				return true
			}
			// 查找Group路径
			groupPath, ok := findGroupPath(x.Rhs[0])
			if ok {
				vars[lhsIdent.Name] = groupPath
			}

		case *ast.CallExpr:
			// 处理HTTP方法调用，如GET、POST
			selExpr, ok := x.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			method := selExpr.Sel.Name
			if !isHTTPMethod(method) {
				return true
			}
			// 获取接收者变量名
			receiver, ok := selExpr.X.(*ast.Ident)
			if !ok {
				return true
			}
			varName := receiver.Name
			groupPath, exists := vars[varName]
			if !exists {
				return true
			}
			// 提取路径
			if len(x.Args) < 1 {
				return true
			}
			pathLit, ok := x.Args[0].(*ast.BasicLit)
			if !ok || pathLit.Kind != token.STRING {
				return true
			}
			path := strings.Trim(pathLit.Value, `"`)
			// 提取处理函数
			var handler string
			if len(x.Args) >= 2 {
				switch arg := x.Args[1].(type) {
				case *ast.SelectorExpr:
					handler = arg.Sel.Name
				case *ast.Ident:
					if arg.Name == "nil" {
						handler = "" // 处理nil情况
					} else {
						handler = arg.Name
					}
				default:
					handler = ""
				}
			}
			routes = append(routes, &Route{
				Method:    method,
				Path:      path,
				Handler:   handler,
				Group:     groupPath,
				GroupName: getGroup,
			})
		}
		return true
	})

	return routes
}

// 判断是否为HTTP方法
func isHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
		return true
	default:
		return false
	}
}

func generateOpenAPIDoc() string {
	conf := config.GetConfig()
	doc := OpenAPIDoc{
		OpenAPI: "3.0.0",
		Info: struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Version     string `json:"version"`
		}{
			Title:       conf.System.Name,
			Version:     conf.System.Version,
			Description: conf.System.Description,
		},
		Tags:  []Tag{},
		Paths: make(map[string]map[string]Operation),
		Components: struct {
			Schemas         map[string]Schema         `json:"schemas"`
			SecuritySchemes map[string]SecurityScheme `json:"securitySchemes"`
		}{
			Schemas:         make(map[string]Schema),
			SecuritySchemes: make(map[string]SecurityScheme),
		},
		Servers: []Server{},
	}

	for _, server := range conf.DocsServers {
		doc.Servers = append(doc.Servers, Server{
			URL:         server.URL,
			Description: server.Description,
		})
	}

	for _, route := range routes {
		if !strings.HasPrefix(route.Group, "/") {
			route.Group = fmt.Sprintf("%s%s", "/", route.Group)
		}
		if strings.HasSuffix(route.Group, "/") {
			route.Group = strings.TrimSuffix(route.Group, "/")
		}

		route.Path = strings.Trim(route.Path, `"`)
		if !strings.HasPrefix(route.Path, "/") {
			route.Path = fmt.Sprintf("%s%s", "/", route.Path)
		}

		path := route.Group + route.Path
		if strings.Contains(path, "/:id") {
			path = strings.ReplaceAll(path, "/:id", "/{id}")
		}
		if _, ok := doc.Paths[path]; !ok {
			doc.Paths[path] = make(map[string]Operation)
		}

		var operation = Operation{
			Summary:     fmt.Sprintf("%s %s", route.Method, path),
			Description: fmt.Sprintf("Handler function is %s", route.Handler),
			Tags:        []string{route.GroupName},
			Parameters:  make([]Parameter, 0),
			Responses: map[string]Response{
				"200": {
					Description: "成功",
					Content: map[string]MediaType{
						"application/json": {
							Schema: schema(),
						},
					},
				},
				"404": {
					Description: "请求有误",
					Content: map[string]MediaType{
						"application/json": {
							Schema: schema(),
						},
					},
				},
			},
			Security: []map[string][]string{},
		}

		findStruct, structName, auth, err := findfield.FindStruct(route.Handler)
		if err != nil {
			log.Fatalln("❌", err.Error())
		}
		if structName == "" {
			structName = route.Handler
		}
		cmmSlice := strings.Split(auth, " ")
		if len(cmmSlice) > 0 && strings.TrimSpace(cmmSlice[len(cmmSlice)-1:][0]) == "true" {
			operation.Security = []map[string][]string{
				{
					"Authorization": {},
				},
			}
		}

		switch route.Method {
		case "GET", "DELETE":
			if strings.Contains(route.Path, "/:id") {
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:        "id",
					In:          "path",
					Description: fmt.Sprintf("%s ID", structName),
					Required:    true,
					Schema: Schema{
						Type: "integer",
					},
				})
			}
			for _, info := range findStruct {
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:        InitialLetterToLower(info.Name),
					In:          "query",
					Description: fmt.Sprintf("%s", info.Comment),
					Required:    false,
					Schema: Schema{
						Type: handleType(info.Type),
					},
				})
			}
		case "POST", "PUT":
			var schema Schema
			pro := map[string]Schema{}
			strslice := make([]string, 0)
			for _, info := range findStruct {
				pro[InitialLetterToLower(info.Name)] = Schema{
					Type:        info.Type,
					Example:     "",
					Description: info.Comment,
				}
				strslice = append(strslice, InitialLetterToLower(info.Name))
			}
			schema = Schema{
				Properties: pro,
				Required:   strslice,
			}
			operation.RequestBody = &RequestBody{
				Content: map[string]MediaType{
					"application/json": {
						Schema: schema,
					},
				},
			}
			// 判断是否包含id
			if strings.Contains(route.Path, "/:id") {
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:        "id",
					In:          "path",
					Description: fmt.Sprintf("%s ID", structName),
					Required:    true,
					Schema: Schema{
						Type: "integer",
					},
				})
			}
		}

		doc.Paths[path][strings.ToLower(route.Method)] = operation
	}

	for _, tag := range group {
		doc.Tags = append(doc.Tags, Tag{
			Name:        tag,
			Description: fmt.Sprintf("%s", tag),
		})
	}

	doc.Components.SecuritySchemes = map[string]SecurityScheme{
		"Authorization": {
			Type: "apiKey",
			Name: "Token",
			In:   "header",
		},
	}

	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		log.Fatal("Error generating JSON:", err)
	}
	return string(jsonData)
}

func schema() Schema {
	return Schema{
		Type: "object",
		Properties: map[string]Schema{
			"code": {
				Type: "integer",
			},
			"success": {
				Type: "boolean",
			},
			"message": {
				Type: "string",
			},
			"data": {
				Type: "object",
			},
			"time": {
				Type: "integer",
			},
		},
		Required: []string{"code", "success", "message", "data", "time"},
	}
}

func InitialLetterToLower(word string) string {
	if word == "ID" || word == "UID" {
		return strings.ToLower(word)
	}
	if len(word) > 1 {
		return strings.ToLower(word[0:1]) + word[1:]
	}
	return strings.ToUpper(word)
}

var OpenapiCmd = &cobra.Command{
	Use:   "openapi",
	Short: "生成 openapi 文档",
	Long:  `生成 openapi 文档`,
	Run: func(cmd *cobra.Command, args []string) {
		genDoc()
	},
}

func handleType(t string) string {
	switch t {
	case "int", "uint", "int64":
		return "integer"
	default:
		return t
	}

}
