package allianceservice

import (
	"errors"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

func (s *service) ApplyChangeAlliance(logger fklog.FKLogI, userID uint64, allianceID int32) error {
	userFamilyModel := familymodel.NewUserFamilyModel(logger, userID)
	familyID, err := userFamilyModel.GetUserFamily(logger)
	if err != nil {
		logger.ErrorWF("ApplyChangeAlliance GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}

	familyInfoModel, err := familymodel.LoadFamilyInfoModel(logger, familyID)
	if err != nil {
		logger.ErrorWF("ApplyChangeAlliance LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	// 检查用户权限是否为族长
	if !familyInfoModel.CheckHaveLeader(logger, []uint64{userID}) {
		logger.ErrorWF("ApplyChangeAlliance User not leader",
			zap.Uint64("userID", userID), zap.Int32("familyID", familyID))
		return errors.New("user not leader")
	}

	// todo 根据后续限定来修改 家族是否能够更改

	// 更改家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(logger, familyID)

	nowAllianceID, err := familyToAllianceModel.GetUserAlliance(logger)
	if err != nil {
		logger.ErrorWF("ApplyChangeAlliance GetUserAlliance err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	if nowAllianceID == allianceID {
		return nil
	}

	if nowAllianceID != 0 {
		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(logger, nowAllianceID)
		if err != nil {
			logger.ErrorWF("ApplyChangeAlliance LoadAllianceInfoModel err",
				zap.Int32("allianceID", nowAllianceID), zap.Error(err))
			return err
		}
		err = allianceInfoModel.RemoveFamilyID(logger, familyID)
		if err != nil {
			logger.ErrorWF("ApplyChangeAlliance RemoveFamilyID err",
				zap.Int32("allianceID", nowAllianceID), zap.Error(err))
			return err
		}
	}

	if allianceID != 0 {
		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(logger, allianceID)
		if err != nil {
			logger.ErrorWF("ApplyChangeAlliance LoadAllianceInfoModel err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			return err
		}
		err = allianceInfoModel.AddFamilyID(logger, familyID)
		if err != nil {
			logger.ErrorWF("ApplyChangeAlliance AddFamilyID err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			return err
		}
	}

	err = familyToAllianceModel.SetUserAlliance(logger, allianceID)
	if err != nil {
		logger.ErrorWF("ApplyChangeAlliance SetUserAlliance err",
			zap.Int32("familyID", familyID), zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}

	return nil
}
