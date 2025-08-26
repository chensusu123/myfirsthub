package costumeservice

import (
	"context"
)

type CostumeService interface {
	// 获取用户装扮
	GetUserCostume(ctx context.Context, userId uint64) (map[int32]int32, error)
	ChangeCostume(ctx context.Context, userId uint64)
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
