package process

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/lib/net/websocket_service"
)

func SendArrivePacket(logger fklog.FKLogI, userID int64, packetType uint16, pack proto.Message) error {
	// data, err := proto.Marshal(pack)
	// if err != nil {
	// 	logger.ErrorWF("SendArrivePacket Marshal err:%v",
	// 		zap.Any("userID", userID),
	// 		zap.Any("packetType", packetType), zap.Error(err))
	// 	return err
	// }
	return websocket_service.SendPacket(logger, uint64(userID), packetType, pack)
}
