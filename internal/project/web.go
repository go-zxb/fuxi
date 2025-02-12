package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	ginstatic "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-zxb/fuxi/internal/cmd/whichFile"
	"github.com/go-zxb/fuxi/middleware"
	"github.com/go-zxb/fuxi/pkg"
	"github.com/go-zxb/fuxi/static"
	"github.com/spf13/cobra"
)

var args = &Args{}

var SSEWebCmd = &cobra.Command{
	Use:   "ui",
	Short: "启动一个web服务 可在网页端创建项目",
	Long:  "启动一个web服务 可在网页端创建项目",
	Run:   SSEWeb,
}

type Info struct {
	Content string `json:"content"`
}

func SSEWeb(cmd *cobra.Command, args []string) {
	r := gin.Default()
	r.Use(middleware.Cors())

	// 处理 POST 请求以初始化 SSE 连接
	r.POST("/start-sse", func(c *gin.Context) {
		// 你可以在这里处理 POST 请求的 body
		var requestBody map[string]interface{}
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "SSE connection initialized"})
	})

	// 处理 SSE 请求
	r.GET("/createProject", createProject)
	r.GET("/addApi", addApi)
	r.GET("/newApi", newApi)
	r.GET("/getModel", getModel)
	r.GET("/getModName", getModName)
	r.NoRoute(func(c *gin.Context) {
		c.File(os.TempDir() + "/fuxi/index.html")
	})

	localDir := os.TempDir()
	localDir += "/fuxi"
	defer os.RemoveAll(localDir)
	static.WriteStaticFiles(localDir)
	r.Use(ginstatic.ServeRoot("/", os.TempDir()+"/fuxi"))
	r.Use(ginstatic.ServeRoot("/assets", os.TempDir()+"/fuxi/assets"))
	log.Println("✅", " web 服务启动成功,运行地址：http://127.0.0.1:8066", "👌")

	// 启动 GinWeb 服务
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "", 8066),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fuxi:启动WebUI服务失败: %s\n", err)
		}
	}()

	go func() {
		// 监控文件变化
		path, _ := os.Getwd()
		whichFile.Which(path, []string{}, false)
	}()

	// 捕获 SIGINT 信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭web服务中...")

	// 优雅关闭服务
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("web服务已被迫关闭:", err)
	}

	log.Println("web服务已退出")
}

type Args struct {
	Api          string `form:"api"  json:"api,omitempty"`
	ApiFunc      string `form:"apiFunc" json:"apiFunc,omitempty"`
	ApiPath      string `form:"apiPath" json:"apiPath,omitempty"`
	Debug        bool   `form:"debug" json:"debug,omitempty"`
	IsByID       bool   `form:"isByID" json:"isByID,omitempty"`
	IsReturnList bool   `form:"isReturnList" json:"isReturnList,omitempty"`
	Method       string `form:"method" json:"method,omitempty"`
	FileName     string `form:"filename" json:"fileName,omitempty"`
	NoParams     bool   `form:"noParams" json:"noParams,omitempty"`
	RepoPath     string `form:"repoPath" json:"repoPath,omitempty"`
	ReturnType   string `form:"returnType" json:"returnType,omitempty"`
	RouterPath   string `form:"routerPath" json:"routerPath,omitempty"`
	ServicePath  string `form:"servicePath" json:"servicePath,omitempty"`
	AddRepo      bool   `form:"addRepo" json:"addRepo"`
}

func getModName(ctx *gin.Context) {
	mod, err := pkg.GetModuleName("go.mod")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": mod})
}

func getModel(ctx *gin.Context) {
	dir := "internal/model"
	var dirs = ""
	entries, err := os.ReadDir(dir)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"data": ""})
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs += entry.Name() + ","
			arr, _ := os.ReadDir(dir + "/" + entry.Name())
			for _, item := range arr {
				if item.IsDir() {
					dirs += entry.Name() + "/" + item.Name() + ","
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": dirs[:len(dirs)-1]})
}

func addApi(ctx *gin.Context) {
	setHeader(ctx)

	if err := ctx.ShouldBindQuery(args); err != nil {
		handleinfo(ctx, pkg.CommandInfo{Message: "Invalid request body"})
		return
	}

	if args.FileName == "" {
		handleinfo(ctx, pkg.CommandInfo{Message: "已有项目名称不能为空"})
		return
	}

	if args.ApiFunc == "" {
		handleinfo(ctx, pkg.CommandInfo{Message: "方法签名不能为空"})
		return
	}

	name = path.Base(args.FileName)
	apiFunc = args.ApiFunc
	method = args.Method
	AddRepo = args.AddRepo
	returnType = args.ReturnType

	if args.IsReturnList {
		isReturnList = "true"
	} else {
		isReturnList = "false"
	}

	if args.IsByID {
		iSByID = "true"
		isReturnList = "false"
	} else {
		iSByID = "false"
	}

	if args.NoParams {
		noParams = "true"
	} else {
		noParams = "false"
	}

	if args.Method == "" {
		method = "GET"
	}

	if args.Api != "" {
		api = args.Api
	}
	if args.ApiPath != "" {
		apiPath = args.ApiPath
	}
	if args.RouterPath != "" {
		routerPath = args.RouterPath
	}
	if args.ServicePath != "" {
		servicePath = args.ServicePath
	}
	if args.RepoPath != "" {
		repoPath = args.RepoPath
	}

	infoChan := make(chan pkg.CommandInfo)
	// 创建一个通道用于接收命令执行信息
	go addApiHandle(infoChan)
	// 从通道接收命令执行信息
	for info := range infoChan {
		handleinfo(ctx, info)
		time.Sleep(time.Millisecond * 500)
	}
	ctx.Writer.Write([]byte("data: ok\n\n"))
	Init()
}

func newApi(ctx *gin.Context) {

	setHeader(ctx)
	if val := ctx.Query("filename"); val == "" {
		handleinfo(ctx, pkg.CommandInfo{Message: "项目名称不能为空"})
		return
	} else {
		filename = val
	}

	empty = ctx.Query("empty")
	localJson = ctx.Query("localJson")
	isOutputJson = ctx.Query("isOutputJson")

	if !isTrue(empty) {
		if val := ctx.Query("question"); val == "" {
			if !isTrue(localJson) {
				handleinfo(ctx, pkg.CommandInfo{Message: "需求不能为空"})
				return
			}
		} else {
			question = val
		}
	}

	if isTrue(ctx.Query("debug")) {
		isWebDebug = true
	}

	if val := ctx.Query("apiPath"); val != "" {
		apiPath = val
	}
	if val := ctx.Query("routerPath"); val != "" {
		routerPath = val
	}
	if val := ctx.Query("servicePath"); val != "" {
		servicePath = val
	}
	if val := ctx.Query("repoPath"); val != "" {
		repoPath = val
	}
	if val := ctx.Query("modelPath"); val != "" {
		modelPath = val
	}

	infoChan := make(chan pkg.CommandInfo)
	// 创建一个通道用于接收命令执行信息
	go createCodeHandle(infoChan)
	// 从通道接收命令执行信息
	for info := range infoChan {
		handleinfo(ctx, info)
		time.Sleep(time.Millisecond * 100)
	}
	ctx.Writer.Write([]byte("data: ok\n\n"))

	Init()
}

func handleinfo(ctx *gin.Context, info pkg.CommandInfo) {
	fmt.Println("本地显示:", info.Message)
	resp := &Info{}
	if info.Error != nil {
		resp.Content = info.Message + " " + info.Error.Error()
		jsonData, _ := json.Marshal(resp)
		_, err := ctx.Writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
		if err != nil {
			log.Println("Writer ❌", err.Error())
			return
		}
		ctx.Writer.Flush()
	} else {
		resp.Content = info.Message
		jsonData, _ := json.Marshal(resp)
		_, err := ctx.Writer.Write([]byte("data: " + string(jsonData) + "\n\n"))
		if err != nil {
			log.Println("Writer ❌", err.Error())
			return
		}
		ctx.Writer.Flush()
	}

}

func createProject(ctx *gin.Context) {
	setHeader(ctx)

	name_ := ctx.Query("name")
	debug_ := ctx.Query("debug")

	if name_ == "" {
		handleinfo(ctx, pkg.CommandInfo{Message: "项目名称不能为空"})
		return
	}
	log.Println("✅", name_, "👌")
	projectName = name_
	if isTrue(debug_) {
		isWebDebug = true
	}

	log.Println("✅", name, "👌")
	infoChan := make(chan pkg.CommandInfo)
	// 创建一个通道用于接收命令执行信息
	go handleGenProjectCode(infoChan)
	// 从通道接收命令执行信息
	for info := range infoChan {
		handleinfo(ctx, info)
		time.Sleep(time.Millisecond * 100)
	}
	ctx.Writer.Write([]byte("data: ok\n\n"))
	Init()
}

func setHeader(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	// 发送初始事件
	resp := Info{Content: "🍲 服务端已收到任务 正在疯狂摸鱼中..."}
	bytes, _ := json.Marshal(resp)
	ctx.Writer.Write([]byte("data: " + string(bytes) + "\n\n"))
	ctx.Writer.Flush()
}
