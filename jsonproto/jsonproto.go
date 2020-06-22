package jsonproto

type JsonMsg struct {
	MsgId   int    `json:"msgid"`
	MsgData string `json:"msgdata"`
}

type CSUserReg struct {
	AccountId string `json:"accountid"`
	Passwd    string `json:"passwd"`
	Phone     string `json:"phone"`
}

type SCUserReg struct {
	ErrorId   int    `json:"errorid"`
	AccountId string `json:"accountid"`
	Passwd    string `json:"passwd"`
	Phone     string `json:"phone"`
}

type CSUserCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
}

type SCUserCall struct {
	ErrorId  int    `json:"errorid"`
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

type SCNotifyBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Token    string `json:"token"`
}

type CSNotifyBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Agree    bool   `json:"agree"`
	Token    string `json:"token"`
}

type SCNotifyCallRing struct {
	ErrorId  int    `json:"errorid"`
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Token    string `json:"token"`
}

//主叫方中止呼叫
type CSTerminalCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Token    string `json:"token"`
}

//服务器通知被叫方挂断
type SCTerminalBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Token    string `json:"token"`
}
