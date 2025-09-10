package barrieritemservice

import (
	"context"
	"fmt"
	"math"
	"maze_game_server/common/constdef"
	"maze_game_server/config/GMazeBariresDropConditionV8Cfg"
	"maze_game_server/config/GMazeBariresDropV8Cfg"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/io/redis/barrierguiditemredis"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"
	"maze_game_server/services/tempbuffservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) FallOffSkillItems(ctx context.Context, userID uint64, barrierID int32, killMonsterNum int32, nowBloodVolume int64,
	allBloodVolume int64, guid int64, pos string) (dropItems []*itemservice.ItemInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "FallOffSkillItems End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),

			zap.String("pos", pos),
		)
	}()
	logger.CtxInfo(ctx, "FallOffSkillItems Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("killMonsterNum", killMonsterNum),
		zap.Int64("nowBloodVolume", nowBloodVolume),
		zap.Int64("allBloodVolume", allBloodVolume),
		zap.String("pos", pos),
	)

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "FallOffSkillItems NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),
			zap.String("pos", pos),
			zap.Error(err),
		)
		return
	}

	logger.CtxInfo(ctx, "FallOffSkillItems GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
	)

	// var reason string
	attrDbs, err := mazecalcattrredis.GetAllMazeCalcAttr(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "FallOffSkillItems BatchGetMazeCalcAttr nil", zap.Uint64("userID", userID))
		return
	}

	tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "FallOffSkillItems GetBarrierTempBuff err", zap.Error(err))
		return nil, err
	}
	for _, buffInfo := range tempBuffInfo.TotalBuff {
		attrDbs[buffInfo.BuffId] += buffInfo.BuffValue
	}

	bloodMap := mazeconfigv8config.GetMazeConfig(ctx, constdef.MazeCfgId951)
	var bloodLimitAttr int32
	var bloodItemID int64
	for k, v := range bloodMap {
		bloodLimitAttr = k
		bloodItemID = v
	}

	bloodCount := 0
	for _, item := range data.Items {
		if item.ItemId == int32(bloodItemID) {
			bloodCount++
		}
	}

	// 先获取当前关卡掉落物品
	barrierCfg := mazebarriesv8config.GetStageConfig(ctx, barrierID)
	for _, itemID := range barrierCfg.Drop_id {
		barrierDrops := GMazeBariresDropV8Cfg.GetAll()
		for _, barrierDrop := range barrierDrops {
			// 指定掉落id
			if barrierDrop.Drop_id != itemID {
				continue
			}

			for dropItemID, dropCount := range barrierDrop.Drop_items {
				// 当前拥有上限
				if dropItemID == int32(bloodItemID) && int64(bloodCount) >= attrDbs[bloodLimitAttr] {
					// reason += fmt.Sprintf("血瓶到上限了 当前:%d 上限:%d\n", data.Items[int64(dropItemID)], row.Value_int)
					logger.CtxInfo(ctx, "FallOffSkillItems BloodBottle Full",
						zap.Uint64("userID", userID),
						zap.Int("bloodCount", bloodCount),
					)
					break
				}

				// 判断杀怪
				if killMonsterNum > barrierDrop.In_barries_kill_max || killMonsterNum < barrierDrop.In_barries_kill_min {
					// reason += fmt.Sprintf("杀怪不满足条件 当前杀怪数:%d 杀怪范围[%d~%d]\n", killMonsterNum, barrierDrop.In_barries_kill_min, barrierDrop.In_barries_kill_max)
					continue
				}

				// 判断掉落cd
				nowTime := time.Now().UnixMilli()
				if _, ok := data.SkillDropTime[itemID]; ok {
					if nowTime-data.SkillDropTime[itemID] < int64(barrierDrop.Drop_cd) {
						// reason += fmt.Sprintf("CD不满足条件 上次掉落时间:%d 当前时间:%d 时间间隔:%d 当前间隔:%d\n", data.SkillDropTime[itemID], nowTime, barrierDrop.Drop_cd, nowTime-data.SkillDropTime[itemID])
						continue
					}
				}

				isCondition := true
				// 判断掉落条件id
				for _, condionID := range barrierDrop.Drop_condition {
					ok, err, _ := checkCondition(ctx, attrDbs, condionID, nowBloodVolume, allBloodVolume, data.SkillsCount[dropItemID])
					if err != nil {
						return nil, err
					}

					if !ok {
						isCondition = false
						// reason += rs
					}
				}

				if !isCondition {
					break
				}

				// 判断百分比 10000
				nowRandNum := fkutil.RandInt(1, 10000)
				if nowRandNum > int(barrierDrop.Drop_ratio_max) || nowRandNum < int(barrierDrop.Drop_ratio_min) {
					// reason += fmt.Sprintf("随机概率不满足条件 当前概率:%d 允许概率[%d~%d]\n", nowRandNum, barrierDrop.Drop_ratio_min, barrierDrop.Drop_ratio_max)
					continue
				}

				// 实际添加物品 并设置此类物品掉落cd
				for i := 1; i <= int(dropCount); i++ {
					itemGuid, err := barrierguiditemredis.IncrNowGuid(ctx, userID, barrierID)
					if err != nil {
						logger.CtxError(ctx, "FallOffSkillItems IncrNowGuid Fail",
							zap.Uint64("userID", userID),
							zap.Int32("barrierID", barrierID),
							zap.Any("data", data),
							zap.Error(err),
						)
					}

					data.Items[itemGuid] = &itemservice.ItemInfo{
						ItemId: int32(dropItemID),
						Count:  1,
						Guid:   itemGuid,
					}

					// 统计战区掉落次数 掉落cd
					data.SkillsCount[dropItemID] += 1
					data.SkillDropTime[itemID] = nowTime

					dropItems = append(dropItems, &itemservice.ItemInfo{
						ItemId: dropItemID,
						Count:  1,
						Guid:   itemGuid,
					})
				}
			}

		}

	}
	// fmt.Println(reason)
	// fmt.Println("------------------------------------------------")
	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "FallOffSkillItems data Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),
			zap.String("pos", pos),
			zap.Any("data", data),
		)
		return
	}

	logger.CtxInfo(ctx, "FallOffSkillItems Add Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("dropItems", dropItems),
	)

	// 推送物品
	if len(dropItems) > 0 {
		err = item.OnSendItemsPack(ctx, userID, dropItems, nil, guid, pos, 1)
		if err != nil {
			logger.CtxWarn(ctx, "FallOffSkillItems OnSendItemsPack Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int64("guid", guid),
				zap.String("pos", pos),
				zap.Any("dropItems", dropItems),
				zap.Error(err),
			)
			return
		}
	}

	return
}

func checkCondition(ctx context.Context, attrDbs map[int32]int64, conditionID int32, nowBloodVolume int64, allBloodVolume int64, nowSkillCount int32) (bool, error, string) {
	barrierDropCondition := GMazeBariresDropConditionV8Cfg.GetWithCtx(ctx, conditionID)
	switch barrierDropCondition.Condition_type {
	case 1:
		return nowBloodVolume*int64(1000000) <= allBloodVolume*(int64(barrierDropCondition.Value)+attrDbs[constdef.BloodBottleProbability]), nil,
			fmt.Sprintf("当前血量不满足条件 当前血量:%d 总血量:%d 低于百分比:%d\n", nowBloodVolume, allBloodVolume, barrierDropCondition.Value)
	case 2:
		// 当前刷怪区掉落上限
		numerator := float64(10000+attrDbs[constdef.BloodBottlesNumber]) / 10000.0
		// fmt.Printf("numerator:%v\n", numerator)
		count := math.Ceil(float64(barrierDropCondition.Value) * numerator)
		// fmt.Printf("count:%v\n", count)
		intCount := int32(count)
		// if intCount == 0 {
		// 	fmt.Printf("conditionID:%d", conditionID)
		// }
		return nowSkillCount <= intCount, nil, fmt.Sprintf("当前使用血瓶上限不满足条件 当前掉落血瓶:%d 允许掉落上限:%d\n", nowSkillCount, intCount)
	case 3:
		return nowSkillCount <= barrierDropCondition.Value, nil, fmt.Sprintf("当前场上同时存在血瓶不满足条件 当前血瓶数:%d 允许存在血瓶数:%d", nowSkillCount, barrierDropCondition.Value)
	case 4:
		return nowBloodVolume*int64(1000000) >= allBloodVolume*(int64(barrierDropCondition.Value)+attrDbs[constdef.BloodBottleProbability]), nil,
			fmt.Sprintf("当前血量不满足条件 当前血量:%d 总血量:%d 高于百分比:%d\n", nowBloodVolume, allBloodVolume, barrierDropCondition.Value)
	}
	return false, nil, ""
}
