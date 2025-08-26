package familyservice

import (
	"context"

	"maze_game_server/model/familymodel"
	"maze_game_server/pb/common/MazeFamily"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// UpgradeFamily 升级家族
func (r *service) UpgradeFamily(logger fklog.FKLogI, familyID int32, targetLevel int32) (*familymodel.FamilyInfoModel, error) {
	familymodel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	familymodel.Upgrade(logger, familyID, targetLevel)
	// 保存
	err = familymodel.Save(logger, familyID)
	if err != nil {
		logger.ErrorWF("UpgradeFamily familymodel.Save err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return familymodel, nil
}

// SendUpgradeFamilyIDPack 发送升级家族id包
func (r *service) SendUpgradeFamilyIDPack(logger fklog.FKLogI, familyID int32) error {
	familymodel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	pack := &MazeFamily.UpgradeFamilyID{
		FamilyId:   proto.Int32(familyID),
		FamilyInfo: familymodel.DataToFamilyInfoPb(logger),
	}

	familyMembers := familymodel.DataToFamilyMembersPb(logger)
	for _, v := range familyMembers {
		online.ClusterPush(context.TODO(), v.GetUserId(), 0, pack)
	}
	return nil
}

// DeductUpgradeFamilyCost 升级家族扣物品
func (r *service) DeductUpgradeFamilyCost(logger fklog.FKLogI, cost map[int32]int64) error {
	// todo 待补充，读表

	return nil
}

// 查询升级家族消耗
func (r *service) QueryUpgradeFamilyCost(logger fklog.FKLogI, familyID, targetLevel int32) (map[int32]int64, error) {
	// todo 待补充，读表

	return nil, nil
}

// CheckUpgradeFamilyCost 检查家族升级消耗
func (r *service) CheckUpgradeFamilyCost(logger fklog.FKLogI, targetLevel int32) (bool, error) {
	// todo 待补充，读表

	return true, nil
}
