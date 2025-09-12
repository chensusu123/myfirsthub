package familymodel

import (
	"context"
	"maze_game_server/io/redis/familyredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type FamilySubscribeModel struct {
	FamilyId int32   `json:"family_id"`
	UserId   uint64  `json:"user_id"`
	GroupIds []int64 `json:"group_ids"`
}

func LoadFamilyScribeModel(ctx context.Context, familyID int32, userId uint64) (*FamilySubscribeModel, error) {
	logger := fklog.ContextAppLogger(ctx)
	r := &FamilySubscribeModel{FamilyId: familyID, UserId: userId}
	if err := r.Load(ctx); err != nil {
		logger.CtxError(ctx, "LoadFamilyModel err",
			zap.Int32("familyID", familyID), zap.Error(err))
		return nil, err
	}
	return r, nil
}

func (m *FamilySubscribeModel) Load(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	groupIDs, err := familyredis.GetSubscribe(ctx, m.UserId)
	if err != nil {
		logger.CtxError(ctx, "GetSubscribe err",
			zap.Int32("familyID", m.FamilyId), zap.Uint64("userID", m.UserId), zap.Error(err))
		return err
	}

	m.GroupIds = groupIDs
	return nil
}

func (m *FamilySubscribeModel) Save(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	err := familyredis.SaveSubscribe(ctx, m.UserId, m.GroupIds)
	if err != nil {
		logger.CtxError(ctx, "SaveSubscribe err",
			zap.Int32("familyID", m.FamilyId), zap.Uint64("userID", m.UserId), zap.Error(err))
		return err
	}
	return nil
}

func (m *FamilySubscribeModel) Delete(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)
	err := familyredis.DelSubscribe(ctx, m.UserId)
	if err != nil {
		logger.CtxError(ctx, "DelSubscribe err",
			zap.Int32("familyID", m.FamilyId), zap.Uint64("userID", m.UserId), zap.Error(err))
		return err
	}
	return nil
}
