package logic

import (
	"errors"
	"wentmin/netmodel"
)

type UserData struct {
	AccountId string
	Passwd    string
	Phone     string
	Session   *netmodel.Session
}

type UserMgr struct {
	UsersMap    map[string]*UserData
	SessUserMap map[int]*UserData
}

var UserMgrInst *UserMgr

func init() {
	UserMgrInst = &UserMgr{UsersMap: make(map[string]*UserData),
		SessUserMap: make(map[int]*UserData)}
}

func (self *UserMgr) AddUser(user *UserData) {
	self.UsersMap[user.AccountId] = user
	sid := user.Session.GetSocketId()
	self.SessUserMap[sid] = user
}

func (self *UserMgr) SetUserOffline(sid int) {
	ud, ok := self.SessUserMap[sid]
	if !ok {
		return
	}
	ud.Session = nil
	delete(self.SessUserMap, sid)
}

func (self *UserMgr) GetUser(accountid string) (*UserData, error) {
	ud, ok := self.UsersMap[accountid]
	if !ok {
		return nil, errors.New("user not found")
	}

	return ud, nil
}

func (self *UserData) IsOnline() bool {
	if self.Session == nil {
		return false
	}

	if self.Session.IsClosed() == true {
		return false
	}

	return true
}

func (self *UserData) GetSession() *netmodel.Session {
	return self.Session
}
