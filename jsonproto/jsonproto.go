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

//一方挂断
type CSTerminalCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Cancel   string `json:"cancel"`
}

//服务器通知另一方挂断
type SCTerminalBeCall struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Cancel   string `json:"cancel"`
}

//主叫方提供offer
type CSCallOffer struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Sdp      string `json:"sdp"`
}

//服务器通知
type SCOfferNotify struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Sdp      string `json:"sdp"`
}

//被叫方回复answer信息
type CSBecallAnswer struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Sdp      string `json:"sdp"`
}

//服务器通知主叫方answer信息
type SCAnswerNotify struct {
	Caller   string `json:"caller"`
	BeCalled string `json:"becalled"`
	Roomid   string `json:"roomid"`
	Sdp      string `json:"sdp"`
}

//主叫方提供ice
type CSCallIce struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"becalled"`
	Roomid    string `json:"roomid"`
	Candidate string `json:"candidate"`
}

//服务器通知被叫方，接收来自主叫方的ice
type SCCallIceNotify struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"becalled"`
	Roomid    string `json:"roomid"`
	Candidate string `json:"candidate"`
}

// 被叫方提供ice
type CSBecallIce struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"becalled"`
	Roomid    string `json:"roomid"`
	Candidate string `json:"candidate"`
}

//服务器通知主叫方，来自被叫方的ice
type SCBecallIceNotify struct {
	Caller    string `json:"caller"`
	BeCalled  string `json:"becalled"`
	Roomid    string `json:"roomid"`
	Candidate string `json:"candidate"`
}
