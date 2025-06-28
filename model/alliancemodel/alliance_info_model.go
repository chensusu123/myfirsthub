package alliancemodel

import (
	"maze_game_server/io/redis/allianceredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/MazeFamily"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	allianceCountLimit = 100
)

type AllianceInfoModel struct {
	AllianceID         int32   `json:"alliance_id"`
	AllianceName       string  `json:"alliance_name"`
	AllianceCountLimit int32   `json:"alliance_count_limit"` // 联盟中家族数量限制
	FamilyIDs          []int32 `json:"family_ids"`           // 联盟中家族ID列表
}

func LoadAllianceInfoModel(logger fklog.FKLogI, allianceID int32) (r *AllianceInfoModel, err error) {
	r = &AllianceInfoModel{}
	if err = r.load(logger, allianceID); err != nil {
		logger.ErrorWF("LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return nil, err
	}
	return
}

func NewAllianceInfoModel(logger fklog.FKLogI, allianceID int32, allianceName string) *AllianceInfoModel {
	return &AllianceInfoModel{
		AllianceID:         allianceID,
		AllianceName:       allianceName,
		AllianceCountLimit: allianceCountLimit,
		FamilyIDs:          []int32{},
	}
}
func (r *AllianceInfoModel) load(logger fklog.FKLogI, allianceID int32) (err error) {
	value, err := allianceredis.GetAllianceInfo(allianceID)
	if err != nil {
		logger.ErrorWF("LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	if value == nil {
		return nil
	}
	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.ErrorWF("LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	return
}

func (r *AllianceInfoModel) Save(logger fklog.FKLogI) (err error) {
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.ErrorWF("Save err", zap.Error(err))
		return err
	}
	return allianceredis.SetAllianceInfo(r.AllianceID, value)
}

func (r *AllianceInfoModel) DataToAllianceInfoPb() *MazeFamily.AllianceInfo {
	return &MazeFamily.AllianceInfo{
		AllianceId:         proto.Int32(r.AllianceID),
		AllianceName:       proto.String(r.AllianceName),
		AllianceCountLimit: proto.Int32(r.AllianceCountLimit),
		FamilyIds:          r.FamilyIDs,
	}
}

func (r *AllianceInfoModel) AddFamilyID(logger fklog.FKLogI, familyID int32) error {
	r.FamilyIDs = append(r.FamilyIDs, familyID)
	return r.Save(logger)
}

func (r *AllianceInfoModel) RemoveFamilyID(logger fklog.FKLogI, familyID int32) error {
	for i, id := range r.FamilyIDs {
		if id == familyID {
			r.FamilyIDs = append(r.FamilyIDs[:i], r.FamilyIDs[i+1:]...)
			return r.Save(logger)
		}
	}
	return nil
}
