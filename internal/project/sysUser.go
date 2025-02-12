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

	"github.com/go-zxb/fuxi/pkg"
	tmpl "github.com/go-zxb/fuxi/template/user"
	"github.com/spf13/cobra"
)

type SysUser struct {
	ModuleName  string
	StructName  string
	isSysUser   string
	FilePath    string
	FileName    string
	GenFilePath string
	debug       string
}

var sysUser = &SysUser{
	ModuleName:  "",
	StructName:  "User",
	FileName:    "user",
	FilePath:    "user",
	GenFilePath: "user",
}

func init() {
	NewSysUserCmd.Flags().StringVarP(&sysUser.debug, "debug", "d", "false", "是否开启debug模式")
	NewSysUserCmd.Flags().StringVarP(&sysUser.StructName, "struct", "s", sysUser.StructName, "结构体名称")
	NewSysUserCmd.Flags().StringVarP(&sysUser.isSysUser, "isSysUser", "i", sysUser.isSysUser, "是否是系统用户模块")
	NewSysUserCmd.Flags().StringVarP(&sysUser.FileName, "filename", "f", sysUser.FileName, "文件名称")

}

var NewSysUserCmd = &cobra.Command{
	Use:   "create:user",
	Short: "自动创建一个用户注册登录模块",
	Long:  "自动创建一个用户注册登录模块",
	Run:   cmdHandleSysUser,
}

var sysUserCodePath []*PathData

func cmdHandleSysUser(cmd *cobra.Command, args []string) {

	infoChan := make(chan pkg.CommandInfo)
	go handleGenSysUserCode(infoChan)
	for info := range infoChan {
		if info.Error != nil {
			log.Fatalln("❌", info.Message, info.Error.Error())
		} else {
			log.Println(info.Message)
		}
	}
}

func handleGenSysUserCode(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	if sysUser.isSysUser == "true" {
		sysUser.FilePath = "system/user"
		sysUser.GenFilePath = "system/user"
	} else {
		sysUser.FilePath = sysUser.FileName
		sysUser.GenFilePath = sysUser.FileName
	}

	addSysUserCodePath("api.go", "internal/api/"+sysUser.FilePath, sysUser.FileName, ".go")
	addSysUserCodePath("repo.go", "internal/repo/"+sysUser.FilePath, sysUser.FileName, ".go")
	addSysUserCodePath("service.go", "internal/service/"+sysUser.FilePath, sysUser.FileName, ".go")
	addSysUserCodePath("router.go", "internal/router/"+sysUser.FilePath, sysUser.FileName, ".go")
	addSysUserCodePath("model.go", "internal/model/"+sysUser.FilePath, sysUser.FileName, ".go")
	addSysUserCodePath("gen.go", "cmd/gorm/gen/"+sysUser.FilePath, sysUser.FileName, ".go")

	infoChan <- pkg.CommandInfo{Message: "🐮🐴正在复制go基础文件....", Error: nil}
	ok := "n"
	isOk := false
	isZip := false
	for _, data := range sysUserCodePath {

		goFilePaht := fmt.Sprintf("%s/%s%s", data.FilePath, data.FileName, data.FileExtension)

		// 文件是否存在
		_, err := os.Stat(goFilePaht)
		if err == nil && isWebDebug == false {
			// debug模式先删掉文件
			if isTrue(sysUser.debug) {
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
						err = zipSysUserCode()
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
				err = zipSysUserCode()
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

		bytes, err := tmpl.TmplSysUserData.ReadFile(data.TmplPath + ".tmpl")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴读取模板代码失败....", Error: err}
		}
		tmplx := template.Must(template.New("demo").Parse(string(bytes)))
		file, err := os.OpenFile(goFilePaht, os.O_CREATE, 0750)
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴创建源码文件失败....", Error: err}
		}
		defer file.Close()

		sysUser.ModuleName, err = pkg.GetModuleName("go.mod")
		if err != nil {
			infoChan <- pkg.CommandInfo{Message: "🐮🐴❗️请先初始化项目: fuxi project -n 项目名称", Error: err}
			return
		}
		if err = tmplx.Execute(file, sysUser); err != nil {
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

	if isTrue(sysUser.isSysUser) {
		InsertInitRouterCode(sysUser.ModuleName, "system/", "user")
		InsertGormGenCode(sysUser.ModuleName, "system/", "user")
		InsertSetDB(sysUser.ModuleName, "system/", "user")
	} else {
		InsertInitRouterCode(sysUser.ModuleName, "", sysUser.FilePath)
		InsertGormGenCode(sysUser.ModuleName, "", sysUser.FilePath)
		InsertSetDB(sysUser.ModuleName, "", sysUser.FilePath)
	}

	err := pkg.RunCommand("go", "get", "gorm.io/gen")
	if err != nil {
		fmt.Println(err)
	}
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

	infoChan <- pkg.CommandInfo{Message: fmt.Sprintf("✅ 创建 %s 系统用户模块成功", filepath.Base(sysUser.ModuleName)), Error: nil}
}

func addSysUserCodePath(tmplPath, filepath, filename string, suffix string) {
	sysUserCodePath = append(sysUserCodePath, &PathData{
		TmplPath:      tmplPath,
		FilePath:      filepath,
		FileName:      filename,
		FileExtension: suffix,
	})
}

func zipSysUserCode() error {
	var goFilePaht []string
	for _, data := range sysUserCodePath {
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
