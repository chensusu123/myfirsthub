package frame

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeRoom"
	frame_service "maze_game_server/services/frame"
)

// 客户端上报房间
func (g *Frame) OnFrameSyncRQ_10539_10540(s *session.Session, req *MazeRoom.MazeFrameSyncRQ) (err error) {
	defer fkprometheus.InfoPMT("OnFrameSyncRQ")()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeFrameSyncRS{}

	logger.InfoWF("OnFrameSyncRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnFrameSyncRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnFrameSyncRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = frame_service.WritePump(logger, userId, req.GetBytes())
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	res.AccessInfo = req.AccessInfo
	res.Bytes = req.Bytes
	return nil
}
