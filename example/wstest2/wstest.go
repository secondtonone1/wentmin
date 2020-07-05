package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"wentmin/common"
	"wentmin/jsonproto"

	"github.com/goinggo/mapstructure"
	"golang.org/x/net/websocket"
)

var gLocker sync.Mutex    //全局锁
var gCondition *sync.Cond //全局条件变量

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

		fmt.Println("notify   is ", string(request))
		jsonmsg := &jsonproto.JsonMsg{}
		err = json.Unmarshal(request[:readLen], jsonmsg)
		if err != nil {
			fmt.Println("err is ", err.Error())
		}
		fmt.Println("receive notify jsondata msg is ", jsonmsg.MsgData)
		becall := &jsonproto.SCNotifyBeCall{}
		mapstructure.Decode(jsonmsg.MsgData, becall)

		fmt.Println("becall is ", becall.BeCalled)
		fmt.Println("caller is ", becall.Caller)
		fmt.Println("roomid is ", becall.Roomid)
		fmt.Println("audioonly is ", becall.IsAudioOnly)

		//同意接通
		jsonmsg = &jsonproto.JsonMsg{}
		jsonmsg.MsgId = common.WEB_REPLY_BECALL
		becallr := &jsonproto.CSNotifyBeCall{}
		becallr.BeCalled = becall.BeCalled
		becallr.Caller = becall.Caller
		becallr.Agree = true
		becallr.Roomid = becall.Roomid

		jsonmsg.MsgData = becallr
		jstmp, _ := json.Marshal(jsonmsg)
		conn.Write(jstmp)

		//等待对方发送中断通话
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

		fmt.Println("receive terminal msg is ", string(request))
		jsonmsg = &jsonproto.JsonMsg{}
		err = json.Unmarshal(request[:readLen], jsonmsg)
		if err != nil {
			fmt.Println("err is ", err.Error())
		}
		fmt.Println("receive notify jsondata msg is ", jsonmsg.MsgData)
		terminalcall := &jsonproto.SCTerminalBeCall{}
		mapstructure.Decode(jsonmsg.MsgData, terminalcall)
		fmt.Println("becall is ", terminalcall.BeCalled)
		fmt.Println("caller is ", terminalcall.Caller)
		fmt.Println("roomid is ", terminalcall.Roomid)
		fmt.Println("cancel is ", terminalcall.Cancel)
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
	regMsg.AccountId = "102"
	regMsg.Passwd = "pwd102"
	regMsg.Phone = "112898988"

	jsMsg := &jsonproto.JsonMsg{}
	jsMsg.MsgId = common.WEB_CS_USER_REG
	jsMsg.MsgData = regMsg

	jsmal, _ := json.Marshal(jsMsg)
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
	fmt.Println("Client2 finish.")
}
