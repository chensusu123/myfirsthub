package barrieritemservice

import (
	"context"
	"maze_game_server/common/errors"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/services/equipdropservice"
	"maze_game_server/services/itemservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) FallOffEquip(ctx context.Context, userID uint64, barrierID int32, equipNum int32) (res []*itemservice.ItemInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "DropEquip Start",
		zap.Uint64("userID", userID),
		zap.Int32("barrierID", barrierID),
	)
	defer func() {
		logger.CtxInfo(ctx, "DropEquip End",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
	}()

	userInfo, err := mazeuserinfo.GetUserInfoV2(ctx, userID)
	if err != nil {
		logger.CtxError(ctx, "DropEquip GetUserInfoV2 fail", zap.Error(err))
		return nil, err
	}
	level := userInfo.Level

	if barrierID != userInfo.Barrier {
		logger.CtxError(ctx, "DropEquip check barrier fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
		return nil, errors.New("请求的关卡id和上报存储的关卡不一致")
	}

	addEquipMap, err := equipdropservice.GlobalEquipDropService.GetNewEquip(ctx, userID, int32(level), barrierID, equipNum)
	if err != nil {
		logger.CtxError(ctx, "DropEquip GetNewEquip fail",
			zap.Uint64("userID", userID),
			zap.Int32("barrierID", barrierID),
		)
		return nil, errors.New("获取数据失败")
	}

	for equipID, equipCount := range addEquipMap {
		res = append(res, &itemservice.ItemInfo{
			ItemId: equipID,
			Count:  int64(equipCount),
		})
	}

	return
}
