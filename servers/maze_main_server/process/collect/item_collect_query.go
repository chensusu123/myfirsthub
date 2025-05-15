package collect

import (
	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecollectredis"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"
	"go.uber.org/zap"
)

// 道具收集查询
func OnMazeCollectInfoQueryRQ(ctx fknet.TCPContext, userId uint64, rq proto.Message, rs proto.Message) (err error) {
	req := rq.(*MazeCollect.MazeCollectInfoQueryRQ)
	res := rs.(*MazeCollect.MazeCollectInfoQueryRS)
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR
	res.QueryType = req.QueryType

	ctx.InfoWF("OnMazeCollectInfoQueryRQ start", zap.Any("req", req))
	defer func() {
		ctx.InfoWF("OnMazeCollectInfoQueryRQ end", zap.Any("res", res))
	}()

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetUserInfoV2", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return err
	}

	// 道具产出信息
	collectInfo, err := mazecollectredis.GetCollectInfo(ctx, userId)
	if err != nil {
		ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	if collectInfo == nil {
		if userInfo.PassBarrier <= 0 {
			return nil
		}
		err = InitMazeCollectLand(ctx, userId, userInfo.PassBarrier)
		if err != nil {
			ctx.ErrorWF("OnMazeCollectInfoQueryRQ InitMazeCollectLand", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
		// 道具产出信息
		collectInfo, err = mazecollectredis.GetCollectInfo(ctx, userId)
		if err != nil {
			ctx.ErrorWF("OnMazeCollectInfoQueryRQ GetCollectInfo", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return err
		}
	}
	if collectInfo == nil {
		return nil
	}
	if IsTimerLoss(collectInfo) {
		// 定时器丢失修复道具产出
		ctx.InfoWF("OnMazeCollectInfoQueryRQ fix ItemCollect start", zap.Any("collectInfo", collectInfo))
		err = ItemCollect(ctx, userId, collectInfo)
		if err != nil {
			ctx.ErrorWF("OnPetCollectInfoQueryRQ fix ItemCollect", zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return
		}
		ctx.InfoWF("OnPetCollectInfoQueryRQ fix ItemCollect end", zap.Any("collectInfo", collectInfo))
	}

	mazeCollectInfoPb, err := MazeCollectToCliPB(ctx, collectInfo, userInfo.PassBarrier)
	if err != nil {
		ctx.ErrorWF("ItemCollect MazeCollectToCliPB err", zap.Any("collectInfo", collectInfo), zap.Error(err))
		return err
	}
	res.MazeCollectInfo = mazeCollectInfoPb
	res.FreshTime = proto.Int64(GetFreshTime(collectInfo))
	return
}
