package logic

import (
	"fmt"
	"wentmin/common"
	"wentmin/netmodel"
	"wentmin/protocol"
)

func init() {
	fmt.Println("reg user req")
	netmodel.GetMsgHandlerIns().RegMsgHandler(common.SYC_CON_CLOSED, OnSessionClosed)
}

func OnSessionClosed(session *netmodel.Session, msgpkg *protocol.MsgPacket) error {
	sid := session.GetSocketId()
	lsid := session.GetLastSocket()
	fmt.Printf("socket id [%d], lastsocketid [%d], session closed \n", sid, lsid)
	//logs.Debug("socket id [%d], lastsocketid [%d], session closed \n", sid, lsid)
	//做应用层连接断开处理
	UserMgrInst.SetUserOffline(lsid)
	return nil
}
