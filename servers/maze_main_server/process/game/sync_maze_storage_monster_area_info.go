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

func (g *Game) OnSyncMazeStorageMonsterAreaInfoRQ_10548_10549(s *session.Session, req *MazeGame.SyncMazeStorageMonsterAreaInfoRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageMonsterAreaInfoRQ")()
	ctx := s.Context()
	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.SyncMazeStorageMonsterAreaInfoRS{}
	logger.InfoWF("OnSyncMazeStorageMonsterAreaInfoRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStorageMonsterAreaInfoRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.MonsterAreaInfo = req.MonsterAreaInfo

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageMonsterAreaInfoRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetMonsterAreaInfo()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_MonsterAreaInfo, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageMonsterAreaInfoRQ SaveSyncMazeStorageInfo failed", zap.String("roleItemData", req.GetMonsterAreaInfo()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	return nil
}
