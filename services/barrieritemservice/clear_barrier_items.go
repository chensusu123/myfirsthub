package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ClearBarrierItems(ctx context.Context, userID uint64, barrierID int32) (int64, int64, error) {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "ClearBarrierItems End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()
	logger.CtxInfo(ctx, "ClearBarrierItems Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "ClearBarrierItems NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return 0, 0, err
	}
	logger.CtxInfo(ctx, "ClearBarrierItems ",
		zap.Any("data", data))
	// fmt.Println("item-------------")
	// for k, v := range data.Items {
	// 	fmt.Printf("k : %d, v : %v\n", k, v)
	// }
	// fmt.Println("Equips-------------")
	// for k, v := range data.Equips {
	// 	fmt.Printf("k : %d, v : %v\n", k, v)
	// }
	// fmt.Println("EquipScore-------------")
	// fmt.Println(data.EquipScore)
	// fmt.Println("ItemsScore-------------")
	// for k, v := range data.ItemsScore {
	// 	fmt.Printf("k : %d, v : %v\n", k, v)
	// }

	data.Items = make(map[int64]*itemservice.ItemInfo)
	data.Equips = make(map[int64]*itemservice.ItemInfo)
	data.EquipScore = 0
	data.ItemsScore = make(map[int32]int32)
	data.SkillsCount = make(map[int32]int32)
	data.SkillDropTime = make(map[int32]int64)
	data.BloodBottleAttr = make(map[int32]int64)

	var bloodlimitAttr int32
	bloodBottleMap := mazeconfigv8config.GetMazeConfig(ctx, constdef.MazeCfgId951)
	bloodBottleCdAttr := mazeconfigv8config.GetMazeValueInt(ctx, constdef.MazeCfgId952)

	for k, v := range bloodBottleMap {
		bloodlimitAttr = k
		_ = v
	}

	// 获取当前血瓶相关属性存储
	attrDbs, err := mazecalcattrredis.BatchGetMazeCalcAttr(ctx, userID, []int32{bloodlimitAttr, int32(bloodBottleCdAttr)})
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr BatchGetMazeCalcAttr nil", zap.Uint64("userID", userID))
		return 0, 0, err
	}

	data.BloodBottleAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
	data.BloodBottleAttr[int32(bloodBottleCdAttr)] = attrDbs[int32(bloodBottleCdAttr)]

	logger.CtxInfo(ctx, "ClearBarrierItems Init Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("data", data),
	)

	return data.BloodBottleAttr[bloodlimitAttr], data.BloodBottleAttr[int32(bloodBottleCdAttr)], data.Save(ctx, userID, barrierID)
}

func (s *service) DelInAdditionToEquips(ctx context.Context, userID uint64, barrierID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	defer func() {
		logger.CtxInfo(ctx, "DelInAdditionToEquips End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()
	logger.CtxInfo(ctx, "DelInAdditionToEquips Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "DelInAdditionToEquips NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return err
	}

	data.Items = make(map[int64]*itemservice.ItemInfo)
	data.ItemsScore = make(map[int32]int32)
	data.SkillDropTime = make(map[int32]int64)
	data.SkillsCount = make(map[int32]int32)

	logger.CtxInfo(ctx, "DelInAdditionToEquips del successful",
		zap.Any("data", data))

	return data.Save(ctx, userID, barrierID)
}
