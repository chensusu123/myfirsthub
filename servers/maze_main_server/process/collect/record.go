package collect

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"maze_game_server/io/kafka/mazecollectrecord"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/server/MazeCollectCache"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func PushDollMazeCollectInfoLog(ctx context.Context, userId uint64, collectInfo *MazeCollectCache.MazeCollectInfo, oldLastTime, collectTimes int64, opType int32, tradeNumber uint64, items []*MazeCommon.MazeItem, retCode int64) error {
	logger := fklog.ContextAppLogger(ctx)
	addItems := make([]string, 0)
	for _, itemInfo := range items {
		addItems = append(addItems, fmt.Sprintf("%d:%d", itemInfo.GetItemId(), itemInfo.GetCount()))
	}
	saleMsg := &mazecollectrecord.MazeCollectChgRecord{
		UserId:        userId,
		OpType:        opType,
		StartTime:     collectInfo.GetStartTime(),
		LastTime:      oldLastTime,
		NewLastTime:   collectInfo.GetLastTime(),
		AvailableTime: collectInfo.GetAvailableTime(),
		EndTime:       collectInfo.GetEndTime(),
		PeriodTime:    collectInfo.GetPeriodTime(),
		CollectTimes:  collectTimes,
		BarrierId:     collectInfo.GetBarrierId(),
		TradeNo:       tradeNumber,
		AddItems:      strings.Join(addItems, ","),
		RemainItems:   CollectItemsToString(collectInfo.Items),
		RetCode:       retCode,
	}
	if err := mazecollectrecord.PushMazeCollectChgRecord(ctx, saleMsg); err != nil {
		logger.CtxError(ctx, "PushDollMazeShopInfoLog PushDollMazeShopRecord err", zap.Error(err))
	}
	return nil
}

func CollectItemsToString(items []*MazeCollectCache.ItemInfo) string {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetId() < items[j].GetId()
	})
	s := ""
	for k, v := range items {
		if v.GetId() <= 0 || v.GetCount() <= 0 {
			continue
		}
		temp := fmt.Sprintf("%d:%d", v.GetId(), v.GetCount())

		if k != len(items)-1 {
			temp += ","
		}
		s += temp
	}
	return s
}
