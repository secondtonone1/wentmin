package weblogic

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

func (cm *ChatMgr) GetRoom(roomid string) *ChatRoom {
	room, ok := cm.ChatRoomMap[roomid]
	if !ok {
		return nil
	}

	return room
}

func (cm *ChatMgr) DelRoomByUser(userid string) {
	for roomid, val := range cm.ChatRoomMap {
		if val.Caller == userid || val.Becalled == userid {
			delete(cm.ChatRoomMap, roomid)
			return
		}
	}
}

func init() {
	ChatMgrInst = &ChatMgr{}
	ChatMgrInst.ChatRoomMap = make(map[string]*ChatRoom)
}
