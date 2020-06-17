package netmodel

import (
	"fmt"
	"sync"
	"wentmin/components"

	"wentmin/protocol"

	"wentmin/common"

	"github.com/astaxie/beego/logs"
)

type MsgQueue struct {
	lock    sync.RWMutex
	running bool
}

var MsgQueueInst *MsgQueue
var MsgQueueClose chan struct{}
var MsgWatiGroup sync.WaitGroup

func NewMsgQueue() {
	MsgQueueInst = &MsgQueue{running: false}
	MsgQueueClose = make(chan struct{}, components.MaxMsgQueNum)
	for i := 0; i < components.MaxMsgQueNum; i++ {
		go MsgQueueInst.ReadFromChan()
		MsgWatiGroup.Add(1)
	}
}

func (mq *MsgQueue) SetStop() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.running == false {
		return
	}
	mq.running = false
}

func (mq *MsgQueue) StartRun() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	if mq.running == true {
		return
	}
	mq.running = true
}

func (mq *MsgQueue) IsRunning() bool {
	mq.lock.RLock()
	defer mq.lock.RUnlock()
	return mq.running == true
}

func (mq *MsgQueue) OnClose() {
	fmt.Println("MsgQueue exit !")
	logs.Debug("MsgQueue exit !")
	mq.SetStop()
	MsgQueueClose <- struct{}{}
	MsgWatiGroup.Done()
}

func (mq *MsgQueue) PutMsgInQue(packet interface{}) error {
	PacketChan <- packet.(*MsgSession)
	return nil
}

func (mq *MsgQueue) PutCloseMsgInQue(session *Session) error {
	if mq.IsRunning() == false {
		fmt.Println("msg queue has been closed")
		logs.Debug("msg queue has been closed")
		return nil
	}

	pkg := new(protocol.MsgPacket)
	pkg.Head.Id = common.SYC_CON_CLOSED
	pkg.Head.Len = 0
	pkg.Body.Data = []byte{}
	msgs := &MsgSession{session, pkg}
	PacketChan <- msgs
	//logs.Debug("PutCloseMsgInQue")
	return nil
}

func (mq *MsgQueue) ReadFromChan() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("msg queue recover from error ", err)
			logs.Debug("msg queue recover from error ", err)
		}
		mq.OnClose()
	}()
	mq.StartRun()
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
					logs.Debug("handle msg failed, msgid is %d, error is %s\n",
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
				logs.Debug("handle msg failed, msgid is %d, error is %s\n",
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
