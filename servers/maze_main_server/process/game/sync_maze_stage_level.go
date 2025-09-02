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

func (g *Game) OnSyncMazeStageLevelRQ_10546_10547(s *session.Session, req *MazeGame.SyncMazeStorageStageLevelRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStageLevelRQ")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeGame.SyncMazeStorageStageLevelRS{}
	logger.CtxInfo(ctx, "OnSyncMazeStageLevelRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSyncMazeStageLevelRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.StageLevel = req.StageLevel

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStageLevelRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetStageLevel()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(ctx, userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_StageLevel, data)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStageLevelRQ SaveSyncMazeStorageInfo", zap.Int32("StageLevel", req.GetStageLevel()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
