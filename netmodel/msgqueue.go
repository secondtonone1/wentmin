package netmodel

import (
	"fmt"
	"sync"
	"wentmin/components"
)

type MsgQueue struct {
	lock sync.Mutex
}

var MsgQueueInst *MsgQueue
var MsgQueueClose chan struct{}
var MsgWatiGroup sync.WaitGroup

func NewMsgQueue() {
	MsgQueueInst = &MsgQueue{}
	MsgQueueClose = make(chan struct{}, components.MaxMsgQueNum)
	for i := 0; i < components.MaxMsgQueNum; i++ {
		go MsgQueueInst.ReadFromChan()
		MsgWatiGroup.Add(1)
	}
}

func (mq *MsgQueue) OnClose() {
	fmt.Println("MsgQueue exit !")
	MsgQueueClose <- struct{}{}
	MsgWatiGroup.Done()
}

func (mq *MsgQueue) PutMsgInQue(packet interface{}) error {
	PacketChan <- packet.(*MsgSession)
	return nil
}

func (mq *MsgQueue) ReadFromChan() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("msg queue recover from error ", err)
		}
		mq.OnClose()
	}()
	for {
		select {
		case <-AcceptClose:
			return
		case msgs := <-PacketChan:
			//多消息队列，多协程处理模式
			if components.MaxMsgQueNum > 1 {
				mq.lock.Lock()
				defer mq.lock.Unlock()
				//handle msg packet
				hdres := GetMsgHandlerIns().HandleMsgPacket(msgs.packet, msgs.session)
				if hdres != nil {
					fmt.Printf("handle msg failed, msgid is %d, error is %s\n",
						msgs.packet.Head.Id, hdres.Error())
					return
				}
				continue
			}
			//单协程模式
			hdres := GetMsgHandlerIns().HandleMsgPacket(msgs.packet, msgs.session)
			if hdres != nil {
				fmt.Printf("handle msg failed, msgid is %d, error is %s\n",
					msgs.packet.Head.Id, hdres.Error())
				return
			}
			continue
		}
	}

}

func (mq *MsgQueue) WaitClose() chan struct{} {
	return MsgQueueClose
}
