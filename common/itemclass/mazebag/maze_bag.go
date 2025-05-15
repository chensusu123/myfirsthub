/*
@Author: xiaobo
@Date: 2025/3/21 11:42
@Description:
*/

package mazebag

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/errors"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeBag"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/additemdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemdefine/constdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemdefine/errdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazebagdb"
	"gitlab.ifreetalk.com/maze/maze_game_server/usecase/mustarrive"
	"go.uber.org/zap"
)

var GlobalMazeBag = &class{}

func Register(reg *additemdefine.RegisterInfo) {
	reg.RegisterByBagType(constdefine.EnterTypeMazeBag, GlobalMazeBag)
}

type class struct{}

// GatherItem 加道具
func (s *class) GatherItem(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items ...*MazeCommon.MazeItem) (addRes *structdefine.AddItemRes, err error) {
	addRes = &structdefine.AddItemRes{}

	var sucItems []*MazeCommon.MazeItem
	for idx, item := range items {
		currItemCount, addErr := mazebagdb.IncrBagItem(userCtx, userCtx.UserID, item.GetItemId(), item.GetCount())
		if addErr != nil {
			userCtx.ErrorWF("MazeBag IncrBagItem err",
				zap.Int32("itemID", item.GetItemId()),
				zap.Int64("count", item.GetCount()),
				zap.Error(addErr),
			)

			if errdefine.IsTimeOut(addErr.Error()) {
				err = errors.ITEM_ADD_TIME_OUT
				addRes.TimeoutItem = append(addRes.TimeoutItem, items[idx])
			} else {
				err = addErr
				addRes.FailItem = append(addRes.FailItem, items[idx])
			}
			continue
		}

		addRes.SucItem = append(addRes.SucItem, items[idx])

		sucItems = append(sucItems, &MazeCommon.MazeItem{
			ItemId: item.ItemId,
			Count:  proto.Int64(currItemCount),
		})
	}

	// 变化通知ID包
	SendBagItemChgID(userCtx, sucItems)

	// 流水 TODO
	return
}

// GatherItemCheck 加道具检查
func (s *class) GatherItemCheck(_ fkserver.UserContext, _ *additemdefine.AddItemOption, _ []*MazeCommon.MazeItem) (*structdefine.AddItemRes, error) {
	return nil, nil
}

// DeductItem 扣道具
func (s *class) DeductItem(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (deductRes *structdefine.AddItemRes, errInfo *MessageType.ErrorInfo) {
	deductRes = &structdefine.AddItemRes{}

	var (
		sucItems       []*MazeCommon.MazeItem // 推id包使用
		removeItemList []int32
	)
	for _, item := range items {
		curCount, subErr := mazebagdb.IncrBagItem(userCtx, userCtx.UserID, item.GetItemId(), -item.GetCount())
		if subErr != nil {
			userCtx.ErrorWF("MazeBag IncrBagItem err", zap.Any("item", item), zap.Error(subErr))
			if errdefine.IsTimeOut(subErr.Error()) {
				errInfo = errors.ITEM_SUB_TIME_OUT.ToInfo()
				deductRes.TimeoutItem = append(deductRes.TimeoutItem, item)
			} else {
				errInfo = errors.DB_SAVE_ERROR.ToInfo()
				deductRes.FailItem = append(deductRes.FailItem, item)
			}
			continue
		}

		if curCount < 0 {
			// 添加修复道具流水 TODO
			removeItemList = append(removeItemList, item.GetItemId())
			deductRes.LessItem = append(deductRes.LessItem, item)
		} else {
			sucItems = append(sucItems, &MazeCommon.MazeItem{
				ItemId: proto.Int32(item.GetItemId()),
				Count:  proto.Int64(curCount),
			})
			deductRes.SucItem = append(deductRes.SucItem, item)
		}
	}

	if len(removeItemList) > 0 {
		err := mazebagdb.BatchDelBagItem(userCtx, userCtx.UserID, removeItemList)
		if err != nil {
			userCtx.ErrorWF("MazeBag BatchDelBagItem err", zap.Any("removeItemList", removeItemList), zap.Error(err))
		}
	}

	// 变化通知ID包
	SendBagItemChgID(userCtx, sucItems)

	// 流水 TODO
	return
}

// DeductItemCheck 扣道具检查
func (s *class) DeductItemCheck(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (checkRes *structdefine.AddItemRes, errInfo *MessageType.ErrorInfo) {
	checkRes = &structdefine.AddItemRes{}

	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.GetItemId())
	}

	dbCountMap, err := mazebagdb.BatchGetBagItem(userCtx, userCtx.UserID, queryIds)
	if err != nil {
		userCtx.ErrorWF("MazeBag BatchGetBagItem err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var bagCount int64
	for _, item := range items {
		bagCount = dbCountMap[item.GetItemId()]

		if bagCount < item.GetCount() {
			userCtx.WarnWF("MazeBag DeductItemCheck item not enough", zap.Int64("bagCount", bagCount), zap.Any("item", item))
			errInfo = errors.ITEM_CHECK_ITEM_NOT_ENOUGH.ToInfo()
			continue
		}
	}
	return
}

// GetItem 查询数据
func (s *class) GetItem(userCtx fkserver.UserContext, _ *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (errInfo *MessageType.ErrorInfo) {
	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.GetItemId())
	}

	dbCountMap, err := mazebagdb.BatchGetBagItem(userCtx, userCtx.UserID, queryIds)
	if err != nil {
		userCtx.ErrorWF("MazeBag BatchGetBagItem err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	for _, item := range items {
		item.Count = proto.Int64(dbCountMap[item.GetItemId()])
	}
	return
}

func SendBagItemChgID(userCtx fkserver.UserContext, items []*MazeCommon.MazeItem) {
	if len(items) == 0 {
		return
	}

	idPack := &MazeBag.MazeBagChgID{
		Token: proto.Int64(time.Now().UnixMilli()),
		Items: make([]*MazeBag.MazeBagItem, 0, len(items)),
	}
	for _, item := range items {
		idPack.Items = append(idPack.Items, itemutil.BuildMazeBagItem(userCtx, item.GetItemId(), item.GetCount()))
	}

	_ = mustarrive.SendArrivePacket(userCtx, int64(userCtx.UserID), 16250, idPack)
}
