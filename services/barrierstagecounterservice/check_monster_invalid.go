package barrierstagecounterservice

import (
	"context"
	"maze_game_server/model/barrierstagecountermodel"
	"maze_game_server/model/passareamodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) CheckMonsterInvalid(ctx context.Context, userID uint64, barrierID int32, stageID int32, areaID int32, areaIndex int32, monsterGuid int64) (res bool, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "CheckMonsterInvalid Start",
		zap.Any("userID", userID),
		zap.Any("barrierID", barrierID),
	)
	defer func() {
		logger.CtxInfo(ctx, "CheckMonsterInvalid End",
			zap.Any("userID", userID),
			zap.Any("barrierID", barrierID),
			zap.Any("res", res),
		)
	}()

	// 检查通过区域
	passInfo, err := passareamodel.NewPassAreaModel(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckMonsterInvalid NewPassAreaModel Fail",
			zap.Any("userID", userID),
			zap.Any("barrierID", barrierID),
			zap.Any("monsterGuid", monsterGuid),
		)
		return
	}

	isPassArea := false
	for _, passArea := range passInfo.PassAreaList {
		if passArea.AreaId == areaID && passArea.AreaIndex == areaIndex {
			isPassArea = true
		}
	}

	// 判断是否处理过
	recordModel, err := barrierstagecountermodel.NewBarrierStageCounterModel(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "AddKillMonsterNum NewBarrierStageCounterModel fail", zap.Error(err))
		return
	}

	isCheck := false
	if _, ok := recordModel.KillMonsterGuidMap[stageID][monsterGuid]; ok {
		isCheck = true
	}

	res = !isPassArea && !isCheck
	return
}
