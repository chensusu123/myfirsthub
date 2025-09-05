package collect

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBarriesOnHookV8Cfg"
	"maze_game_server/pb/common/MazeCollect"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeCollectCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func MazeCollectToCliPB(ctx context.Context, collectInfo *MazeCollectCache.MazeCollectInfo, passBarrier int32) (*MazeCollect.MazeCollectInfo, error) {
	logger := fklog.ContextAppLogger(ctx)
	res := &MazeCollect.MazeCollectInfo{}
	res.StartTime = proto.Int64(collectInfo.GetStartTime())
	res.EndTime = proto.Int64(collectInfo.GetEndTime())
	res.AvailableTime = proto.Int64(collectInfo.GetAvailableTime())
	res.StageId = proto.Int32(passBarrier)
	cfg := GMazeBarriesOnHookV8Cfg.GetWithCtx(ctx, collectInfo.GetBarrierId())
	if cfg == nil {
		logger.CtxError(ctx, "MazeCollectToCliPB GMazeBarriesOnHookV8Cfg error",
			zap.Any("barrierId", collectInfo.GetBarrierId()))
		return nil, errors.New("装备详情配置不存在")
	}
	userItems, _ := GetUserItemsAndRemains(ctx, collectInfo.GetItems())
	for _, item := range userItems {
		if _, ok := cfg.Cycle_award_2[item.GetItemId()]; ok {
			res.RareItems = append(res.RareItems, item)
		} else {
			res.Items = append(res.Items, item)
		}
	}
	periodMultiple := 3600 / cfg.Cycle_time
	for k, v := range cfg.Cycle_award {
		count := v * int64(periodMultiple)
		res.CycleItems = append(res.CycleItems, &MazeCommon.MazeItem{
			ItemId: proto.Int32(k),
			Count:  proto.Int64(count),
		})
	}
	for k, v := range cfg.Cycle_award_2 {
		count := v * int64(periodMultiple)
		res.CycleItems = append(res.CycleItems, &MazeCommon.MazeItem{
			ItemId: proto.Int32(k),
			Count:  proto.Int64(count),
		})
	}
	return res, nil
}
