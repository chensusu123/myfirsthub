package frame

import (
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeRoom"
	frame_service "maze_game_server/services/frame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 创建房间
func (g *Frame) OnCreateRoomRQ_10531_10532(s *session.Session, req *MazeRoom.MazeCreateRoomRQ) (err error) {
	defer fkprometheus.InfoPMT("OnCreateRoomRQ")()

	ctx := s.Context()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeCreateRoomRS{}

	logger.CtxInfo(ctx, "OnCreateRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnCreateRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnCreateRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	room, err := frame_service.CreateRoom(ctx, userId, req.OpenIdList, req.GetGameTick(), req.GetStartPercent(), req.GetGameLastTime(), req.GetUdpReliabilityStrategy(), req.GetRoomExtInfo(), req.GetNeedGameSeed())
	if err != nil {
		logger.CtxError(ctx, "OnCreateRoomRQ CreateRoomRQ fail", zap.Error(err))
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}

	res.AccessInfo = proto.String(room.GameAccessInfo)
	return nil
}
