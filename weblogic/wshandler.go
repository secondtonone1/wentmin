package weblogic

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"wentmin/common"
	"wentmin/jsonproto"

	"golang.org/x/net/websocket"
)

var WebMsgMap map[int]func(*websocket.Conn, string) error
var WebLogicLock sync.RWMutex

func init() {
	WebMsgMap = make(map[int]func(*websocket.Conn, string) error)
	RegMsgHandler(common.WEB_CS_USER_REG, UserReg)
	RegMsgHandler(common.WEB_CS_USER_CALL, UserCall)
	RegMsgHandler(common.WEB_REPLY_BECALL, UserCallReply)
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
	fmt.Println("receive user reg req , msgdata is ", string(msgdata))
	csreg := &jsonproto.CSUserReg{}
	err := json.Unmarshal([]byte(msgdata), csreg)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return common.ErrJsonUnMarshal
	}

	ud := &UserData{}
	ud.AccountId = csreg.AccountId
	ud.Conn = conn
	ud.Passwd = csreg.Passwd
	ud.Phone = csreg.Phone
	UserMgrInst.AddUser(ud)

	screg := &jsonproto.SCUserReg{}
	screg.AccountId = csreg.AccountId
	screg.ErrorId = common.RSP_SUCCESS
	screg.Passwd = csreg.Passwd
	screg.Phone = csreg.Phone

	jsreg, _ := json.Marshal(screg)
	jsmsg := jsonproto.JsonMsg{MsgId: common.WEB_SC_USER_REG, MsgData: string(jsreg)}
	jsrt, _ := json.Marshal(jsmsg)

	nw, err := conn.Write(jsrt)
	if err != nil {
		fmt.Println("write failed")
		return common.ErrJsonUnMarshal
	}
	if nw == 0 {
		fmt.Println("peer connection closed ")
		return common.ErrWebSocketClosed
	}
	return nil
}

func UserCall(conn *websocket.Conn, msgdata string) error {
	fmt.Println("receive user call req , msgdata is ", string(msgdata))
	cscall := &jsonproto.CSUserCall{}
	err := json.Unmarshal([]byte(msgdata), cscall)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return common.ErrJsonUnMarshal
	}

	fmt.Println("caller is ", cscall.Caller)
	fmt.Println("becalled is ", cscall.BeCalled)
	becall, err := UserMgrInst.GetUser(cscall.BeCalled)
	//没找到用户
	if err != nil {
		sccall := &jsonproto.SCUserCall{}
		sccall.Caller = cscall.Caller
		sccall.BeCalled = cscall.BeCalled
		sccall.ErrorId = common.RSP_USER_NOT_FOUND
		sccallms, _ := json.Marshal(sccall)
		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = string(sccallms)
		jmsgms, _ := json.Marshal(jmsg)
		nw, err := conn.Write(jmsgms)
		if err != nil {
			fmt.Println("write failed")
			return common.ErrWebSocketClosed
		}
		if nw == 0 {
			fmt.Println("peer connection closed ")
			return common.ErrWebSocketClosed
		}
		return nil
	}
	//用户不在线
	if !becall.IsOnline() {
		sccall := &jsonproto.SCUserCall{}
		sccall.Caller = cscall.Caller
		sccall.BeCalled = cscall.BeCalled
		sccall.ErrorId = common.RSP_USER_NOT_ONLINE
		sccall.Phone = becall.Phone
		sccallms, _ := json.Marshal(sccall)
		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = string(sccallms)
		jmsgms, _ := json.Marshal(jmsg)
		nw, err := conn.Write(jmsgms)
		if err != nil {
			fmt.Println("write failed")
			return common.ErrWebSocketClosed
		}
		if nw == 0 {
			fmt.Println("peer connection closed ")
			return common.ErrWebSocketClosed
		}
		return nil
	}

	//被呼叫人在线，通知被呼叫人

	notifybecall := &jsonproto.SCNotifyBeCall{}
	notifybecall.Caller = cscall.Caller
	notifybecall.BeCalled = cscall.BeCalled

	notifyms, _ := json.Marshal(notifybecall)
	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_NOTIFY_BECALL
	jmsg.MsgData = string(notifyms)
	jmsgms, _ := json.Marshal(jmsg)
	fmt.Println("send notify msg is ", string(jmsgms))
	nw, err := becall.GetConn().Write(jmsgms)
	if err != nil {
		fmt.Println("write failed")
		return common.ErrWebSocketClosed
	}
	if nw == 0 {
		fmt.Println("peer connection closed ")
		return common.ErrWebSocketClosed
	}

	return nil
}

func UserCallReply(conn *websocket.Conn, msgdata string) error {
	fmt.Println("receive user call reply , msgdata is ", string(msgdata))
	cscall := &jsonproto.CSNotifyBeCall{}
	err := json.Unmarshal([]byte(msgdata), cscall)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return common.ErrJsonUnMarshal
	}

	fmt.Println("caller is ", cscall.Caller)
	fmt.Println("becalled is ", cscall.BeCalled)
	fmt.Println("agree is ", cscall.Agree)

	usrcaller, err := UserMgrInst.GetUser(cscall.Caller)
	if err != nil {
		return nil
	}
	if !usrcaller.IsOnline() {
		return nil
	}
	//被叫方不同意接听
	if !cscall.Agree {

		sccall := &jsonproto.SCUserCall{}
		sccall.Caller = cscall.Caller
		sccall.BeCalled = cscall.BeCalled
		sccall.ErrorId = common.RSP_USER_NOT_AGREE
		sccallms, _ := json.Marshal(sccall)
		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = string(sccallms)
		jmsgms, _ := json.Marshal(jmsg)
		nw, err := usrcaller.GetConn().Write(jmsgms)
		if err != nil {
			fmt.Println("write failed")
			return common.ErrWebSocketClosed
		}
		if nw == 0 {
			fmt.Println("peer connection closed ")
			return common.ErrWebSocketClosed
		}
		return nil
	}

	timestr := time.Now().Format("2006-01-02 15:04:05")
	tokenstr := fmt.Sprintf("%x", md5.Sum([]byte(cscall.Caller+cscall.BeCalled+timestr)))
	fmt.Println("token str is ", tokenstr)

	sccall := &jsonproto.SCUserCall{}
	sccall.Caller = cscall.Caller
	sccall.BeCalled = cscall.BeCalled
	sccall.ErrorId = common.RSP_SUCCESS
	sccall.Token = tokenstr
	sccallms, _ := json.Marshal(sccall)
	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_USER_CALL
	jmsg.MsgData = string(sccallms)
	jmsgms, _ := json.Marshal(jmsg)
	nw, err := usrcaller.GetConn().Write(jmsgms)
	if err != nil {
		fmt.Println("write failed")
		return common.ErrWebSocketClosed
	}
	if nw == 0 {
		fmt.Println("peer connection closed ")
		return common.ErrWebSocketClosed
	}
	return nil

}
