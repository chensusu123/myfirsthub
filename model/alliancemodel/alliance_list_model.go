package alliancemodel

import (
	"context"
	"maze_game_server/io/redis/allianceredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/MazeFamily"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type AllianceListModel struct {
	Alliances []int32 `json:"alliances,omitempty"`
}

func LoadAllianceListModel(ctx context.Context) (r *AllianceListModel, err error) {
	logger := fklog.ContextAppLogger(ctx)
	r = &AllianceListModel{}
	if err = r.load(ctx); err != nil {
		logger.CtxError(ctx, "LoadAllianceListModel err", zap.Error(err))
		return nil, err
	}
	return
}

func (r *AllianceListModel) load(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := allianceredis.GetAllianceList()
	if err != nil {
		logger.CtxError(ctx, "LoadAllianceListModel err", zap.Error(err))
		return err
	}
	if value == nil {
		return nil
	}
	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.CtxError(ctx, "LoadAllianceListModel err", zap.Error(err))
		return err
	}
	return
}

func (r *AllianceListModel) Save(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.CtxError(ctx, "Save err", zap.Error(err))
		return err
	}
	return allianceredis.SetAllianceList(value)
}

func (r *AllianceListModel) Delete(ctx context.Context) (err error) {
	return allianceredis.DelAllianceList()
}

func (r *AllianceListModel) AddAlliance(ctx context.Context, allianceID int32) error {
	r.Alliances = append(r.Alliances, allianceID)
	return r.Save(ctx)
}

func (r *AllianceListModel) RemoveAlliance(ctx context.Context, allianceID int32) error {
	for i, v := range r.Alliances {
		if v == allianceID {
			r.Alliances = append(r.Alliances[:i], r.Alliances[i+1:]...)
			return r.Save(ctx)
		}
	}
	return nil
}

func (r *AllianceListModel) DataToAllianceListPb(ctx context.Context) []*MazeFamily.AllianceInfo {
	logger := fklog.ContextAppLogger(ctx)
	allianceList := make([]*MazeFamily.AllianceInfo, 0)
	for _, allianceID := range r.Alliances {
		allianceInfo, err := LoadAllianceInfoModel(ctx, allianceID)
		if err != nil {
			logger.CtxError(ctx, "DataToAllianceListPb LoadAllianceInfoModel err", zap.Error(err))
			continue
		}
		allianceList = append(allianceList, allianceInfo.DataToAllianceInfoPb())
	}
	return allianceList
}

func (r *AllianceListModel) GetAllianceID() int32 {
	return allianceredis.GetAllianceID()
}
