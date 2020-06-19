package weblogic

import (
	"encoding/json"
	"fmt"

	"wentmin/jsonproto"

	"net/http"

	"golang.org/x/net/websocket"
)

//  ws://localhost:9527/wsmsg
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
