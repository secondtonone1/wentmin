package jsonproto

type JsonMsg struct {
	MsgId   int         `json:"msgId"`
	MsgData interface{} `json:"msgData"`
}

type CSUserLogin struct {
	AccountId string `json:"accountId"`
	Passwd    string `json:"passwd"`
	Phone     string `json:"phone"`
}

type SCUserLogin struct {
	ErrorId int    `json:"errorId"`
	Token   string `json:"token"`
}

type CSUserCall struct {
	Caller      string `json:"caller"`
	BeCalled    string `json:"beCalled"`
	IsAudioOnly bool   `json:"isAudioOnly"`
}

type SCUserCall struct {
	ErrorId  int    `json:"errorId"`
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Phone    string `json:"phone"`
	Roomid   string `json:"roomId"`
}

type SCNotifyBeCall struct {
	Caller      string `json:"caller"`
	BeCalled    string `json:"beCalled"`
	Roomid      string `json:"roomId"`
	IsAudioOnly bool   `json:"isAudioOnly"`
}

type CSNotifyBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Agree    bool   `json:"agree"`
	Roomid   string `json:"roomId"`
}

type SCNotifyCallRing struct {
	ErrorId  int    `json:"errorId"`
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
}

//一方挂断
type CSTerminalCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Cancel   string `json:"cancel"`
}

//服务器通知另一方挂断
type SCTerminalBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Cancel   string `json:"cancel"`
}

//主叫方提供offer
type CSCallOffer struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Sdp      string `json:"sdp"`
}

//服务器通知
type SCOfferNotify struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Sdp      string `json:"sdp"`
}

//被叫方回复answer信息
type CSBecallAnswer struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Sdp      string `json:"sdp"`
}

//服务器通知主叫方answer信息
type SCAnswerNotify struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"beCalled"`
	Roomid   string `json:"roomId"`
	Sdp      string `json:"sdp"`
}

//主叫方提供ice
type CSCallIce struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"beCalled"`
	Roomid    string `json:"roomId"`
	Candidate string `json:"candidate"`
}

//服务器通知被叫方，接收来自主叫方的ice
type SCCallIceNotify struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"beCalled"`
	Roomid    string `json:"roomId"`
	Candidate string `json:"candidate"`
}

// 被叫方提供ice
type CSBecallIce struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"beCalled"`
	Roomid    string `json:"roomId"`
	Candidate string `json:"candidate"`
}

//服务器通知主叫方，来自被叫方的ice
type SCBecallIceNotify struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"beCalled"`
	Roomid    string `json:"roomId"`
	Candidate string `json:"candidate"`
}
