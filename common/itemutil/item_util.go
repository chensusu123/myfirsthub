package itemutil

import (
	"context"

	"maze_game_server/config/GMazeBagOrderV8Cfg"
	"maze_game_server/pb/common/Common"
	"maze_game_server/pb/common/MazeBag"
	"maze_game_server/pb/common/MazeCommon"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// WrapUserContext 封装UserContext信息
func WrapUserContext(ctx context.Context, uid uint64, logger fklog.FKLogI, header *Common.PacketHeader) fkserver.UserContext {
	if header == nil {
		header = &Common.PacketHeader{}
	}

	userCtx := fkserver.NewUserContext(ctx, uid, logger)
	// userCtx.Header = header
	return userCtx
}

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

// CheckAndMergeItem 检查并合并物品
func CheckAndMergeItem(input []*MazeCommon.MazeItem) (output []*MazeCommon.MazeItem) {
	if len(input) == 0 {
		return
	}

	itemMap := make(map[int32]*MazeCommon.MazeItem, len(input))
	output = make([]*MazeCommon.MazeItem, 0, len(input))

	var (
		tmpItem *MazeCommon.MazeItem
		isExist bool
	)
	for _, item := range input {
		if item.GetItemId() <= 0 || item.GetCount() <= 0 {
			continue
		}

		if tmpItem, isExist = itemMap[item.GetItemId()]; !isExist {
			itemMap[item.GetItemId()] = &MazeCommon.MazeItem{
				ItemId: proto.Int32(item.GetItemId()),
				Count:  proto.Int64(item.GetCount()),
			}
			continue
		}

		itemMap[item.GetItemId()].Count = proto.Int64(tmpItem.GetCount() + item.GetCount())
	}

	for _, item := range itemMap {
		output = append(output, item)
	}
	return
}

func BuildItemIds(input []*MazeCommon.MazeItem) (output []int32) {
	if len(input) == 0 {
		return
	}

	output = make([]int32, 0, len(input))
	for _, item := range input {
		output = append(output, item.GetItemId())
	}
	return
}

// MergeQueryItems 合并查询物品
func MergeQueryItems(input []*MazeCommon.MazeItem) (output []*MazeCommon.MazeItem) {
	if len(input) == 0 {
		return
	}

	itemMap := make(map[int32]struct{}, len(input))
	for _, item := range input {
		if item.GetItemId() <= 0 {
			continue
		}

		itemMap[item.GetItemId()] = struct{}{}
	}

	for id := range itemMap {
		output = append(output, &MazeCommon.MazeItem{
			ItemId: proto.Int32(id),
		})
	}
	return
}
