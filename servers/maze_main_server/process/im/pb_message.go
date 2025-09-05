package im

import (
	"maze_game_server/io/redis/im/msgstore"
	"maze_game_server/pb/common/MazeIM"

	"google.golang.org/protobuf/proto"
)

func PbMessage(messageList msgstore.Message) *MazeIM.Message {
	pbMsg := &MazeIM.Message{
		MsgId:      proto.Uint64(messageList.MessageID),
		Type:       proto.Int32(messageList.Type),
		Content:    []byte(messageList.Content),
		Sender:     proto.Uint64(messageList.UserID),
		CreateTime: proto.Int64(messageList.CreateTime),
	}
	return pbMsg
}
