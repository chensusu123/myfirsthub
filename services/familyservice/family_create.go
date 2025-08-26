package familyservice

import (
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

// CreateFamily 创建家族
func (r *service) CreateFamily(logger fklog.FKLogI, userID uint64, allianceID int32, familyName string,
	familySetting int32, userInfo familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error) {
	// 创建新家族
	familyInfoModel, err := familymodel.NewFamilyInfoModel(logger, familyName, familySetting)
	if err != nil {
		logger.ErrorWF("CreateFamily NewFamilyInfoModel failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting))
		return nil, err
	}
	// 加入成员
	familyInfoModel.AddMember(logger, userInfo)
	// 保存家族信息
	err = familyInfoModel.Save(logger, familyInfoModel.FamilyID)
	if err != nil {
		logger.ErrorWF("CreateFamily model.Save failed",
			zap.Any("familyName", familyName), zap.Any("familySetting", familySetting),
			zap.Error(err))
		return nil, err
	}

	// 把家族id加入家族列表
	familyListModel, err := familymodel.LoadFamilyListModel(logger)
	if err != nil {
		logger.ErrorWF("CreateFamily LoadFamilyListModel err",
			zap.Error(err))
		return nil, err
	}
	familyListModel.AddFamily(logger, familyInfoModel.FamilyID)
	err = familyListModel.Save(logger)
	if err != nil {
		logger.ErrorWF("CreateFamily familyListModel.Save err",
			zap.Error(err))
		return nil, err
	}

	// 设置家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(logger, familyInfoModel.FamilyID)
	err = familyToAllianceModel.SetUserAlliance(logger, allianceID)
	if err != nil {
		logger.ErrorWF("CreateFamily familyToAllianceModel.SetUserAlliance err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	// 联盟中增加该家族id
	allianceModel, err := alliancemodel.LoadAllianceInfoModel(logger, allianceID)
	if err != nil {
		logger.ErrorWF("CreateFamily LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID),
			zap.Error(err))
		return nil, err
	}

	err = allianceModel.AddFamilyID(logger, familyInfoModel.FamilyID)
	if err != nil {
		logger.ErrorWF("CreateFamily allianceModel.AddFamilyID err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", familyInfoModel.FamilyID),
			zap.Error(err))
		return nil, err
	}

	return familyInfoModel, nil
}

// 设置家族群组id
func (r *service) SetFamilyGroupID(logger fklog.FKLogI, familyID int32, groupID int32) error {
	familymodel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("UpgradeFamily familymodel.LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	familymodel.SetFamilyGroupID(logger, groupID)
	familymodel.Save(logger, familyID)
	return nil
}

// DeductCreateFamilyCost 创建家族扣物品
func (r *service) DeductCreateFamilyCost(logger fklog.FKLogI, cost map[int32]int64) error {
	// todo 待补充，读表

	return nil
}

// 查询创建家族消耗
func (r *service) QueryCreateFamilyCost(logger fklog.FKLogI) (map[int32]int64, error) {
	// todo 待补充，读表

	return nil, nil
}

// CheckCreateFamilyCost 检查创建家族消耗
func (r *service) CheckCreateFamilyCost(logger fklog.FKLogI) (bool, error) {
	// todo 待补充，读表

	return true, nil
}
