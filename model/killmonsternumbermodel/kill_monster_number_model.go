package passareamodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/passarearedis"
	"maze_game_server/lib/serialize"
)

type KillMonsterNumberAreaInfo struct {
	KillMonsterNumber int32 `json:"kill_monster_number,omitempty"`
	AreaIndex         int32 `json:"area_index,omitempty"`
}

type KillMonsterNumberModel struct {
	PassAreaList []*KillMonsterNumberAreaInfo `json:"pass_area_list,omitempty"`
}

func NewKillMonsterNumberModel(logger fklog.FKLogI, userID uint64, stageId int32) (*PassAreaModel, error) {
	passArea := &PassAreaModel{}
	if err := passArea.load(logger, userID, stageId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *KillMonsterNumberModel) load(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
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

func (p *KillMonsterNumberModel) Save(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := serialize.Marshal(p)
	if err != nil {
		logger.ErrorWF("TempBuff save Marshal failed", zap.Error(err), zap.Uint64("userID", userID), zap.Int32("stageId", stageId))
		return err
	}
	return passarearedis.SetBarrierPassArea(logger, userID, stageId, bytes)
}

func (p *KillMonsterNumberModel) Del(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	return passarearedis.DelBarrierPassArea(logger, userID, stageId)
}
