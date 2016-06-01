package main

import (
	"runtime"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/smtc/glog"
	"flag"
	"os"
	"os/signal"
	"syscall"
)
var(
	configFn                                    = flag.String("config", "./config.json", "config file path")
	debugFlag                                   = flag.Bool("d", false, "debug mode")
)
/**
主程序函数
创建人:邵炜
创建时间:2016年6月1日09:38:36
 */
func main() {
	//判断进程是否存在
	if checkPid() {
		return
	}

	flag.Parse()

	serverRun(*configFn, *debugFlag)

	c := make(chan os.Signal, 1)
	writePid()
	// 信号处理
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// 等待信号
	<-c

	serverExit()
	rmPidFile()
	glog.Close()
	os.Exit(0)
}

func serverRun(cfn string, debug bool) {
	config.ReadCfg(cfn)
	logInit(debug)

	// 初始化
	deferinit.InitAll()
	glog.Info("init all module successfully.\n")

	// 设置多cpu运行
	runtime.GOMAXPROCS(runtime.NumCPU())

	deferinit.RunRoutines()
	glog.Info("run routines successfully.\n")

	readFile("./urlList.txt",true)
	go loadFileDB("./afterProcess.txt")
}

// 结束进程
func serverExit() {
	// 结束所有go routine
	deferinit.StopRoutines()
	glog.Info("stop routine successfully.\n")

	deferinit.FiniAll()
	glog.Info("fini all modules successfully.\n")
}