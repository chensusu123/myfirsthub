/*
@Author: xiaobo
@Date: 2025/3/22 14:52
@Description:
*/

package mazecommonvalue

import (
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommonValueSvr"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/additemdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemdefine/constdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/itemdefine/errdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/game/common_value"
	"fmt"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/mazecommonvaluedb"
)

var GlobalMazeCommonValue = &class{}

func Register(reg *additemdefine.RegisterInfo) {
	reg.RegisterByItemID(constdefine.MazeSilverCoin, GlobalMazeCommonValue)
	reg.RegisterByItemID(constdefine.MazeGoldCoin, GlobalMazeCommonValue)
}

type class struct {
}

// GatherItem 加道具
func (s *class) GatherItem(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items ...*MazeCommon.MazeItem) (addRes *structdefine.AddItemRes, err error) {
	addRes = &structdefine.AddItemRes{}

	rpcRq := &MazeCommonValueSvr.MazeCommonValueAddRQ{
		UserId:      proto.Uint64(userCtx.UserID),
		AddItems:    items,
		OpType:      proto.Int32(option.OpType),
		TradeNumber: proto.Uint64(option.TradeNumber),
		Header:      userCtx.Header,
	}
	rpcRs := &MazeCommonValueSvr.MazeCommonValueAddRS{}
	// err = mazecommonvaluerpc.MazeCommonValueAddRQ(userCtx, rpcRq, rpcRs)
	err = common_value.MazeCommonValueAddRQ(userCtx, int64(userCtx.UserID), rpcRq, rpcRs)
	if err == nil && rpcRs.GetErrInfo() != nil && rpcRs.GetErrInfo().GetErrCode() != errors.NO_ERROR_CODE {
		err = fmt.Errorf("MazeCommonValue failed, %v", rpcRs.GetErrInfo().GetErrMsg())
	}
	if err != nil {
		userCtx.ErrorWF("MazeCommonValue MazeCommonValueAddRQ err", zap.Any("rpcReq", rpcRq), zap.Any("rpcRs", rpcRs), zap.Error(err))
		if errdefine.IsTimeOut(err.Error()) {
			err = errors.ITEM_ADD_TIME_OUT
			addRes.TimeoutItem = items
		} else {
			addRes.FailItem = items
		}
		return
	}

	addRes.SucItem = append(addRes.SucItem, items...)
	return
}

// GatherItemCheck 加道具检查
func (s *class) GatherItemCheck(_ fkserver.UserContext, _ *additemdefine.AddItemOption, _ []*MazeCommon.MazeItem) (*structdefine.AddItemRes, error) {
	return nil, nil
}

// DeductItem 扣道具
func (s *class) DeductItem(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (deductRes *structdefine.AddItemRes, errInfo *MessageType.ErrorInfo) {
	deductRes = &structdefine.AddItemRes{}

	rpcRq := &MazeCommonValueSvr.MazeCommonValueSubRQ{
		UserId:      proto.Uint64(userCtx.UserID),
		SubItems:    items,
		OpType:      proto.Int32(option.OpType),
		TradeNumber: proto.Uint64(option.TradeNumber),
		Header:      userCtx.Header,
	}
	rpcRs := &MazeCommonValueSvr.MazeCommonValueSubRS{}

	// err := mazecommonvaluerpc.MazeCommonValueSubRQ(userCtx, rpcRq, rpcRs)
	err := common_value.MazeCommonValueSubRQ(userCtx, int64(userCtx.UserID), rpcRq, rpcRs)
	if err != nil {
		userCtx.WarnWF("ExpireClass DeductItem err", zap.Any("rpcRq", rpcRq), zap.Error(err))
		if errdefine.IsTimeOut(err.Error()) {
			errInfo = errors.ITEM_SUB_TIME_OUT.Wrap("操作超时")
			deductRes.TimeoutItem = items
		} else {
			errInfo = errors.NewCommonCodeError(err.Error())
			deductRes.FailItem = items
		}
		return
	}

	deductRes.SucItem = items
	return
}

// DeductItemCheck 扣道具检查
func (s *class) DeductItemCheck(userCtx fkserver.UserContext, _ *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (checkRes *structdefine.AddItemRes, errInfo *MessageType.ErrorInfo) {
	checkRes = &structdefine.AddItemRes{}

	queryIds := make([]int32, 0, len(items))
	for _, item := range items {
		queryIds = append(queryIds, item.GetItemId())
	}

	dbCountMap, err := mazecommonvaluedb.BatchGetItem(userCtx, userCtx.UserID, queryIds)
	if err != nil {
		userCtx.ErrorWF("MazeCommonValue BatchGetBagItem err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	var bagCount int64
	for _, item := range items {
		bagCount = dbCountMap[item.GetItemId()]

		if bagCount < item.GetCount() {
			userCtx.WarnWF("MazeCommonValue DeductItemCheck item not enough", zap.Int64("bagCount", bagCount), zap.Any("item", item),
			)
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

	dbCountMap, err := mazecommonvaluedb.BatchGetItem(userCtx, userCtx.UserID, queryIds)
	if err != nil {
		userCtx.ErrorWF("MazeCommonValue BatchGetItem err", zap.Int32s("queryIds", queryIds), zap.Error(err))
		errInfo = errors.MODULE_ERROR.ToInfo()
		return
	}

	for _, item := range items {
		item.Count = proto.Int64(dbCountMap[item.GetItemId()])
	}
	return
}
