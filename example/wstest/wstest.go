package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"videocall/common"
	"videocall/jsonproto"

	"golang.org/x/net/websocket"
)

var gLocker sync.Mutex    //全局锁
var gCondition *sync.Cond //全局条件变量

/*
var origin = "http://81.68.86.146:9528/"
var url = "ws://81.68.86.146:9528/wsmsg"
*/

var origin = "http://127.0.0.1:9527/"
var url = "ws://127.0.0.1:9527/wsmsg"

//错误处理函数
func checkErr(err error, extra string) bool {
	if err != nil {
		formatStr := " Err : %s\n"
		if extra != "" {
			formatStr = extra + formatStr
		}

		fmt.Fprintf(os.Stderr, formatStr, err.Error())
		return true
	}

	return false
}

//连接处理函数
func clientConnHandler(conn *websocket.Conn) {
	gLocker.Lock()
	defer gLocker.Unlock()
	defer conn.Close()
	request := make([]byte, 1280)
	for {
		readLen, err := conn.Read(request)
		if checkErr(err, "Read") {
			gCondition.Signal()
			break
		}

		//socket被关闭了
		if readLen == 0 {
			fmt.Println("Server connection close!")

			//条件变量同步通知
			gCondition.Signal()
			break
		}

		fmt.Println("reg rsp  is ", string(request))

		//呼叫102
		callMsg := &jsonproto.CSUserCall{}
		callMsg.Caller = "101"
		callMsg.BeCalled = "102"

		jsMsg := &jsonproto.JsonMsg{}
		jsMsg.MsgId = common.WEB_CS_USER_CALL
		jsmal, _ := json.Marshal(callMsg)
		jsMsg.MsgData = string(jsmal)

		jsmal, _ = json.Marshal(jsMsg)
		fmt.Println("send data is ", string(jsmal))
		_, err = conn.Write(jsmal)

		//等待接收响铃通知
		request = make([]byte, 1280)

		readLen, err = conn.Read(request)
		if checkErr(err, "Read") {
			gCondition.Signal()
			break
		}

		//socket被关闭了
		if readLen == 0 {
			fmt.Println("Server connection close!")

			//条件变量同步通知
			gCondition.Signal()
			break
		}

		fmt.Println("ring rsp  is ", string(request))

		jsonmsg := &jsonproto.JsonMsg{}
		err = json.Unmarshal(request[:readLen], jsonmsg)
		if err != nil {
			fmt.Println("err is ", err.Error())
		}
		fmt.Println("receive notify jsondata msg is ", jsonmsg.MsgData)
		jsondata := []byte(jsonmsg.MsgData)

		ringmsg := &jsonproto.SCNotifyCallRing{}
		_ = json.Unmarshal(jsondata, ringmsg)
		fmt.Println("becall is ", ringmsg.BeCalled)
		fmt.Println("caller is ", ringmsg.Caller)
		fmt.Println("errorid is ", ringmsg.ErrorId)
		fmt.Println("errorid is ", ringmsg.Roomid)

		request = make([]byte, 1280)

		readLen, err = conn.Read(request)
		if checkErr(err, "Read") {
			gCondition.Signal()
			break
		}

		//socket被关闭了
		if readLen == 0 {
			fmt.Println("Server connection close!")

			//条件变量同步通知
			gCondition.Signal()
			break
		}

		fmt.Println("reg rsp  is ", string(request))
		jsonmsg = &jsonproto.JsonMsg{}
		err = json.Unmarshal(request[:readLen], jsonmsg)
		if err != nil {
			fmt.Println("err is ", err.Error())
		}
		fmt.Println("receive notify jsondata msg is ", jsonmsg.MsgData)
		jsondata = []byte(jsonmsg.MsgData)
		becallss := &jsonproto.SCUserCall{}
		_ = json.Unmarshal(jsondata, becallss)
		fmt.Println("becall is ", becallss.BeCalled)
		fmt.Println("caller is ", becallss.Caller)
		fmt.Println("errorid is ", becallss.ErrorId)
		fmt.Println("room id is ", becallss.Roomid)
		fmt.Println("begin chat .....")
		time.Sleep(time.Second * 2)

		//中断聊天
		terminalCallMsg := &jsonproto.CSTerminalCall{}
		terminalCallMsg.Caller = "101"
		terminalCallMsg.BeCalled = "102"
		terminalCallMsg.Cancel = "101"
		terminalCallMsg.Roomid = becallss.Roomid
		jsMsg = &jsonproto.JsonMsg{}
		jsMsg.MsgId = common.WEB_CS_TERMINAL_CALL
		jsmal, _ = json.Marshal(terminalCallMsg)
		jsMsg.MsgData = string(jsmal)

		jsmal, _ = json.Marshal(jsMsg)
		fmt.Println("send data is ", string(jsmal))
		_, err = conn.Write(jsmal)

		fmt.Println("end chat .....")
		//条件变量同步通知
		gCondition.Signal()
		break
	}
}

func main() {
	conn, err := websocket.Dial(url, "", origin)
	if checkErr(err, "Dial") {
		return
	}

	gLocker.Lock()
	gCondition = sync.NewCond(&gLocker)

	regMsg := &jsonproto.CSUserLogin{}
	regMsg.AccountId = "101"
	regMsg.Passwd = "pwd101"
	regMsg.Phone = "18301152008"

	jsMsg := &jsonproto.JsonMsg{}
	jsMsg.MsgId = common.WEB_CS_USER_REG
	jsmal, _ := json.Marshal(regMsg)
	jsMsg.MsgData = string(jsmal)

	jsmal, _ = json.Marshal(jsMsg)
	fmt.Println("send data is ", string(jsmal))
	_, err = conn.Write(jsmal)
	go clientConnHandler(conn)

	//主线程阻塞，等待Singal结束
	for {
		//条件变量同步等待
		gCondition.Wait()
		break
	}
	gLocker.Unlock()
	fmt.Println("Client1 finish.")
}
