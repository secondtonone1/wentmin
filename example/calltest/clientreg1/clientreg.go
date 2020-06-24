package main

import (
	"fmt"
	"protobuf/proto"
	"wentmin/common"
	"wentmin/netmodel"
	wtproto "wentmin/proto"
	"wentmin/protocol"
)

func main() {

	cs, err := netmodel.Dial("tcp4", "127.0.0.1:9902")
	if err != nil {
		return
	}
	packet := new(protocol.MsgPacket)
	packet.Head.Id = common.USER_REG_CS
	csusereg := &wtproto.CSUserLogin{
		Accountid: "101",
		Passwd:    "pawd101",
		Phone:     "15110024987",
	}

	//protobuf编码
	pData, err := proto.Marshal(csusereg)
	if err != nil {
		fmt.Println(common.ErrProtobuffMarshal.Error())
		return
	}
	packet.Head.Len = (uint16)(len(pData))
	packet.Body.Data = pData
	cs.Send(packet)
	packetrsp, err := cs.Recv()
	if err != nil {
		fmt.Println("receive error")
		return
	}

	datarsp := packetrsp.(*protocol.MsgPacket)
	fmt.Println("packet id is", datarsp.Head.Id)
	fmt.Println("packet len is", datarsp.Head.Len)
	scusereg := &wtproto.SCUserLogin{}

	error2 := proto.Unmarshal(datarsp.Body.Data, scusereg)
	if error2 != nil {
		fmt.Println(common.ErrProtobuffUnMarshal.Error())
		return
	}

	if scusereg.Errid != common.RSP_SUCCESS {
		fmt.Println("user reg failed ")
		return
	}

	fmt.Println("user reg success ")
	fmt.Println("user token is ", scusereg.Token)

	//开始呼叫 102
	packetcall := new(protocol.MsgPacket)
	packetcall.Head.Id = common.CS_USER_CALL
	csusercall := &wtproto.CSUserCall{
		Caller:   "101",
		Becalled: "102",
	}

	//protobuf编码
	pDatacall, err := proto.Marshal(csusercall)
	if err != nil {
		fmt.Println(common.ErrProtobuffMarshal.Error())
		return
	}
	packetcall.Head.Len = (uint16)(len(pDatacall))
	packetcall.Body.Data = pDatacall
	cs.Send(packetcall)

	callrspt, err := cs.Recv()
	if err != nil {
		fmt.Println("receive error")
		return
	}

	callrsp := callrspt.(*protocol.MsgPacket)
	scusercall := &wtproto.SCUserCall{}

	error2 = proto.Unmarshal(callrsp.Body.Data, scusercall)
	if error2 != nil {
		fmt.Println(common.ErrProtobuffUnMarshal.Error())
		return
	}

	if scusercall.Errid != common.RSP_SUCCESS {
		fmt.Println("user call res is ", scusercall.Errid)
		fmt.Println("user call failed ")
		return
	}

	fmt.Println("callrsp.Head.Id is ", callrsp.Head.Id)

	//等待服务器通知
	notifyrsps, err := cs.Recv()
	if err != nil {
		fmt.Println("receive notifyrsp failed ")
		return
	}
	fmt.Println("..........................")
	notifyrsp := notifyrsps.(*protocol.MsgPacket)
	notifychat := &wtproto.SCNotifyChat{}
	error2 = proto.Unmarshal(notifyrsp.Body.Data, notifychat)
	if error2 != nil {
		fmt.Println(common.ErrProtobuffUnMarshal.Error())
		return
	}

	fmt.Println("notify chat caller is ", notifychat.Caller)
	fmt.Println("notify chat becalled is ", notifychat.Becalled)

	terminateCall := &wtproto.CSTerminateChat{}
	terminateCall.Roomid = notifychat.Roomid
	terminateCall.Caller = notifychat.Caller
	terminateCall.Becalled = notifychat.Becalled

	cstermcall := &protocol.MsgPacket{}
	cstermcall.Head.Id = common.CS_TERMINAL_CHAT
	cstermcall.Body.Data, _ = proto.Marshal(terminateCall)
	cstermcall.Head.Len = uint16(len(cstermcall.Body.Data))
	cs.Send(cstermcall)

	terminalrt := &wtproto.SCTerminateChat{}
	it, err := cs.Recv()
	sctermcall := it.(*protocol.MsgPacket)

	proto.Unmarshal(sctermcall.Body.Data, terminalrt)
	fmt.Println("receive terminate reply")
	fmt.Println("terminalrt.Caller", terminalrt.Caller)
	fmt.Println("terminalrt.Becalled", terminalrt.Becalled)
	cs.Close()
}
