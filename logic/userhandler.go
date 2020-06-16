package logic

import (
	"fmt"
	"protobuf/proto"
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
}

func UserReg(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	//fmt.Printf("socket id [%d], receive userreg msg \n", session.GetSocketId())
	//logs.Debug("socket id [%d], receive userreg msg \n", session.GetSocketId())
	userreg := &wtproto.CSUserReg{}
	err := proto.Unmarshal(msgpkg.Body.Data, userreg)
	if err != nil {
		fmt.Println("userreg proto unmarshal failed")
		return err
	}

	fmt.Println("user account id is ", userreg.Accountid)
	fmt.Println("user passwd is ", userreg.Passwd)
	logs.Debug("user account id is ", userreg.Accountid)
	logs.Debug("user passwd is ", userreg.Passwd)
	logs.Debug("user phone is ", userreg.Phone)
	fmt.Println("user phone is ", userreg.Phone)
	ud := &UserData{AccountId: userreg.Accountid, Passwd: userreg.Passwd,
		Phone: userreg.Phone, Session: session}
	UserMgrInst.AddUser(ud)
	userregrsp := &wtproto.SCUserReg{}
	userregrsp.Errid = common.RSP_SUCCESS
	userregrsp.Passwd = userreg.Passwd
	/*
		timestr := time.Now().Format("2006-01-02 15:04:05")

		tokenstr := fmt.Sprintf("%x", md5.Sum([]byte(userreg.Accountid+timestr)))
		fmt.Println("token str is ", tokenstr)
		logs.Debug("token str is ", tokenstr)
		userregrsp.Token = tokenstr
	*/
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
	//fmt.Printf("socket id [%d], receive UserCall msg \n", session.GetSocketId())
	//logs.Debug("socket id [%d], receive UserCall msg \n", session.GetSocketId())
	uc := &wtproto.CSUserCall{}
	err := proto.Unmarshal(msgpkg.Body.Data, uc)
	if err != nil {
		logs.Debug("userreg proto unmarshal failed")
		return err
	}

	logs.Debug("caller is ", uc.Caller)
	logs.Debug("becalled is ", uc.Becalled)

	scusercall := &wtproto.SCUserCall{}
	defer func() {
		msgrsp := &protocol.MsgPacket{}
		msgrsp.Head.Id = common.SC_USER_CALL
		msgrsp.Body.Data, _ = proto.Marshal(scusercall)
		msgrsp.Head.Len = uint16(len(msgrsp.Body.Data))
	}()

	//account id 不准确
	ud, err := UserMgrInst.GetUser(uc.Becalled)
	if err != nil {
		logs.Debug("user not found %s", ud.AccountId)
		scusercall.Errid = common.RSP_USER_NOT_FOUND
		return nil
	}

	//判断对方是否在线
	if !ud.IsOnline() {
		logs.Debug("user [%s] not online ", ud.AccountId)
		scusercall.Errid = common.RSP_USER_NOT_ONLINE
		scusercall.Phone = ud.Phone
		return nil
	}

	scusercall.Errid = common.RSP_SUCCESS
	scusercall.Phone = ud.Phone
	scusercall.Caller = uc.Caller
	scusercall.Becalled = uc.Becalled

	return nil
}
