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

func (g *Game) OnSyncMazeStorageRoleItemDataRQ_10515_10516(s *session.Session, req *MazeGame.SyncMazeStorageRoleItemDataRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageRoleItemDataRQ")()
	ctx := s.Context()
	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.SyncMazeStorageRoleItemDataRS{}
	logger.InfoWF("OnSyncMazeStorageRoleItemDataRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStorageRoleItemDataRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.RoleItemData = req.RoleItemData

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageRoleItemDataRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetRoleItemData()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_RoleItemData, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageRoleItemDataRQ GetCollectInfo", zap.String("roleItemData", req.GetRoleItemData()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	return nil
}
