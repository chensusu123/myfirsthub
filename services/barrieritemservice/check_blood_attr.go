package barrieritemservice

import (
	"context"
	"maze_game_server/common/constdef"
	"maze_game_server/excel/mazeconfigv8config"
	"maze_game_server/io/redis/mazecalcattrredis"
	"maze_game_server/model/barrieritemsmodel"
	"maze_game_server/servers/maze_main_server/process/item"
	"maze_game_server/services/itemservice"
	"maze_game_server/services/tempbuffservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) CheckBloodAttr(ctx context.Context, userID uint64, barrierID int32) (bloodBottleLimit, bloodBottleCd int64, err error) {
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
		return 0, 0, err
	}

	var bloodlimitAttr int32
	var bloodID int64
	bloodBottleMap := mazeconfigv8config.GetMazeConfig(ctx, constdef.MazeCfgId951)
	bloodBottleCdAttr := mazeconfigv8config.GetMazeValueInt(ctx, constdef.MazeCfgId952)

	for k, v := range bloodBottleMap {
		bloodlimitAttr = k
		bloodID = v
	}

	// 获取当前血瓶相关属性存储
	attrDbs, err := mazecalcattrredis.BatchGetMazeCalcAttr(ctx, userID, []int32{bloodlimitAttr, int32(bloodBottleCdAttr)})
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr BatchGetMazeCalcAttr nil", zap.Uint64("userID", userID))
		return 0, 0, err
	}

	tempBuffInfo, err := tempbuffservice.GlobalTempBuffService.GetTempBuffInfo(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr GetBarrierTempBuff err", zap.Error(err))
		return 0, 0, err
	}

	for _, buffInfo := range tempBuffInfo.TotalBuff {
		if buffInfo.BuffId == bloodlimitAttr || buffInfo.BuffId == int32(bloodBottleCdAttr) {
			attrDbs[buffInfo.BuffId] += buffInfo.BuffValue
		}
	}

	logger.CtxInfo(ctx, "CheckBloodAttr GetData Successful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("nowBloodBottleCount", data.Items[bloodID]),
		zap.Any("lastBloodBottleLimit", data.BloodBottleAttr[bloodlimitAttr]),
		zap.Any("lastBloodBottleCd", data.BloodBottleAttr[int32(bloodBottleCdAttr)]),
		zap.Any("nowBloodBottleLimit", attrDbs[bloodlimitAttr]),
		zap.Any("nowBloodBottleCd", attrDbs[int32(bloodBottleCdAttr)]),
	)

	if _, ok := data.BloodBottleAttr[bloodlimitAttr]; !ok {
		// 第一次初始化
		data.BloodBottleAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
		data.BloodBottleAttr[int32(bloodBottleCdAttr)] = attrDbs[int32(bloodBottleCdAttr)]

		logger.CtxInfo(ctx, "CheckBloodAttr BloodBottleAttr init",
			zap.Any("nowAttr", data.BloodBottleAttr),
		)

		err = data.Save(ctx, userID, barrierID)
		if err != nil {
			logger.CtxError(ctx, "CheckBloodAttr data Init Fail",
				zap.Uint64("userID", userID),
				zap.Int32("barrierID", barrierID),
				zap.Any("nowdata", data),
				zap.Any("nowBloodBottleCount", data.Items[bloodID]),
				zap.Any("lastBloodBottleLimit", data.BloodBottleAttr[bloodlimitAttr]),
				zap.Any("lastBloodBottleCd", data.BloodBottleAttr[int32(bloodBottleCdAttr)]),
				zap.Any("nowBloodBottleLimit", attrDbs[bloodlimitAttr]),
				zap.Any("nowBloodBottleCd", attrDbs[int32(bloodBottleCdAttr)]),
			)
			return data.BloodBottleAttr[bloodlimitAttr], data.BloodBottleAttr[int32(bloodBottleCdAttr)], err
		}

		return data.BloodBottleAttr[bloodlimitAttr], data.BloodBottleAttr[int32(bloodBottleCdAttr)], nil
	}

	dropItems := make([]*itemservice.ItemInfo, 0)
	// 初始化过了
	// 判断属性是否有变化
	// 血瓶上限 变化则增加对应血瓶
	if data.BloodBottleAttr[bloodlimitAttr] != attrDbs[bloodlimitAttr] {
		bloodBottleLimit = attrDbs[bloodlimitAttr]
		data.Items[bloodID] += (attrDbs[bloodlimitAttr] - data.BloodBottleAttr[bloodlimitAttr])
		// 推包
		dropItems = append(dropItems, &itemservice.ItemInfo{
			ItemId: int32(bloodID),
			Count:  1,
		})

		logger.CtxInfo(ctx, "CheckBloodAttr Drop BloodBottle",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)

		data.BloodBottleAttr[bloodlimitAttr] = attrDbs[bloodlimitAttr]
	}

	if data.BloodBottleAttr[int32(bloodBottleCdAttr)] != attrDbs[int32(bloodBottleCdAttr)] {
		bloodBottleCd = attrDbs[int32(bloodBottleCdAttr)]
		data.BloodBottleAttr[int32(bloodBottleCdAttr)] = attrDbs[int32(bloodBottleCdAttr)]
	}

	err = data.Save(ctx, userID, barrierID)
	if err != nil {
		logger.CtxError(ctx, "CheckBloodAttr Save fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Any("newdata", data),
		)
		return bloodBottleLimit, bloodBottleCd, err
	}

	err = item.OnSendItemsPack(ctx, userID, dropItems, nil, 0, "", 2)
	if err != nil {
		logger.CtxWarn(ctx, "AddScoreItem OnSendItemsPack Fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
			zap.Any("dropItems", dropItems),
			zap.Error(err),
		)
		return
	}

	logger.CtxInfo(ctx, "CheckBloodAttr Save Succesful",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
		zap.Any("nowdata", data),
		zap.Any("nowBloodBottleCount", data.Items[bloodID]),
		zap.Any("nowBloodBottleLimit", data.BloodBottleAttr[bloodlimitAttr]),
		zap.Any("nowBloodBottleCd", data.BloodBottleAttr[int32(bloodBottleCdAttr)]),
	)

	return data.BloodBottleAttr[bloodlimitAttr], data.BloodBottleAttr[int32(bloodBottleCdAttr)], nil
}
