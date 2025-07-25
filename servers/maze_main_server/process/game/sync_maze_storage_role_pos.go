package game

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"
)

func (g *Game) OnSyncMazeStorageRolePosRQ_10519_10520(s *session.Session, req *MazeGame.SyncMazeStorageRolePosRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageRolePosRQ")()

	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.SyncMazeStorageRolePosRS{}
	logger.InfoWF("OnSyncMazeStorageRolePosRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStorageRolePosRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.RolePos = req.RolePos

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageRolePosRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetRolePos()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_RolePos, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageRolePosRQ GetCollectInfo", zap.String("rolePos", data), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
