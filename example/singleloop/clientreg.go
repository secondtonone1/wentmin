package main

import (
	"fmt"
	"protobuf/proto"
	"strconv"
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

	for i := 0; i < 1000; i++ {
		accountid := strconv.Itoa(i)

		packet := new(protocol.MsgPacket)
		packet.Head.Id = common.USER_REG_CS
		csusereg := &wtproto.CSUserReg{
			Accountid: accountid,
			Passwd:    "pawd" + accountid,
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
	}

	for i := 0; i < 2050; i++ {
		accountid := strconv.Itoa(i)

		packet := new(protocol.MsgPacket)
		packet.Head.Id = common.USER_REG_CS
		csusereg := &wtproto.CSUserReg{
			Accountid: accountid,
			Passwd:    "pawd" + accountid,
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

	}

	for i := 0; i < 2050; i++ {
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
		fmt.Println("user token is ", scusereg.Token)
	}

}
