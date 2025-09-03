package barrieritemservice

import (
	"context"
	"fmt"
	"maze_game_server/config/GMazeBariresDropConditionV8Cfg"
	"maze_game_server/config/GMazeBariresDropV8Cfg"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil"
	"go.uber.org/zap"
)

func (s *service) FallOffSkillItems(ctx context.Context, userID uint64, barrierID int32, killMonsterNum int32, nowBloodVolume int64, allBloodVolume int64, monsterGuid int64, monsterPos string) error {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "FallOffSkillItems End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),
			zap.String("monsterPos", monsterPos),
		)
	}()
	logger.CtxInfo(ctx, "FallOffSkillItems Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Int32("killMonsterNum", killMonsterNum),
		zap.Int64("nowBloodVolume", nowBloodVolume),
		zap.Int64("allBloodVolume", allBloodVolume),
		zap.String("monsterPos", monsterPos),
	)

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "ClearBarrierItems NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),
			zap.String("monsterPos", monsterPos),
			zap.Error(err),
		)
		return err
	}

	dropItems := make([]*itemservice.ItemInfo, 0)

	// 先获取当前关卡掉落物品
	barrierCfg := mazebarriesv8config.GetStageConfig(ctx, barrierID)
	for _, itemID := range barrierCfg.Drop_id {
		barrierDrops := GMazeBariresDropV8Cfg.GetAll()
		for _, barrierDrop := range barrierDrops {
			// 指定掉落id
			if barrierDrop.Drop_id != itemID {
				fmt.Println("指定掉落id过滤", itemID, barrierDrop.Drop_id)
				continue
			}

			for dropItemID, dropCount := range barrierDrop.Drop_items {
				// 判断掉落条件id
				for _, condionID := range barrierDrop.Drop_condition {
					if !checkCondition(ctx, condionID, nowBloodVolume, allBloodVolume, data.SkillsCount[dropItemID], data.Items[int64(dropItemID)]) {
						fmt.Println("条件过滤", itemID, dropItemID)
						continue
					}

					// 判断杀怪
					if killMonsterNum > barrierDrop.In_barries_kill_max || killMonsterNum < barrierDrop.In_barries_kill_min {
						fmt.Println("判断杀怪过滤", itemID, dropItemID)
						continue
					}

					// 判断百分比
					nowRandNum := fkutil.RandInt(1, 10000)
					if nowRandNum > int(barrierDrop.Drop_ratio_max) || nowRandNum < int(barrierDrop.Drop_ratio_min) {
						fmt.Println("判断百分比过滤", itemID, dropItemID)
						continue
					}

					// 判断掉落cd
					nowTime := time.Now().UnixMilli()
					if _, ok := data.SkillDropTime[dropItemID]; ok {
						if data.SkillDropTime[dropItemID]-nowTime < int64(barrierDrop.Drop_cd) {
							fmt.Println("判断掉落cd过滤", data.SkillDropTime[dropItemID], nowTime)
							continue
						}
					}

					// 实际添加物品 并设置此类物品掉落cd
					data.Items[int64(dropItemID)] += dropCount
					data.SkillDropTime[dropItemID] = nowTime
					dropItems = append(dropItems, &itemservice.ItemInfo{
						ItemId: dropItemID,
						Count:  dropCount,
					})
				}

			}

		}
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "FallOffSkillItems data Save Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Int32("killMonsterNum", killMonsterNum),
			zap.Int64("nowBloodVolume", nowBloodVolume),
			zap.Int64("allBloodVolume", allBloodVolume),
			zap.String("monsterPos", monsterPos),
			zap.Any("data", data),
		)
		return err
	}

	// 推送物品
	if len(dropItems) > 0 {
		err = item.OnSendItemsPack(ctx, userID, dropItems, nil, monsterGuid, monsterPos)
		if err != nil {
			logger.CtxWarn(ctx, "AddEquipScore OnSendItemsPack Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Int64("monsterGuid", monsterGuid),
				zap.String("monsterPos", monsterPos),
				zap.Any("dropItems", dropItems),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}

func checkCondition(ctx context.Context, conditionID int32, nowBloodVolume int64, allBloodVolume int64, skillCount int32, nowSkillCount int64) bool {
	barrierDropCondition := GMazeBariresDropConditionV8Cfg.GetWithCtx(ctx, conditionID)
	switch barrierDropCondition.Condition_type {
	case 1:
		return nowBloodVolume*int64(1000000) <= allBloodVolume*int64(barrierDropCondition.Value)
	case 2:
		return skillCount <= barrierDropCondition.Value
	case 3:
		return nowSkillCount <= int64(barrierDropCondition.Value)
	case 4:
		return nowBloodVolume*int64(1000000) >= allBloodVolume*int64(barrierDropCondition.Value)
	}
	return false
}
