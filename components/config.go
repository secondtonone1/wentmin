package components

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
)

var BConfig config.Configer = nil
var ServerPort = 9091
var MaxMsgLen uint16 = 2048
var MaxMsgQueLen int = 2048
var MaxMsgQueNum int = 1
var OutputQueNum int = 2
var OutputQueLen int = 2048
var MaxMsgId uint16 = 9999
var WebPort = 9527

func init() {
	var err error
	BConfig, err = config.NewConfig("ini", "config/server.conf")
	if err != nil {
		panic("config init error")
	}
}

func InitTcpCfg() {
	maxlines, lerr := BConfig.Int64("log::maxlines")
	if lerr != nil {
		maxlines = 1000
	}

	logConf := make(map[string]interface{})
	logConf["filename"] = BConfig.String("log::log_path")
	level, _ := BConfig.Int("log::log_level")
	logConf["level"] = level
	logConf["maxlines"] = maxlines

	confStr, err := json.Marshal(logConf)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(confStr))
	logs.SetLogFuncCall(true)

	ServerPort, err = BConfig.Int("server::port")
	if err != nil {
		fmt.Println("server port error is ", err)
		return
	}

	maxmsglen, err := BConfig.Int("server::max_msg_len")
	if err != nil {
		fmt.Println("server max msg len read failed , err is ", err)
		return
	}
	MaxMsgLen = uint16(maxmsglen)

	MaxMsgQueLen, err = BConfig.Int("server::max_msg_queue_len")
	if err != nil {
		fmt.Println("server max msg queue len read failed ")
		return
	}

	MaxMsgQueNum, err = BConfig.Int("server::max_msg_queue_num")
	if err != nil {
		fmt.Println("server max_msg_queue_num read failed ")
		return
	}

	OutputQueNum, err = BConfig.Int("server::output_queue_num")
	if err != nil {
		fmt.Println("server server::output_queue_num read failed ")
		return
	}

	OutputQueLen, err = BConfig.Int("server::output_queue_len")
	if err != nil {
		fmt.Println("server server::output_queue_len read failed ")
		return
	}

	msgid, err := BConfig.Int("server::max_msg_id")
	if err != nil {
		fmt.Println("server server::max_msg_id read failed ")
		return
	}

	MaxMsgId = uint16(msgid)

}

func InitWebCfg() {
	var err error
	WebPort, err = BConfig.Int("server::webport")
	if err != nil {
		fmt.Println("server server::webport read failed ")
		return
	}
}
