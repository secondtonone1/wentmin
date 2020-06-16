package logic

import (
	"fmt"
	"wentmin/common"
	"wentmin/netmodel"
	"wentmin/protocol"

	"github.com/astaxie/beego/logs"
)

func init() {
	fmt.Println("reg user req")
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.SYC_CON_CLOSED, OnSessionClosed)
}

func OnSessionClosed(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	fmt.Printf("socket id [%d], lastsocketid [%d], session closed \n", session.GetSocketId(), session.GetLastSocket())
	logs.Debug("socket id [%d], lastsocketid [%d], session closed \n", session.GetSocketId(), session.GetLastSocket())
	//做应用层

	return nil
}
