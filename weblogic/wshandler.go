package weblogic

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
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
	RegMsgHandler(common.WEB_CS_USER_REG, UserLogin)
	RegMsgHandler(common.WEB_CS_USER_CALL, UserCall)
	RegMsgHandler(common.WEB_REPLY_BECALL, UserCallReply)
	RegMsgHandler(common.WEB_CS_TERMINAL_CALL, UserTerminalCall)
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

func UserLogin(conn *websocket.Conn, msgdata string) error {
	fmt.Println("receive user reg req , msgdata is ", string(msgdata))
	csreg := &jsonproto.CSUserLogin{}
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

	screg := &jsonproto.SCUserLogin{}
	screg.ErrorId = common.RSP_SUCCESS
	timestr := time.Now().Format("20060102150405")
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	randres := rand.Intn(100)
	screg.Token = csreg.AccountId + "-" + timestr + "-" + strconv.Itoa(randres)

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
	fmt.Println("isaudioonly is ", cscall.IsAudioOnly)
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

	//被呼叫人在线，回复主叫人唤起响铃
	callring := &jsonproto.SCNotifyCallRing{}
	callring.Caller = cscall.Caller
	callring.BeCalled = cscall.BeCalled

	rand.Seed(time.Now().UnixNano())
	randres := rand.Intn(100)

	timestr := time.Now().Format("20060102150405")
	roomidstr := cscall.Caller + "-" + cscall.BeCalled + "-" + timestr + "-" + strconv.Itoa(randres)

	fmt.Println("roomid str is ", roomidstr)
	callring.Roomid = roomidstr
	ringms, _ := json.Marshal(callring)
	ringmsg := &jsonproto.JsonMsg{}
	ringmsg.MsgId = common.WEB_NOTIFY_CALLRING
	ringmsg.MsgData = string(ringms)
	ringmsgs, _ := json.Marshal(ringmsg)

	fmt.Println("send call ring msg is ", string(ringmsgs))
	ringnw, err := conn.Write(ringmsgs)
	if err != nil {
		fmt.Println("write failed")
		return common.ErrWebSocketClosed
	}
	if ringnw == 0 {
		fmt.Println("peer connection closed ")
		return common.ErrWebSocketClosed
	}

	//将两个人放入房间
	chatroot := &ChatRoom{Roomid: roomidstr, Caller: cscall.Caller, Becalled: cscall.BeCalled}
	ChatMgrInst.AddRoom(chatroot)

	//被呼叫人在线，通知被呼叫人

	notifybecall := &jsonproto.SCNotifyBeCall{}
	notifybecall.Caller = cscall.Caller
	notifybecall.BeCalled = cscall.BeCalled
	notifybecall.Roomid = roomidstr
	notifybecall.IsAudioOnly = cscall.IsAudioOnly

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
	fmt.Println("roomid is ", cscall.Roomid)

	usrcaller, err := UserMgrInst.GetUser(cscall.Caller)
	if err != nil {
		return nil
	}
	if !usrcaller.IsOnline() {
		return nil
	}
	//被叫方不同意接听
	if !cscall.Agree {

		//将用户从房间中
		ChatMgrInst.DelRoom(cscall.Roomid)
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

	//判断房间信息是否存在，可能主叫方此时已经挂断
	room := ChatMgrInst.GetRoom(cscall.Roomid)
	if room == nil {
		fmt.Println("peer terminal call")
		return nil
	}

	sccall := &jsonproto.SCUserCall{}
	sccall.Caller = cscall.Caller
	sccall.BeCalled = cscall.BeCalled
	sccall.ErrorId = common.RSP_SUCCESS
	sccall.Roomid = cscall.Roomid
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

func UserTerminalCall(conn *websocket.Conn, msgdata string) error {
	fmt.Println("receive user terminal call  , msgdata is ", string(msgdata))
	cstcall := &jsonproto.CSTerminalCall{}
	err := json.Unmarshal([]byte(msgdata), cstcall)
	if err != nil {
		fmt.Println("json unmarshal failed")
		return common.ErrJsonUnMarshal
	}

	fmt.Println("caller is ", cstcall.Caller)
	fmt.Println("becalled is ", cstcall.BeCalled)
	fmt.Println("roomid is ", cstcall.Roomid)
	fmt.Println("cancel is ", cstcall.Cancel)

	//清除房间信息，
	room := ChatMgrInst.GetRoom(cstcall.Roomid)
	if room == nil {
		fmt.Println("not found room data, may be dismissed")
	} else {
		ChatMgrInst.DelRoom(cstcall.Roomid)
	}

	//被通知的另一方，另一方被挂断,假设是被呼叫人
	beterminated := cstcall.BeCalled
	//如果是被呼叫人挂断，将被中断的一方设置为呼叫人
	if cstcall.Cancel == cstcall.BeCalled {
		beterminated = cstcall.Caller
	}

	beterminal, err := UserMgrInst.GetUser(beterminated)
	if err != nil {
		fmt.Println("not found beterminal ")
		return nil
	}

	if !beterminal.IsOnline() {
		fmt.Println("beterminal is not online")
		return nil
	}

	beterminalConn := beterminal.GetConn()

	//同时发送消息推送给另一方挂断信息
	scbecall := &jsonproto.SCTerminalBeCall{}
	scbecall.BeCalled = cstcall.BeCalled
	scbecall.Caller = cstcall.Caller
	scbecall.Roomid = cstcall.Roomid
	scbecall.Cancel = cstcall.Cancel

	scbecallms, _ := json.Marshal(scbecall)

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
	jmsg.MsgData = string(scbecallms)
	jmsgms, _ := json.Marshal(jmsg)
	nw, err := beterminalConn.Write(jmsgms)
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
