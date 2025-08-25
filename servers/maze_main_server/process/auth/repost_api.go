package auth

import (
	"context"

	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
)

func SendArrivePacket(logger fklog.FKLogI, userID int64, packetType uint16, pack proto.Message) error {
	return online.ClusterPush(context.TODO(), uint64(userID), packetType, pack)
}

func SendArrivePacketWithContext(ctx context.Context, logger fklog.FKLogI, userID int64, packetType uint16, pack proto.Message) error {
	return online.PushWithContext(ctx, logger, uint64(userID), packetType, pack)
}
