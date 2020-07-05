package weblogic

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"wentmin/common"
	"wentmin/jsonproto"

	"github.com/goinggo/mapstructure"
	"golang.org/x/net/websocket"
)

func UserLogin(conn *websocket.Conn, msgdata interface{}) error {
	cslogin := &jsonproto.CSUserLogin{}
	if err := mapstructure.Decode(msgdata, cslogin); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSUserLogin data is :%v\n", cslogin)

	ud := &UserData{}
	ud.AccountId = cslogin.AccountId
	ud.Conn = conn
	ud.Passwd = cslogin.Passwd
	ud.Phone = cslogin.Phone
	UserMgrInst.AddUser(ud)

	screg := &jsonproto.SCUserLogin{}
	screg.ErrorId = common.RSP_SUCCESS
	timestr := time.Now().Format("20060102150405")
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	randres := rand.Intn(100)
	screg.Token = cslogin.AccountId + "-" + timestr + "-" + strconv.Itoa(randres)

	jsmsg := &jsonproto.JsonMsg{}
	jsmsg.MsgId = common.WEB_SC_USER_REG
	jsmsg.MsgData = screg

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

func UserCall(conn *websocket.Conn, msgdata interface{}) error {
	cscall := &jsonproto.CSUserCall{}
	if err := mapstructure.Decode(msgdata, cscall); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSUserCall data is :%v\n", cscall)

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
		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = sccall
		jmsgms, _ := json.Marshal(jmsg)
		return SendData(conn, jmsgms)
	}
	//用户不在线
	if !becall.IsOnline() {
		sccall := &jsonproto.SCUserCall{}
		sccall.Caller = cscall.Caller
		sccall.BeCalled = cscall.BeCalled
		sccall.ErrorId = common.RSP_USER_NOT_ONLINE
		sccall.Phone = becall.Phone

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = sccall
		jmsgms, _ := json.Marshal(jmsg)
		return SendData(conn, jmsgms)
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

	ringmsg := &jsonproto.JsonMsg{}
	ringmsg.MsgId = common.WEB_NOTIFY_CALLRING
	ringmsg.MsgData = callring
	ringmsgs, _ := json.Marshal(ringmsg)

	fmt.Println("send call ring msg is ", string(ringmsgs))
	err = SendData(conn, ringmsgs)
	if err != nil {
		fmt.Println("write failed")
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

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_NOTIFY_BECALL
	jmsg.MsgData = notifybecall
	jmsgms, _ := json.Marshal(jmsg)
	fmt.Println("send notify msg is ", string(jmsgms))
	return SendData(becall.GetConn(), jmsgms)
}

func UserCallReply(conn *websocket.Conn, msgdata interface{}) error {
	cscall := &jsonproto.CSNotifyBeCall{}
	if err := mapstructure.Decode(msgdata, cscall); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSNotifyBeCall data is :%v\n", cscall)

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

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_USER_CALL
		jmsg.MsgData = sccall
		jmsgms, _ := json.Marshal(jmsg)
		return SendData(usrcaller.GetConn(), jmsgms)
	}

	//判断房间信息是否存在，可能主叫方此时已经挂断
	//也可能是返回的房间号不正确
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

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_USER_CALL
	jmsg.MsgData = sccall
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

func UserTerminalCall(conn *websocket.Conn, msgdata interface{}) error {
	cstcall := &jsonproto.CSTerminalCall{}
	if err := mapstructure.Decode(msgdata, cstcall); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSTerminalCall data is :%v\n", cstcall)

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

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
	jmsg.MsgData = scbecall
	jmsgms, _ := json.Marshal(jmsg)
	return SendData(beterminalConn, jmsgms)

}

func SendData(conn *websocket.Conn, msgdata []byte) error {
	if conn == nil {
		return common.ErrConnInvalid
	}

	nw, err := conn.Write(msgdata)
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

//收到主叫方offer
func ReceiveOffer(conn *websocket.Conn, msgdata interface{}) error {
	cscalloffer := &jsonproto.CSCallOffer{}
	if err := mapstructure.Decode(msgdata, cscalloffer); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSCallOffer data is :%v\n", cscalloffer)

	fmt.Println("caller is ", cscalloffer.Caller)
	fmt.Println("becalled is ", cscalloffer.BeCalled)
	fmt.Println("roomid is ", cscalloffer.Roomid)
	fmt.Println("sdp is ", cscalloffer.Sdp)

	//判断被叫方是否存在
	becall, err := UserMgrInst.GetUser(cscalloffer.BeCalled)
	if err != nil {
		fmt.Printf("user %s not found \n", cscalloffer.BeCalled)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = cscalloffer.BeCalled
		scbecall.Caller = cscalloffer.Caller
		scbecall.Roomid = cscalloffer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//被叫不在线，通知主叫方挂断
	if !becall.IsOnline() {
		fmt.Printf("user %s is not online \n", cscalloffer.BeCalled)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = cscalloffer.BeCalled
		scbecall.Caller = cscalloffer.Caller
		scbecall.Roomid = cscalloffer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//判断房间是否解散
	roomdata := ChatMgrInst.GetRoom(cscalloffer.Roomid)
	//房间信息不存在，则通知主叫方挂断
	if roomdata == nil {
		fmt.Printf("room %s is not exist, maybe dismissed \n", roomdata)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = cscalloffer.BeCalled
		scbecall.Caller = cscalloffer.Caller
		scbecall.Roomid = cscalloffer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}
	//将offer信息转发给被叫方
	offernotify := &jsonproto.SCOfferNotify{}
	offernotify.BeCalled = cscalloffer.BeCalled
	offernotify.Caller = cscalloffer.Caller
	offernotify.Roomid = cscalloffer.Roomid
	offernotify.Sdp = cscalloffer.Sdp

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_OFFER_NOTIFY
	jmsg.MsgData = offernotify
	jmsgms, _ := json.Marshal(jmsg)
	return SendData(becall.GetConn(), jmsgms)
}

//收到被叫方offer
func ReceiveAnswer(conn *websocket.Conn, msgdata interface{}) error {
	csanswer := &jsonproto.CSBecallAnswer{}
	if err := mapstructure.Decode(msgdata, csanswer); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSBecallAnswer data is :%v\n", csanswer)

	fmt.Println("caller is ", csanswer.Caller)
	fmt.Println("becalled is ", csanswer.BeCalled)
	fmt.Println("roomid is ", csanswer.Roomid)
	fmt.Println("sdp is ", csanswer.Sdp)

	//判断主叫方是否存在
	caller, err := UserMgrInst.GetUser(csanswer.Caller)
	if err != nil {
		fmt.Printf("user id %s not found ", csanswer.Caller)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = csanswer.BeCalled
		scbecall.Caller = csanswer.Caller
		scbecall.Roomid = csanswer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//通知被叫方挂断
		return SendData(conn, jmsgms)
	}

	//判断主叫方是否在线
	if !caller.IsOnline() {
		fmt.Printf("caller %s is not online \n", csanswer.Caller)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = csanswer.BeCalled
		scbecall.Caller = csanswer.Caller
		scbecall.Roomid = csanswer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//通知被叫方挂断
		return SendData(conn, jmsgms)
	}

	//判断房间是否存在
	//判断房间是否解散
	roomdata := ChatMgrInst.GetRoom(csanswer.Roomid)
	//房间信息不存在，则通知主叫方挂断
	if roomdata == nil {
		fmt.Printf("room %s is not exist, maybe dismissed \n", roomdata)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = csanswer.BeCalled
		scbecall.Caller = csanswer.Caller
		scbecall.Roomid = csanswer.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//将answer转发给主叫方
	answernotify := &jsonproto.SCAnswerNotify{}

	answernotify.BeCalled = csanswer.BeCalled
	answernotify.Caller = csanswer.Caller
	answernotify.Roomid = csanswer.Roomid
	answernotify.Sdp = csanswer.Sdp

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_ANSWER_NOTIFY
	jmsg.MsgData = answernotify
	jmsgms, _ := json.Marshal(jmsg)
	return SendData(caller.GetConn(), jmsgms)
}

//收到主叫方ice_candidate
func ReceiveCallIce(conn *websocket.Conn, msgdata interface{}) error {

	callice := &jsonproto.CSCallIce{}
	if err := mapstructure.Decode(msgdata, callice); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSCallIce data is :%v\n", callice)

	fmt.Println("caller is ", callice.Caller)
	fmt.Println("becalled is ", callice.BeCalled)
	fmt.Println("roomid is ", callice.Roomid)
	fmt.Println("Candidate is ", callice.Candidate)

	//判断被叫方是否存在
	becall, err := UserMgrInst.GetUser(callice.BeCalled)
	if err != nil {
		fmt.Printf("user %s not found \n", callice.BeCalled)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = callice.BeCalled
		scbecall.Caller = callice.Caller
		scbecall.Roomid = callice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//被叫不在线，通知主叫方挂断
	if !becall.IsOnline() {
		fmt.Printf("user %s is not online \n", callice.BeCalled)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = callice.BeCalled
		scbecall.Caller = callice.Caller
		scbecall.Roomid = callice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//判断房间是否解散
	roomdata := ChatMgrInst.GetRoom(callice.Roomid)
	//房间信息不存在，则通知主叫方挂断
	if roomdata == nil {
		fmt.Printf("room %s is not exist, maybe dismissed \n", roomdata)
		//服务器通知主叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = callice.BeCalled
		scbecall.Caller = callice.Caller
		scbecall.Roomid = callice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}
	//将ice信息转发给被叫方
	ICEnotify := &jsonproto.SCCallIceNotify{}
	ICEnotify.BeCalled = callice.BeCalled
	ICEnotify.Caller = callice.Caller
	ICEnotify.Roomid = callice.Roomid
	ICEnotify.Candidate = callice.Candidate

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_CALL_ICE_NOTIFY
	jmsg.MsgData = ICEnotify
	jmsgms, _ := json.Marshal(jmsg)
	return SendData(becall.GetConn(), jmsgms)
}

//收到被叫方ice_candidate
func ReceiveBeCallIce(conn *websocket.Conn, msgdata interface{}) error {

	becallice := &jsonproto.CSBecallIce{}
	if err := mapstructure.Decode(msgdata, becallice); err != nil {
		fmt.Println("map to struct failed, err is ", err)
		return err
	}
	fmt.Printf("CSBecallIce data is :%v\n", becallice)

	fmt.Println("caller is ", becallice.Caller)
	fmt.Println("becalled is ", becallice.BeCalled)
	fmt.Println("roomid is ", becallice.Roomid)
	fmt.Println("candidate is ", becallice.Candidate)

	//判断主叫方是否存在
	caller, err := UserMgrInst.GetUser(becallice.Caller)
	if err != nil {
		fmt.Printf("user id %s not found ", becallice.Caller)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = becallice.BeCalled
		scbecall.Caller = becallice.Caller
		scbecall.Roomid = becallice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//通知被叫方挂断
		return SendData(conn, jmsgms)
	}

	//判断主叫方是否在线
	if !caller.IsOnline() {
		fmt.Printf("caller %s is not online \n", becallice.Caller)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = becallice.BeCalled
		scbecall.Caller = becallice.Caller
		scbecall.Roomid = becallice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//通知被叫方挂断
		return SendData(conn, jmsgms)
	}

	//判断房间是否存在
	//判断房间是否解散
	roomdata := ChatMgrInst.GetRoom(becallice.Roomid)
	//房间信息不存在，则通知主叫方挂断
	if roomdata == nil {
		fmt.Printf("room %s is not exist, maybe dismissed \n", roomdata)
		//服务器通知被叫方挂断
		scbecall := &jsonproto.SCTerminalBeCall{}
		scbecall.BeCalled = becallice.BeCalled
		scbecall.Caller = becallice.Caller
		scbecall.Roomid = becallice.Roomid
		scbecall.Cancel = "server"

		jmsg := &jsonproto.JsonMsg{}
		jmsg.MsgId = common.WEB_SC_TERMINAL_BECALL
		jmsg.MsgData = scbecall
		jmsgms, _ := json.Marshal(jmsg)
		//给主叫方推送挂断消息
		return SendData(conn, jmsgms)
	}

	//将answer转发给主叫方
	becallicenotify := &jsonproto.SCBecallIceNotify{}

	becallicenotify.BeCalled = becallice.BeCalled
	becallicenotify.Caller = becallice.Caller
	becallicenotify.Roomid = becallice.Roomid
	becallicenotify.Candidate = becallice.Candidate

	jmsg := &jsonproto.JsonMsg{}
	jmsg.MsgId = common.WEB_SC_BECALL_ICE_NOTIFY
	jmsg.MsgData = becallicenotify
	jmsgms, _ := json.Marshal(jmsg)
	return SendData(caller.GetConn(), jmsgms)
}
