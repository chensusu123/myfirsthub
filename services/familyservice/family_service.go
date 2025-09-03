package familyservice

import (
	"context"
	"maze_game_server/model/familymodel"
)

type FamilyService interface {
	// 创建家族
	CreateFamily(ctx context.Context, userID uint64, allianceID int32, familyName string,
		familySetting int32, userInfo familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 创建家族扣除物品
	DeductCreateFamilyCost(ctx context.Context, cost map[int32]int64) error
	// 查询创建家族消耗
	QueryCreateFamilyCost(ctx context.Context) (map[int32]int64, error)
	// 检查创建家族消耗
	CheckCreateFamilyCost(ctx context.Context) (bool, error)
	// 升级家族
	UpgradeFamily(ctx context.Context, familyID int32, targetLevel int32) (*familymodel.FamilyInfoModel, error)
	// 升级家族扣除物品
	DeductUpgradeFamilyCost(ctx context.Context, cost map[int32]int64) error
	// 查询升级家族消耗
	QueryUpgradeFamilyCost(ctx context.Context, familyID, targetLevel int32) (map[int32]int64, error)
	// 检查升级家族消耗
	CheckUpgradeFamilyCost(ctx context.Context, targetLevel int32) (bool, error)
	// 发送升级家族id包
	SendUpgradeFamilyIDPack(ctx context.Context, familyID int32) error
	// 获取家族详情
	GetFamilyInfo(ctx context.Context, familyID int32) (*familymodel.FamilyInfoModel, error)
	// 获取家族列表
	GetFamilyList(ctx context.Context) (familymodel.FamilysInfoModel, error)
	// 申请加入家族
	ApplyFamily(ctx context.Context, familyID int32, applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 同意加入家族
	AgreeApplyFamily(ctx context.Context, familyID int32, userID uint64,
		applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 拒绝加入家族
	RefuseApplyFamily(ctx context.Context, familyID int32, userID uint64,
		applyUser familymodel.FamilyMember) (*familymodel.FamilyInfoModel, error)
	// 踢出家族
	KickFamily(ctx context.Context, userID uint64, familyID int32,
		kickUsers []uint64) error
	// 发送踢出家族id包
	SendKickFamilyID(ctx context.Context, userID uint64, familyID int32,
		kickUsers []uint64) error
	// 退出家族
	ExitFamily(ctx context.Context, userID uint64, familyID int32) error
	// 发送退出家族id包
	SendExitFamilyID(ctx context.Context, userID uint64, familyID int32) error
	// 检查家族申请是否合理
	CheckApplyFamily(ctx context.Context, familyID int32, applyUser *familymodel.FamilyMember) (bool, error)
	// 设置用户绑定家族
	SetUserFamily(ctx context.Context, userID uint64, familyID int32) error
	// 获取用户绑定家族
	GetUserFamily(ctx context.Context, userID uint64) (int32, error)
	// 删除用户绑定家族
	DeleteUserFamily(ctx context.Context, userID uint64) error
	// 设置用户最后离开家族时间
	SetUserLastLeaveFamilyTime(ctx context.Context, userID uint64) error
	// 获取用户最后离开家族时间
	GetUserLastLeaveFamilyTime(ctx context.Context, userID uint64) (int64, error)
	// 检查用户最后离开家族时间是否超过一天
	CheckUserLastLeaveFamilyTime(ctx context.Context, userID uint64) (bool, error)
	// 检查用户是否有家族
	CheckUserHaveFamily(ctx context.Context, userID uint64) (bool, error)
	// 检查用户是否在家族申请列表
	CheckUserInApplyList(ctx context.Context, familyID int32, userID uint64) (bool, error)
	// 检查用户是否在家族中
	CheckUserInFamily(ctx context.Context, familyID int32, userID uint64) (bool, error)
	// 修改家族某玩家权限
	OperatePrivilege(ctx context.Context, userID uint64, familyID int32,
		privilegeLevel int32, operateUsers []uint64) error
	// 解散家族
	DissolutionFamily(ctx context.Context, userID uint64, familyID int32) error
	//设置家族群组id
	SetFamilyGroupID(ctx context.Context, familyID int32, groupID int64) error
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
