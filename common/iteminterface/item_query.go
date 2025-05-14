/*
@Author: xiaobo
@Date: 2023/11/30 19:28
@Description: 查询物品接口
*/

package iteminterface

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/additemdefine"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze-plate/protodef/MessageType"

	"go.uber.org/zap"
)

func QueryItems(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (resItems []*MazeCommon.MazeItem, errInfo *MessageType.ErrorInfo) {
	if option == nil {
		errInfo = errors.NewCommonCodeError("option nil")
		return
	}

	classifiedItems, err := option.RegIns.GetClassifiedItems(userCtx, items...)
	if err != nil {
		return items, errors.NewCommonCodeError(err.Error())
	}

	for handleProcess, queryItems := range classifiedItems {
		if handleProcess == nil || len(queryItems) == 0 {
			userCtx.ErrorWF("QueryItems not fount gatherItemFunc", zap.Any("queryItems", queryItems))
			errInfo = errors.NewCommonCodeError("handle class nil")
			return
		}
		errInfo = handleProcess.GetItem(userCtx, option, queryItems)
		if errInfo != nil && errInfo.GetErrCode() != errors.NO_ERROR.GetErrCode() {
			userCtx.WarnWF("QueryItems fail", zap.Any("items", items), zap.Any("errInfo", errInfo))
			return
		}

		if len(queryItems) > 0 {
			resItems = append(resItems, queryItems...)
		}
	}
	return
}
