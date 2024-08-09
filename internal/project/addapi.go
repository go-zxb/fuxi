package project

import (
	"fmt"
	"github.com/go-zxb/fuxi/internal/ast/add/api"
	addRepo "github.com/go-zxb/fuxi/internal/ast/add/repo"
	addService "github.com/go-zxb/fuxi/internal/ast/add/service"
	"github.com/go-zxb/fuxi/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var gushi = `程序优化我心欢😊，
携手同行路更宽🛣️；
你我共勉心相连🤝，
共创辉煌乐无边🎉。`

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

func Init() {
	//Print()
	name = "great"
	method = "get"
	api = ""
	apiFunc = ""
	apiPath = "internal/api"
	routerPath = "internal/router"
	servicePath = "internal/service"
	repoPath = "internal/repo"
	iSByID = "true"
	isReturnList = "false"
	noParams = "false"
	returnType = ""
	debug = ""
	filename = ""
	modelPath = "internal/model"
	gormGenPath = "cmd/gorm/gen"
	question = ""
	isOutputJson = ""
	localJson = "false"
	empty = "false"
}

func Print() {
	fmt.Println(
		"name: ", name,
		"method: ", method,
		"api: ", api,
		"apiFunc: ", apiFunc,
		"apiPath: ", apiPath,
		"routerPath: ", routerPath,
		"servicePath: ", servicePath,
		"repoPath: ", repoPath,
		"iSByID: ", iSByID,
		"isReturnList: ", isReturnList,
		"noParams: ", noParams,
		"returnType: ", returnType,
		"debug: ", debug,
		"filename: ", filename,
		"modelPath: ", modelPath,
		"gormGenPath: ", gormGenPath,
		"question: ", question,
		"isOutputJson: ", isOutputJson,
		"localJson: ", localJson,
		"empty: ", empty)
}

var AddApiCmd = &cobra.Command{
	Use:   "api:add",
	Short: "单独添加一个接口",
	Long:  "单独添加一个接口",
	Example: "fuxi api:add -n (模块名称) -m (请求方法) -a (api路由路径) -f (方法签名)\n" +
		"fuxi api:add -n user -m get -a userNickname -f userNickname",
	Run: func(cmd *cobra.Command, args []string) {
		infoChan := make(chan pkg.CommandInfo)
		go addApiHandle(infoChan)
		for info := range infoChan {
			if info.Error != nil {
				log.Fatalln("❌", info.Message, info.Error.Error())
			} else {
				log.Println(info.Message)
			}
		}
	},
}

func addApiHandle(infoChan chan<- pkg.CommandInfo) {
	defer close(infoChan) // 确保在函数返回时关闭通道
	if name == "" || method == "" || apiFunc == "" {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴args is empty... 什么? 需要帮助? fuxi api:add -h 可以帮到你!", Error: nil}
		return
	}

	if api == "" {
		api = strings.ToLower(apiFunc)
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
	moduleName, err := pkg.GetModuleName("go.mod")
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 获取项目名称错误", Error: err}
		return
	}
	zipFileName := fmt.Sprintf(".fuxi/%s/%s/%sBeforeCode/%s.zip", moduleName, name, apiFunc, time.Now().Format("20060102150405"))
	err = BeforeCodeToZip(zipFileName, apiAdd.FilePath, fmt.Sprintf("%s/%s/%s.go", apiPath, name, name))
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
		return
	}

	err = apiAdd.InsertRouter()
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 添加路由代码错误", Error: err}
	} else {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅ 添加路由代码成功", Error: err}
	}

	apiAdd.FilePath = fmt.Sprintf("%s/%s/%s.go", apiPath, name, name)
	err = apiAdd.InsertApiHandle()
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 添加ApiHandle代码错误", Error: err}
	} else {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅ 添加ApiHandle代码成功", Error: err}
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

	err = BeforeCodeToZip(zipFileName, sv.FilePath)
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
		return
	}

	err = sv.InsertService()
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 添加service代码错误", Error: err}
	} else {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅ 添加service代码成功", Error: err}
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

	err = BeforeCodeToZip(zipFileName, repo.FilePath)
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴🚶‍♀任务中断🚶 因为备份数据失败🎒....", Error: err}
		return
	}

	err = repo.InsertRepo()
	if err != nil {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 添加repo代码错误", Error: err}
		infoChan <- pkg.CommandInfo{Message: "🐮🐴❎ 似乎啥也没干 就下班了!!!", Error: nil}
	} else {
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅ 添加repo代码成功", Error: nil}
		infoChan <- pkg.CommandInfo{Message: "🐮🐴✅ 生成" + apiFunc + "代码添加成功👌！", Error: nil}
		infoChan <- pkg.CommandInfo{Message: gushi, Error: nil}
	}
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

func BeforeCodeToZip(zipFileName string, apiCodePath ...string) error {

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

	return pkg.FilesToZip(zipFileName, apiCodePath)
}
