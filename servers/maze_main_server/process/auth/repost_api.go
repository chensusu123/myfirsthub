package auth

import (
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"google.golang.org/protobuf/proto"
)

func SendArrivePacket(logger fklog.FKLogI, userID int64, packetType uint16, pack proto.Message) error {
	// data, err := proto.Marshal(pack)
	// if err != nil {
	// 	logger.ErrorWF("SendArrivePacket Marshal ",
	// 		zap.Any("userID", userID),
	// 		zap.Any("packetType", packetType), zap.Error(err))
	// 	return err
	// }
	// return websocket_service.SendPacket(logger, uint64(userID), packetType, pack)
	return online.Push(logger, uint64(userID), packetType, pack)
}
