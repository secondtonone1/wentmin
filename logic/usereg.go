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
}

func UserReg(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Printf("socket id [%d], receive userreg msg \n", session.GetSocketId())
	logs.Debug("socket id [%d], receive userreg msg \n", session.GetSocketId())
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

	userregrsp := &wtproto.SCUserReg{}
	userregrsp.Errid = common.RSP_SUCCESS
	userregrsp.Passwd = userreg.Passwd

	timestr := time.Now().Format("2006-01-02 15:04:05")

	tokenstr := fmt.Sprintf("%x", md5.Sum([]byte(userreg.Accountid+timestr)))
	fmt.Println("token str is ", tokenstr)
	logs.Debug("token str is ", tokenstr)
	userregrsp.Token = tokenstr
	userregrsp.Accountid = userreg.Accountid

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
