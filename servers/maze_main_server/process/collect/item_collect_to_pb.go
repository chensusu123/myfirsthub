package collect

import (
	"gitlab.ifreetalk.com/plate/excel/auto/GMazeBarriesOnHookV8Cfg"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollect"
	"gitlab.ifreetalk.com/plate/protodef/MazeCollectCache"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"go.uber.org/zap"
)

func MazeCollectToCliPB(logger fklog.FKLogI, collectInfo *MazeCollectCache.MazeCollectInfo, passBarrier int32) (*MazeCollect.MazeCollectInfo, error) {
	res := &MazeCollect.MazeCollectInfo{}
	res.StartTime = proto.Int64(collectInfo.GetStartTime())
	res.EndTime = proto.Int64(collectInfo.GetEndTime())
	res.AvailableTime = proto.Int64(collectInfo.GetAvailableTime())
	res.StageId = proto.Int32(passBarrier)
	cfg := GMazeBarriesOnHookV8Cfg.Get(collectInfo.GetBarrierId())
	if cfg == nil {
		logger.ErrorWF("MazeCollectToCliPB GMazeBarriesOnHookV8Cfg error",
			zap.Any("barrierId", collectInfo.GetBarrierId()))
		return nil, errors.New("装备详情配置不存在")
	}
	userItems, _ := GetUserItemsAndRemains(logger, collectInfo.GetItems())
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
