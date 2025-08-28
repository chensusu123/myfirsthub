package allianceservice

import (
	"context"
	"errors"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) ApplyChangeAlliance(ctx context.Context, userID uint64, allianceID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	userFamilyModel := familymodel.NewUserFamilyModel(ctx, userID)
	familyID, err := userFamilyModel.GetUserFamily(ctx)
	if err != nil {
		logger.CtxError(ctx, "ApplyChangeAlliance GetUserFamily err",
			zap.Uint64("userID", userID), zap.Error(err))
		return err
	}

	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "ApplyChangeAlliance LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	// 检查用户权限是否为族长
	if !familyInfoModel.CheckHaveLeader(ctx, []uint64{userID}) {
		logger.CtxError(ctx, "ApplyChangeAlliance User not leader",
			zap.Uint64("userID", userID), zap.Int32("familyID", familyID))
		return errors.New("user not leader")
	}

	// todo 根据后续限定来修改 家族是否能够更改

	// 更改家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, familyID)

	nowAllianceID, err := familyToAllianceModel.GetUserAlliance(ctx)
	if err != nil {
		logger.CtxError(ctx, "ApplyChangeAlliance GetUserAlliance err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	if nowAllianceID == allianceID {
		return nil
	}

	if nowAllianceID != 0 {
		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, nowAllianceID)
		if err != nil {
			logger.CtxError(ctx, "ApplyChangeAlliance LoadAllianceInfoModel err",
				zap.Int32("allianceID", nowAllianceID), zap.Error(err))
			return err
		}
		err = allianceInfoModel.RemoveFamilyID(ctx, familyID)
		if err != nil {
			logger.CtxError(ctx, "ApplyChangeAlliance RemoveFamilyID err",
				zap.Int32("allianceID", nowAllianceID), zap.Error(err))
			return err
		}
	}

	if allianceID != 0 {
		allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
		if err != nil {
			logger.CtxError(ctx, "ApplyChangeAlliance LoadAllianceInfoModel err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			return err
		}
		err = allianceInfoModel.AddFamilyID(ctx, familyID)
		if err != nil {
			logger.CtxError(ctx, "ApplyChangeAlliance AddFamilyID err",
				zap.Int32("allianceID", allianceID), zap.Error(err))
			return err
		}
	}

	err = familyToAllianceModel.SetUserAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "ApplyChangeAlliance SetUserAlliance err",
			zap.Int32("familyID", familyID), zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}

	return nil
}
