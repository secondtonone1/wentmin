package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
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
		Accountid: "102",
		Passwd:    "pawd102",
		Phone:     "18301152098",
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

	//等待被呼叫
	scnbca, _ := cs.Recv()
	notifyBeCall := scnbca.(*protocol.MsgPacket)
	notifybc := &wtproto.SCNotifyBeCalled{}
	proto.Unmarshal(notifyBeCall.Body.Data, notifybc)

	fmt.Println(" becalled is ", notifybc.Becalled)
	fmt.Println(" caller is ", notifybc.Caller)

	//假设同意,发送同意请求
	replybc := &wtproto.CSReplyBeCalled{}
	replybc.Caller = notifybc.Caller
	replybc.Becalled = notifybc.Becalled
	replybc.Agree = true

	replybcp := &protocol.MsgPacket{}
	replybcp.Head.Id = common.CS_NOTIFY_REPLY
	replybcp.Body.Data, _ = proto.Marshal(replybc)
	replybcp.Head.Len = uint16(len(replybcp.Body.Data))
	cs.Send(replybcp)

	//等待服务器通知
	notifyrsps, _ := cs.Recv()

	notifyrsp := notifyrsps.(*protocol.MsgPacket)
	notifychat := &wtproto.SCNotifyChat{}
	error2 = proto.Unmarshal(notifyrsp.Body.Data, notifychat)
	if error2 != nil {
		fmt.Println(common.ErrProtobuffUnMarshal.Error())
		return
	}

	fmt.Println("notify chat caller is ", notifychat.Caller)
	fmt.Println("notify chat becalled is ", notifychat.Becalled)

	cs.Close()
}
