package tempbuffservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/pb/common/MazeCommon"
)

// TempBuffService 临时buff服务可用接口
type TempBuffService interface {
	// GetMazeTempBuffList 获取临时buff列表
	GetMazeTempBuffList(logger fklog.FKLogI, userId uint64, stageId int32) ([]*BuffInfo, error)

	// GetOptionalTempBuffList 获取可选临时buff列表
	GetOptionalTempBuffList(logger fklog.FKLogI, userId uint64, stageId, level, buffType, areaId, areaIndex int32) (*OptionalBuffInfo, error)

	// RefreshOptionalMazeTempBuffList 刷新可选临时buff列表
	RefreshOptionalMazeTempBuffList(logger fklog.FKLogI, userId uint64, stageId, level, areaId int32, cost []*MazeCommon.MazeItem) (*OptionalBuffInfo, error)

	// SelectMazeTempBuffRQ 选择临时buff
	SelectMazeTempBuff(logger fklog.FKLogI, userId uint64, stageId, level, buffId, buffType int32) ([]*BuffInfo, error)

	// DelTempBuff 删除所有临时buff
	DelTempBuff(logger fklog.FKLogI, userID uint64, stageId int32) error

	// GetTempBuffAttr 获取临时buff的实际属性
	GetTempBuffAttr(logger fklog.FKLogI, userID uint64, stageId int32) (map[int32]int64, error)

	// 进入关卡前检查关卡的buff情况，因为可能会有清除部分buff的情况
	CheckTempBuff(logger fklog.FKLogI, userId uint64, barrierId int32) (*tempbuffmodel.TempBuffInfoModel, error)

	// 获取临时buff信息
	GetTempBuffInfo(logger fklog.FKLogI, userID uint64, stageId int32) (*tempbuffmodel.TempBuffInfoModel, error)

	// 获取buff权重
	GetOptionBuffWeightInfo(logger fklog.FKLogI, buffId int32, selectedBuffMap, selectedBuffGroupMap map[int32]int32) *WeightInfo

	// 根据选择的buff获取全部buff属性
	GetTotalBuff(logger fklog.FKLogI, buffList []*tempbuffmodel.SelectedBuffInfo) (map[int32]int64, []*tempbuffmodel.TotalBuffInfo)
}

// GlobalTempBuffService 临时buff可用全局唯一对象
var GlobalTempBuffService TempBuffService

func init() {
	GlobalTempBuffService = newTempBuffService()
}

type service struct {
}

func newTempBuffService() TempBuffService {
	return &service{}
}

type BuffInfo struct {
	BuffId      int32  // buffId
	Value       int64  // 当前带的buff数量
	Name        string // buff名称
	Decs        string // 描述
	IsRecommend int32  // 是否推荐 0=否 1=是
}
type OptionalBuffInfo struct {
	SelectBuffList []*BuffInfo // 可选buff列表
	IsRefresh      int32       // 是否可以刷新 0=不可以 1=可以
	Cost           []*Item     // 刷新消耗
	SelectBuffTime int32       // 选择buff时间配置
}
type Item struct {
	ItemId    int32  //物品Id
	Count     int64  //物品数量
	Guid      int64  //guid
	ItemName  string //物品名称
	AtlasName string //图集名
	IconName  string // 图标名
	Quality   int32  // 品质
}
