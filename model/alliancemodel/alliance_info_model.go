package alliancemodel

import (
	"context"
	"maze_game_server/io/redis/allianceredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/MazeFamily"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
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
	AllianceGroupID    int32   `json:"alliance_group_id"`    // 联盟所属组ID
}

func LoadAllianceInfoModel(ctx context.Context, allianceID int32) (r *AllianceInfoModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	r = &AllianceInfoModel{}
	if err = r.load(ctx, allianceID); err != nil {
		logger.CtxError(ctx, "LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return nil, err
	}
	return
}

func NewAllianceInfoModel(ctx context.Context, allianceID int32, allianceName string) *AllianceInfoModel {
	return &AllianceInfoModel{
		AllianceID:         allianceID,
		AllianceName:       allianceName,
		AllianceCountLimit: allianceCountLimit,
		FamilyIDs:          []int32{},
		AllianceGroupID:    0,
	}
}
func (r *AllianceInfoModel) load(ctx context.Context, allianceID int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := allianceredis.GetAllianceInfo(allianceID)
	if err != nil {
		logger.CtxError(ctx, "LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	if value == nil {
		return nil
	}
	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.CtxError(ctx, "LoadAllianceInfoModel err",
			zap.Int32("allianceID", allianceID), zap.Error(err))
		return err
	}
	return
}

func (r *AllianceInfoModel) Save(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.CtxError(ctx, "Save err", zap.Error(err))
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

func (r *AllianceInfoModel) AddFamilyID(ctx context.Context, familyID int32) error {
	r.FamilyIDs = append(r.FamilyIDs, familyID)
	return r.Save(ctx)
}

func (r *AllianceInfoModel) RemoveFamilyID(ctx context.Context, familyID int32) error {
	for i, id := range r.FamilyIDs {
		if id == familyID {
			r.FamilyIDs = append(r.FamilyIDs[:i], r.FamilyIDs[i+1:]...)
			return r.Save(ctx)
		}
	}
	return nil
}

func (r *AllianceInfoModel) SetAllianceGroupID(ctx context.Context, allianceGroupID int32) error {
	if r.AllianceGroupID != 0 {
		return nil
	}
	r.AllianceGroupID = allianceGroupID
	return r.Save(ctx)
}
