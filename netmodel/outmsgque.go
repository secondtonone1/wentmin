package netmodel

import (
	"fmt"
	"sync"
	"wentmin/components"
	"wentmin/protocol"

	"github.com/astaxie/beego/logs"
)

var OutMsgQueInst *OutMsgQue
var OutputWaitGroup sync.WaitGroup

type OutMsgQue struct {
	MsgChanMap  map[int]chan *MsgSession
	CloseNotify chan struct{}
	RWLock      sync.RWMutex //用来控制map的互斥访问
}

func (oq *OutMsgQue) GetMsgChanByIndex(index int) chan *MsgSession {
	oq.RWLock.RLock()
	defer oq.RWLock.RUnlock()
	msgchan, ok := oq.MsgChanMap[index]
	if !ok {
		fmt.Println("not found output msgchan by index ", index)
		logs.Debug("not found output msgchan by index ", index)
		return nil
	}
	return msgchan
}

//连接关闭通知server退出
func (oq *OutMsgQue) OnGoroutinClose() {
	if err := recover(); err != nil {
		fmt.Println("out msg queue goroutine recover from error ", err)
		logs.Debug("out msg queue goroutine recover from error ", err)
	}
	oq.CloseNotify <- struct{}{}
	OutputWaitGroup.Done()
	fmt.Println("out msg queue goroutine exited ")
	logs.Debug("out msg queue goroutine exited ")
}

func (oq *OutMsgQue) ReadFromOutQue(index int) {
	msgchan := oq.GetMsgChanByIndex(index)
	if msgchan == nil {
		return
	}
	defer oq.OnGoroutinClose()
	for {
		select {
		case msgs := <-msgchan:
			//考虑消息发送
			msgs.session.Write(msgs.packet)
		case <-AcceptClose:
			//考虑处理server退出逻辑
			return
		}
	}
}

func (oq *OutMsgQue) PostMsgtoOutQue(sess *Session, msgpkg *protocol.MsgPacket) {

	msgse := new(MsgSession)
	msgse.session = sess
	msgse.packet = msgpkg

	msgchan := oq.GetMsgChanByIndex(msgse.session.GetSocketId() % components.OutputQueNum)
	if msgchan == nil {
		return
	}
	msgchan <- msgse
}

func (oq *OutMsgQue) WaitClose() chan struct{} {
	return oq.CloseNotify
}

func NewOutMsgQues() {
	OutMsgQueInst = new(OutMsgQue)
	OutMsgQueInst.MsgChanMap = make(map[int]chan *MsgSession)
	//根据输出队列数量，创建chan写入map，每个chan对应一类socket发送
	//从而实现非锁并发
	for i := 0; i < components.OutputQueNum; i++ {
		msgchan := make(chan *MsgSession, components.OutputQueLen)
		OutMsgQueInst.MsgChanMap[i] = msgchan
	}

	//根据输出队列数量，创建关闭回传的通知chan大小
	OutMsgQueInst.CloseNotify = make(chan struct{}, components.OutputQueNum)

	//启动n个协程并发处理输出队列消息
	for i := 0; i < components.OutputQueNum; i++ {
		go OutMsgQueInst.ReadFromOutQue(i)
		OutputWaitGroup.Add(1)
	}
}

func PostMsgOut(sess *Session, msgpkg *protocol.MsgPacket) {
	OutMsgQueInst.PostMsgtoOutQue(sess, msgpkg)
}
