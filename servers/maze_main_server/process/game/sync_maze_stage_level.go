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

func (g *Game) OnSyncMazeStageLevelRQ_10546_10547(s *session.Session, req *MazeGame.SyncMazeStorageStageLevelRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStageLevelRQ")()
	logger := log.Clone("Game", uint64(s.UID()), 0)
	ctx := s.Context()
	res := &MazeGame.SyncMazeStorageStageLevelRS{}
	logger.InfoWF("OnSyncMazeStageLevelRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStageLevelRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.StageLevel = req.StageLevel

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStageLevelRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetStageLevel()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_StageLevel, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStageLevelRQ SaveSyncMazeStorageInfo", zap.Int32("StageLevel", req.GetStageLevel()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
