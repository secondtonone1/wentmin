package weblogic

import (
	"fmt"
	"time"
	"videocall/common"

	"golang.org/x/net/websocket"
)

type UserData struct {
	AccountId string
	Passwd    string
	Phone     string
	Conn      *websocket.Conn
}

type UserMgr struct {
	UsersMap   map[string]*UserData
	SessionMap map[string]*UserData
}

var UserMgrInst *UserMgr
var WebsocketClose chan struct{}
var HeartClose chan struct{}

func init() {
	UserMgrInst = new(UserMgr)
	UserMgrInst.UsersMap = make(map[string]*UserData)
	UserMgrInst.SessionMap = make(map[string]*UserData)
	WebsocketClose = make(chan struct{})
	HeartClose = make(chan struct{})
	go AliveCheck()
}

func CloseHeartG() {
	close(WebsocketClose)
}

func WaitHeartClose() chan struct{} {
	return HeartClose
}

func AliveCheck() {
	t1 := time.NewTimer(10 * time.Second)
	//后期检测用户心跳时间，做清理处理
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("alive recover from panic ")
		}
		fmt.Println("alive exit")
		t1.Stop()
		close(HeartClose)
	}()

	//做检测心跳时间，后期补充
	//to do
	for {
		select {
		case <-WebsocketClose:
			return
		case <-t1.C:
			//fmt.Println("timer tick now")
			//UserMgrInst.ClearDeadSession()
			t1.Reset(10 * time.Second)
			continue

		}
	}

}

func (um *UserMgr) ClearDeadSession() {
	//清理dead session
	// 后期补充
	WebLogicLock.Lock()
	defer WebLogicLock.Unlock()
}

func (um *UserMgr) AddUser(ud *UserData) {

	um.UsersMap[ud.AccountId] = ud
	wskey := ud.Conn.Request().Header.Get("Sec-Websocket-Key")
	//fmt.Println("wskey is ", wskey)
	_, ok := um.SessionMap[wskey]
	if ok {
		fmt.Printf("web socket key %s is exsit\n", wskey)
		return
	}

	um.SessionMap[wskey] = ud

}

func (um *UserMgr) GetUser(accountid string) (*UserData, error) {
	ud, ok := um.UsersMap[accountid]
	if !ok {
		return nil, common.ErrAccountNotFound
	}

	return ud, nil
}

func (um *UserMgr) OnOffline(conn *websocket.Conn) {
	wskey := conn.Request().Header.Get("Sec-Websocket-Key")
	ud, ok := um.SessionMap[wskey]
	if !ok {
		return
	}
	ud.Conn = nil
	ChatMgrInst.DelRoomByUser(ud.AccountId)
}

func (ud *UserData) IsOnline() bool {
	if ud.Conn == nil {
		return false
	}

	return true
}

func (ud *UserData) GetConn() *websocket.Conn {
	return ud.Conn
}
