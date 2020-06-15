package netmodel

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"wentmin/common"
	"wentmin/components"
)

var PacketChan chan *MsgSession
var AcceptClose chan struct{}
var SocketIndex int
var TcpServerInst *WtServer = nil

func init() {
	PacketChan = make(chan *MsgSession, components.MaxMsgQueLen)
	AcceptClose = make(chan struct{})
	//创建消息处理队列
	NewMsgQueue()
	//创建消息发送队列
	NewOutMsgQues()
}
func NewTcpServer() (*WtServer, error) {
	address := "0.0.0.0:" + strconv.Itoa(components.ServerPort)
	listenert, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listen failed !!!")
		return nil, common.ErrListenFailed
	}

	TcpServerInst = &WtServer{listener: listenert,
		once: &sync.Once{}, sessionGroup: &sync.WaitGroup{},
		notifyMain: make(chan struct{}), sessionMap: make(map[int]*Session)}
	return TcpServerInst, nil
}

type WtServer struct {
	listener     net.Listener
	once         *sync.Once
	sessionGroup *sync.WaitGroup
	notifyMain   chan struct{}
	sessionMap   map[int]*Session
	sessionLock  sync.Mutex
}

//主协程主动关闭accept
func (wt *WtServer) Close() {
	wt.once.Do(func() {
		if wt.listener != nil {
			defer wt.listener.Close()
		}
	})
}

func (wt *WtServer) OnSessConnect(se *Session) {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	wt.sessionGroup.Add(1)
	wt.sessionMap[se.SocketId] = se
}

func (wt *WtServer) acceptLoop() error {

	tcpConn, err := wt.listener.Accept()
	if err != nil {
		fmt.Println("Accept error!, err is ", err.Error())
		return common.ErrAcceptFailed
	}
	SocketIndex++
	newsess := NewSession(tcpConn, SocketIndex)
	fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
	newsess.Start()
	wt.OnSessConnect(newsess)
	return nil

}

func (wt *WtServer) AcceptLoop() {
	fmt.Println("Server begin accept ...")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("server recover from err , err is ", err)
		}
		wt.ClearSessions()
		close(AcceptClose)
		wt.sessionGroup.Wait()
		MsgWatiGroup.Wait()
		OutputWaitGroup.Wait()
		close(wt.notifyMain)
		fmt.Println("main io goroutin exit ")
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

//服务器关闭所有连接
func (wt *WtServer) ClearSessions() {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	for id, ss := range wt.sessionMap {
		ss.Close()
		delete(wt.sessionMap, id)
		wt.sessionGroup.Done()
	}
}

//服务器主动关闭session
func (wt *WtServer) CloseSession(sid int) {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	session, ok := wt.sessionMap[sid]
	if !ok {
		fmt.Println("not found session by id ", sid)
		return
	}
	session.Close()
	delete(wt.sessionMap, sid)
	wt.sessionGroup.Done()
	fmt.Printf("session id %d closed successfully", sid)
}

//连接断开回调函数
func (wt *WtServer) OnSessionClosed(sid int) {
	if err := recover(); err != nil {
		fmt.Println(" recover from error ", err)
	}

	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	session, ok := wt.sessionMap[sid]
	if !ok {
		fmt.Printf("not found session by %d , maybe it has been closed \n", sid)
		return
	}
	session.Close()
	delete(wt.sessionMap, sid)
	wt.sessionGroup.Done()
	fmt.Printf("session id %d closed successfully", sid)
}
