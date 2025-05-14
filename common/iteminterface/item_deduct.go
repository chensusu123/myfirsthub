/*
@Author: xiaobo
@Date: 2023/11/30 17:40
@Description: 扣物品接口
*/

package iteminterface

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/additemdefine"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structdefine"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"
	"go.uber.org/zap"
)

// DeductItems 扣物品
func DeductItems(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items ...*MazeCommon.MazeItem) (deductRes *structdefine.AddItemRes, errInfo *MessageType.ErrorInfo) {
	if len(items) == 0 { // 这里不会走到，方法调用放会检查 itemId 和count
		return
	}

	deductRes = &structdefine.AddItemRes{}

	classifiedItems, err := option.RegIns.GetClassifiedItems(userCtx.FKLogI, items...)
	if err != nil {
		errInfo = errors.NewCommonCodeError(err.Error())
		deductRes.FailItem = append(deductRes.FailItem, items...)
		return
	}

	var checkDuctRes *structdefine.AddItemRes
	for handleProcess, opItems := range classifiedItems {
		if handleProcess == nil {
			userCtx.WarnWF("DeductItems cannot find deduct check func", zap.Any("opItems", opItems), zap.Any("option", option))
			errInfo = errors.ITEM_CHECK_ERROR.ToInfo()
			return
		}

		checkDuctRes, errInfo = handleProcess.DeductItemCheck(userCtx, option, opItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR_CODE {
			userCtx.WarnWF("DeductItems check deduce item fail", zap.Any("option", option), zap.Any("opItems", opItems), zap.Any("errInfo", errInfo))
			deductRes.Merge(checkDuctRes)
			return
		}
	}

	var deductItems *structdefine.AddItemRes
	for handProcess, opItems := range classifiedItems {
		deductItems, errInfo = handProcess.DeductItem(userCtx, option, opItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR.GetErrCode() {
			userCtx.WarnWF("DeductItems cost fail", zap.Any("items", opItems), zap.Any("errInfo", errInfo))
			deductRes.FailItem = append(deductRes.FailItem, deductItems.FailItem...)
			deductRes.TimeoutItem = append(deductRes.TimeoutItem, deductItems.TimeoutItem...)
			deductRes.LessItem = append(deductRes.LessItem, deductItems.LessItem...)
			return
		}

		// 扣除成功的道具
		deductRes.SucItem = append(deductRes.SucItem, opItems...)
	}
	return
}
