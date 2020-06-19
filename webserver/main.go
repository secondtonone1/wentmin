package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wentmin/components"
	"wentmin/netmodel"
	"wentmin/weblogic"

	"github.com/astaxie/beego/logs"
)

func main() {
	components.InitWebCfg()
	weblogic.RegWSHandlers()
	weblogic.RegHttpHandlers()
	wbs := netmodel.NewWebServer()
	go wbs.Start()
	stopsignal := make(chan os.Signal) // 接收系统中断信号
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(stopsignal, shutdownSignals...)
	select {
	case sign := <-stopsignal:
		fmt.Println("catch stop signal, ", sign)
		logs.Debug("catch stop signal, ", sign)
		wbs.Stop()
		weblogic.CloseHeartG()

	case <-weblogic.WaitHeartClose():
		fmt.Println("catch hearbeat goroutine closed")
		wbs.Stop()
	}

	logs.Debug("wt server closed successfully ")
}
