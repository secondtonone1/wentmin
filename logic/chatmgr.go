package logic

type ChatRoom struct {
	Roomid   string
	Caller   string
	Becalled string
}

type ChatMgr struct {
	ChatRoomMap map[string]*ChatRoom
}

var ChatMgrInst *ChatMgr

func (cm *ChatMgr) AddRoom(cr *ChatRoom) {
	cm.ChatRoomMap[cr.Roomid] = cr
}

func (cm *ChatMgr) DelRoom(roomid string) {
	delete(cm.ChatRoomMap, roomid)
}

func init() {
	ChatMgrInst = &ChatMgr{}
	ChatMgrInst.ChatRoomMap = make(map[string]*ChatRoom)
}
