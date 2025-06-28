package alliancemodel

import (
	"maze_game_server/io/redis/allianceredis"
	"maze_game_server/lib/serialize"
	"maze_game_server/pb/common/MazeFamily"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
	"go.uber.org/zap"
)

type AllianceListModel struct {
	Alliances []int32 `json:"alliances,omitempty"`
}

func LoadAllianceListModel(logger fklog.FKLogI) (r *AllianceListModel, err error) {
	r = &AllianceListModel{}
	if err = r.load(logger); err != nil {
		logger.ErrorWF("LoadAllianceListModel err", zap.Error(err))
		return nil, err
	}
	return
}

func (r *AllianceListModel) load(logger fklog.FKLogI) (err error) {
	value, err := allianceredis.GetAllianceList()
	if err != nil {
		logger.ErrorWF("LoadAllianceListModel err", zap.Error(err))
		return err
	}
	if value == nil {
		return nil
	}
	err = serialize.Unmarshal(value, r)
	if err != nil {
		logger.ErrorWF("LoadAllianceListModel err", zap.Error(err))
		return err
	}
	return
}

func (r *AllianceListModel) Save(logger fklog.FKLogI) (err error) {
	value, err := serialize.Marshal(r)
	if err != nil {
		logger.ErrorWF("Save err", zap.Error(err))
		return err
	}
	return allianceredis.SetAllianceList(value)
}

func (r *AllianceListModel) Delete(logger fklog.FKLogI) (err error) {
	return allianceredis.DelAllianceList()
}

func (r *AllianceListModel) AddAlliance(logger fklog.FKLogI, allianceID int32) error {
	r.Alliances = append(r.Alliances, allianceID)
	return r.Save(logger)
}

func (r *AllianceListModel) RemoveAlliance(logger fklog.FKLogI, allianceID int32) error {
	for i, v := range r.Alliances {
		if v == allianceID {
			r.Alliances = append(r.Alliances[:i], r.Alliances[i+1:]...)
			return r.Save(logger)
		}
	}
	return nil
}

func (r *AllianceListModel) DataToAllianceListPb(logger fklog.FKLogI) []*MazeFamily.AllianceInfo {
	allianceList := make([]*MazeFamily.AllianceInfo, 0)
	for _, allianceID := range r.Alliances {
		allianceInfo, err := LoadAllianceInfoModel(logger, allianceID)
		if err != nil {
			logger.ErrorWF("DataToAllianceListPb LoadAllianceInfoModel err", zap.Error(err))
			continue
		}
		allianceList = append(allianceList, allianceInfo.DataToAllianceInfoPb())
	}
	return allianceList
}

func (r *AllianceListModel) GetAllianceID() int32 {
	return allianceredis.GetAllianceID()
}
