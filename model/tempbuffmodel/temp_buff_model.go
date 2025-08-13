package tempbuffmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/tempbuffredis"
	"maze_game_server/lib/serialize"
)

type TempBuffInfoModel struct {
	BuffSequence *BuffSequence       `json:"buff_sequence,omitempty"`
	TotalBuff    []*TotalBuffInfo    `json:"total_buff,omitempty"`
	SelectedBuff []*SelectedBuffInfo `json:"selected_buff,omitempty"`
}
type BuffSequence struct {
	Level            int32   `json:"level,omitempty"`            // 当前buff等级
	RefreshCount     int32   `json:"refresh_count,omitempty"`    // 已刷新次数
	OptionalBuffList []int32 `json:"select_buff_list,omitempty"` // 可选buff列表
	AreaId           int32   `json:"area_id,omitempty"`          // 选择buff的区域
	AreaIndex        int32   `json:"area_index,omitempty"`       // 选择buff的子区域
}
type TotalBuffInfo struct {
	BuffId    int32 `json:"buff_id,omitempty"`
	BuffValue int64 `json:"buff_value,omitempty"`
}
type SelectedBuffInfo struct {
	BuffId    int32 `json:"buff_id,omitempty"`
	Level     int32 `json:"level,omitempty"`      // buff对应的等级
	Type      int32 `json:"type,omitempty"`       // buff选择类型  0等级选择 1道具选择
	AreaId    int32 `json:"area_id,omitempty"`    // 选择buff的区域
	AreaIndex int32 `json:"area_index,omitempty"` // 选择buff的子区域
}

func NewTempBuffInfoModel(logger fklog.FKLogI, userID uint64, stageId int32) (*TempBuffInfoModel, error) {
	tempBuff := &TempBuffInfoModel{}
	if err := tempBuff.load(logger, userID, stageId); err != nil {
		return nil, err
	}
	return tempBuff, nil
}

func (tb *TempBuffInfoModel) load(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := tempbuffredis.GetMazeTempBuff(logger, userID, stageId)
	if err != nil {
		return err
	}
	if bytes == nil {
		return nil
	}
	err = serialize.Unmarshal(bytes, tb)
	if err != nil {
		logger.ErrorWF("TempBuff load Unmarshal failed", zap.Error(err), zap.Int32("stageId", stageId))
		return err
	}
	return
}

func (tb *TempBuffInfoModel) Save(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	bytes, err := serialize.Marshal(tb)
	if err != nil {
		logger.ErrorWF("TempBuff save Marshal failed", zap.Error(err), zap.Int32("stageId", stageId))
		return err
	}
	return tempbuffredis.SetMazeTempBuff(logger, userID, stageId, bytes)
}

func (tb *TempBuffInfoModel) Del(logger fklog.FKLogI, userID uint64, stageId int32) (err error) {
	return tempbuffredis.DelMazeTempBuff(logger, userID, stageId)
}
