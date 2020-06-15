package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "wentmin/logic"
	"wentmin/netmodel"
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
		wt.Close()
	case <-wt.WaitClose():
		fmt.Println("catch wt server closed ")
	case <-netmodel.MsgQueueInst.WaitClose():
		fmt.Println("catch msgqueue closed")
		wt.Close()
	case <-netmodel.OutMsgQueInst.WaitClose():
		fmt.Println("catch output msg queue closed ")
		wt.Close()
	}
	<-wt.WaitClose()
	fmt.Println("wt server closed successfully ")
}
