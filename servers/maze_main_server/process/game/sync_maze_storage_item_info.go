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

func (g *Game) OnSyncMazeStorageItemInfoRQ_10521_10522(s *session.Session, req *MazeGame.SyncMazeStorageItemInfoRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageItemInfoRQ")()
	ctx := s.Context()
	logger := log.Clone("Game", uint64(s.UID()), 0)

	res := &MazeGame.SyncMazeStorageItemInfoRS{}
	logger.InfoWF("OnSyncMazeStorageItemInfoRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnSyncMazeStorageItemInfoRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.StorageItemInfo = req.StorageItemInfo

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageItemInfoRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetStorageItemInfo()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_StorageItemInfo, data)
	if err != nil {
		logger.ErrorWF("OnSyncMazeStorageItemInfoRQ SaveSyncMazeStorageInfo", zap.Any("StorageItemInfo", req.GetStorageItemInfo()), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
