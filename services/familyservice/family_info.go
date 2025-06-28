package familyservice

import (
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

// GetFamilyInfo 查询家族详情
func (r *service) GetFamilyInfo(logger fklog.FKLogI, familyID int32) (*familymodel.FamilyInfoModel, error) {
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("GetFamilyInfo LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}

	return familyInfoModel, nil
}

// GetFamilyList 获取家族列表
func (r *service) GetFamilyList(logger fklog.FKLogI) (familymodel.FamilysInfoModel, error) {
	// 查询当前家族id列表
	familyListModel, err := familymodel.LoadFamilyListModel(logger)
	if err != nil {
		logger.ErrorWF("GetFamilyList LoadFamilyListModel err",
			zap.Error(err))
		return nil, err
	}
	// 查询家族列表
	familysInfoModel, err := familymodel.LoadFamilyListInfoModel(logger, familyListModel.Familys)
	if err != nil {
		logger.ErrorWF("GetFamilyInfo LoadFamilyInfoModel err",
			zap.Error(err))
		return nil, err
	}
	return familysInfoModel, nil
}
