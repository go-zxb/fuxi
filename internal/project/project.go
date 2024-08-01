package project

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-zxb/fuxi/pkg"
	tmpl "github.com/go-zxb/fuxi/template/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/template"
)

var (
	projectName = ""
	debug       = ""
)

type Project struct {
	ProjectName string
	PackageName string
	Description string
	ModuleName  string
	Version     string
}

func init() {
	NewProjectCmd.Flags().StringVarP(&projectName, "name", "n", "great", "项目名 如:fuxi")
	NewProjectCmd.Flags().StringVarP(&debug, "debug", "d", "false", "是否开启debug模式")
}

var NewProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "创建一个go项目工程",
	Long:  "创建一个go项目工程",
	Run:   handleGenPkgCode,
}

var projectCodePath []*PathData

func handleGenPkgCode(cmd *cobra.Command, args []string) {
	if pkg.HasChinese(projectName) {
		log.Fatalln("------名称不能包含中文--------")
	}

	addProjectCodePath("mod", "", "go", ".mod")
	addProjectCodePath("main.go", "", "main", ".go")
	addProjectCodePath("core.go", "core", "core", ".go")
	addProjectCodePath("dockerfile", "", "dockerfile", "")
	addProjectCodePath(".gitignore", "", ".gitignore", "")
	addProjectCodePath("aes.go", "pkg/encrypt", "aes", ".go")
	addProjectCodePath("config.go", "config", "config", ".go")
	addProjectCodePath("mysql.go", "core/data", "mysql", ".go")
	addProjectCodePath("config.yaml", "config", "config", ".yaml")
	addProjectCodePath("strconv.go", "pkg/strconv", "strconv", ".go")
	addProjectCodePath("main.gorm.go", "cmd/gorm/gen", "main", ".go")
	addProjectCodePath("httpRequest.go", "pkg/httputil", "request", ".go")
	addProjectCodePath("response.go", "pkg/response", "response", ".go")
	addProjectCodePath("json_arm.go", "pkg/json", "json_arm", ".go")
	addProjectCodePath("json_default.go", "pkg/json", "json_default", ".go")
	addProjectCodePath("redis.go", "core/data", "redis", ".go")
	addProjectCodePath("cors.go", "middleware", "cors", ".go")

	log.Println("正在复制go基础文件....")
	ok := "n"
	isOk := false
	for _, data := range projectCodePath {

		goFilePaht := ""
		if data.FilePath != "" {
			goFilePaht = fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension)
		} else {
			goFilePaht = fmt.Sprintf("%s%s", data.FileName, data.FileExtension)
		}

		//文件是否存在
		_, err := os.Stat(goFilePaht)
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
				continue
			}
		}

		if data.FilePath != "" {
			err = os.MkdirAll(data.FilePath, os.ModePerm)
			if err != nil {
				log.Println(data.FilePath, err)
			}
		}

		//log.Println(data.FileName+data.FileExtension, "文件路径:", goFilePaht)
		bytes, err := tmpl.GoCode.ReadFile(data.TmplPath + ".tmpl")
		if err != nil {
			log.Fatalln(err)
		}
		tmplx := template.Must(template.New("demo").Parse(string(bytes)))
		file, err := os.OpenFile(goFilePaht, os.O_CREATE, 0750)
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		if err = tmplx.Execute(file, Project{
			ProjectName: projectName,
			PackageName: projectName,
			ModuleName:  projectName,
			Version:     gin.Version,
			Description: "代码生成工具",
		}); err != nil {
			log.Println("⚠️❎ 代码渲染失败", err)
		}
		if data.FileExtension == ".go" {
			_ = FormatGoCode(goFilePaht)
		}
		isOk = true
		log.Println("✅ ", goFilePaht, "文件创建成功👌！")
	}

	if !isOk {
		log.Fatalln("❎ 项目还是哪个项目, 未做任何改变 😊 ~略略略略~")
	}

	log.Println("🎁٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 正在拉取依赖包...")
	err := pkg.RunCommandNoOutput("go", "mod", "tidy")
	if err != nil {
		log.Fatalln("❎ 依赖文件拉取失败:", err)
	}
	log.Println("✅ 拉取依赖包成功")
	log.Println("✅ 创建", projectName, "项目成功")
	fmt.Println(`程序流畅心自喜😊，
助人为乐情更怡🤗；
共享成功欢声起🎊，
同心协力福缘齐🤝。`)
}

func addProjectCodePath(tmplPath, filepath, filename string, suffix string) {
	projectCodePath = append(projectCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}
