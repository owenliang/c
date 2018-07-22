package main

import (
	"flag"
	"runtime"
	"github.com/owenliang/c/master"
	"fmt"
	"os"
	"time"
)

var (
	confFile string		// 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// --config master.json
	flag.StringVar(&confFile, "config", "master.json", "master.json配置文件路径")
	flag.Parse()
}

// 初始化线程数量 == CPU数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// 解析命令行参数
	initArgs()

	// 初始化运行环境
	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	// 任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// API接口服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	// 启动失败
ERR:
	fmt.Fprintln(os.Stderr, err)
	return
}