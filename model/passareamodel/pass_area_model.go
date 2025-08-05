package passareamodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/passarearedis"
	"maze_game_server/lib/serialize"
)

type PassAreaInfo struct {
	AreaId    int32 `json:"areaId,omitempty"`
	AreaIndex int32 `json:"areaIndex,omitempty"`
}

type PassAreaModel struct {
	PassAreaList []*PassAreaInfo `json:"pass_area_list,omitempty"`
}

func NewPassAreaModel(logger fklog.FKLogI, userID uint64, stageId int32) (*PassAreaModel, error) {
	passArea := &PassAreaModel{}
	if err := passArea.load(logger, userID, stageId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *PassAreaModel) load(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := passarearedis.GetBarrierPassArea(logger, userID, stageId)
	if err != nil {
		return err
	}
	if bytes == nil {
		p.PassAreaList = make([]*PassAreaInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, p)
	if err != nil {
		logger.ErrorWF("TempBuff load Unmarshal failed", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("stageId", stageId))
		return err
	}
	return
}

func (p *PassAreaModel) Save(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := serialize.Marshal(p)
	if err != nil {
		logger.ErrorWF("TempBuff save Marshal failed", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("stageId", stageId))
		return err
	}
	return passarearedis.SetBarrierPassArea(logger, userID, stageId, bytes)
}

func (p *PassAreaModel) Del(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	return passarearedis.DelBarrierPassArea(logger, userID, stageId)
}
