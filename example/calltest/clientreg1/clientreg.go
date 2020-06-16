package main

import (
	"fmt"
	"protobuf/proto"
	"time"
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
	csusereg := &wtproto.CSUserReg{
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
	scusereg := &wtproto.SCUserReg{}

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
	fmt.Println("user account is ", scusereg.Accountid)
	fmt.Println("user passwd is ", scusereg.Passwd)
	fmt.Println("user phone is ", scusereg.Phone)

	time.Sleep(time.Second * 2)

	//开始呼叫 102
	packetcall := new(protocol.MsgPacket)
	packetcall.Head.Id = common.CS_USER_CALL
	csusercall := &wtproto.CSUserCall{
		Caller:"101",
		Becalled:"102"
	}

	//protobuf编码
	pDatacall, err := proto.Marshal(csusercall)
	if err != nil {
		fmt.Println(common.ErrProtobuffMarshal.Error())
		return
	}
	packet.Head.Len = (uint16)(len(pDatacall))
	packet.Body.Data = pDatacall
	cs.Send(packet)
	callrsp, err := cs.Recv()
	if err != nil {
		fmt.Println("receive error")
		return
	}
	

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

	

	cs.Close()
}
