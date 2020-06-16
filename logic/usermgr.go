package logic

type UserData struct {
	AccountId string
	Passwd    string
	Phone     string
}

type UserMgr struct {
	UsersMap map[string]*UserData
}

var UserMgrInst *UserMgr

func init() {
	UserMgrInst = &UserMgr{UsersMap: make(map[string]*UserData)}
}

func (self *UserMgr) AddUser(user *UserData) {
	self.UsersMap[user.AccountId] = user
}
