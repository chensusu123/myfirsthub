package familyservice

import (
	"context"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// CreateFamily 创建家族
func (r *service) CreateFamily(ctx context.Context, userID uint64, allianceID int32, familyName string,
	familySetting int32, userInfo familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 创建新家族
	familyInfoModel, err := familymodel.NewFamilyInfoModel(ctx, familyName, familySetting)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily NewFamilyInfoModel failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting))
		return nil, err
	}
	// 加入成员
	familyInfoModel.AddMember(ctx, userInfo)
	// 保存家族信息
	err = familyInfoModel.Save(ctx, familyInfoModel.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily model.Save failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting),
			zap.Error(err))
		return nil, err
	}

	// 把家族id加入家族列表
	familyListModel, err := familymodel.LoadFamilyListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily LoadFamilyListModel err",
			zap.Error(err))
		return nil, err
	}
	familyListModel.AddFamily(ctx, familyInfoModel.FamilyID)
	err = familyListModel.Save(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily familyListModel.Save err",
			zap.Error(err))
		return nil, err
	}

	// 设置家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, familyInfoModel.FamilyID)
	err = familyToAllianceModel.SetUserAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily familyToAllianceModel.SetUserAlliance err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	// 联盟中增加该家族id
	allianceModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID),
			zap.Error(err))
		return nil, err
	}

	err = allianceModel.AddFamilyID(ctx, familyInfoModel.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "CreateFamily allianceModel.AddFamilyID err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	return familyInfoModel, nil
}

// 设置家族群组id
func (r *service) SetFamilyGroupID(ctx context.Context, familyID int32, groupID int32) error {
	familymodel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	logger := fklog.ContextAppLogger(ctx)
	if err != nil {
		logger.ErrorWF("UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	familymodel.SetFamilyGroupID(ctx, groupID)
	familymodel.Save(ctx, familyID)
	return nil
}

// DeductCreateFamilyCost 创建家族扣物品
func (r *service) DeductCreateFamilyCost(ctx context.Context, cost map[int32]int64) error {
	// todo 待补充，读表

	return nil
}

// 查询创建家族消耗
func (r *service) QueryCreateFamilyCost(ctx context.Context) (map[int32]int64, error) {
	// todo 待补充，读表

	return nil, nil
}

// CheckCreateFamilyCost 检查创建家族消耗
func (r *service) CheckCreateFamilyCost(ctx context.Context) (bool, error) {
	// todo 待补充，读表

	return true, nil
}
