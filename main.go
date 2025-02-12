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
)

func init() {
	runCmd.PersistentFlags().StringSliceVarP(&excludePaths, "exclude", "e", []string{}, "需要排除监听的go文件目录")
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
	rootCmd.AddCommand(project.NewRoleCmd)
	rootCmd.AddCommand(openapi.OpenapiCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(genCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
