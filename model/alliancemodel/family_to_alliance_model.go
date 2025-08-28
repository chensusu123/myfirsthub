package alliancemodel

import (
	"context"
	"maze_game_server/io/redis/allianceredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type FamilyToAllianceModel struct {
	FamilyID int32 `json:"family_id,omitempty"`
}

func NewFamilyToAllianceModel(ctx context.Context, familyID int32) *FamilyToAllianceModel {
	return &FamilyToAllianceModel{
		FamilyID: familyID,
	}
}

func (r *FamilyToAllianceModel) GetUserAlliance(ctx context.Context) (int32, error) {
	logger := fklog.ContextAppLogger(ctx)
	allianceID, err := allianceredis.GetFamilyToAlliance(r.FamilyID)
	if err != nil {
		logger.CtxError(ctx, "GetUserAlliance GetUserAlliance err",
			zap.Int32("familyID", r.FamilyID), zap.Error(err))
		return 0, err
	}
	return allianceID, nil
}

func (r *FamilyToAllianceModel) SetUserAlliance(ctx context.Context, allianceID int32) error {
	return allianceredis.SetFamilyToAlliance(r.FamilyID, allianceID)
}

func (r *FamilyToAllianceModel) DelUserAlliance(ctx context.Context) error {
	return allianceredis.DelFamilyToAlliance(r.FamilyID)
}
