package logic

import (
	"crypto/md5"
	"fmt"
	"protobuf/proto"
	"time"
	"wentmin/common"
	"wentmin/netmodel"
	wtproto "wentmin/proto"
	"wentmin/protocol"

	"github.com/astaxie/beego/logs"
)

func init() {
	fmt.Println("reg user req")
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.USER_REG_CS, UserReg)
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.CS_USER_CALL, UserCall)
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.CS_NOTIFY_REPLY, BeCallReply)
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.CS_TERMINAL_CHAT, TerminateChat)
}

func UserReg(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Println("receive user reg msg")
	userreg := &wtproto.CSUserReg{}
	err := proto.Unmarshal(msgpkg.Body.Data, userreg)
	if err != nil {
		fmt.Println("userreg proto unmarshal failed")
		return err
	}

	fmt.Println("user account id is ", userreg.Accountid)
	fmt.Println("user passwd is ", userreg.Passwd)
	fmt.Println("user phone is ", userreg.Phone)
	ud := &UserData{AccountId: userreg.Accountid, Passwd: userreg.Passwd,
		Phone: userreg.Phone, Session: session}
	UserMgrInst.AddUser(ud)
	userregrsp := &wtproto.SCUserReg{}
	userregrsp.Errid = common.RSP_SUCCESS
	userregrsp.Passwd = userreg.Passwd

	userregrsp.Accountid = userreg.Accountid
	userregrsp.Phone = userreg.Phone
	pData, err := proto.Marshal(userregrsp)
	if err != nil {
		fmt.Println(common.ErrProtobuffMarshal.Error())
		return common.ErrProtobuffMarshal
	}

	msgpkg.Head.Id = common.USER_REG_SC
	msgpkg.Head.Len = uint16(len(pData))
	msgpkg.Body.Data = pData
	netmodel.PostMsgOut(session, msgpkg)

	return nil
}

func UserCall(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Println("receive user call msg")
	uc := &wtproto.CSUserCall{}
	err := proto.Unmarshal(msgpkg.Body.Data, uc)
	if err != nil {
		logs.Debug("userreg proto unmarshal failed")
		return err
	}

	fmt.Println("caller is ", uc.Caller)
	fmt.Println("becalled is ", uc.Becalled)

	ud, err := UserMgrInst.GetUser(uc.Becalled)
	//user account id 没找到
	if err != nil {
		fmt.Printf("user not found %s\n", uc.Becalled)
		scusercall := &wtproto.SCUserCall{}
		scusercall.Errid = common.RSP_USER_NOT_FOUND
		msgrsp := &protocol.MsgPacket{}
		msgrsp.Head.Id = common.SC_USER_CALL
		msgrsp.Body.Data, _ = proto.Marshal(scusercall)
		msgrsp.Head.Len = uint16(len(msgrsp.Body.Data))
		netmodel.PostMsgOut(session, msgrsp)
		return nil
	}

	//判断对方是否在线
	if !ud.IsOnline() {
		//logs.Debug("user [%s] not online ", ud.AccountId)
		logs.Debug("user [%s] not online ", uc.Becalled)
		scusercall := &wtproto.SCUserCall{}
		scusercall.Errid = common.RSP_USER_NOT_ONLINE
		scusercall.Phone = ud.Phone
		msgrsp := &protocol.MsgPacket{}
		msgrsp.Head.Id = common.SC_USER_CALL
		msgrsp.Body.Data, _ = proto.Marshal(scusercall)
		msgrsp.Head.Len = uint16(len(msgrsp.Body.Data))
		netmodel.PostMsgOut(session, msgrsp)
		return nil
	}

	//发送会话通知给被呼叫方
	notifyBc := &wtproto.SCNotifyBeCalled{Caller: uc.Caller, Becalled: uc.Becalled}
	notifyms, _ := proto.Marshal(notifyBc)
	msgnotify := &protocol.MsgPacket{}
	msgnotify.Head.Id = common.SC_NOTIFY_BECALL
	msgnotify.Head.Len = uint16(len(notifyms))
	msgnotify.Body.Data = notifyms
	netmodel.PostMsgOut(ud.GetSession(), msgnotify)

	return nil
}

func BeCallReply(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Println("receive becall reply  msg")

	csreply := &wtproto.CSReplyBeCalled{}
	err := proto.Unmarshal(msgpkg.Body.Data, csreply)
	if err != nil {
		return common.ErrProtobuffUnMarshal
	}
	//获取呼叫人
	ud, err := UserMgrInst.GetUser(csreply.Caller)
	//不做任何处理，不通知呼叫方结果
	if err != nil {
		return nil
	}

	//判断呼叫方在线
	if online := ud.IsOnline(); !online {
		return nil
	}

	//通知呼叫方呼叫结果(是否同意)
	if !csreply.Agree {
		scusercall := &wtproto.SCUserCall{}
		scusercall.Errid = common.RSP_USER_NOT_AGREE
		scusercall.Caller = csreply.Caller
		scusercall.Becalled = csreply.Becalled
		msgrsp := &protocol.MsgPacket{}
		msgrsp.Head.Id = common.SC_USER_CALL
		msgrsp.Body.Data, _ = proto.Marshal(scusercall)
		msgrsp.Head.Len = uint16(len(msgrsp.Body.Data))
		netmodel.PostMsgOut(ud.GetSession(), msgrsp)
		return nil
	}

	//被呼叫方同意，先给呼叫方回包，告诉他呼叫结果
	scusercall := &wtproto.SCUserCall{}
	scusercall.Errid = common.RSP_SUCCESS
	scusercall.Caller = csreply.Caller
	scusercall.Becalled = csreply.Becalled
	msgrsp := &protocol.MsgPacket{}
	msgrsp.Head.Id = common.SC_USER_CALL
	msgrsp.Body.Data, _ = proto.Marshal(scusercall)
	msgrsp.Head.Len = uint16(len(msgrsp.Body.Data))
	netmodel.PostMsgOut(ud.GetSession(), msgrsp)

	//被呼叫方同意，则通知双方建立会话消息
	scnotify := &wtproto.SCNotifyChat{}
	scnotify.Caller = csreply.Caller
	scnotify.Becalled = csreply.Becalled

	timestr := time.Now().Format("2006-01-02 15:04:05")
	tokenstr := fmt.Sprintf("%x", md5.Sum([]byte(csreply.Caller+csreply.Becalled+timestr)))
	fmt.Println("token str is ", tokenstr)
	//logs.Debug("token str is ", tokenstr)
	scnotify.Token = tokenstr

	notifyChat := &protocol.MsgPacket{}
	notifyChat.Head.Id = common.SC_NOTIFY_CHAT
	notifyChat.Body.Data, _ = proto.Marshal(scnotify)
	notifyChat.Head.Len = uint16(len(notifyChat.Body.Data))

	netmodel.PostMsgOut(ud.GetSession(), notifyChat)
	netmodel.PostMsgOut(session, notifyChat)
	cr := new(ChatRoom)
	cr.Token = tokenstr
	cr.Caller = csreply.Caller
	cr.Becalled = csreply.Becalled
	ChatMgrInst.AddRoom(cr)
	return nil
}

func TerminateChat(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Println("receive TerminateChat  msg")
	cstermchat := &wtproto.CSTerminateChat{}
	err := proto.Unmarshal(msgpkg.Body.Data, cstermchat)
	if err != nil {
		return common.ErrProtobuffUnMarshal
	}

	ChatMgrInst.DelRoom(cstermchat.Token)
	sctermchat := &wtproto.SCTerminateChat{}
	sctermchat.Errid = common.RSP_SUCCESS
	sctermchat.Caller = cstermchat.Caller
	sctermchat.Becalled = cstermchat.Becalled
	sctermchat.Token = cstermchat.Token

	sctermpkg := &protocol.MsgPacket{}
	sctermpkg.Head.Id = common.SC_TERMINAL_CHAT
	sctermpkg.Body.Data, _ = proto.Marshal(sctermchat)
	sctermpkg.Head.Len = uint16(len(sctermpkg.Body.Data))
	netmodel.PostMsgOut(session, sctermpkg)
	return nil
}
