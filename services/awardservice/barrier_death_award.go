package awardservice

import (
	"context"
	"errors"
	"maze_game_server/config/GMazeBarriesV8Cfg"
	"maze_game_server/config/GMazeBoxV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/config/GMazeShopV8Cfg"
	"maze_game_server/io/redis/mazeboxredis"
	"maze_game_server/module/calsweepbarrier"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/services/barrierscorerewardservice"
	"maze_game_server/services/equipdropservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) GetBarrierDeathAward(ctx context.Context, userId uint64, barrier int32) (realItemMap map[int32]int64, showItemMap map[int32]int64,
	realEquipMap map[int32]int32, showEquipMap map[int32]int32, expCount int64, err error) {

	logger := fklog.ContextAppLogger(ctx)
	// 获取当前关卡的扫荡奖励
	equipItem, ohterItem, expItem, equipNum, err := calsweepbarrier.GetSweepBarrierAward(ctx, userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierDeathAward GetSweepBarrierAward err", zap.Error(err))
		return nil, nil, nil, nil, 0, err
	}
	logger.CtxInfo(ctx, "GetBarrierDeathAward GetSweepBarrierAward", zap.Any("equipItem", equipItem), zap.Any("ohterItem", ohterItem), zap.Any("expItem", expItem), zap.Any("equipNum", equipNum))

	// 获取存储的当前关卡的奖励数据
	nowBarrierEquipList, nowBarrierItemList, err := barrierscorerewardservice.GlobalScoreRewardService.GetBarrierScoreReward(context.TODO(), userId, barrier)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierDeathAward GetBarrierScoreReward err", zap.Error(err),
			zap.Any("barrier", barrier),
			zap.Any("userId", userId),
			zap.Any("nowBarrierEquipList", nowBarrierEquipList),
			zap.Any("nowBarrierItemList", nowBarrierItemList),
		)
		return nil, nil, nil, nil, 0, err
	}
	logger.CtxInfo(ctx, "GetBarrierDeathAward GetBarrierScoreReward", zap.Any("nowBarrierEquipList", nowBarrierEquipList), zap.Any("nowBarrierItemList", nowBarrierItemList))

	// 先获取当前用户关卡内打开过的宝箱
	barrierCfg := GMazeBarriesV8Cfg.Get(barrier)
	if barrierCfg == nil {
		logger.CtxError(ctx, "GetBarrierDeathAward get box cfg fail", zap.Any("barrier", barrier))
		return nil, nil, nil, nil, 0, errors.New("box cfg nil")
	}

	boxEquipCount := 0

	boxCfg := GMazeBoxV8Cfg.Get(barrierCfg.Box_id)
	if boxCfg == nil {
		logger.CtxError(ctx, "GetBarrierDeathAward get box cfg fail", zap.Any("boxId", barrierCfg.Box_id))
		return nil, nil, nil, nil, 0, errors.New("box cfg nil")
	}

	opened, err := mazeboxredis.IsOpenedBox(ctx, userId, barrier, barrierCfg.Box_id)
	if err != nil {
		logger.CtxError(ctx, "GetBarrierDeathAward IsOpenedBox fail", zap.Error(err), zap.Any("boxId", barrierCfg.Box_id), zap.Any("barrierId", barrier))
		return nil, nil, nil, nil, 0, err
	}

	// 首次打开宝箱奖励 转化为非首次
	if opened <= 0 {
		// 物品奖励
		for k, v := range boxCfg.Drop_item {
			_, ok := ohterItem[k]
			if k > 0 && ok {
				ohterItem[k] -= v
			} else if k > 0 && !ok {
				logger.CtxError(ctx, "GetBarrierDeathAward get box cfg fail", zap.Any("boxId", barrierCfg.Box_id), zap.Any("itemId", k))
				return nil, nil, nil, nil, 0, errors.New("item not found")
			}
		}

		// 装备奖励
		for _, v := range boxCfg.Award_equip {
			if v > 0 {
				equipItem[int32(v)] -= 1
				boxEquipCount += 1
			}
		}
	}

	// 装备随机取整
	equipNum -= int32(boxEquipCount)
	sum := int64(equipNum) * int64(GMazeConfigV8Cfg.Get(912).Value_int)
	nowequipNum := int32(sum / 10000)
	probability := sum % 10000
	if probability > 0 {
		tmp := fkutil.RandInt(1, 10000)
		if tmp <= int(probability) {
			nowequipNum += 1
		}
	}

	if nowequipNum > 0 {
		// 获取用户信息
		userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userId)
		if err != nil {
			logger.CtxError(ctx, "CalUserSweepBarrierAward GetUserInfoV2 fail", zap.Error(err))
			return nil, nil, nil, nil, 0, err
		}

		calLv := equipdropservice.GlobalEquipDropService.GetMazeBarrierLv(int32(userInfo.Level), barrier)
		shopCfg := GMazeShopV8Cfg.Get(calLv)
		if shopCfg == nil {
			logger.CtxError(ctx, "GetSweepBarrierAward get shop cfg fail", zap.Any("calLv", calLv))
			err = errors.New("shop cfg nil")
			return nil, nil, nil, nil, 0, err
		}

		addEquipMap, err := equipdropservice.GlobalEquipDropService.GetNewEquip(ctx, userId, calLv, barrier, nowequipNum)
		if err != nil {
			logger.CtxError(ctx, "GetSweepBarrierAward GetNewEquip fail", zap.Error(err), zap.Any("barrier", barrier), zap.Any("calLv", calLv))
			return nil, nil, nil, nil, 0, err
		}

		for k, v := range addEquipMap {
			equipItem[k] += v
		}
	}

	realEquipMap = make(map[int32]int32)
	// 记录失败新增装备奖励
	for k, v := range equipItem {
		realEquipMap[k] = v
	}

	for k, v := range nowBarrierEquipList {
		equipItem[k] += v
	}

	// 道具
	for k, v := range ohterItem {
		tmpSum := int64(v) * int64(GMazeConfigV8Cfg.Get(912).Value_int)
		ohterItem[k] = tmpSum / 10000
	}

	realItemMap = make(map[int32]int64)
	// 记录失败新增物品奖励
	for k, v := range ohterItem {
		realItemMap[k] += v
	}

	for k, v := range nowBarrierItemList {
		ohterItem[k] += v
	}

	// 过滤
	for k, v := range realItemMap {
		if v == 0 {
			delete(realItemMap, k)
		}
	}

	for k, v := range ohterItem {
		if v == 0 {
			delete(ohterItem, k)
		}
	}

	for k, v := range realEquipMap {
		if v == 0 {
			delete(realEquipMap, k)
		}
	}

	for k, v := range equipItem {
		if v == 0 {
			delete(equipItem, k)
		}
	}

	logger.CtxInfo(ctx, "GetBarrierDeathAward", zap.Any("realItemMap", realItemMap),
		zap.Any("ohterItem", ohterItem),
		zap.Any("realEquipMap", realEquipMap),
		zap.Any("equipItem", equipItem),
		zap.Any("expItem", expItem),
	)

	return realItemMap, ohterItem, realEquipMap, equipItem, expItem.GetCount(), nil
}
