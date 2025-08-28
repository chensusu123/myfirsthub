package familyservice

import (
	"context"
	"errors"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) DissolutionFamily(ctx context.Context, uid uint64, familyID int32) error {
	logger := fklog.ContextAppLogger(ctx)
	familyModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyModel err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	// 检查用户权限
	if !familyModel.CheckHaveLeader(ctx, []uint64{uid}) {
		logger.CtxError(ctx, "DissolutionFamily CheckHaveLeader err", zap.Int32("familyID", familyID), zap.Uint64("uid", uid))
		return errors.New("用户没有权限解散家族")
	}

	// 删除家族
	familyModel.Delete(ctx, familyID)

	// 删除家族列表
	familyListModel, err := familymodel.LoadFamilyListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "LoadFamilyListModel err", zap.Error(err))
		return err
	}
	familyListModel.RemoveFamily(ctx, familyID)

	// 删除联盟内该家族ID
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, familyID)
	allianceID, err := familyToAllianceModel.GetUserAlliance(ctx)
	if err != nil {
		logger.CtxError(ctx, "DissolutionFamily GetUserAlliance err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}
	allianceInfoModel, err := alliancemodel.LoadAllianceInfoModel(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "DissolutionFamily LoadAllianceInfoModel err", zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	allianceInfoModel.RemoveFamilyID(ctx, familyID)

	// 删除家族对应联盟
	err = familyToAllianceModel.DelUserAlliance(ctx)
	if err != nil {
		logger.CtxError(ctx, "DissolutionFamily DelUserAlliance err", zap.Int32("familyID", familyID), zap.Error(err))
		return err
	}

	return nil
}
