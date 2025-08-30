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

func (g *Game) OnSyncMazeStoragePassLevelRQ_10517_10518(s *session.Session, req *MazeGame.SyncMazeStoragePassLevelRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStoragePassLevelRQ")()
	ctx := s.Context()

	logger := fklog.ContextAppLogger(ctx)

	res := &MazeGame.SyncMazeStoragePassLevelRS{}
	logger.CtxInfo(ctx, "OnSyncMazeStoragePassLevelRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSyncMazeStoragePassLevelRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.PassLevel = req.PassLevel

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStoragePassLevelRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetPassLevel()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_PassLevel, data)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStoragePassLevelRQ SaveSyncMazeStorageInfo", zap.Int64("passLevel", req.GetPassLevel()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
