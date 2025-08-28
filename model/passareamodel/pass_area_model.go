package passareamodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
)

type PassAreaInfo struct {
	AreaId    int32 `json:"areaId,omitempty"`
	AreaIndex int32 `json:"areaIndex,omitempty"`
}

type PassAreaModel struct {
	PassAreaList []*PassAreaInfo `json:"pass_area_list,omitempty"`
}

func getKey(userId uint64, barrierId int32) string {
	return fmt.Sprintf("area:u:%d:barrier:%d", userId, barrierId)
}

func NewPassAreaModel(ctx context.Context, userID uint64, stageId int32) (*PassAreaModel, error) {
	passArea := &PassAreaModel{
		PassAreaList: make([]*PassAreaInfo, 0),
	}
	if err := passArea.load(ctx, userID, stageId); err != nil {
		return nil, err
	}
	return passArea, nil
}

func (p *PassAreaModel) load(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.LoadSvrData(ctx, getKey(userID, barrierId), p)
}

func (p *PassAreaModel) Save(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.SaveSvrData(ctx, getKey(userID, barrierId), p)
}

func (p *PassAreaModel) Del(ctx context.Context, userID uint64, barrierId int32) (err error) {
	return io.DeleteSvrData(ctx, getKey(userID, barrierId))
}
