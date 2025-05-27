package collect

import (
	"maze_game_server/common/errors"
	"maze_game_server/io/redis/mazecollectredis"
	"maze_game_server/lib/log"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeCollect"

	"github.com/lonng/nano/session"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 道具收集查询
func (c *Collect) OnMazeCollectInfoQueryRQ_10465_10466(s *session.Session, req *MazeCollect.MazeCollectInfoQueryRQ) (err error) {

	logger := log.Clone("Collect", uint64(s.UID()), 0)
	res := &MazeCollect.MazeCollectInfoQueryRS{}
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.QueryType = req.QueryType

	logger.InfoWF("OnMazeCollectInfoQueryRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeCollectInfoQueryRQ end", zap.Any("res", res))
	}()

	userId := uint64(s.UID())

	userInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeCollectInfoQueryRQ GetUserInfoV2", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(logger, userId)
	if err != nil {
		logger.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if collectInfo == nil {
		if userInfo.PassBarrier <= 0 {
			return nil
		}
		err = InitMazeCollectLand(logger, userId, userInfo.PassBarrier)
		if err != nil {
			logger.ErrorWF("OnMazeCollectInfoQueryRQ InitMazeCollectLand", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
		// 道具产出信息
		collectInfo, err = mazecollectredis.GetCollectInfo(logger, userId)
		if err != nil {
			logger.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
	}
	if collectInfo == nil {
		return nil
	}
	if IsTimerLoss(collectInfo) {
		// 定时器丢失修复道具产出
		logger.InfoWF("OnMazeCollectInfoQueryRQ fix ItemCollect start", zap.Any("collectInfo", collectInfo))
		err = ItemCollect(logger, userId, collectInfo)
		if err != nil {
			logger.ErrorWF("OnPetCollectInfoQueryRQ fix ItemCollect", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		logger.InfoWF("OnPetCollectInfoQueryRQ fix ItemCollect end", zap.Any("collectInfo", collectInfo))
	}

	mazeCollectInfoPb, err := MazeCollectToCliPB(logger, collectInfo, userInfo.PassBarrier)
	if err != nil {
		logger.ErrorWF("ItemCollect MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return err
	}
	res.MazeCollectInfo = mazeCollectInfoPb
	res.FreshTime = proto.Int64(GetFreshTime(collectInfo))
	return
}
