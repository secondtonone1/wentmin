package logic

import (
	"fmt"
	"protobuf/proto"
	"wentmin/common"
	"wentmin/netmodel"
	wtproto "wentmin/proto"
	"wentmin/protocol"
)

func init() {
	fmt.Println("reg user req")
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.USER_REG_CS, netmodel.CallBackFunc(UserReg))
}

func UserReg(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Printf("socket id [%d], receive userreg msg \n", session.GetSocketId())
	userreg := &wtproto.CSUserReg{}
	err := proto.Unmarshal(msgpkg.Body.Data, userreg)
	if err != nil {
		fmt.Println("userreg proto unmarshal failed")
		return err
	}

	fmt.Println("user account id is ", userreg.Accountid)
	fmt.Println("user passwd is ", userreg.Passwd)
	return nil
}
