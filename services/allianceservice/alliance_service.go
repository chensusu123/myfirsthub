package allianceservice

import (
	"maze_game_server/model/alliancemodel"

	"gitlab.ifreetalk.com/nano-ecosystem/fklog"
)

type AllianceService interface {
	// 查询联盟信息
	QueryAllianceInfo(logger fklog.FKLogI, allianceID int32) (*alliancemodel.AllianceInfoModel, error)
	// 查询联盟列表
	QueryAllianceList(logger fklog.FKLogI) (*alliancemodel.AllianceListModel, error)
	// 查询用户所在联盟
	QueryUserAlliance(logger fklog.FKLogI, userID uint64) (allianceID int32, err error)
	// 申请更改联盟
	ApplyChangeAlliance(logger fklog.FKLogI, userID uint64, allianceID int32) error
	// 增加联盟
	AddAlliance(logger fklog.FKLogI, allianceName string) error
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
