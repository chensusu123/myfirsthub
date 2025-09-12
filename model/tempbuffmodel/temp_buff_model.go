package tempbuffmodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
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

func getKey(userId uint64, barrierId int32) string {
	return fmt.Sprintf("tempbuff:u:%d:barrier:%d", userId, barrierId)
}

func NewTempBuffInfoModel(ctx context.Context, userID uint64, barrierId int32) (*TempBuffInfoModel, error) {
	tempBuff := &TempBuffInfoModel{}
	if err := tempBuff.load(ctx, userID, barrierId); err != nil {
		return nil, err
	}
	return tempBuff, nil
}

func (t *TempBuffInfoModel) load(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierId), t)
}

func (t *TempBuffInfoModel) Save(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierId), t)
}

func (t *TempBuffInfoModel) Del(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierId))
}
