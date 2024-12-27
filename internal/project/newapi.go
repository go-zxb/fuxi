package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/go-zxb/fuxi/internal/gpt"
	"github.com/go-zxb/fuxi/internal/model"
	"github.com/go-zxb/fuxi/internal/project/base"
	"github.com/go-zxb/fuxi/pkg"
	templatex "github.com/go-zxb/fuxi/template"
	emptytemplatex "github.com/go-zxb/fuxi/template/empty"
	"github.com/spf13/cobra"
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
	question     = ""
	isOutputJson = ""
	localJson    = "false"
	empty        = "false"
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
	NewApiCmd.Flags().StringVarP(&localJson, "localJson", "l", "false", "使用本地json数据生成 ")
	NewApiCmd.Flags().StringVarP(&empty, "empty", "e", "false", "是否生成空代码文件")
}

var NewApiCmd = &cobra.Command{
	Use:     "api:new",
	Short:   "一键生成xxCRUD代码",
	Long:    "一键生成xxCRUD代码",
	Example: "生成一个用户CRUD→ fuxi api:new user",
	Run: func(cmd *cobra.Command, args []string) {
		infoChan := make(chan pkg.CommandInfo)
		go createCodeHandle(infoChan)
		for info := range infoChan {
			if info.Error != nil {
				log.Fatalln("❌", info.Message, info.Error.Error())
			} else {
				log.Println(info.Message)
			}
		}
	},
}

func handleArgs(infoChan chan<- pkg.CommandInfo) *model.CodeModel {
	jsonStr := ""
	var err error

	var genCode = &model.CodeModel{
		Fields: make([]*model.GenCodeStruct, 0),
	}

	fmt.Println(filename != "" && localJson == "true")
	fmt.Println(filename, localJson)
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
			infoChan <- pkg.CommandInfo{Message: "🐮🐴 ------❎ 名称或需求不可为空❎--------", Error: nil}
			return nil
		}
	}

	if len(genCode.Fields) == 0 {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴 别急...正在疯狂的和Ai生产数据...", Error: nil}
		log.Println("➡️你的需求是: ", question)
		genCode, jsonStr, err = gpt.GenCode(question)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴 ------❎ Ai造数据失败了--------", Error: err}
			return nil
		}
		if filename == "" {
			filename = strings.ToLower(genCode.StructName)
		}
		// 导出json数据
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
		infoChan <- pkg.CommandInfo{Message: "🐮🐴 读取本地json数据成功", Error: nil}
		if filename == "" {
			filename = strings.ToLower(genCode.StructName)
		}
	}
	return genCode
}

// createCode 创建代码
func createCodeHandle(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	packagename, err := pkg.GetModuleName("go.mod")
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️请先初始化项目: fuxi project -n 项目名称", Error: err}
		return
	}

	if pkg.HasChinese(filename) {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❗------❎ 名称不能包含中文❎--------", Error: nil}
		return
	}

	var genCode = &model.CodeModel{
		Fields: make([]*model.GenCodeStruct, 0),
	}

	if !isTrue(empty) {
		genCode = handleArgs(infoChan)
		if len(genCode.Fields) == 0 {
			return
		}
	} else {
		if filename == "" {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗------❎ 名称不可为空❎--------", Error: nil}
			return
		}
		genCode.Table = pkg.CamelToSnake(filename)
	}

	// 对genCode进行单词长短排序😊
	sort.Slice(genCode.Fields, func(i, j int) bool {
		return len(genCode.Fields[i].Name) < len(genCode.Fields[j].Name)
	})

	// 首字母小写
	filename = pkg.InitialLetterToLower(filename)
	addApiCodePath("api.go", apiPath, filename, ".go")
	addApiCodePath("router.go", routerPath, filename, ".go")
	addApiCodePath("service.go", servicePath, filename, ".go")
	addApiCodePath("repo.go", repoPath, filename, ".go")
	addApiCodePath("model.go", modelPath, filename, ".go")
	addApiCodePath("gormGen.go", gormGenPath, filename, ".go")
	var slicePath = make([]string, 0)
	var ok = "n"
	var isOK = false // 判断是否有改动文件
	isZip := false   // 是否压缩啦
	for _, data := range apiCodePath {
		goFilePaht := fmt.Sprintf("%s/%s/%s%s", data.FilePath, data.FileName, data.FileName, data.FileExtension)
		// 文件是否存在
		_, err = os.Stat(goFilePaht)
		if err == nil && isWebDebug == false {
			// debug模式先删掉文件
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
					// 删除之前先把能涉及到的文件备份压缩
					if !isZip {
						err = zipModelCode()
						if err != nil {
							infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
							return
						}
						isZip = true
					}
					_ = os.Remove(goFilePaht)
				} else {
					log.Fatalln("🚶‍♀️告辞🚶")
				}
			} else {
				// 跳过 继续执行下一个文件
				continue
			}
		} else if err == nil && isWebDebug {
			// 删除之前先把能涉及到的文件备份压缩
			if !isZip {
				err = zipModelCode()
				if err != nil {
					infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
					return
				}
				isZip = true
			}
			_ = os.Remove(goFilePaht)
		}

		// 创建文件夹
		exists, err := pkg.PathExists(fmt.Sprintf("%s/%s", data.FilePath, data.FileName))
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️创建所需文件夹失败", Error: err}
			return
		}

		if !exists {
			err = os.MkdirAll(fmt.Sprintf("%s/%s", data.FilePath, data.FileName), 0750)
			if err != nil {
				infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️创建所需文件夹失败", Error: err}
				return
			}
		}

		// log.Println("🫎文件路径:", goFilePaht)
		slicePath = append(slicePath, goFilePaht)
		// 生成代码
		err = TmplExecute(packagename, goFilePaht, data, genCode)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️生成代码渲染失败", Error: err}
			// 如果失败一个 就删除已生成的文件
			for _, s := range slicePath {
				_ = os.Remove(s)
			}
			return
		}
		if data.FileExtension == ".go" {
			_ = FormatGoCode(goFilePaht)
		}
		isOK = true
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅" + goFilePaht + "文件创建成功👌！", Error: nil}
	}

	if !isOK {
		// 没有生成任何文件
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❗⚠️❎ 相关代码文件已存在，生成失败！", Error: nil}
		return
	}

	InsertInitRouterCode(packagename, "", filename)
	InsertGormGenCode(packagename, "", filename)
	InsertSetDB(packagename, "", filename)

	_ = pkg.RunCommandNoOutput("fuxi", "gen")
	_ = pkg.RunCommandNoOutput("fuxi", "openapi")

	// 运行mod tidy
	infoChan <- pkg.CommandInfo{Message: "🐮🐴 ✅ :添加成功正在 go mod tidy...", Error: nil}
	err = pkg.RunCommand("go", "mod", "tidy")
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️拉取依赖错误", Error: err}
		return
	}
	infoChan <- pkg.CommandInfo{Message: "🐮🐴 ✅ :执行 go mod tidy 成功👌", Error: nil}

	// 运行gormGen生成gen代码
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	modChan := make(chan pkg.CommandInfo)
	go pkg.RunCommandChannel(ctx, modChan, "go", "run", gormGenPath+"/main.go")
	for info := range modChan {
		if info.Error != nil {
			infoChan <- pkg.CommandInfo{Message: "❌ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 行gormGenBuild生成gen代码失败 ,请手动执行" + info.Message, Error: info.Error}
			cancel()
			return
		}
		if info.Message == "命令执行成功" {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗✅ :执行" + gormGenPath + "main.go 生成gen代码成功👌！", Error: nil}
		} else {
			infoChan <- pkg.CommandInfo{Message: info.Message, Error: info.Error}
		}
	}

	gushi := `代码精妙我自豪😎，
分享经验乐陶陶😄；
你我同欢心相照👫，
共筑梦想乐逍遥🌟。`
	infoChan <- pkg.CommandInfo{Message: gushi, Error: nil}
}

// TmplExecute 模板渲染
func TmplExecute(packageName, goFilePath string, data *PathData, genCode *model.CodeModel) error {
	// 读取模板
	var bytes []byte
	var err error
	if isTrue(empty) {
		bytes, err = emptytemplatex.EmptyTmplData.ReadFile(data.TmplPath + ".tmpl")
	} else {
		bytes, err = templatex.TmplData.ReadFile(data.TmplPath + ".tmpl")
	}
	if err != nil {
		return err
	}

	// 创建tmpl渲染对象
	tmpl, err := template.New("demo").Funcs(base.FuncMap).Parse(string(bytes))
	file, err := os.OpenFile(goFilePath, os.O_CREATE, 0750)
	if err != nil {
		return err
	}
	defer file.Close()

	// 渲染模板
	if err = tmpl.Execute(file, model.TmplData{
		ModuleName:  packageName,
		Table:       genCode.Table,
		StructName:  filename,
		GormGenPath: fmt.Sprintf("%s/%s/query", modelPath, filename),
		Data:        genCode.Fields,
		Import:      handleImport(genCode.Fields),
		Desc:        genCode.Desc,
	}); err != nil {
		return err
	}
	return nil
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

func zipModelCode() error {
	var goFilePaht []string
	for _, data := range apiCodePath {
		goFilePaht = append(goFilePaht, fmt.Sprintf("%s/%s/%s%s", data.FilePath, data.FileName, data.FileName, data.FileExtension))
	}
	moduleName, err := pkg.GetModuleName("go.mod")
	if err != nil {
		return err
	}

	zipFileName := fmt.Sprintf(".fuxi/%s/%s/%s.zip", moduleName, filename, time.Now().Format("20060102150405"))
	exists, err := pkg.PathExists(path.Dir(zipFileName))
	if err != nil {
		return err
	}

	if !exists {
		err = os.MkdirAll(path.Dir(zipFileName), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return pkg.FilesToZip(zipFileName, goFilePaht)
}
