package jsonproto

type JsonMsg struct {
	MsgId   int    `json:"msgid"`
	MsgData string `json:"msgdata"`
}

type CSUserLogin struct {
	AccountId string `json:"accountid"`
	Passwd    string `json:"passwd"`
	Phone     string `json:"phone"`
}

type SCUserLogin struct {
	ErrorId int    `json:"errorid"`
	Token   string `json:"token"`
}

type CSUserCall struct {
	Caller      string `json:"caller"`
	BeCalled    string `json:"becalled"`
	IsAudioOnly bool   `json:"isaudioonly"`
}

type SCUserCall struct {
	ErrorId  int    `json:"errorid"`
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Phone    string `json:"phone"`
	Roomid   string `json:"roomid"`
}

type SCNotifyBeCall struct {
	Caller      string `json:"caller"`
	BeCalled    string `json:"becalled"`
	Roomid      string `json:"roomid"`
	IsAudioOnly bool   `json:"isaudioonly"`
}

type CSNotifyBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Agree    bool   `json:"agree"`
	Roomid   string `json:"roomid"`
}

type SCNotifyCallRing struct {
	ErrorId  int    `json:"errorid"`
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
}

//主叫方中止呼叫
type CSTerminalCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
}

//服务器通知被叫方挂断
type SCTerminalBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
}
