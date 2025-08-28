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

// 退出房间
func (g *Frame) OnExistRoomRQ_10541_10542(s *session.Session, req *MazeRoom.MazeExistRoomRQ) (err error) {
	defer fkprometheus.InfoPMT("OnExistRoomRQ")()
	ctx := s.Context()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeExistRoomRS{}

	logger.InfoWF("OnExistRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnExistRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnExistRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = frame_service.ExistRoom(logger, userId)
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return
	}
	return nil
}
