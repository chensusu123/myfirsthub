package familyservice

import (
	"context"
	"maze_game_server/model/alliancemodel"
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) SysCreateFamily(ctx context.Context, allianceID int32) (r *familymodel.FamilyInfoModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "SysCreateFamily Start",
		zap.Any("allianceID", allianceID),
	)
	defer func() {
		logger.CtxInfo(ctx, "SysCreateFamily End",
			zap.Any("allianceID", allianceID),
			zap.Any("familymodel", r),
			zap.Any("err", err),
		)
	}()

	r, err = familymodel.SysCreateFamilyInfoModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "SysCreateFamily SysCreateFamilyInfoModel failed",
			zap.Error(err),
			zap.Any("allianceID", allianceID),
		)
		return nil, err
	}
	// 加入成员
	r.AddMember(ctx, familymodel.FamilyMember{
		UserID:         1,
		NickName:       "admin",
		Sex:            1,
		Avatar:         "",
		Level:          1,
		PrivilegeLevel: 10,
	})

	err = r.Save(ctx, r.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "SysCreateFamily model.Save failed",
			zap.Error(err))
		return nil, err
	}

	// 设置家族对应联盟
	familyToAllianceModel := alliancemodel.NewFamilyToAllianceModel(ctx, r.FamilyID)
	err = familyToAllianceModel.SetUserAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "SysCreateFamily familyToAllianceModel.SetUserAlliance err",
			zap.Int32("allianceID", allianceID),
			zap.Int32("familyID", r.FamilyID),
			zap.Error(err))
		return nil, err
	}

	return
}
