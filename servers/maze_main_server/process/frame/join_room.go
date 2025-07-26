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

// 加入房间
func (g *Frame) OnJoinRoomRQ_10539_10540(s *session.Session, req *MazeRoom.MazeJoinRoomRQ) (err error) {
	defer fkprometheus.InfoPMT("OnJoinRoomRQ")()

	logger := log.Clone("Frame", uint64(s.UID()), 0)
	res := &MazeRoom.MazeJoinRoomRS{}

	logger.InfoWF("OnJoinRoomRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnJoinRoomRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.AccessInfo = req.AccessInfo

	userId := uint64(s.UID())

	_, err = mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnJoinRoomRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	err = frame_service.JoinRoom(logger, userId, req.GetAccessInfo())
	if err != nil {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap(err.Error())
		return err
	}

	return nil
}
