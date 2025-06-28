package familyservice

import (
	"errors"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func (s *service) DissolutionFamily(logger fklog.FKLogI, uid uint64, familyID int32) error {
	familyModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("LoadFamilyModel err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	// 检查用户权限
	if !familyModel.CheckHaveLeader(logger, []uint64{uid}) {
		logger.ErrorWF("DissolutionFamily CheckHaveLeader err", zap.Int32("familyID", familyID), zap.Uint64("uid", uid))
		return errors.New("用户没有权限解散家族")
	}

	// 删除家族
	familyModel.Delete(logger, familyID)

	// 删除家族列表
	familyListModel, err := familymodel.LoadFamilyListModel(logger)
	if err != nil {
		logger.ErrorWF("LoadFamilyListModel err", zap.Error(err))
		return err
	}
	familyListModel.RemoveFamily(logger, familyID)

	// 删除联盟内该家族ID
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(logger, familyID)
	allianceID, err := familyToAllianceModel.GetUserAlliance(logger)
	if err != nil {
		logger.ErrorWF("DissolutionFamily GetUserAlliance err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(logger, allianceID)
	if err != nil {
		logger.ErrorWF("DissolutionFamily LoadAllianceInfoModel err", zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	allianceInfoModel.RemoveFamilyID(logger, familyID)

	// 删除家族对应联盟
	err = familyToAllianceModel.DelUserAlliance(logger)
	if err != nil {
		logger.ErrorWF("DissolutionFamily DelUserAlliance err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	return nil
}
