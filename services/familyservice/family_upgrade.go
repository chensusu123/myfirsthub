package familyservice

import (
	"context"

	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// UpgradeFamily 升级家族
func (r *service) UpgradeFamily(ctx context.Context, familyID int32, targetLevel int32) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	familymodel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	familymodel.Upgrade(ctx, familyID, targetLevel)
	// 保存
	err = familymodel.Save(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "UpgradeFamily familymodel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familymodel, nil
}

// SendUpgradeFamilyIDPack 发送升级家族id包
func (r *service) SendUpgradeFamilyIDPack(ctx context.Context, familyID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	familymodel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.UpgradeFamilyID{
		FamilyId:   proto.Int32(familyID),
		FamilyInfo: familymodel.DataToFamilyInfoPb(ctx),
	}

	familyMembers := familymodel.DataToFamilyMembersPb(ctx)
	for _, v := range familyMembers {
		online.ClusterPush(ctx, v.GetUserId(), 0, pack)
	}
	return nil
}

// DeductUpgradeFamilyCost 升级家族扣物品
func (r *service) DeductUpgradeFamilyCost(ctx context.Context, cost map[int32]int64) error {
	// todo 待补充，读表

	return nil
}

// 查询升级家族消耗
func (r *service) QueryUpgradeFamilyCost(ctx context.Context, familyID, targetLevel int32) (map[int32]int64, error) {
	// todo 待补充，读表

	return nil, nil
}

// CheckUpgradeFamilyCost 检查家族升级消耗
func (r *service) CheckUpgradeFamilyCost(ctx context.Context, targetLevel int32) (bool, error) {
	// todo 待补充，读表

	return true, nil
}
