package game

import (
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/syncmazestorageinforedis"
	"maze_game_server/lib/nano/session"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeGame"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"go.uber.org/zap"
)

func (g *Game) OnSyncMazeStorageMonsterAreaInfoRQ_10548_10549(s *session.Session, req *MazeGame.SyncMazeStorageMonsterAreaInfoRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageMonsterAreaInfoRQ")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeGame.SyncMazeStorageMonsterAreaInfoRS{}
	logger.CtxInfo(ctx, "OnSyncMazeStorageMonsterAreaInfoRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSyncMazeStorageMonsterAreaInfoRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.MonsterAreaInfo = req.MonsterAreaInfo

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStorageMonsterAreaInfoRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetMonsterAreaInfo()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_MonsterAreaInfo, data)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStorageMonsterAreaInfoRQ SaveSyncMazeStorageInfo failed", zap.String("roleItemData", req.GetMonsterAreaInfo()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	return nil
}
