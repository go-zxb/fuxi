package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/cmd/openapi"
	"github.com/go-zxb/fuxi/internal/cmd/whichFile"
	"github.com/go-zxb/fuxi/internal/project"
	"github.com/go-zxb/fuxi/pkg"
	"github.com/spf13/cobra"
)

var (
	excludePaths []string
	output       string
	osType       string
)

func init() {
	runCmd.PersistentFlags().StringSliceVarP(&excludePaths, "exclude", "e", []string{}, "需要排除监听的go文件目录")
	buildGoCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "输出目录")
	buildGoCmd.PersistentFlags().StringVarP(&osType, "osType", "t", "win", "编译平台")
}

var rootCmd = &cobra.Command{
	Use:     "fuxi",
	Short:   "fuxi",
	Long:    "fuxi",
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
🫣Hi~欢迎使用 FuXi Gin Web开发辅助工具! 
👻祝你编码无BUG🐛,流畅上线稳定运行😶‍🌫
🫨微信交流:GinWeb
`)
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成gorm gen代码(生成api CRUD代码后才能执行)",
	Long:  "生成gorm gen代码(生成api CRUD代码后才能执行)",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.RunCommand("go", "run", "cmd/gorm/gen/main.go")
		if err != nil {
			log.Panicln(err)
			return
		}
		log.Println("✅ gen success 👌")
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "运行当前目录下的go项目并监听go文件变化实现自动重新运行",
	Long:  "运行当前目录下的go项目并监听go文件变化实现自动重新运行",
	Run: func(cmd *cobra.Command, args []string) {
		// 运行并监控文件变化
		path, _ := os.Getwd()
		whichFile.Which(path, excludePaths, false)
	},
}

var buildGoCmd = &cobra.Command{
	Use:   "build:go",
	Short: "运行go build -ldflags=\"-s -w\"",
	Long:  "运行go build -ldflags=\"-s -w\"",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.RunCommand("go", "build", `-ldflags=-s -w`, `-o`, output)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("✅ run success 👌")
	},
}

var buildFlutterCmd = &cobra.Command{
	Use:   "build:flt",
	Short: "运行flutter build windows ",
	Long:  "运行flutter build windows ",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.RunCommand("flutter", "build", `windows`)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("✅ run success 👌")
	},
}

func init() {
	path, _ := os.Getwd()
	log.Println("💻 当前运行目录:", path)
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "build:go", "build:flt", "gen", "project":
		default:
			// 读取配置文件
			_, err := config.NewConfig("")
			if err != nil {
				log.Fatalln("📒 配置文件读取错误", err)
			}
		}
	}

	rootCmd.AddCommand(project.AddApiCmd)
	rootCmd.AddCommand(project.NewProjectCmd)
	rootCmd.AddCommand(project.NewApiCmd)
	rootCmd.AddCommand(project.SSEWebCmd)
	rootCmd.AddCommand(project.NewSysUserCmd)
	rootCmd.AddCommand(project.NewCasbinCmd)
	rootCmd.AddCommand(project.NewApisCmd)
	rootCmd.AddCommand(openapi.OpenapiCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(buildGoCmd)
	rootCmd.AddCommand(buildFlutterCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
