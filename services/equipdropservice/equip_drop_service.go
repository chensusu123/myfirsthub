package equipdropservice

import (
	"context"
	"maze_game_server/model/equipdropmodel"
)

type EquipDropService interface {
	//获取装备特殊掉落信息
	GetMazeEquipSpecialDropInfo(ctx context.Context, userId uint64) (*equipdropmodel.EquipSpecialDropModel, error)
	//装备掉落
	GetNewEquip(ctx context.Context, userId uint64, mazeLevel int32, barrierId int32, equipNum int32) (newEquip map[int32]int32, err error)

	GetMazeBarrierLv(level int32, barrier int32) int32

	GmDelete(ctx context.Context, userId uint64) error
}

var GlobalEquipDropService EquipDropService

func init() {
	GlobalEquipDropService = newEquipDropService()
}

type service struct {
}

func newEquipDropService() EquipDropService {
	return &service{}
}

const equipSpecialDropRate = 10000
