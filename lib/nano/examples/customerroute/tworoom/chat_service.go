package tworoom

import (
	"log"

	"maze_game_server/lib/nano"
	"maze_game_server/lib/nano/component"
	"maze_game_server/lib/nano/examples/cluster/protocol"
	"maze_game_server/lib/nano/session"
)

type ChatRoomService struct {
	component.Base
	group *nano.Group
}

func newChatRoomService() *ChatRoomService {
	return &ChatRoomService{
		group: nano.NewGroup("all-users"),
	}
}

func (rs *ChatRoomService) JoinRoom(s *session.Session, msg *protocol.JoinRoomRequest) error {
	return rs.group.Add(s)
}

type SyncMessage struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (rs *ChatRoomService) SyncMessage(s *session.Session, msg *SyncMessage) error {
	// Sync message to all members in this room
	return rs.group.Broadcast("onMessage", msg)
}

func (rs *ChatRoomService) userDisconnected(s *session.Session) {
	if err := rs.group.Leave(s); err != nil {
		log.Println("Remove user from group failed", s.UID(), err)
		return
	}
	log.Println("User session disconnected", s.UID())
}
