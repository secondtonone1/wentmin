package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "wentmin/logic"
	"wentmin/netmodel"

	"github.com/astaxie/beego/logs"
)

func main() {
	//logs.Debug("server port is %d", components.ServerPort)
	wt, err := netmodel.NewTcpServer()
	if err != nil {
		panic("new tcp server failed")

	}
	go wt.AcceptLoop()
	stopsignal := make(chan os.Signal) // 接收系统中断信号
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(stopsignal, shutdownSignals...)
	select {
	case sign := <-stopsignal:
		fmt.Println("catch stop signal, ", sign)
		logs.Debug("catch stop signal, ", sign)
		wt.Close()
	case <-wt.WaitClose():
		fmt.Println("catch wt server closed ")
		logs.Debug("catch wt server closed ")
	case <-netmodel.MsgQueueInst.WaitClose():
		fmt.Println("catch msgqueue closed")
		logs.Debug("catch msgqueue closed")
		wt.Close()
	case <-netmodel.OutMsgQueInst.WaitClose():
		fmt.Println("catch output msg queue closed ")
		logs.Debug("catch output msg queue closed ")
		wt.Close()
	case <-netmodel.WaitAliveClose():
		fmt.Println("catch alive watcher goroutine close")
		logs.Debug("catch alive watcher goroutine close")
		wt.Close()
	}
	<-wt.WaitClose()
	fmt.Println("wt server closed successfully ")
	logs.Debug("wt server closed successfully ")
}
