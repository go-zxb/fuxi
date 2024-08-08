package openapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/cmd/openapi/findfield"
	"github.com/spf13/cobra"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	Method  string
	Path    string
	Handler string
	Group   string
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
					walkToGetRoute(path)
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

// 获取路由信息
func walkToGetRoute(path string) {
	// 解析 Go 文件
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok {
					if sel.Sel.Name == "GET" || sel.Sel.Name == "POST" || sel.Sel.Name == "PUT" || sel.Sel.Name == "DELETE" {
						method := sel.Sel.Name
						path_ := x.Args[0].(*ast.BasicLit).Value
						handler := ""
						if val, ok := x.Args[1].(*ast.SelectorExpr); ok {
							handler = val.Sel.Name
						} else {
							handler = x.Args[1].(*ast.Ident).Name
						}
						group[ident.Name] = ident.Name
						routes = append(routes, &Route{Method: method, Path: path_, Handler: handler, Group: ident.Name})
					}
				}

			}
		case *ast.CommentGroup:

		default:

		}
		return true
	})
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
		Servers: []Server{
			{
				URL:         "http://127.0.0.1:8080/api/v1",
				Description: "开发环境",
			},
			{
				URL:         "http://127.0.0.1:8081/api/v1",
				Description: "测试环境",
			},
			{
				URL:         "http://127.0.0.1:8082/api/v1",
				Description: "生产环境",
			},
		},
	}

	for _, route := range routes {
		path := strings.Trim(route.Path, `"`)
		if _, ok := doc.Paths[path]; !ok {
			doc.Paths[path] = make(map[string]Operation)
		}

		var operation = Operation{
			Summary:     fmt.Sprintf("%s %s", route.Method, path),
			Description: fmt.Sprintf("Handler function is %s", route.Handler),
			Tags:        []string{route.Group},
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
					In:          "query",
					Description: fmt.Sprintf("%s ID", structName),
					Required:    true,
					Schema: Schema{
						Type: "integer",
					},
				})
			}
		case "POST", "PUT":
			var schema Schema
			pro := map[string]Schema{}
			strslice := make([]string, 0)
			for _, info := range findStruct {
				pro[InitialLetter(info.Name)] = Schema{
					Type:        info.Type,
					Example:     "",
					Description: info.Comment,
				}
				strslice = append(strslice, InitialLetter(info.Name))
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
			if strings.Contains(route.Path, "/:id") {
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:        "id",
					In:          "query",
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

func InitialLetter(word string) string {
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
