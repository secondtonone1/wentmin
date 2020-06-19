package common

const (
	//系统级别消息从1~1000
	SYS_CON_CONNECT = 1
	SYC_CON_CLOSED  = 2
	SYS_HEART_BEAT  = 3 //tcp心跳,后期扩充
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

const (
	WEB_HEART_BEAT    = 9    //web心跳，后期可以扩充
	WEB_CS_USER_REG   = 2001 //用户注册身份信息
	WEB_SC_USER_REG   = 2002 //服务器回复注册结果
	WEB_CS_USER_CALL  = 2003 //用户呼叫请求
	WEB_SC_USER_CALL  = 2004 //服务器回复呼叫结果
	WEB_NOTIFY_BECALL = 2005 //服务器通知被呼叫方有通话请求
	WEB_REPLY_BECALL  = 2006 //客户端回复服务器，是否同意接听
)
