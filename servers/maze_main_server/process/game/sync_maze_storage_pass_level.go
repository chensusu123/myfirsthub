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

func (g *Game) OnSyncMazeStoragePassLevelRQ_10517_10518(s *session.Session, req *MazeGame.SyncMazeStoragePassLevelRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStoragePassLevelRQ")()
	ctx := s.Context()

	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.SyncMazeStoragePassLevelRS{}
	logger.InfoWF("OnSyncMazeStoragePassLevelRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStoragePassLevelRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.PassLevel = req.PassLevel

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStoragePassLevelRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetPassLevel()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_PassLevel, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStoragePassLevelRQ SaveSyncMazeStorageInfo", zap.Int64("passLevel", req.GetPassLevel()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
