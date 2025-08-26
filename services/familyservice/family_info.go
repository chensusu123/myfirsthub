package familyservice

import (
	"context"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// GetFamilyInfo 查询家族详情
func (r *service) GetFamilyInfo(ctx context.Context, familyID int32) (*familymodel.FamilyInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	familyInfoModel, err := familymodel.LoadFamilyInfoModel(ctx, familyID)
	if err != nil {
		logger.CtxError(ctx, "GetFamilyInfo LoadFamilyInfoModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}

	return familyInfoModel, nil
}

// GetFamilyList 获取家族列表
func (r *service) GetFamilyList(ctx context.Context) (familymodel.FamilysInfoModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	// 查询当前家族id列表
	familyListModel, err := familymodel.LoadFamilyListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "GetFamilyList LoadFamilyListModel err",
			zap.Error(err))
		return nil, err
	}
	// 查询家族列表
	familysInfoModel, err := familymodel.LoadFamilyListInfoModel(ctx, familyListModel.Familys)
	if err != nil {
		logger.CtxError(ctx, "GetFamilyInfo LoadFamilyInfoModel err",
			zap.Error(err))
		return nil, err
	}
	return familysInfoModel, nil
}
