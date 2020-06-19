package weblogic

type ChatRoom struct {
	Token    string
	Caller   string
	Becalled string
}

type ChatMgr struct {
	ChatRoomMap map[string]*ChatRoom
}

var ChatMgrInst *ChatMgr

func (cm *ChatMgr) AddRoom(cr *ChatRoom) {
	cm.ChatRoomMap[cr.Token] = cr
}

func (cm *ChatMgr) DelRoom(token string) {
	delete(cm.ChatRoomMap, token)
}

func init() {
	ChatMgrInst = &ChatMgr{}
	ChatMgrInst.ChatRoomMap = make(map[string]*ChatRoom)
}
