package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-zxb/fuxi/pkg"
	tmpl "github.com/go-zxb/fuxi/template/pkg"
	"github.com/spf13/cobra"
)

var (
	projectName = ""
	debug       = ""
	isWebDebug  = false
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
	Run:   cmdHandle,
}

var projectCodePath []*PathData

func cmdHandle(cmd *cobra.Command, args []string) {
	if pkg.HasChinese(projectName) {
		log.Fatalln("------名称不能包含中文--------")
	}
	infoChan := make(chan pkg.CommandInfo)
	go handleGenProjectCode(infoChan)
	for info := range infoChan {
		if info.Error != nil {
			log.Fatalln("❌", info.Message, info.Error.Error())
		} else {
			log.Println(info.Message)
		}
	}
}

func handleGenProjectCode(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	addProjectCodePath("mod", "", "go", ".mod")
	addProjectCodePath("main.go", "", "main", ".go")
	addProjectCodePath("core.go", "core", "core", ".go")
	addProjectCodePath("dockerfile", "", "dockerfile", "")
	addProjectCodePath(".gitignore", "", ".gitignore", "")
	addProjectCodePath("aes.go", "pkg/encrypt", "aes", ".go")
	addProjectCodePath("md5.go", "pkg/encrypt", "md5", ".go")
	addProjectCodePath("pwd.go", "pkg/encrypt", "pwd", ".go")
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
	addProjectCodePath("docs.go", "docs/openapi", "docs", ".go")
	addProjectCodePath("consts.go", "utils", "consts", ".go")
	addProjectCodePath("jwt.go", "pkg/jwt", "jwt", ".go")
	addProjectCodePath("middlewareJwt.go", "middleware", "jwt", ".go")
	addProjectCodePath("modelBase.go", "internal/model/base", "base", ".go")
	addProjectCodePath("log.go", "log", "log", ".go")
	addProjectCodePath("README.md", "", "README", ".md")

	infoChan <- pkg.CommandInfo{Message: "🐮🐴正在复制go基础文件....", Error: nil}
	ok := "n"
	isOk := false
	isZip := false
	for _, data := range projectCodePath {

		goFilePaht := ""
		if data.FilePath != "" {
			goFilePaht = fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension)
		} else {
			goFilePaht = fmt.Sprintf("%s%s", data.FileName, data.FileExtension)
		}

		// 文件是否存在
		_, err := os.Stat(goFilePaht)
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
						err = zipProjectCode()
						if err != nil {
							log.Fatalln("🚶‍♀️告辞🚶 备份失败🎒")
						}
						isZip = true
					}
					_ = os.Remove(goFilePaht)
				} else {
					log.Fatalln("🚶‍♀️告辞🚶")
				}
			} else {
				continue
			}
		} else if err == nil && isWebDebug {
			// 删除之前先把能涉及到的文件备份压缩
			if !isZip {
				err = zipProjectCode()
				if err != nil {
					infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
					return
				}
				isZip = true
			}
			_ = os.Remove(goFilePaht)
		}

		if data.FilePath != "" {
			// 创建文件夹
			exists, err := pkg.PathExists(data.FilePath)
			if err != nil {
				infoChan <- pkg.CommandInfo{Message: "🐮🐴创建目录失败....", Error: err}
				return
			}

			if !exists {
				err = os.MkdirAll(data.FilePath, os.ModePerm)
				if err != nil {
					infoChan <- pkg.CommandInfo{Message: "🐮🐴创建目录失败....", Error: err}
					return
				}
			}

		}

		// log.Println(data.FileName+data.FileExtension, "文件路径:", goFilePaht)
		bytes, err := tmpl.GoCode.ReadFile(data.TmplPath + ".tmpl")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴读取模板代码失败....", Error: err}
		}
		tmplx := template.Must(template.New("demo").Parse(string(bytes)))
		file, err := os.OpenFile(goFilePaht, os.O_CREATE, 0750)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴创建源码文件失败....", Error: err}
		}
		defer file.Close()
		if err = tmplx.Execute(file, Project{
			ProjectName: projectName,
			PackageName: projectName,
			ModuleName:  projectName,
			Version:     gin.Version,
			Description: "代码生成工具",
		}); err != nil {
			infoChan <- pkg.CommandInfo{Message: "⚠️❎🐮🐴代码渲染失败....", Error: err}
		}
		if data.FileExtension == ".go" {
			_ = FormatGoCode(goFilePaht)
		}
		isOk = true
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  %s 文件创建成功👌！", goFilePaht), Error: nil}
	}

	if !isOk {
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  项目还是哪个项目, 未做任何改变 😊 ~略略略略~"), Error: nil}
	}

	infoChan <- pkg.CommandInfo{Message: "🎁٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 正在拉取依赖包...", Error: nil}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	modChan := make(chan pkg.CommandInfo)
	go pkg.RunCommandChannel(ctx, modChan, "go", "mod", "tidy")
	for info := range modChan {
		if info.Error != nil {
			infoChan <- pkg.CommandInfo{Message: "❌ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 依赖文件拉取失败... "}
			cancel()
			return
		}
		if info.Message == "命令执行成功" {
			infoChan <- pkg.CommandInfo{Message: "🎁٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 拉取依赖包成功...", Error: nil}
		} else {
			infoChan <- pkg.CommandInfo{Message: info.Message, Error: info.Error}
		}
	}
	infoChan <- pkg.CommandInfo{Message: "🎁٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 正在创建相关模块文件...", Error: nil}

	err := pkg.RunCommand("fuxi", "create:casbin")
	infoChan <- pkg.CommandInfo{Message: "✅ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 生成casbin代码成功...", Error: err}

	err = pkg.RunCommand("fuxi", "create:user", "-i", "true", "-s", "User")
	infoChan <- pkg.CommandInfo{Message: "✅ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 生成user代码成功...", Error: err}

	err = pkg.RunCommand("fuxi", "create:apis")
	infoChan <- pkg.CommandInfo{Message: "✅ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 生成apis代码成功...", Error: err}

	err = pkg.RunCommand("fuxi", "create:role")
	infoChan <- pkg.CommandInfo{Message: "✅ ٩(•̤̀ᵕ•̤́๑)ᵒᵏᵎᵎᵎᵎ 生成role代码成功...", Error: err}

	_ = pkg.RunCommandNoOutput("fuxi", "openapi")
	infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅ 创建 %s 项目成功", projectName), Error: nil}
	si := `程序流畅心自喜😊,
助人为乐情更怡🤗;
共享成功欢声起🎊,
同心协力福缘齐🤝.`
	infoChan <- pkg.CommandInfo{Message: si, Error: nil}
}

func addProjectCodePath(tmplPath, filepath, filename string, suffix string) {
	projectCodePath = append(projectCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}

func zipProjectCode() error {
	var goFilePaht []string
	for _, data := range projectCodePath {
		if data.FilePath != "" {
			goFilePaht = append(goFilePaht, fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension))
		} else {
			goFilePaht = append(goFilePaht, fmt.Sprintf("%s%s", data.FileName, data.FileExtension))
		}
	}
	zipFileName := fmt.Sprintf(".fuxi/%s/%s.zip", projectName, time.Now().Format("20060102150405"))

	exists, _ := pkg.PathExists(path.Dir(zipFileName))

	if !exists {
		err := os.MkdirAll(path.Dir(zipFileName), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return pkg.FilesToZip(zipFileName, goFilePaht)
}
