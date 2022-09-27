package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/internal/config"
	"github.com/super-l/nproxy/internal/consts"
	"github.com/super-l/nproxy/internal/loader"
	"github.com/super-l/nproxy/services/rpc"
	"github.com/super-l/nproxy/task"
	"github.com/super-l/nproxy/utils"
	"runtime"
	"strings"
	"sync"
)

func showBanner() {
	name := fmt.Sprintf("%s (v.%s)", consts.Name, consts.Version)
	banner := `

     __  ___                     
  /\ \ \/ _ \_ __ _____  ___   _ 
 /  \/ / /_)/ '__/ _ \ \/ / | | |
/ /\  / ___/| | | (_) >  <| |_| |
\_\ \/\/    |_|  \___/_/\_\\__, |
                           |___/ 

	`
	// Shell width
	all_lines := strings.Split(banner, "\n")
	w := len(all_lines[1])

	// Print Centered
	color.Green(banner)
	color.Red(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(name))/2, name)))

	fmt.Println()
}

func initDir() {
	dirs := []string{
		"data",
		"data/logs",
	}
	for _, dir := range dirs {
		isExist, _ := utils.File.FileOrPathIsExists(dir)
		if !isExist {
			utils.File.CreateDir(dir)
		}
	}
}

func main() {
	showBanner()
	maxProcess := runtime.NumCPU()
	if maxProcess > 1 {
		maxProcess -= 1
	}
	runtime.GOMAXPROCS(maxProcess)

	var wg sync.WaitGroup
	wg.Add(1)

	initDir()

	// 初始化配置
	var err error
	err = config.InitConfig()
	if err != nil {
		internal.SLogger.StdoutLogger.Errorf("load config.yaml exception, %s", err.Error())
		return
	}

	// 初始化日志组件
	err = internal.SLogger.InitLogger()
	if err != nil {
		internal.SLogger.StdoutLogger.Errorf("log component initialization exception, %s", err.Error())
		return
	}

	// 检查数据库是否正常链接
	db := internal.GetDbInstance()
	if db == nil {
		internal.SLogger.StdoutLogger.Error("database is connect failed!")
		return
	}

	// 初始化数据库数据
	err = loader.InitDbLoader()
	if err != nil {
		internal.SLogger.StdoutLogger.Errorf("connect db failed! %s", err.Error())
		return
	}

	// API代理功能支持
	go task.ProxyTaskSercice.Start()

	// 运行RPC服务
	go rpc.InitRpcServer()
	internal.SLogger.StdoutLogger.Infof("rpc server start successful! port: %s", config.GetConfig().Rpc.Port)
	wg.Wait()
}
