package netmodel

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	"wentmin/common"
	"wentmin/components"

	"github.com/astaxie/beego/logs"
)

var PacketChan chan *MsgSession
var AcceptClose chan struct{}

var TcpServerInst *WtServer = nil
var AliveClose chan struct{}

func init() {
	PacketChan = make(chan *MsgSession, components.MaxMsgQueLen)
	AcceptClose = make(chan struct{})
	AliveClose = make(chan struct{})
}
func NewTcpServer() (*WtServer, error) {
	address := "0.0.0.0:" + strconv.Itoa(components.ServerPort)
	fmt.Println("address ", address)
	listenert, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listen failed !!!")
		logs.Debug("listen failed !!!")
		fmt.Println(err.Error())
		return nil, common.ErrListenFailed
	}

	TcpServerInst = &WtServer{listener: listenert,
		once: &sync.Once{}, sessionGroup: &sync.WaitGroup{},
		notifyMain: make(chan struct{}), sessionMap: make(map[int]*Session),
		SocketIndex: 0, UnUsedSocketMap: make(map[int]bool)}

	//创建消息处理队列
	NewMsgQueue()
	//创建消息发送队列
	NewOutMsgQues()
	//启动心跳监听协程
	go TcpServerInst.OnCheckAlive()
	return TcpServerInst, nil
}

type WtServer struct {
	listener        net.Listener
	once            *sync.Once
	sessionGroup    *sync.WaitGroup
	notifyMain      chan struct{}
	sessionMap      map[int]*Session
	sessionLock     sync.Mutex
	SocketIndex     int
	UnUsedSocketMap map[int]bool
}

//主协程主动关闭accept
func (wt *WtServer) Close() {
	wt.once.Do(func() {
		if wt.listener != nil {
			defer wt.listener.Close()
		}
	})
}

func (wt *WtServer) GenerateSocket() int {
	//fmt.Println("unused socket map is ", wt.UnUsedSocketMap)
	if len(wt.UnUsedSocketMap) < 20 {
		wt.SocketIndex++
		return wt.SocketIndex
	}

	value := 0
	for k, _ := range wt.UnUsedSocketMap {
		value = k
		delete(wt.UnUsedSocketMap, k)
	}
	return value
}

func (wt *WtServer) RecycleSocket(socket int) {
	wt.UnUsedSocketMap[socket] = true
	//fmt.Println("after RecycleSocket, unusedsocket map is ", wt.UnUsedSocketMap)
}

func (wt *WtServer) OnSessConnect(se *Session) {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	wt.sessionGroup.Add(1)
	se.SocketId = wt.GenerateSocket()
	wt.sessionMap[se.SocketId] = se
	se.Start()
}

func (wt *WtServer) acceptLoop() error {

	tcpConn, err := wt.listener.Accept()
	if err != nil {
		fmt.Println("Accept error!, err is ", err.Error())
		logs.Debug("Accept error!, err is ", err.Error())
		return common.ErrAcceptFailed
	}

	newsess := NewSession(tcpConn)
	fmt.Println("A client connected :" + tcpConn.RemoteAddr().String())
	logs.Debug("A client connected :" + tcpConn.RemoteAddr().String())
	wt.OnSessConnect(newsess)
	return nil

}

func (wt *WtServer) AcceptLoop() {
	fmt.Println("Server begin accept ...")
	logs.Debug("Server begin accept ...")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("server recover from err , err is ", err)
			logs.Debug("server recover from err , err is ", err)
		}
		//先清除连接
		wt.ClearSessions()
		close(AcceptClose)
		wt.sessionGroup.Wait()
		MsgWatiGroup.Wait()
		OutputWaitGroup.Wait()
		<-WaitAliveClose()
		close(wt.notifyMain)
		fmt.Println("main io goroutin exit ")
		logs.Debug("main io goroutin exit ")
	}()
	for {
		if err := wt.acceptLoop(); err != nil {
			fmt.Println("went server accept failed!! ")
			logs.Debug("went server accept failed!! ")
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
		//清除所有连接，意味着服务器可能要官服，那么就不投递连接断开消息
		//MsgQueueInst.PutCloseMsgInQue(session)
		delete(wt.sessionMap, id)
		wt.RecycleSocket(id)
		wt.sessionGroup.Done()
		fmt.Printf("session id %d closed successfully\n", id)
		logs.Debug("session id %d closed successfully\n", id)
	}
}

//服务器主动关闭session
func (wt *WtServer) CloseSession(sid int) {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	session, ok := wt.sessionMap[sid]
	if !ok {
		fmt.Println("not found session by id ", sid)
		logs.Debug("not found session by id ", sid)
		return
	}
	session.Close()
	MsgQueueInst.PutCloseMsgInQue(session)
	delete(wt.sessionMap, sid)
	wt.RecycleSocket(sid)
	wt.sessionGroup.Done()
	fmt.Printf("session id %d closed successfully\n", sid)
	logs.Debug("session id %d closed successfully\n", sid)
}

//连接断开回调函数
func (wt *WtServer) OnSessionClosed(sid int) {
	if err := recover(); err != nil {
		fmt.Println(" recover from error ", err)
		logs.Debug(" recover from error ", err)
	}

	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	session, ok := wt.sessionMap[sid]
	if !ok {
		fmt.Printf("not found session by %d , maybe it has been closed \n", sid)
		logs.Debug("not found session by %d , maybe it has been closed \n", sid)
		return
	}
	session.Close()
	MsgQueueInst.PutCloseMsgInQue(session)
	delete(wt.sessionMap, sid)
	wt.RecycleSocket(sid)
	wt.sessionGroup.Done()
	fmt.Printf("session id %d closed successfully\n", sid)
	logs.Debug("session id %d closed successfully\n", sid)
}

func (wt *WtServer) ClearDeadSession() {
	wt.sessionLock.Lock()
	defer wt.sessionLock.Unlock()
	cur := time.Now().Unix()
	for _, session := range wt.sessionMap {
		if !session.CheckAlive(cur) {
			fmt.Printf("session id %d not alive ", session.SocketId)
			session.Close()
			MsgQueueInst.PutCloseMsgInQue(session)
			delete(wt.sessionMap, session.SocketId)
			wt.RecycleSocket(session.SocketId)
			wt.sessionGroup.Done()
			fmt.Printf("session id %d closed successfully\n", session.SocketId)
			logs.Debug("session id %d closed successfully\n", session.SocketId)
		}
	}
}

func (wt *WtServer) OnCheckAlive() {
	t1 := time.NewTimer(10 * time.Second)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("timer recover from error ", err)
			logs.Debug("timer recover from error ", err)
		}
		t1.Stop()
		close(AliveClose)
	}()
	for {
		select {
		case <-AcceptClose:
			return
		case <-t1.C:
			//fmt.Println("timer tick now")
			//logs.Debug("timer tick now")
			wt.ClearDeadSession()
			t1.Reset(10 * time.Second)
			continue
		}
	}
}

func WaitAliveClose() chan struct{} {
	return AliveClose
}
