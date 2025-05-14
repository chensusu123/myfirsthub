package collect

import (
	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCollect"
	"go.uber.org/zap"
)

// 道具收集查询
func (c *Collect) OnMazeCollectInfoQueryRQ(s *session.Session, req *MazeCollect.MazeCollectInfoQueryRQ) (err error) {
	logger := fklog.AppLogger().Clone("collect")
	res := &MazeCollect.MazeCollectInfoQueryRS{}

	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.QueryType = req.QueryType

	userId := uint64(s.UID())

	logger.InfoWF("OnMazeCollectInfoQueryRQ start", zap.Any("req", req))
	defer func() {
		err = s.Response(res)
		logger.InfoWF("OnMazeCollectInfoQueryRQ end", zap.Any("res", res))
	}()

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
