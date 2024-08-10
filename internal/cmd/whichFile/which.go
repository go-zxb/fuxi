package whichFile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	config2 "github.com/go-zxb/fuxi/config"
	"github.com/go-zxb/fuxi/log"
	"github.com/go-zxb/fuxi/pkg"
)

var (
	startTime = time.Now()
)

func Which(watchPath string, excludeDirs []string, allFiles bool) {
	var timer *time.Timer
	restartApp()
	for {
		filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if path == ".git" && info.IsDir() {
				return filepath.SkipDir
			}
			for _, x := range excludeDirs {
				if x == path {
					return filepath.SkipDir
				}
			}

			if filepath.Base(path)[0] == '.' {
				return nil
			}

			if (allFiles || filepath.Ext(path) == ".go") && info.ModTime().After(startTime) {
				if timer != nil {
					timer.Stop()
				}

				log.InfoWithPath(path, " 源文件变动3秒后重新运行...")
				timer = time.AfterFunc(3*time.Second, restartApp)
				startTime = time.Now()
				return errors.New("done")
			}

			return nil
		})
		time.Sleep(500 * time.Millisecond)
	}
}

var ctx context.Context
var cancel context.CancelFunc

func restartApp() {
	if cancel != nil {
		cancel()
		exists, err := pkg.PathExists("config/config.yaml")
		if err != nil {
			log.Error(err)
			return
		}
		if !exists {
			log.Info("读取项目配置文件失败,项目配置文件不存在")
			return
		}

		conf, err := config2.NewConfig2("config/config.yaml")
		if err != nil {
			log.Info(err.Error())
			return
		}
		// 获取进程pid杀死正在运行的进程
		pid, err := pkg.PIDByPort(conf.Gin.Port)
		if err != nil {
			log.Info(err.Error())
			return
		}
		err = pkg.RunCommand("taskkill", "/PID", fmt.Sprintf("%s", pid), "/F")
		if err != nil {
			log.Info(err.Error())
		}
		time.Sleep(500 * time.Millisecond)
	}

	ctx, cancel = context.WithCancel(context.Background())
	go pkg.RunCommandWithCtx(ctx, "go", "run", "main.go", "server")

	log.Info("开始监听文件变动...")
}
