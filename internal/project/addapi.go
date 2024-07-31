package project

import (
	"fmt"
	"github.com/go-zxb/fuxi/internal/ast/add/api"
	addRepo "github.com/go-zxb/fuxi/internal/ast/add/repo"
	addService "github.com/go-zxb/fuxi/internal/ast/add/service"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

var (
	name         = ""
	method       = "get"
	api          = ""
	apiFunc      = ""
	apiPath      = "internal/api"
	routerPath   = "internal/router"
	servicePath  = "internal/service"
	repoPath     = "internal/repo"
	iSByID       = "true"
	isReturnList = "false"
	noParams     = "false"
	returnType   = ""
)

var AddApiCmd = &cobra.Command{
	Use:   "api:add",
	Short: "单独添加一个接口",
	Long:  "单独添加一个接口",
	Example: "fuxi api:add -n (模块名称) -m (请求方法) -a (api路由路径) -f (方法签名)\n" +
		"fuxi api:add -n user -m get -a userNickname -f userNickname",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" || method == "" || apiFunc == "" {
			cobra.CheckErr("args is empty... 什么? 需要帮助? fuxi api:add -h 可以帮到你!")
			return
		}

		//fmt.Println(iSByID, "-------------", isReturnList, noParams)

		if api == "" {
			api = apiFunc
		}

		apiAdd := addapi.AddApi{
			Name:         name,
			FilePath:     fmt.Sprintf("%s/%s/%s.go", routerPath, name, name),
			Method:       strings.ToUpper(method),
			Api:          api,
			ApiFunc:      apiFunc,
			ISByID:       isTrue(iSByID),
			IsReturnList: isTrue(isReturnList),
			NoParams:     isTrue(noParams),
		}
		if apiAdd.NoParams {
			apiAdd.ISByID = false
		}
		fmt.Println(apiAdd.ISByID)
		err := apiAdd.InsertRouter()
		if err != nil {
			log.Fatalln("❎ 添加api代码错误", err)
		}

		apiAdd.FilePath = fmt.Sprintf("%s/%s/%s.go", apiPath, name, name)
		err = apiAdd.InsertApiHandle()
		if err != nil {
			log.Fatalln("❎ 添加代码错误", err)
		}

		//添加service
		sv := addService.AddService{
			Name:         name,
			FilePath:     fmt.Sprintf("%s/%s/%s.go", servicePath, name, name),
			Method:       strings.ToUpper(method),
			Api:          api,
			ApiFunc:      apiFunc,
			ReturnType:   returnType,
			ISByID:       isTrue(iSByID),
			IsReturnList: isTrue(isReturnList),
			NoParams:     isTrue(noParams),
		}
		if sv.NoParams {
			sv.ISByID = false
		}
		err = sv.InsertService()
		if err != nil {
			log.Fatalln("❎ 添加service代码错误", err)
		}

		repo := addRepo.AddRepo{
			Name:         name,
			FilePath:     fmt.Sprintf("%s/%s/%s.go", repoPath, name, name),
			Method:       strings.ToUpper(method),
			Api:          api,
			ApiFunc:      apiFunc,
			ReturnType:   returnType,
			ISByID:       isTrue(iSByID),
			IsReturnList: isTrue(isReturnList),
			NoParams:     isTrue(noParams),
		}
		if repo.NoParams {
			repo.ISByID = false
		}
		err = repo.InsertRepo()
		if err != nil {
			log.Fatalln("❎ 添加repo代码错误", err)
		}
		log.Println("✅ ", apiFunc, "生成代码添加成功👌！")
		fmt.Println(`程序优化我心欢😊，
携手同行路更宽🛣️；
你我共勉心相连🤝，
共创辉煌乐无边🎉。`)
	},
}

func init() {
	AddApiCmd.Flags().StringVarP(&name, "name", "n", name, "模块名称")
	AddApiCmd.Flags().StringVarP(&method, "method", "m", method, "请求方法:GET,POST,DELETE,PUT")
	AddApiCmd.Flags().StringVarP(&api, "api", "a", api, "api路由路径")
	AddApiCmd.Flags().StringVarP(&apiFunc, "apiFunc", "f", apiFunc, "方法签名")

	AddApiCmd.Flags().StringVarP(&apiPath, "apiPath", "p", apiPath, "api文件路径")
	AddApiCmd.Flags().StringVarP(&routerPath, "routerPath", "r", routerPath, "路由文件路径")
	AddApiCmd.Flags().StringVarP(&servicePath, "servicePath", "s", servicePath, "service文件路径")
	AddApiCmd.Flags().StringVarP(&repoPath, "repoPath", "d", repoPath, "repo文件路径")

	AddApiCmd.Flags().StringVarP(&iSByID, "isByID", "i", "true", "入参为ID ")
	AddApiCmd.Flags().StringVarP(&isReturnList, "list", "l", "false", "是否返回list列表 ")
	AddApiCmd.Flags().StringVarP(&noParams, "noParams", "o", "false", "是否无参数 ")
	AddApiCmd.Flags().StringVarP(&returnType, "returnType", "t", returnType, "返回类型支持 int,uint,int64,float64,string,以及对应的结构体（默认返回结构体）")
}

func isTrue(v string) bool {
	if v == "true" {
		return true
	}
	return false
}
