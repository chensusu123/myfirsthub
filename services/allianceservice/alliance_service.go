package allianceservice

import (
	"context"
	"maze_game_server/model/alliancemodel"
)

type AllianceService interface {
	// 查询联盟信息
	QueryAllianceInfo(ctx context.Context, allianceID int32) (*alliancemodel.AllianceInfoModel, error)
	// 查询联盟列表
	QueryAllianceList(ctx context.Context) (*alliancemodel.AllianceListModel, error)
	// 查询用户所在联盟
	QueryUserAlliance(ctx context.Context, userID uint64) (allianceID int32, err error)
	// 申请更改联盟
	ApplyChangeAlliance(ctx context.Context, userID uint64, allianceID int32) error
	// 增加联盟
	CreateAlliance(ctx context.Context, allianceName string) error
	//订阅联盟
	SubscribeAllianceChat(ctx context.Context, userID uint64, allianceID int32, group_ids []int64) error
}

var GlobalAllianceService AllianceService

func init() {
	GlobalAllianceService = newAllianceService()
}

type service struct {
}

func newAllianceService() AllianceService {
	return &service{}
}
