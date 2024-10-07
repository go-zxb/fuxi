package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	newapi "github.com/go-zxb/fuxi/internal/ast/new"
	"github.com/go-zxb/fuxi/internal/model"
	"github.com/go-zxb/fuxi/pkg"
	casbinTmpl "github.com/go-zxb/fuxi/template/casbin"
	"github.com/spf13/cobra"
)

type Casbin struct {
	ModuleName  string
	StructName  string
	FilePath    string
	FileName    string
	GenFilePath string
	debug       string
}

var casbin = &Casbin{}

func init() {
	NewCasbinCmd.Flags().StringVarP(&casbin.debug, "debug", "d", "false", "是否开启debug模式")
}

var NewCasbinCmd = &cobra.Command{
	Use:   "create:casbin",
	Short: "自动创建Casbin模块",
	Long:  "自动创建Casbin模块",
	Run:   cmdHandleCasbin,
}

var casbinCodePath []*PathData

func cmdHandleCasbin(cmd *cobra.Command, args []string) {

	infoChan := make(chan pkg.CommandInfo)
	go handleGenCasbinCode(infoChan)
	for info := range infoChan {
		if info.Error != nil {
			log.Fatalln("❌", info.Message, info.Error.Error())
		} else {
			log.Println(info.Message)
		}
	}
}

func handleGenCasbinCode(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	casbin.FilePath = "system/casbin"
	casbin.GenFilePath = "system/casbin"
	casbin.FileName = "casbin"

	addCasbinCodePath("api.go", "internal/api/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("repo.go", "internal/repo/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("service.go", "internal/service/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("router.go", "internal/router/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("model.go", "internal/model/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("gen.go", "cmd/gorm/gen/"+casbin.FilePath, casbin.FileName, ".go")
	addCasbinCodePath("middleware.go", "middleware", casbin.FileName, ".go")

	infoChan <- pkg.CommandInfo{Message: "🐮🐴正在复制go基础文件....", Error: nil}
	ok := "n"
	isOk := false
	isZip := false
	for _, data := range casbinCodePath {

		goFilePaht := fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension)

		// 文件是否存在
		_, err := os.Stat(goFilePaht)
		if err == nil && isWebDebug == false {
			// debug模式先删掉文件
			if isTrue(casbin.debug) {
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
						err = zipCasbinCode()
						if err != nil {
							if strings.Contains(err.Error(), "The system cannot find the file specified") {
								log.Println("⚠️ 老家伙,文件不存在, 继续吧....")
							} else {
								log.Fatalln("🚶‍♀️告辞🚶 备份失败🎒", err)
							}
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
				err = zipCasbinCode()
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

		bytes, err := casbinTmpl.TmplCasbinData.ReadFile(data.TmplPath + ".tmpl")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴读取模板代码失败....", Error: err}
		}
		tmplx := template.Must(template.New("demo").Parse(string(bytes)))
		file, err := os.OpenFile(goFilePaht, os.O_CREATE, 0750)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴创建源码文件失败....", Error: err}
		}
		defer file.Close()

		casbin.ModuleName, err = pkg.GetModuleName("go.mod")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️请先初始化项目: fuxi project -n 项目名称", Error: err}
			return
		}
		if err = tmplx.Execute(file, casbin); err != nil {
			infoChan <- pkg.CommandInfo{Message: "⚠️❎🐮🐴代码渲染失败....", Error: err}
		}
		if data.FileExtension == ".go" {
			_ = FormatGoCode(goFilePaht)
		}
		isOk = true
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  %s 文件创建成功👌！", goFilePaht), Error: nil}
	}

	if !isOk {
		infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅  模块还是哪个模块, 未做任何改变 😊 ~略略略略~"), Error: nil}
		return
	}

	InsertInitRouterCode(casbin.ModuleName, "system/", "casbin")
	InsertGormGenCode(casbin.ModuleName, "system/", "casbin")
	InsertSetDB(casbin.ModuleName, "system/", "casbin")

	e := newapi.Enforcer{
		AST: model.AST{
			Name:        "casbin",
			ImportsName: []string{"casbinService"},
			Imports:     []string{fmt.Sprintf("%s/internal/service/system/%s", casbin.ModuleName, "casbin")},
			FilePath:    "core/core.go",
		},
	}
	e.Enforcer()

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

	infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅ 创建 %s 系统用户模块成功", filepath.Base(casbin.ModuleName)), Error: nil}
}

func addCasbinCodePath(tmplPath, filepath, filename string, suffix string) {
	casbinCodePath = append(casbinCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}

func zipCasbinCode() error {
	var goFilePaht []string
	for _, data := range casbinCodePath {
		if data.FilePath != "" {
			goFilePaht = append(goFilePaht, fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension))
		} else {
			goFilePaht = append(goFilePaht, fmt.Sprintf("%s%s", data.FileName, data.FileExtension))
		}
	}
	zipFileName := fmt.Sprintf(".fuxi/%s/%s.zip", projectName, time.Now().Format("20060102150405"))

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
