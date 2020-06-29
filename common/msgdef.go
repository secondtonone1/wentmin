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
	WEB_HEART_BEAT           = 9    //web心跳，后期可以扩充
	WEB_CS_USER_REG          = 2001 //用户注册身份信息
	WEB_SC_USER_REG          = 2002 //服务器回复注册结果
	WEB_CS_USER_CALL         = 2003 //用户呼叫请求
	WEB_SC_USER_CALL         = 2004 //服务器回复呼叫结果
	WEB_NOTIFY_BECALL        = 2005 //服务器通知被呼叫方有通话请求
	WEB_REPLY_BECALL         = 2006 //客户端回复服务器，是否同意接听
	WEB_NOTIFY_CALLRING      = 2007 //通知主叫人唤起响铃
	WEB_CS_TERMINAL_CALL     = 2008 //主叫方终止呼叫，或者通话中任意一方挂断
	WEB_SC_TERMINAL_BECALL   = 2009 //服务器通知另一方终止通话
	WEB_CS_CALL_OFFER        = 2010 //主叫方给服务器发送offer
	WEB_SC_OFFER_NOTIFY      = 2011 //服务器通知被叫方offer信息
	WEB_CS_BECALL_ANSWER     = 2012 //被叫方将answer信息发送给服务器
	WEB_SC_ANSWER_NOTIFY     = 2013 //服务器通知主叫方answer信息
	WEB_CS_CALL_ICE          = 2014 //主叫方将ICE_CANDIDATE信息发送给服务器
	WEB_SC_CALL_ICE_NOTIFY   = 2015 //服务器通知被叫方 来自主叫方的ICE信息
	WEB_CS_BECALL_ICE        = 2016 //被叫方将ICE_CANDIDATE信息发送给服务器
	WEB_SC_BECALL_ICE_NOTIFY = 2017 //服务器通知主叫方 来自被叫方的ICE信息
)
