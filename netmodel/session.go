package netmodel

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"wentmin/protocol"
)

type Session struct {
	conn        net.Conn
	closed      int32                  //session是否关闭，-1未开启，0未关闭，1关闭
	protocol    protocol.ProtocolInter //字节序和自己处理器
	lock        sync.Mutex             //协程锁
	SocketId    int
	closeNotify chan struct{} //当accept主动关闭session协程通知
}

func NewSession(connt net.Conn,
	soId int) *Session {
	sess := &Session{
		conn:        connt,
		closed:      -1,
		protocol:    new(protocol.ProtocolImpl),
		SocketId:    soId,
		closeNotify: make(chan struct{}),
	}
	tcpConn := sess.conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)
	tcpConn.SetReadBuffer(64 * 1024)
	tcpConn.SetWriteBuffer(64 * 1024)
	return sess
}

func (se *Session) RawConn() *net.TCPConn {
	return se.conn.(*net.TCPConn)
}

func (se *Session) Start() {
	if atomic.CompareAndSwapInt32(&se.closed, -1, 0) {
		go se.recvLoop()
	}
}

// Close the session, destory other resource.
func (se *Session) Close() error {
	if atomic.CompareAndSwapInt32(&se.closed, 0, 1) {
		se.conn.Close()
	}
	return nil
}

func (se *Session) recvLoop() {
	defer TcpServerInst.OnSessionClosed(se.SocketId)
	var packet interface{}
	var err error
	for {

		select {
		case <-se.closeNotify:
			return
		case <-AcceptClose:
			return
		default:
			{
				packet, err = se.protocol.ReadPacket(se.conn)
				if packet == nil || err != nil {
					fmt.Println("Read packet error ", err.Error())
					return
				}

				//handle msg packet
				hdres := GetMsgHandlerIns().HandleMsgPacket(packet, se)
				if hdres != nil {
					fmt.Println(hdres.Error())
					return
				}
			}

		}

	}
}
