package alliancemodel

import (
	"maze_game_server/io/redis/allianceredis"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

type FamilyToAllianceModel struct {
	FamilyID int32 `json:"family_id,omitempty"`
}

func NewFamilyToAllianceModel(logger fklog.FKLogI, familyID int32) *FamilyToAllianceModel {
	return &FamilyToAllianceModel{
		FamilyID: familyID,
	}
}

func (r *FamilyToAllianceModel) GetUserAlliance(logger fklog.FKLogI) (int32, error) {
	allianceID, err := allianceredis.GetFamilyToAlliance(r.FamilyID)
	if err != nil {
		logger.ErrorWF("GetUserAlliance GetUserAlliance err",
			zap.Int32("familyID", r.FamilyID), zap.Error(err))
		return 0, err
	}
	return allianceID, nil
}

func (r *FamilyToAllianceModel) SetUserAlliance(logger fklog.FKLogI, allianceID int32) error {
	return allianceredis.SetFamilyToAlliance(r.FamilyID, allianceID)
}

func (r *FamilyToAllianceModel) DelUserAlliance(logger fklog.FKLogI) error {
	return allianceredis.DelFamilyToAlliance(r.FamilyID)
}
