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

// 客户端上报房间
func (g *Frame) OnFrameSyncRQ_10539_10540(s *session.Session, req *MazeRoom.MazeFrameSyncRQ) (err error) {
	defer fkprometheus.InfoPMT("OnFrameSyncRQ")()
	ctx := s.Context()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeFrameSyncRS{}

	logger.CtxInfo(ctx, "OnFrameSyncRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnFrameSyncRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnFrameSyncRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = frame_service.WritePump(ctx, userId, req.GetBytes())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	res.AccessInfo = req.AccessInfo
	res.Bytes = req.Bytes
	return nil
}
