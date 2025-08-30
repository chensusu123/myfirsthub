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

func (g *Game) OnSyncMazeStorageRolePosRQ_10519_10520(s *session.Session, req *MazeGame.SyncMazeStorageRolePosRQ) (err error) {
	defer fkprometheus.InfoPMT("OnSyncMazeStorageRolePosRQ")()
	ctx := s.Context()
	logger := fklog.ContextAppLogger(ctx)

	res := &MazeGame.SyncMazeStorageRolePosRS{}
	logger.CtxInfo(ctx, "OnSyncMazeStorageRolePosRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.CtxInfo(ctx, "OnSyncMazeStorageRolePosRQ end", zap.Any("res", res))
	}()

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.RolePos = req.RolePos

	userId := uint64(s.UID())
	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStorageRolePosRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	data := req.GetRolePos()
	// 道具产出信息
	err = syncmazestorageinforedis.SaveSyncMazeStorageInfo(userId, userInfo.Barrier, syncmazestorageinforedis.StorageInfo_RolePos, data)
	if err != nil {
		logger.CtxError(ctx, "OnSyncMazeStorageRolePosRQ GetCollectInfo", zap.String("rolePos", data), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	return nil
}
