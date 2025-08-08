package costumeservice

import "gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

type CostumeService interface {
	// 获取用户装扮
	GetUserCostume(logger fklog.FKLogI, userId uint64) (map[int32]int32, error)
	ChangeCostume(logger fklog.FKLogI, userId uint64)
}

var GlobalCostumeService CostumeService

func init() {
	GlobalCostumeService = newCostumeService()
}

type service struct {
}

func newCostumeService() CostumeService {
	return &service{}
}
