package collect

import (
	"sort"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/constdef"
	"go.uber.org/zap"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCollectCache"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
)

// 获取道具刷新时间
func GetFreshTime(collectInfo *MazeCollectCache.MazeCollectInfo) int64 {
	if IsLastCollect(collectInfo) { // 产出已结束
		if time.Now().Unix() < collectInfo.GetEndTime() {
			// 产出停止时间早于挂机结束时间，让客户端多查一次，虽然没有道具更新了
			// 因为客户端逻辑是收到刷新时间为0会停止挂机进度条，返回非0刷新时间避免客户端挂机进度条停止
			return collectInfo.GetEndTime()
		} else {
			return 0
		}
	} else {
		freshTime := collectInfo.GetLastTime() + int64(collectInfo.GetPeriodTime())
		return freshTime + 1
	}
}

// 本次挂机道具产出是否已结束
func IsLastCollect(collectInfo *MazeCollectCache.MazeCollectInfo) bool {
	return collectInfo.GetLastTime()+int64(collectInfo.GetPeriodTime()) > collectInfo.GetEndTime()
}

// 是否可正确领取道具
func CheckReceiveItems(collectInfo *MazeCollectCache.MazeCollectInfo) bool {
	if IsLastCollect(collectInfo) {
		return true
	}
	if time.Now().Unix()-collectInfo.GetLastTime() > int64(collectInfo.GetPeriodTime()) {
		// 还有周期未结算
		return false
	}
	return true
}

// 定时器是否丢失
func IsTimerLoss(collectInfo *MazeCollectCache.MazeCollectInfo) bool {
	capSeconds := time.Now().Unix() - collectInfo.GetLastTime()
	if !IsLastCollect(collectInfo) &&
		capSeconds > int64(collectInfo.GetPeriodTime())+TimerDelaySeconds {
		return true
	}
	return false
}

func Map2Common(m map[int32]int64) (items []*MazeCommon.MazeItem) {
	for id, count := range m {
		items = append(items, &MazeCommon.MazeItem{
			ItemId: proto.Int32(id),
			Count:  proto.Int64(count),
		})
	}
	return
}

// 获取用户可领取道具和留存道具
func GetUserItemsAndRemains(logger fklog.FKLogI, items []*MazeCollectCache.ItemInfo) (
	userItems []*MazeCommon.MazeItem, remainItems []*MazeCollectCache.ItemInfo) {
	mergeItems := make(map[int32]int64) // 合并相同道具
	for _, item := range items {
		id := item.GetId()
		entireCount := item.GetCount() / int64(constdef.ItemProducePercent) // 用户可见道具数
		debrisCount := item.GetCount() % int64(constdef.ItemProducePercent) // 用户不可见道具数

		if entireCount != 0 {
			mergeItems[id] += entireCount
		}
		if debrisCount != 0 {
			remainItem := &MazeCollectCache.ItemInfo{
				Id:    proto.Int32(id),
				Count: proto.Int64(debrisCount),
			}
			remainItems = append(remainItems, remainItem)
		}
	}

	for id, count := range mergeItems {
		userItems = append(userItems, &MazeCommon.MazeItem{
			ItemId: proto.Int32(id),
			Count:  proto.Int64(count),
		})
	}

	sort.Slice(userItems, func(i, j int) bool {
		return userItems[i].GetItemId() < userItems[j].GetItemId()
	})

	logger.InfoWF("GetUserItems end",
		zap.Any("inItems", items), zap.Any("userItems", userItems), zap.Any("remainItems", remainItems))
	return
}
