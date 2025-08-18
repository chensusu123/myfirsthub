package barrierservice

import (
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeBoxV8Cfg"
	"maze_game_server/io/redis/mazebarrieropstatusredis"
	"maze_game_server/io/redis/mazeboxredis"
	"maze_game_server/pb/common/MessageType"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// OpenBox implements BarrierService.
func (b *barrier) OpenBox(logger fklog.FKLogI, userID uint64, barrierID int32, boxID int32) (
	kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo) {
	boxCfg := GMazeBoxV8Cfg.Get(boxID)
	if boxCfg == nil {
		logger.ErrorWF("OpenBox get box cfg fail", zap.Any("boxId", boxID))
		return 0, nil, nil, errors.CONFIG_NOT_FOUND.ToInfo()
	}

	// // 更新box表格 读表校验宝箱对应的关卡id
	// if fkutil.ToInt64(boxCfg.Level_id) != int64(barrierID) {
	// 	logger.ErrorWF("OpenBox barrier and box not match", zap.Any("boxId", boxID), zap.Any("barrierId", barrierID))
	// 	return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("宝箱关卡信息不匹配")
	// }

	// 防重复操作校验
	triggered, triggerFn, err := mazebarrieropstatusredis.IsTriggered(logger, userID, barrierID, fmt.Sprintf("openbox:%d", boxID))
	if err != nil {
		logger.ErrorWF("OpenBox IsTriggered fail", zap.Error(err), zap.Any("boxId", boxID), zap.Any("barrierId", barrierID))
		return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("数据校验失败")
	}
	if triggered {
		logger.WarnWF("OpenBox already opened", zap.Any("boxId", boxID), zap.Any("barrierId", barrierID))
		return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("宝箱已打开")
	} else {
		defer triggerFn()
	}
	opened, err := mazeboxredis.IsOpenedBox(logger, userID, barrierID, boxID)
	if err != nil {
		logger.ErrorWF("OpenBox IsOpenedBox fail", zap.Error(err), zap.Any("boxId", boxID), zap.Any("barrierId", barrierID))
		return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("宝箱打开失败")
	}

	// 奖励通关值
	kongfu = 0 //boxCfg.Add_kongfu

	var (
		awardEquip []int32
	)
	if opened > 0 {
		awardEquip = boxCfg.Award_equip
		items = boxCfg.Drop_item
	} else {
		awardEquip = boxCfg.Award_equip_first
		items = boxCfg.Drop_item_first
	}

	// 宝箱掉落装备
	equips = make(map[int32]int32)
	for _, v := range awardEquip {
		if v > 0 {
			equips[v] += 1
		}
	}

	// 标记宝箱已打开过
	err = mazeboxredis.SetOpenBoxTime(logger, userID, barrierID, boxID)
	if err != nil {
		logger.ErrorWF("OpenBox SetOpenBoxTime fail", zap.Error(err), zap.Any("boxId", boxID), zap.Any("barrierId", barrierID))
	}

	return kongfu, equips, items, errors.NO_ERROR
}
