package common

const (
	//系统级别消息从1~1000
	SYS_CON_CONNECT = 1
	SYC_CON_CLOSED  = 2
	//用户级别的消息从1001开始
	USER_REG_CS      = 1001
	USER_REG_SC      = 1002
	CS_USER_CALL     = 1003
	SC_USER_CALL     = 1004
	SC_NOTIFY_BECALL = 1005 //通知被叫方会话
	CS_NOTIFY_REPLY  = 1006 //被叫方回复服务器
	SC_NOTIFY_CHAT   = 1007 //回复chat
	CS_TERMINAL_CHAT = 1008 //客户端终止chat
	SC_TERMINAL_CHAT = 1009 //服务器回复终止chat
)
