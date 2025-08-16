package barrierservice

import (
	"fmt"
	"maze_game_server/common/errors"
	"maze_game_server/config/GMazeFoeV8Cfg"
	"maze_game_server/io/redis/mazebarrieropstatusredis"
	"maze_game_server/pb/common/MessageType"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// GuardDeath implements BarrierService.
func (b *barrier) GuardDeath(logger fklog.FKLogI, userID uint64, barrierID int32, monsterID int32, monsterGuid int32) (kongfu int32, equips map[int32]int32, items map[int32]int64, errinfo *MessageType.ErrorInfo) {
	foeCfg := GMazeFoeV8Cfg.Get(monsterID)
	if foeCfg == nil {
		logger.ErrorWF("GuardDeath get foe cfg fail", zap.Any("monsterID", monsterID))
		return 0, nil, nil, errors.CONFIG_NOT_FOUND.ToInfo()
	}

	// 防重复操作校验
	triggered, triggerFn, err := mazebarrieropstatusredis.IsTriggered(logger, userID, barrierID, fmt.Sprintf("monsterid:%d", monsterGuid))
	if err != nil {
		logger.ErrorWF("GuardDeath IsTriggered fail", zap.Error(err), zap.Any("MonsterId", monsterID), zap.Any("barrierId", barrierID))
		return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("数据校验失败")
	}
	if triggered {
		logger.WarnWF("GuardDeath already killed", zap.Any("MonsterId", monsterID), zap.Any("barrierId", barrierID))
		return 0, nil, nil, errors.COMMON_ERROR_TIPS.Wrap("怪物已击杀")
	} else {
		defer triggerFn()
	}

	// 通关值
	kongfu = foeCfg.Kongfu

	// 怪物掉落装备
	equips = make(map[int32]int32)
	for _, v := range foeCfg.Drop_equip {
		if v > 0 {
			equips[v] += 1
		}
	}

	// 怪物掉落道具
	items = make(map[int32]int64)
	for itemID, count := range foeCfg.Drop_item {
		if itemID > 0 && count > 0 {
			items[itemID] = count
		}
	}

	return kongfu, equips, items, errors.NO_ERROR
}
