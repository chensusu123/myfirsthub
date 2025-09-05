package barrieritemservice

import (
	"context"
	"math"
	"maze_game_server/common/constdef"
	"maze_game_server/config/GMazeBariresDropConditionV8Cfg"
	"maze_game_server/config/GMazeBariresDropV8Cfg"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/excel/mazebarriesv8config"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"
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
		logger.CtxError(ctx, "ClearBarrierItems NewBarrierItems Fail",
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

	attrDbs, err := mazecalcattrredis.BatchGetMazeCalcAttr(ctx, userID, []int32{constdef.BloodBottleProbability, constdef.BloodBottlesNumber})
	if err != nil {
		logger.CtxError(ctx, "checkCondition GetAllMazeCalcAttr nil", zap.Uint64("userID", userID))
		return
	}

	row := GMazeConfigV8Cfg.GetWithCtx(ctx, constdef.MazeCfgId951)

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
				if dropItemID == constdef.BloodBottleID && data.Items[int64(dropItemID)] == row.Value_int {
					continue
				}

				isCondition := true
				// 判断掉落条件id
				for _, condionID := range barrierDrop.Drop_condition {
					ok, err := checkCondition(ctx, attrDbs, condionID, nowBloodVolume, allBloodVolume, data.SkillsCount[dropItemID], data.Items[int64(dropItemID)])
					if err != nil {
						return nil, err
					}

					if !ok {
						isCondition = false
					}
				}

				if !isCondition {
					break
				}

				// 判断杀怪
				if killMonsterNum > barrierDrop.In_barries_kill_max || killMonsterNum < barrierDrop.In_barries_kill_min {
					continue
				}

				// 判断百分比
				nowRandNum := fkutil.RandInt(1, 10000)

				if nowRandNum > int(barrierDrop.Drop_ratio_max) || nowRandNum < int(barrierDrop.Drop_ratio_min) {
					continue
				}

				// 判断掉落cd
				nowTime := time.Now().UnixMilli()
				if _, ok := data.SkillDropTime[itemID]; ok {
					if data.SkillDropTime[itemID]-nowTime < int64(barrierDrop.Drop_cd) {
						continue
					}
				}

				// 实际添加物品 并设置此类物品掉落cd
				data.Items[int64(dropItemID)] += dropCount
				data.SkillDropTime[itemID] = nowTime
				dropItems = append(dropItems, &itemservice.ItemInfo{
					ItemId: dropItemID,
					Count:  dropCount,
				})
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
			zap.String("pos", pos),
			zap.Any("data", data),
		)
		return
	}

	// 推送物品
	if len(dropItems) > 0 {
		err = item.OnSendItemsPack(ctx, userID, dropItems, nil, guid, pos)
		if err != nil {
			logger.CtxWarn(ctx, "AddEquipScore OnSendItemsPack Fail",
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

func checkCondition(ctx context.Context, attrDbs map[int32]int64, conditionID int32, nowBloodVolume int64, allBloodVolume int64, skillCount int32, nowSkillCount int64) (bool, error) {
	barrierDropCondition := GMazeBariresDropConditionV8Cfg.GetWithCtx(ctx, conditionID)
	switch barrierDropCondition.Condition_type {
	case 1:
		return nowBloodVolume*int64(1000000) <= allBloodVolume*(int64(barrierDropCondition.Value)+attrDbs[constdef.BloodBottleProbability]), nil
	case 2:
		numerator := float64(10000+attrDbs[constdef.BloodBottlesNumber]) / 10000.0
		count := math.Ceil(float64(barrierDropCondition.Value) * numerator)
		intCount := int32(count)
		return skillCount <= intCount, nil
	case 3:
		return nowSkillCount <= int64(barrierDropCondition.Value), nil
	case 4:
		return nowBloodVolume*int64(1000000) >= allBloodVolume*(int64(barrierDropCondition.Value)+attrDbs[constdef.BloodBottleProbability]), nil
	}
	return false, nil
}
