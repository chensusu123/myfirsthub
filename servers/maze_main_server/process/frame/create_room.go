package frame

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeRoom"
	frame_service "maze_game_server/services/frame"
)

// 创建房间
func (g *Frame) OnCreateRoomRQ_10531_10532(s *session.Session, req *MazeRoom.MazeCreateRoomRQ) (err error) {
	defer fkprometheus.InfoPMT("OnCreateRoomRQ")()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeCreateRoomRS{}

	logger.InfoWF("OnCreateRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnCreateRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnCreateRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	room, err := frame_service.CreateRoom(logger, userId, req.OpenIdList, req.GetGameTick(), req.GetStartPercent(), req.GetGameLastTime(), req.GetUdpReliabilityStrategy(), req.GetRoomExtInfo(), req.GetNeedGameSeed())
	if err != nil {
		logger.ErrorWF("OnCreateRoomRQ CreateRoomRQ fail", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}

	res.AccessInfo = proto.String(room.GameAccessInfo)
	return nil
}
