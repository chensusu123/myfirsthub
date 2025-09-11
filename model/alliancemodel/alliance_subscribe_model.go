package alliancemodel

import (
	"context"
	"maze_game_server/io/redis/allianceredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type AllianceSubscribeModel struct {
	AllianceID int32   `json:"alliance_id"`
	UserID     uint64  `json:"user_id"`
	GroupIDs   []int64 `json:"group_ids"`
}

func LoadFamilyScribeModel(ctx context.Context, allianceID int32, userId uint64) (*AllianceSubscribeModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	r := &AllianceSubscribeModel{AllianceID: allianceID, UserID: userId}
	if err := r.Load(ctx); err != nil {
		logger.CtxError(ctx, "LoadFamilyModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return nil, err
	}
	return r, nil
}

func (m *AllianceSubscribeModel) Load(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	groupIDs, err := allianceredis.GetSubscribe(ctx, m.UserID)
	if err != nil {
		logger.CtxError(ctx, "GetSubscribe err",
			zap.Int32("allianceID", m.AllianceID), zap.Uint64("userID", m.UserID), zap.Error(err))
		return err
	}

	m.GroupIDs = groupIDs
	return nil
}

func (m *AllianceSubscribeModel) Save(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	err := allianceredis.SaveSubscribe(ctx, m.UserID, m.GroupIDs)
	if err != nil {
		logger.CtxError(ctx, "SaveSubscribe err",
			zap.Int32("allianceID", m.AllianceID), zap.Uint64("userID", m.UserID), zap.Error(err))
		return err
	}
	return nil
}

func (m *AllianceSubscribeModel) Delete(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	err := allianceredis.DelSubscribe(ctx, m.UserID)
	if err != nil {
		logger.CtxError(ctx, "DelSubscribe err",
			zap.Int32("allianceID", m.AllianceID), zap.Uint64("userID", m.UserID), zap.Error(err))
		return err
	}
	return nil
}
