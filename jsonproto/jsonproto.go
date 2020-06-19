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
	ErrorId   string `json:"errorid"`
	AccountId string `json:"accountid"`
	Passwd    string `json:"passwd"`
	Phone     string `json:"phone"`
}
