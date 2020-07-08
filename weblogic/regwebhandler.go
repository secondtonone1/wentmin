package weblogic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"

	"wentmin/jsonproto"

	"net/http"

	"golang.org/x/net/websocket"
)

//  ws://localhost:9528/wsmsg
func RegWSHandlers() {
	http.Handle("/wsmsg", svrConnHandler)
}

/*
var svrConnHandler websocket.Handler = func(conn *websocket.Conn) {
	request := make([]byte, components.MaxMsgLen*3)
	defer conn.Close()
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			fmt.Println(config.ErrWebSocketRead.Error())
			return
		}

		//socket被关闭了
		if readLen == 0 {
			fmt.Println(config.ErrWebSocketClosed.Error())
			//执行逻辑断开处理
			return
		}

		fmt.Println(string(request[:readLen]))
		//后期这个改在回调函数中
		conn.Write([]byte("Recieve Hello World Msg!"))
		request = make([]byte, components.MaxMsgLen*3)
	}
}
*/

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

var svrConnHandler websocket.Handler = func(ws *websocket.Conn) {
	var err error
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("web socket logic goroutine recover from panic")
		}

		//处理断线逻辑
		WebLogicLock.Lock()
		defer WebLogicLock.Unlock()
		UserMgrInst.OnOffline(ws)
	}()
	fmt.Println("new connection arraived")
	fmt.Println("cur goroutine is ", GetGID())
	/*
		//校验头部信息
		var token = ws.Request().Header.Get("token")
		if token == "" {
			fmt.Println("connection token empty ")
			return
		}
	*/
	for {

		var reply string

		//websocket接受信息

		if err = websocket.Message.Receive(ws, &reply); err != nil {

			fmt.Println("receive failed:", err)

			break

		}

		fmt.Println("receive msg is ", reply)
		jsonMsg := &jsonproto.JsonMsg{}

		err = json.Unmarshal([]byte(reply), jsonMsg)
		if err != nil {
			fmt.Println("json unmarshal failed , error is ", err.Error())
			continue
		}

		err = HandleWebMsg(ws, jsonMsg.MsgId, jsonMsg.MsgData)
		if err != nil {
			fmt.Printf("handle web msg[%d] failed, error is %s", jsonMsg.MsgId, err.Error())
			continue
		}

	}
}
