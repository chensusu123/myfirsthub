package allianceservice

import (
	"context"

	"maze_game_server/app"
	"maze_game_server/services/familyservice"

	"maze_game_server/model/alliancemodel"

	grouppkg "maze_game_server/io/redis/im/group"

	"maze_game_server/common/constdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) CreateAlliance(ctx context.Context, allianceName string) (int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	allianceListModel, err := alliancemodel.LoadAllianceListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance LoadAllianceListModel err", zap.Error(err))
		return 0, err
	}
	allianceID := allianceListModel.GetAllianceID(ctx)

	allianceInfoModel := alliancemodel.NewAllianceInfoModel(ctx, allianceID, allianceName)

	r, err := familyservice.GlobalFamilyService.SysCreateFamily(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance SysCreateFamily Fail",
			zap.Any("allianceID", allianceID),
			zap.Any("allianceName", allianceName),
		)
		return 0, err
	}

	allianceInfoModel.SetDeafultFamily(r.FamilyID)

	err = allianceInfoModel.Save(ctx)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance Save allianceInfoModel err", zap.Error(err))
		return 0, err
	}
	// 创建联盟聊天组
	_, err = grouppkg.CreateGroup(ctx, app.Maze.ID(), int64(0), allianceInfoModel.AllianceGroupID, constdef.GroupTypeLeague, make([]int64, 0))
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance CreateGroup err", zap.Error(err))
		return 0, err
	}

	err = allianceListModel.AddAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "CreateAlliance AddAlliance err", zap.Error(err))
		return 0, err
	}
	return allianceID, nil
}
