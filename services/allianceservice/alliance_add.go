package allianceservice

import (
	"context"

	"maze_game_server/app"

	"maze_game_server/model/alliancemodel"
	"time"

	grouppkg "maze_game_server/io/redis/im/group"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

func (s *service) AddAlliance(ctx context.Context, allianceName string) error {
	logger := fklog.ContextAppLogger(ctx)
	allianceListModel, err := alliancemodel.LoadAllianceListModel(ctx)
	if err != nil {
		logger.CtxError(ctx, "AddAlliance LoadAllianceListModel err", zap.Error(err))
		return err
	}
	allianceID := allianceListModel.GetAllianceID(ctx)

	allianceInfoModel := alliancemodel.NewAllianceInfoModel(ctx, allianceID, allianceName)
	allianceInfoModel.SetAllianceGroupID(ctx, int32(time.Now().Unix()))
	err = allianceInfoModel.Save(ctx)
	if err != nil {
		logger.CtxError(ctx, "AddAlliance Save allianceInfoModel err", zap.Error(err))
		return err
	}
	// 创建联盟聊天组
	_, err = grouppkg.CreateGroup(ctx, app.Maze.ID(), uint64(0), allianceInfoModel.AllianceGroupID, make([]uint64, 0))
	if err != nil {
		logger.CtxError(ctx, "AddAlliance CreateGroup err", zap.Error(err))
		return err
	}

	err = allianceListModel.AddAlliance(ctx, allianceID)
	if err != nil {
		logger.CtxError(ctx, "AddAlliance AddAlliance err", zap.Error(err))
		return err
	}
	return nil
}
