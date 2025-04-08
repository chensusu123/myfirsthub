// @Author pangchenyang 2025/3/22 15:35:00
// @Desc: 
package process

import (
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
)

func ItemsMapToList(itemsMap map[int32]int64) []*Common.Item {
	itemList := make([]*Common.Item, 0)
	for itemId, count := range itemsMap {
		itemList = append(itemList, &Common.Item{
			ItemId: proto.Int32(itemId),
			Count:  proto.Int64(count),
		})
	}
	return itemList
}
