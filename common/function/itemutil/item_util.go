package itemutil

import (
	"fmt"
	"sort"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/uniqueid"
	"gitlab.ifreetalk.com/maze-plate/protodef/Common"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
)

func Common2Map(attrs []*Common.Attr) (m map[int32]int64) {
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

func GetTradeNum() (tradeNum uint64) {
	m := uniqueid.NewTradeNoMaker(uint64(fkconfig.GetServerConfig().ServerID))
	if m == nil {
		return
	}
	tradeNum = m.MakeTradeNo()
	return
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
