package netmodel

import (
	"fmt"
	"sync"
	"wentmin/components"
)

type MsgQueue struct {
	lock sync.Mutex
	once sync.Once
}

var MsgQueueInst *MsgQueue
var MsgQueueClose chan struct{}
var MsgWatiGroup sync.WaitGroup

func NewMsgQueue() {
	MsgQueueInst = &MsgQueue{}
	MsgQueueClose = make(chan struct{}, components.MaxMsgQueNum)
}

func (mq *MsgQueue) CloseMsgQue() {
	mq.lock.Lock()
	defer mq.lock.Unlock()
	fmt.Println("MsgQueue exit !")
	MsgQueueClose <- struct{}{}
}

func (mq *MsgQueue) ReadFromChan() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("msg queue recover from error ", err)
		}
		mq.CloseMsgQue()
	}()
	for {
		select {
		case <-AcceptClose:
			return
		case _ = <-PacketChan:

		}
	}

}

func (mq *MsgQueue) WaitClose() chan struct{} {
	return MsgQueueClose
}
