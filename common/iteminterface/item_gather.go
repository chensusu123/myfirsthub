/*
@Author: xiaobo
@Date: 2023/11/30 17:44
@Description: 加物品接口
*/

package iteminterface

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"go.uber.org/zap"
	"maze_game_server/common/additemdefine"
	"maze_game_server/common/structdefine"
	"maze_game_server/pb/common/MazeCommon"
)

// CheckAddItems 加检查
func CheckAddItems(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items []*MazeCommon.MazeItem) (checkItemRes *structdefine.AddItemRes, err error) {
	checkItemRes = &structdefine.AddItemRes{}

	classifiedItems, err := option.RegIns.GetClassifiedItems(userCtx.FKLogI, items...)
	if err != nil {
		checkItemRes.FailItem = append(checkItemRes.FailItem, items...)
		return
	}

	var checkRes *structdefine.AddItemRes
	for handleProcess, opItems := range classifiedItems {
		checkRes, err = handleProcess.GatherItemCheck(userCtx, option, opItems)
		if err != nil {
			userCtx.WarnWF("CheckAddItems GatherItem error", zap.Any("items", opItems),
				zap.Any("checkRes", checkRes), zap.Error(err),
			)
			checkItemRes.Merge(checkRes)
			return
		}
	}

	return
}

// GatherItems 检查并添加物品
func GatherItems(userCtx fkserver.UserContext, option *additemdefine.AddItemOption, items ...*MazeCommon.MazeItem) (addRes *structdefine.AddItemRes, err error) {
	addRes = &structdefine.AddItemRes{}

	classifiedItems, err := option.RegIns.GetClassifiedItems(userCtx.FKLogI, items...)
	if err != nil {
		addRes.FailItem = append(addRes.FailItem, items...)
		return
	}

	var checkRes *structdefine.AddItemRes
	checkResMap := make(map[additemdefine.ItemClassProcess]*structdefine.AddItemRes, len(classifiedItems))
	for handleProcess, opItems := range classifiedItems {
		checkRes, err = handleProcess.GatherItemCheck(userCtx, option, opItems)
		if err != nil {
			userCtx.WarnWF("GatherItems GatherItem error", zap.Any("items", opItems),
				zap.Any("checkRes", checkRes), zap.Error(err),
			)
			addRes.Merge(checkRes)
			return
		}
		checkResMap[handleProcess] = checkRes
	}

	option.CheckResMap = checkResMap

	var (
		gatherItems *structdefine.AddItemRes
		addErr      error
	)
	for handleProcess, opItems := range classifiedItems {
		gatherItems, addErr = handleProcess.GatherItem(userCtx, option, opItems...)
		if addErr != nil {
			userCtx.WarnWF("GatherItems GatherItem error", zap.Any("items", opItems), zap.Any("gatherRes", gatherItems), zap.Error(addErr))
			err = addErr
		}
		addRes.Merge(gatherItems)
	}

	return
}
