package familyservice

import (
	"maze_game_server/model/familymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type FamilyService interface {
	// 创建家族
	CreateFamily(logger fklog.FKLogI, userID uint64, allianceID int32, familyName string,
		familySetting int32, userInfo familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 创建家族扣除物品
	DeductCreateFamilyCost(logger fklog.FKLogI, cost map[int32]int64) error
	// 查询创建家族消耗
	QueryCreateFamilyCost(logger fklog.FKLogI) (map[int32]int64, error)
	// 检查创建家族消耗
	CheckCreateFamilyCost(logger fklog.FKLogI) (bool, error)
	// 升级家族
	UpgradeFamily(logger fklog.FKLogI, familyID int32, targetLevel int32) (*familymodel.FamilyInfoModel, error)
	// 升级家族扣除物品
	DeductUpgradeFamilyCost(logger fklog.FKLogI, cost map[int32]int64) error
	// 查询升级家族消耗
	QueryUpgradeFamilyCost(logger fklog.FKLogI, familyID, targetLevel int32) (map[int32]int64, error)
	// 检查升级家族消耗
	CheckUpgradeFamilyCost(logger fklog.FKLogI, targetLevel int32) (bool, error)
	// 发送升级家族id包
	SendUpgradeFamilyIDPack(logger fklog.FKLogI, familyID int32) error
	// 获取家族详情
	GetFamilyInfo(logger fklog.FKLogI, familyID int32) (*familymodel.FamilyInfoModel, error)
	// 获取家族列表
	GetFamilyList(logger fklog.FKLogI) (familymodel.FamilysInfoModel, error)
	// 申请加入家族
	ApplyFamily(logger fklog.FKLogI, familyID int32, applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 同意加入家族
	AgreeApplyFamily(logger fklog.FKLogI, familyID int32, userID uint64,
		applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 拒绝加入家族
	RefuseApplyFamily(logger fklog.FKLogI, familyID int32, userID uint64,
		applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 踢出家族
	KickFamily(logger fklog.FKLogI, userID uint64, familyID int32,
		kickUsers []uint64) error
	// 发送踢出家族id包
	SendKickFamilyID(logger fklog.FKLogI, userID uint64, familyID int32,
		kickUsers []uint64) error
	// 退出家族
	ExitFamily(logger fklog.FKLogI, userID uint64, familyID int32) error
	// 发送退出家族id包
	SendExitFamilyID(logger fklog.FKLogI, userID uint64, familyID int32) error
	// 检查家族申请是否合理
	CheckApplyFamily(logger fklog.FKLogI, familyID int32, applyUser *familymodel.FamilyMember) (bool, error)
	// 设置用户绑定家族
	SetUserFamily(logger fklog.FKLogI, userID uint64, familyID int32) error
	// 获取用户绑定家族
	GetUserFamily(logger fklog.FKLogI, userID uint64) (int32, error)
	// 删除用户绑定家族
	DeleteUserFamily(logger fklog.FKLogI, userID uint64) error
	// 设置用户最后离开家族时间
	SetUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) error
	// 获取用户最后离开家族时间
	GetUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) (int64, error)
	// 检查用户最后离开家族时间是否超过一天
	CheckUserLastLeaveFamilyTime(logger fklog.FKLogI, userID uint64) (bool, error)
	// 检查用户是否有家族
	CheckUserHaveFamily(logger fklog.FKLogI, userID uint64) (bool, error)
	// 检查用户是否在家族申请列表
	CheckUserInApplyList(logger fklog.FKLogI, familyID int32, userID uint64) (bool, error)
	// 检查用户是否在家族中
	CheckUserInFamily(logger fklog.FKLogI, familyID int32, userID uint64) (bool, error)
	// 修改家族某玩家权限
	OperatePrivilege(logger fklog.FKLogI, userID uint64, familyID int32,
		privilegeLevel int32, operateUsers []uint64) error
	// 解散家族
	DissolutionFamily(logger fklog.FKLogI, userID uint64, familyID int32) error
}

var GlobalFamilyService FamilyService

func init() {
	GlobalFamilyService = newFamilyService()
}

type service struct {
}

func newFamilyService() FamilyService {
	return &service{}
}
