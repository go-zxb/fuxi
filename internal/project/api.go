package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-zxb/fuxi/pkg"
	apisTmpl "github.com/go-zxb/fuxi/template/api"
	"github.com/spf13/cobra"
)

type Apis struct {
	ModuleName  string
	StructName  string
	FilePath    string
	FileName    string
	GenFilePath string
	debug       string
}

var apis = &Apis{}

func init() {
	NewApisCmd.Flags().StringVarP(&apis.debug, "debug", "d", "false", "是否开启debug模式")
}

var NewApisCmd = &cobra.Command{
	Use:   "create:apis",
	Short: "自动创建Apis模块",
	Long:  "自动创建Apis模块",
	Run:   cmdHandleApis,
}

var apisCodePath []*PathData

func cmdHandleApis(cmd *cobra.Command, args []string) {

	infoChan := make(chan pkg.CommandInfo)
	go handleGenApisCode(infoChan)
	for info := range infoChan {
		if info.Error != nil {
			log.Fatalln("❌", info.Message, info.Error.Error())
		} else {
			log.Println(info.Message)
		}
	}
}

func handleGenApisCode(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	apis.FilePath = "system/api"
	apis.GenFilePath = "system/api"
	apis.FileName = "api"

	addApisCodePath("api.go", "internal/api/"+apis.FilePath, apis.FileName, ".go")
	addApisCodePath("repo.go", "internal/repo/"+apis.FilePath, apis.FileName, ".go")
	addApisCodePath("service.go", "internal/service/"+apis.FilePath, apis.FileName, ".go")
	addApisCodePath("router.go", "internal/router/"+apis.FilePath, apis.FileName, ".go")
	addApisCodePath("model.go", "internal/model/"+apis.FilePath, apis.FileName, ".go")
	addApisCodePath("gen.go", "cmd/gorm/gen/"+apis.FilePath, apis.FileName, ".go")

	infoChan <- pkg.CommandInfo{Message: "🐮🐴正在复制go基础文件....", Error: nil}
	isOk := false
	for _, data := range apisCodePath {

		goFilePaht := fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension)

		// 文件是否存在
		_, err := os.Stat(goFilePaht)
		if err == nil && apis.debug == "true" {
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

		bytes, err := apisTmpl.TmplApisData.ReadFile(data.TmplPath + ".tmpl")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴读取模板代码失败....", Error: err}
		}
		tmplx := template.Must(template.New("demo").Parse(string(bytes)))
		file, err := os.OpenFile(goFilePaht, os.O_CREATE, 0750)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴创建源码文件失败....", Error: err}
		}
		defer file.Close()

		apis.ModuleName, err = pkg.GetModuleName("go.mod")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️请先初始化项目: fuxi project -n 项目名称", Error: err}
			return
		}
		if err = tmplx.Execute(file, apis); err != nil {
			infoChan <- pkg.CommandInfo{Message: "⚠️❎🐮🐴代码渲染失败....", Error: err}
		}
		if data.FileExtension == ".go" {
			_ = pkg.FormatGoCode(goFilePaht)
		}
		isOk = true
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  %s 文件创建成功👌！", goFilePaht), Error: nil}
	}

	if !isOk {
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  模块还是哪个模块, 未做任何改变 😊 ~略略略略~"), Error: nil}
		return
	}

	InsertInitRouterCode(apis.ModuleName, "system/", "api")
	InsertGormGenCode(apis.ModuleName, "system/", "api")
	InsertSetDB(apis.ModuleName, "system/", "api")

	_ = pkg.RunCommandNoOutput("fuxi", "gen")
	_ = pkg.RunCommandNoOutput("fuxi", "openapi")

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

	infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅ 创建 %s 系统用户模块成功", filepath.Base(apis.ModuleName)), Error: nil}
}

func addApisCodePath(tmplPath, filepath, filename string, suffix string) {
	apisCodePath = append(apisCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}
