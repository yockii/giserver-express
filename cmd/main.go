package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/controller"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/config"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/server"
	"github.com/yockii/giserver-express/pkg/util"
)

var VERSION = ""
var (
	daemon bool
)

func init() {
	config.DefaultInstance.SetDefault("server.port", 8080)
	config.DefaultInstance.SetDefault("database.driver", "sqlite")
	config.DefaultInstance.SetDefault("database.host", "./conf/data.db")
	config.DefaultInstance.SetDefault("database.prefix", "t_")
	config.DefaultInstance.SetDefault("logger.level", "debug")

	flag.BoolVar(&daemon, "daemon", false, "以守护进程方式启动")
	flag.Parse()

	// 写入配置文件
	if err := config.DefaultInstance.WriteConfig(); err != nil {
		logger.Errorln(err)
	}

	util.InitNode(1)
}

func main() {
	if daemon {
		runDaemon(os.Args)
		return
	}

	logger.Infoln("当前应用版本: " + VERSION)

	// 检查数据库文件是否存在
	if config.GetString("database.driver") == "sqlite" {
		_, err := os.Stat(config.GetString("database.host"))
		if err != nil && os.IsNotExist(err) {
			// 不存在
			f, _ := os.Create(config.GetString("database.host"))
			f.Close()
		}
	}
	database.Initial()
	database.DB.Sync2(model.SyncModels...)

	startWeb()
}

// 以守护进程方式启动
func runDaemon(args []string) {
	fmt.Printf("pid:%d ppid: %d, arg: %s \n", os.Getpid(), os.Getppid(), os.Args)
	// 去除--daemon参数，启动主程序
	for i := 0; i < len(args); {
		if args[i] == "--daemon" && i != len(args)-1 {
			args = append(args[:i], args[i+1:]...)
		} else if args[i] == "--daemon" && i == len(args)-1 {
			args = args[:i]
		} else {
			i++
		}
	}
	// 启动子进程
	for {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "启动失败, Error: %s \n", err)
			return
		}
		fmt.Printf("守护进程模式启动, pid:%d ppid: %d, arg: %s \n", cmd.Process.Pid, os.Getpid(), args)
		cmd.Wait()
	}
}

func startWeb() {
	controller.InitRouter()
	logger.Error(server.Start())
}
