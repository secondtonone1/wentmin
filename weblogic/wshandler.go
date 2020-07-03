package weblogic

import (
	"sync"
	"wentmin/common"

	"golang.org/x/net/websocket"
)

var WebMsgMap map[int]func(*websocket.Conn, string) error
var WebLogicLock sync.RWMutex

func init() {
	WebMsgMap = make(map[int]func(*websocket.Conn, string) error)
	RegMsgHandler(common.WEB_CS_USER_REG, UserLogin)
	RegMsgHandler(common.WEB_CS_USER_CALL, UserCall)
	RegMsgHandler(common.WEB_REPLY_BECALL, UserCallReply)
	RegMsgHandler(common.WEB_CS_TERMINAL_CALL, UserTerminalCall)
	RegMsgHandler(common.WEB_CS_CALL_OFFER, ReceiveOffer)
	RegMsgHandler(common.WEB_CS_BECALL_ANSWER, ReceiveAnswer)
	RegMsgHandler(common.WEB_CS_CALL_ICE, ReceiveCallIce)
	RegMsgHandler(common.WEB_CS_BECALL_ICE, ReceiveBeCallIce)
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
