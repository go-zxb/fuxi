package project

import (
	"encoding/json"
	"fmt"
	newapi "github.com/go-zxb/fuxi/internal/ast/new"
	"github.com/go-zxb/fuxi/internal/gpt"
	"github.com/go-zxb/fuxi/internal/model"
	"github.com/go-zxb/fuxi/internal/project/base"
	"github.com/go-zxb/fuxi/pkg"
	"github.com/go-zxb/fuxi/template"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
)

// PathData 路径数据
type PathData struct {
	TmplPath      string
	FilePath      string
	FileName      string
	FileExtension string
}

var apiCodePath []*PathData

var (
	filename     = ""
	modelPath    = "internal/model"
	gormGenPath  = "cmd/gorm/gen"
	openapi      = "docs/openapi"
	question     = ""
	isOutputJson = ""
	localJson    = "false"
)

// init 初始化
func init() {
	apiCodePath = make([]*PathData, 0)
	NewApiCmd.Flags().StringVarP(&apiPath, "apiPath", "a", apiPath, "apiHandle文件路径 ")
	NewApiCmd.Flags().StringVarP(&routerPath, "routerPath", "r", routerPath, "router文件路径 ")
	NewApiCmd.Flags().StringVarP(&servicePath, "servicePath", "s", servicePath, "service文件路径 ")
	NewApiCmd.Flags().StringVarP(&repoPath, "repoPath", "p", repoPath, "repo文件路径 ")
	NewApiCmd.Flags().StringVarP(&filename, "filename", "n", "", "file name ")
	NewApiCmd.Flags().StringVarP(&modelPath, "modelPath", "m", modelPath, "model文件路径 ")
	NewApiCmd.Flags().StringVarP(&question, "question", "q", question, "你的设计需求: 比如帮我设计一个用户数据表...()")
	NewApiCmd.Flags().StringVarP(&gormGenPath, "gormGenPath", "g", gormGenPath, "gormGen文件路径 ")
	NewApiCmd.Flags().StringVarP(&isOutputJson, "json", "j", "false", "是否输出json数据")
	NewApiCmd.Flags().StringVarP(&debug, "debug", "d", "false", "是否开启debug模式 ")
	NewApiCmd.Flags().StringVarP(&localJson, "skip", "l", "false", "是否跳过生成代码步骤 ")
}

var NewApiCmd = &cobra.Command{
	Use:     "api:new",
	Short:   "一键生成xxCRUD代码",
	Long:    "一键生成xxCRUD代码",
	Example: "生成一个用户CRUD→ fuxi api:new user",
	Run:     createCode,
}

// createCode 创建代码
func createCode(cmd *cobra.Command, args []string) {
	packagename, err := base.GetModuleName("go.mod")
	if err != nil {
		log.Fatalln("❗️请先初始化项目: fuxi project -n 项目名称")
	}

	if pkg.HasChinese(filename) {
		log.Fatalln("------❎ 名称不能包含中文❎--------")
	}
	jsonStr := ""
	var genCode = &model.CodeModel{}
	if filename != "" && localJson == "true" {
		_, err = os.Stat("docs/json/" + filename + ".json")
		if err == nil {
			file, err := os.ReadFile("docs/json/" + filename + ".json")
			if err != nil {
				log.Println(err)
			} else {
				jsonStr = string(file)
				_ = json.Unmarshal(file, &genCode)
			}
		}
	} else {
		if question == "" {
			log.Fatalln("------❎ 名称或需求不可为空❎--------")
		}
	}

	if len(genCode.Fields) == 0 {
		log.Println("➡️你的需求是: ", question)
		genCode, jsonStr, err = gpt.GenCode(question)
		if err != nil {
			log.Fatalln(err)
		}
		if filename == "" {
			filename = strings.ToLower(genCode.StructName)
		}
		//导出json数据
		if isTrue(isOutputJson) {
			err = os.MkdirAll("docs/json", 0750)
			if err != nil {
				log.Fatalln(err)
			}
			_ = os.Remove("docs/json/" + filename + ".json")
			tmpl := template.Must(template.New("").Parse(jsonStr))
			file, err := os.OpenFile("docs/json/"+filename+".json", os.O_CREATE, 0750)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			_ = tmpl.Execute(file, nil)
		}
	} else {
		if filename == "" {
			filename = strings.ToLower(genCode.StructName)
		}
	}

	//对genCode进行单词长短排序😊
	sort.Slice(genCode.Fields, func(i, j int) bool {
		return len(genCode.Fields[i].Name) < len(genCode.Fields[j].Name)
	})

	addApiCodePath("api.go", apiPath, filename, ".go")
	addApiCodePath("router.go", routerPath, filename, ".go")
	addApiCodePath("service.go", servicePath, filename, ".go")
	addApiCodePath("repo.go", repoPath, filename, ".go")
	addApiCodePath("model.go", modelPath, filename, ".go")
	addApiCodePath("gormGen.go", gormGenPath, filename, ".go")
	addApiCodePath("openapi", openapi, filename, ".json")
	var slicePath = make([]string, 0)
	var ok = "n"
	var isOK = false //判断是否有改动文件
	for _, data := range apiCodePath {
		goFilePaht := fmt.Sprintf("%s/%s/%s%s", data.FilePath, data.FileName, data.FileName, data.FileExtension)
		//文件是否存在
		_, err = os.Stat(goFilePaht)
		if err == nil {
			//debug模式先删掉文件
			if isTrue(debug) {
				if ok == "n" {
					log.Println(goFilePaht, "🍵 Hi 文件已存在...")
					path, _ := os.Getwd()
					log.Println("⚠️ 老家伙你现在在的目录是:", path)
					fmt.Print("🧒老伙计请确认是否继续,继续会覆盖相关代码文件 (Y/n)⬇️:")
					scan, err := fmt.Scan(&ok)
					if err != nil {
						log.Fatalln(scan, err)
					}
				}

				if ok == "Y" {
					_ = os.Remove(goFilePaht)
				} else {
					log.Fatalln("🚶‍♀️告辞🚶")
				}
			} else {
				//跳过 继续执行下一个文件
				continue
			}
		}

		//创建文件夹
		err = os.MkdirAll(fmt.Sprintf("%s/%s", data.FilePath, data.FileName), 0750)
		if err != nil {
			log.Println(err)
		}
		//log.Println("🫎文件路径:", goFilePaht)
		slicePath = append(slicePath, goFilePaht)
		//生成代码
		err = TmplExecute(packagename, goFilePaht, data, genCode)
		if err != nil {
			log.Println(err)
			//如果失败一个 就删除已生成的文件
			for _, s := range slicePath {
				_ = os.Remove(s)
			}
			return
		}
		if data.FileExtension == ".go" {
			_ = FormatGoCode(goFilePaht)
		}
		isOK = true
		log.Println("✅", goFilePaht, "文件创建成功👌！")
	}

	if !isOK {
		//没有生成任何文件
		log.Fatalln("⚠️❎ 相关代码文件已存在，生成失败！")
	}

	InsertInitRouterCode(packagename)
	InsertGormGenCode(packagename)
	InsertSetDB(packagename)

	//运行mod tidy
	err = pkg.RunCommand("go", "mod", "tidy")
	if err != nil {
		log.Fatalln("failed to mod tidy:", err)
	}
	log.Println("✅ :执行 go mod tidy 成功👌")

	//运行gormGen生成gen代码
	err = pkg.RunCommand("go", "run", gormGenPath+"/main.go")
	if err != nil {
		log.Fatalln("❎ Error ", err)
	}
	log.Println("✅ :执行", gormGenPath, "/main.go", "生成gen代码成功👌！")
	fmt.Println(`代码精妙我自豪😎，
分享经验乐陶陶😄；
你我同欢心相照👫，
共筑梦想乐逍遥🌟。`)
}

// TmplExecute 模板渲染
func TmplExecute(packageName, goFilePath string, data *PathData, genCode *model.CodeModel) error {
	//读取模板
	bytes, err := templatex.TmplData.ReadFile(data.TmplPath + ".tmpl")
	if err != nil {
		return err
	}

	//创建tmpl渲染对象
	tmpl, err := template.New("demo").Funcs(base.FuncMap).Parse(string(bytes))
	file, err := os.OpenFile(goFilePath, os.O_CREATE, 0750)
	if err != nil {
		return err
	}
	defer file.Close()

	//渲染模板
	if err = tmpl.Execute(file, model.TmplData{
		ModuleName:  packageName,
		Table:       genCode.Table,
		StructName:  strings.ToLower(filename),
		GormGenPath: fmt.Sprintf("%s/%s/query", modelPath, strings.ToLower(filename)),
		Data:        genCode.Fields,
		Import:      handleImport(genCode.Fields),
		Desc:        genCode.Desc,
	}); err != nil {
		return err
	}
	return nil
}

// InsertInitRouterCode 插入初始化路由代码
func InsertInitRouterCode(packageName string) {
	//自动插入初始化路由代码
	a := newapi.ASTRouter{
		Name:     filename,
		Imports:  []string{fmt.Sprintf("%s/internal/router/%s", packageName, filename)},
		FilePath: "core/core.go",
	}
	//自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		log.Println(err)
	}

}

// InsertGormGenCode 插入gormGen代码
func InsertGormGenCode(packageName string) {
	//自动插入初始化路由代码
	a := newapi.ASTGormGen{
		Name:     filename,
		Imports:  []string{fmt.Sprintf("%s/%s/%s", packageName, gormGenPath, filename)},
		FilePath: gormGenPath + "/main.go",
	}
	//自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		fmt.Println(err)
		return
	}
}

func InsertSetDB(packageName string) {
	//自动插入初始化路由代码
	a := newapi.ASTSetDB{
		Name: filename,
		Imports: map[string]string{
			"query": fmt.Sprintf("%s/%s/%s/query", packageName, modelPath, filename),
			"model": fmt.Sprintf("%s/%s/%s", packageName, modelPath, filename)},
		FilePath: "core/data/mysql.go",
	}
	//自动插入初始化路由代码
	if err := a.InsetCode(); err != nil {
		fmt.Println(err)
		return
	}
}

// add 添加路径
func addApiCodePath(tmplPath, filepath, filename, suffix string) {
	apiCodePath = append(apiCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}

// handleImport 处理导入包
func handleImport(data []*model.GenCodeStruct) map[string]string {
	var imports = make(map[string]string)
	for _, i2 := range data {
		switch i2.StructType {
		case "time.Time":
			imports["time"] = "time"
		case "io.Reader":
			imports["io"] = "io"
		case "io.Writer":
			imports["io"] = "io"
		case "io.Closer":
			imports["io"] = "io"
		case "os.File":
			imports["os"] = "os"
		case "os.Process":
			imports["os"] = "os"
		case "sync.Mutex":
			imports["sync"] = "sync"
		case "sync.RWMutex":
			imports["sync"] = "sync"
		case "sync.WaitGroup":
			imports["sync"] = "sync"
		case "sync.Once":
			imports["sync"] = "sync"
		case "context.Context":
			imports["context"] = "context"
		case "net.IP":
			imports["net"] = "net"
		case "net.IPNet":
			imports["net"] = "net"
		case "net.URL":
			imports["net"] = "net"
		case "net.Addr":
			imports["net"] = "net"
		case "filepath.WalkFunc":
			imports["apiCodePath/filepath"] = "apiCodePath/filepath"
		}
	}
	return imports
}

// FormatGoCode 格式化Go代码
func FormatGoCode(filePath string) error {
	return pkg.RunCommand("gofmt", "-w", filePath)
}
