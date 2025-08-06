package equipdropservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/equipdropmodel"
)

type EquipDropService interface {
	//获取装备特殊掉落信息
	GetMazeEquipSpecialDropInfo(logger fklog.FKLogI, userId uint64) (*equipdropmodel.EquipSpecialDropModel, error)
	//装备掉落
	GetNewEquip(logger fklog.FKLogI, userId uint64, mazeLevel int32, barrierId int32, equipNum int32) (newEquip map[int32]int32, err error)

	GetMazeBarrierLv(level int32, barrier int32) int32
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
