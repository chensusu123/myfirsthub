package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/services/tempbuffservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) CheckBloodAttr(ctx context.Context, userID uint64, barrierID int32) (bloodBottleCount int64, nowAttr map[int32]int64, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "CheckBloodAttr Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)

	defer func() {
		logger.CtxInfo(ctx, "CheckBloodAttr End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()

	data, err := barrieritemsmodel.NewBarrierItems(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr NewBarrierItems Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Error(err),
		)
		return 0, nil, err
	}

	logger.CtxInfo(ctx, "CheckBloodAttr GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
	)

	var bloodlimitAttr int32
	var bloodID int64
	bloodBottleMap := mazeconfigv8config.GetMazeConfig(ctx, constdef.MazeCfgId951)
	bloodBottleCd := mazeconfigv8config.GetMazeValueInt(ctx, constdef.MazeCfgId952)

	for k, v := range bloodBottleMap {
		bloodlimitAttr = k
		bloodID = v
	}

	attrDbs, err := mazecalcattrredis.BatchGetMazeCalcAttr(ctx, userID, []int32{bloodlimitAttr, int32(bloodBottleCd)})
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr BatchGetMazeCalcAttr nil", zap.Uint64("userID", userID))
		return data.Items[bloodID], nil, err
	}

	tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr GetBarrierTempBuff err", zap.Error(err))
		return data.Items[bloodID], nil, err
	}
	for _, buffInfo := range tempBuffInfo.TotalBuff {
		attrDbs[buffInfo.BuffId] += buffInfo.BuffValue
	}

	if _, ok := data.BloodBottleAttr[bloodlimitAttr]; !ok {
		// 第一次初始化
		data.BloodBottleAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
		data.BloodBottleAttr[int32(bloodBottleCd)] = attrDbs[int32(bloodBottleCd)]

		logger.CtxInfo(ctx, "CheckBloodAttr BloodBottleAttr init",
			zap.Any("nowAttr", data.BloodBottleAttr),
		)

		return data.Items[bloodID], nil, nil
	}

	// 初始化过了
	// 判断属性是否有变化
	// 血瓶上限 变化则增加对应血瓶
	if data.BloodBottleAttr[bloodlimitAttr] != attrDbs[bloodlimitAttr] {
		nowAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
		data.Items[bloodID] += (attrDbs[bloodlimitAttr] - data.BloodBottleAttr[bloodlimitAttr])
		data.BloodBottleAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
	}

	if data.BloodBottleAttr[int32(bloodBottleCd)] != attrDbs[int32(bloodBottleCd)] {
		nowAttr[int32(bloodBottleCd)] = attrDbs[int32(bloodBottleCd)]
		data.BloodBottleAttr[int32(bloodBottleCd)] = attrDbs[int32(bloodBottleCd)]
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr Save fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Any("newdata", data),
		)
		return data.Items[bloodID], nil, err
	}

	logger.CtxInfo(ctx, "CheckBloodAttr Save Succesful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("nowAttr", nowAttr),
	)

	return data.Items[bloodID], nowAttr, nil
}
