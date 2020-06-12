package netmodel

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"wentmin/common"
	"wentmin/components"
	"wentmin/protocol"
)

var PacketChan chan *protocol.MsgPacket
var AcceptClose chan struct{}

func init() {
	PacketChan = make(chan *protocol.MsgPacket, components.MaxMsgQueLen)
	AcceptClose = make(chan struct{})
	NewMsgQueue()
	for i := 0; i < components.MaxMsgQueNum; i++ {
		go MsgQueueInst.ReadFromChan()
		MsgWatiGroup.Add(1)
	}

}
func NewTcpServer() (*WtServer, error) {
	address := "0.0.0.0:" + strconv.Itoa(components.ServerPort)
	listenert, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listen failed !!!")
		return nil, common.ErrListenFailed
	}

	return &WtServer{listener: listenert,
		once: &sync.Once{}, sessionGroup: &sync.WaitGroup{}, notifyMain: make(chan struct{})}, nil
}

type WtServer struct {
	listener     net.Listener
	once         *sync.Once
	sessionGroup *sync.WaitGroup
	notifyMain   chan struct{}
}

//主协程主动关闭accept
func (wt *WtServer) Close() {
	wt.once.Do(func() {
		if wt.listener != nil {
			defer wt.listener.Close()
		}
	})
}

func (wt *WtServer) acceptLoop() error {

	tcpConn, err := wt.listener.Accept()
	if err != nil {
		fmt.Println("Accept error!, err is ", err.Error())
		return common.ErrAcceptFailed
	}

	newsess := NewSession(tcpConn, wt.sessionGroup)
	fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
	newsess.Start()
	wt.sessionGroup.Add(1)
	return nil

}

func (wt *WtServer) AcceptLoop() {
	fmt.Println("Server begin accept ...")
	defer func() {
		fmt.Println("main io goroutin exit ")
		if err := recover(); err != nil {
			fmt.Println("server recover from err , err is ", err)
		}
		close(AcceptClose)
		wt.sessionGroup.Wait()
		close(wt.notifyMain)
	}()
	for {
		if err := wt.acceptLoop(); err != nil {
			fmt.Println("went server accept failed!! ")
			return
		}
	}
}

func (wt *WtServer) WaitClose() chan struct{} {
	return wt.notifyMain
}
