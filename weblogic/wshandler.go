package weblogic

import (
	"fmt"
	"sync"
	"wentmin/common"

	"golang.org/x/net/websocket"
)

var WebMsgMap map[int]func(*websocket.Conn, string) error
var WebLogicLock sync.RWMutex

func init() {
	WebMsgMap = make(map[int]func(*websocket.Conn, string) error)
	RegMsgHandler(common.WEB_CS_USER_CALL, UserReg)
}

func RegMsgHandler(msgid int, handler func(*websocket.Conn, string) error) {
	WebMsgMap[msgid] = handler
}

func HandleWebMsg(conn *websocket.Conn, msgid int, msgdata string) error {
	handler, ok := WebMsgMap[msgid]
	if !ok {
		return common.ErrMsgIDNotReg
	}
	WebLogicLock.Lock()
	defer WebLogicLock.Unlock()
	return handler(conn, msgdata)
}

func UserReg(conn *websocket.Conn, msgdata string) error {
	fmt.Println("receive user reg req , msgdata is ", msgdata)

	//这里是发送消息
	/*
		if err := websocket.Message.Send(ws, msg); err != nil {

			fmt.Println("send failed:", err)

		}*/
	return nil
}
