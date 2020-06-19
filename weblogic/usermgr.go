package weblogic

import (
	"sync"

	"golang.org/x/net/websocket"
)

type UserData struct {
	AccountId string          `json:"accountid"`
	Passwd    string          `json:"passwd"`
	Phone     string          `json:"phone"`
	Conn      *websocket.Conn `json:"-"`
	RWlock    sync.RWMutex    `json:"-"`
}

type UserMgr struct {
	UsersMap    map[string]*UserData
	SessUserMap map[string]*UserData
}

var UserMgrInst *UserMgr
