package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/golang/protobuf/proto"
	"strconv"
	"sync"
	"syscall"
	"time"
	"wentmin/common"
	"wentmin/netmodel"
	wtproto "wentmin/proto"
	"wentmin/protocol"

	"github.com/astaxie/beego/logs"
)

var wg sync.WaitGroup
var stopsignal chan os.Signal
var exitsignal chan struct{}

func CreateClient(id int) {
	cs, err := netmodel.Dial("tcp4", "127.0.0.1:9902")
	if err != nil {
		return
	}

	defer func() {
		cs.Close()
	}()

	for {
		select {
		case <-exitsignal:
			fmt.Println("socket get exit signal")
			logs.Debug("socket get exit signal")
			return
		default:
			accountid := strconv.Itoa(id)
			timestr := strconv.FormatInt(time.Now().Unix(), 10)
			accountid = accountid + "timerstr:" + timestr
			packet := new(protocol.MsgPacket)
			packet.Head.Id = common.USER_REG_CS
			csusereg := &wtproto.CSUserLogin{
				Accountid: accountid,
				Passwd:    "pawd" + accountid,
			}

			//protobuf编码
			pData, err := proto.Marshal(csusereg)
			if err != nil {
				fmt.Println(common.ErrProtobuffMarshal.Error())
				return
			}
			packet.Head.Len = (uint16)(len(pData))
			packet.Body.Data = pData
			cs.Send(packet)
			packetrsp, err := cs.Recv()
			if err != nil {
				fmt.Println("receive error")
				return
			}

			datarsp := packetrsp.(*protocol.MsgPacket)
			fmt.Println("packet id is", datarsp.Head.Id)
			fmt.Println("packet len is", datarsp.Head.Len)
			scusereg := &wtproto.SCUserLogin{}

			error2 := proto.Unmarshal(datarsp.Body.Data, scusereg)
			if error2 != nil {
				fmt.Println(common.ErrProtobuffUnMarshal.Error())
				return
			}

			if scusereg.Errid != common.RSP_SUCCESS {
				fmt.Println("user reg failed ")
				return
			}

			fmt.Println("user reg success ")
			fmt.Println("user token is ", scusereg.Token)
			time.Sleep(time.Second)
		}

	}

}

func main() {
	stopsignal = make(chan os.Signal) // 接收系统中断信号
	var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	signal.Notify(stopsignal, shutdownSignals...)
	exitsignal = make(chan struct{})
	go func() {
		select {
		case sign := <-stopsignal:
			fmt.Println("catch stop signal, ", sign)
			logs.Debug("catch stop signal, ", sign)
			close(exitsignal)
		}
	}()

	for i := 1000; i < 3000; i++ {
		go CreateClient(i)
		time.Sleep(time.Millisecond * 10)

	}
	<-exitsignal
	fmt.Println("main goroutine exit ")
	logs.Debug("main goroutine exit ")
}
