package netmodel

import (
	"fmt"
	"net"
	"sync"
	"time"
	"wentmin/protocol"

	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	conn       net.Conn
	closed     int32                  //session是否关闭，-1未开启，0未关闭，1关闭
	protocol   protocol.ProtocolInter //字节序和自己处理器
	RWLock     sync.RWMutex           //协程锁
	SocketId   string                 // 当前socket
	AliveTime  int64
	LastSocket string // 上次socket
}

type MsgSession struct {
	session *Session
	packet  *protocol.MsgPacket
}

func NewSession(connt net.Conn,
) *Session {
	uid := uuid.NewV4()
	sess := &Session{
		conn:      connt,
		closed:    -1,
		protocol:  new(protocol.ProtocolImpl),
		AliveTime: time.Now().Unix(),
		SocketId:  uid.String(),
	}
	tcpConn := sess.conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)
	tcpConn.SetReadBuffer(64 * 1024)
	tcpConn.SetWriteBuffer(64 * 1024)
	return sess
}

func (se *Session) GetSocketId() string {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return se.SocketId
}

func (se *Session) GetSocketLen() int {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return len(se.SocketId)
}

func (se *Session) GetLastSocket() string {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return se.LastSocket
}

func (se *Session) RawConn() *net.TCPConn {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return se.conn.(*net.TCPConn)
}

func (se *Session) CheckAlive(now int64) bool {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return ((now - se.AliveTime) < 60*60)
}

func (se *Session) UpdateAlive(now int64) {
	se.RWLock.Lock()
	defer se.RWLock.Unlock()
	se.AliveTime = now
}

func (se *Session) IsClosed() bool {
	se.RWLock.RLock()
	defer se.RWLock.RUnlock()
	return se.closed == 1
}

func (se *Session) Start() {
	se.RWLock.Lock()
	defer se.RWLock.Unlock()
	if se.closed != -1 {
		return
	}

	se.closed = 0
	go se.recvLoop()

}

// Close the session, destory other resource.
func (se *Session) Close() error {
	se.RWLock.Lock()
	defer se.RWLock.Unlock()
	if se.closed != 0 {
		return nil
	}
	se.LastSocket = se.SocketId
	se.SocketId = ""
	se.closed = 1
	se.conn.Close()
	return nil
}

func (se *Session) Write(msgpkg *protocol.MsgPacket) {
	se.protocol.WritePacket(se.conn, msgpkg)
	fmt.Println("send msg success , msg id is ", msgpkg.Head.Id)
}

func (se *Session) Read() (interface{}, error) {
	packet, err := se.protocol.ReadPacket(se.conn)
	return packet, err
}

func (se *Session) recvLoop() {
	defer TcpServerInst.OnSessionClosed(se.SocketId)
	var packet interface{}
	var err error
	for {

		select {
		case <-AcceptClose:
			return
		default:
			{
				packet, err = se.Read()
				if packet == nil || err != nil {
					//fmt.Println("Read packet error ", err.Error())
					return
				}
				cur := time.Now().Unix()
				se.UpdateAlive(cur)
				msgs := new(MsgSession)
				msgs.packet = packet.(*protocol.MsgPacket)
				msgs.session = se
				err = MsgQueueInst.PutMsgInQue(msgs)
				if err != nil {
					fmt.Println("put msg into queue failed , err is ", err.Error())
					logs.Debug("put msg into queue failed , err is ", err.Error())
					return
				}

			}

		}

	}
}
