package main

import (
	"fmt"
	"github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/internal/cmd/openapi"
	"github.com/go-zxb/fuxi/internal/project"
	"github.com/go-zxb/fuxi/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

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
	Short: "运行当前目录下的go项目",
	Long:  "运行当前目录下的go项目",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.RunCommand("go", "run", "main.go", "server")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("✅ run success 👌")
	},
}

func init() {
	path, _ := os.Getwd()
	log.Println("💻 当前运行目录:", path)
	// 读取配置文件
	_, err := config.NewConfig("")
	if err != nil {
		log.Fatalln("📒 配置文件读取错误", err)
	}
	//fmt.Println(v)

	rootCmd.AddCommand(project.AddApiCmd)
	rootCmd.AddCommand(project.NewProjectCmd)
	rootCmd.AddCommand(project.NewApiCmd)
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
