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
)

// 加入房间
func (g *Frame) OnJoinRoomRQ_10539_10540(s *session.Session, req *MazeRoom.MazeJoinRoomRQ) (err error) {
	defer fkprometheus.InfoPMT("OnJoinRoomRQ")()
	ctx := s.Context()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeJoinRoomRS{}

	logger.CtxInfo(ctx, "OnJoinRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnJoinRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.AccessInfo = req.AccessInfo

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnJoinRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = frame_service.JoinRoom(ctx, userId, req.GetAccessInfo())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return err
	}

	return nil
}
