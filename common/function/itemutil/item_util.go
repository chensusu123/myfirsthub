package itemutil

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/config/GMazeBagOrderV8Cfg"
	"maze_game_server/pb/common/MazeBag"
	"maze_game_server/services/itemservice"
	"sort"

	"maze_game_server/pb/common/MazeCommon"

	"google.golang.org/protobuf/proto"
)

func BuildMazeBagItem(logger fklog.FKLogI, id int32, count int64) (bagItem *MazeBag.MazeBagItem) {
	bagItem = &MazeBag.MazeBagItem{
		ItemId: proto.Int32(id),
		Count:  proto.Int64(count),
	}

	cfg := GMazeBagOrderV8Cfg.Get(id)
	if cfg == nil {
		logger.WarnWF("BuildMazeBagItem GMazeBagOrderV8Cfg nil", zap.Int32("id", id))
		return
	}

	bagItem.Quality = proto.Int32(cfg.Quality)
	bagItem.OrderType = proto.Int32(cfg.Order_type_2)
	return
}
func Common2Map(attrs []*MazeCommon.Attr) (m map[int32]int64) {
	m = make(map[int32]int64)
	for _, attr := range attrs {
		if attr.GetAttrId() <= 0 {
			continue
		}
		m[attr.GetAttrId()] += int64(attr.GetAttrValue())
	}
	return
}

func Map2Common(m map[int32]int64) (items []*MazeCommon.MazeItem) {
	for id, count := range m {
		if id <= 0 || count <= 0 {
			continue
		}
		items = append(items, &MazeCommon.MazeItem{
			ItemId: proto.Int32(id),
			Count:  proto.Int64(count),
		})
	}
	return
}

func CheckItemMap(m map[int32]int64) map[int32]int64 {
	validM := make(map[int32]int64)
	for k, v := range m {
		if k == 0 || v == 0 {
			continue
		}
		validM[k] = v
	}
	return validM
}

func CheckItemMatch(items []*MazeCommon.MazeItem, cost []*MazeCommon.MazeItem) bool {
	if len(items) != len(cost) {
		return false
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetItemId() < items[j].GetItemId()
	})
	sort.Slice(cost, func(i, j int) bool {
		return cost[i].GetItemId() < cost[j].GetItemId()
	})

	for i := range items {
		if items[i].GetItemId() == 0 || items[i].GetCount() == 0 ||
			cost[i].GetItemId() == 0 || cost[i].GetCount() == 0 {
			return false
		}
		if items[i].GetItemId() != cost[i].GetItemId() || items[i].GetCount() != cost[i].GetCount() {
			return false
		}
	}

	return true
}

func CheckItemMatchEx(items []*MazeCommon.MazeItem, cost []*MazeCommon.MazeItem) bool {
	if len(items) != len(cost) {
		return false
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetItemId() < items[j].GetItemId()
	})
	sort.Slice(cost, func(i, j int) bool {
		return cost[i].GetItemId() < cost[j].GetItemId()
	})

	for i := range items {
		if items[i].GetItemId() == 0 || cost[i].GetItemId() == 0 {
			return false
		}

		if items[i].GetItemId() != cost[i].GetItemId() || items[i].GetCount() != cost[i].GetCount() {
			return false
		}
	}

	return true
}

// 道具转为字符串
func CommonItemsToString(items []*MazeCommon.MazeItem) string {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetItemId() < items[j].GetItemId()
	})

	s := ""

	for k, v := range items {
		if v.GetItemId() <= 0 || v.GetCount() <= 0 {
			continue
		}
		temp := fmt.Sprintf("%d:%d", v.GetItemId(), v.GetCount())

		if k != len(items)-1 {
			temp += ","
		}
		s += temp
	}
	return s
}

func Map2ItemInfo(m map[int32]int64) (items []*itemservice.ItemInfo) {
	for id, count := range m {
		if id <= 0 || count <= 0 {
			continue
		}
		items = append(items, &itemservice.ItemInfo{
			ItemId: id,
			Count:  count,
		})
	}
	return
}

func ItemPb2ItemInfo(items []*MazeCommon.MazeItem) []*itemservice.ItemInfo {
	res := make([]*itemservice.ItemInfo, 0, len(items))
	for _, i := range items {
		res = append(res, &itemservice.ItemInfo{
			ItemId: i.GetItemId(),
			Count:  i.GetCount(),
		})
	}
	return res
}
