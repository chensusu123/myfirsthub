package npclibrarymodel

import (
	"context"
	"fmt"
	"maze_game_server/io"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 图鉴来源类型
type LibFromType int32

const (
	LibFromTypeCliRefresh          LibFromType = 1    // 战斗抓取
	LibFromTypeNpcValetDrop        LibFromType = 2    // 跟班打工掉落
	LibFromTypeTreasureBoxDrop     LibFromType = 3    // 普通抽卡掉落
	LibFromTypeLoot                LibFromType = 4    // 顺
	LibFromTypeFreeTreasureBoxDrop LibFromType = 5    // 免费开宝箱掉落
	LibFromTypeNewUserGuideAward   LibFromType = 6    // 新手引导掉落
	LibFromTypeVipTreasureBoxDrop  LibFromType = 7    // vip抽卡
	LibFromTypeClientBuyNpc        LibFromType = 8    // 客户端购买
	LibFromTypeVipReceiveAwards    LibFromType = 100  // vip领取奖励
	LibFromTypeNewRegisterAwards   LibFromType = 9    // 新注册赠送星卡
	LibFromTypeMiniAppsNewRegister LibFromType = 10   // h5版小程序用户到客户端登录注册迁移星卡
	LibFromTypeMallBuy             LibFromType = 11   // 商城购买
	LibFromTypePhpActive           LibFromType = 12   // PHP活动
	LibFromTypeLootSlot            LibFromType = 13   // 领取槽位奖励
	LibFromTypeRobSlot             LibFromType = 14   // 打劫槽位奖励
	LibFromTypeExploreAward        LibFromType = 15   // 游戏探索奖励
	LibFromTypeExchangeNpc         LibFromType = 270  // 提现页兑换星卡
	LibFromTypeBackground          LibFromType = 2000 // 后台操作
	LibFromTypeUpProgress          LibFromType = 2001 // 使用碎片
	LibFromTypeUpStar              LibFromType = 2002 // 升星
	LibFromTypeFullStarClean       LibFromType = 2003 // 满星后清除碎片
	LibFromTypeDayTask             LibFromType = 3000 // 日常任务
	LibFromTypeSectionTask         LibFromType = 3001 // 章节任务
	LibFromTypeAchievement         LibFromType = 3002 // 成就
	LibFromTypeResourceBattle      LibFromType = 3003 // 资源站奖励星卡碎片
	LibFromTypeDiamondShop         LibFromType = 3004 // 钻石商店发送
	LibFromTypeWebShop             LibFromType = 3005 // web商店-礼包
	LibFromTypeRedpacketShop       LibFromType = 3006 // 红包商城
	LibFromTypeSendSysUserID       LibFromType = 3007 // 后台造账号
	LibFromTypePeiyuAI             LibFromType = 3008 // 培育ai账号
	LibFromTypeCritAddValue        LibFromType = 9000 // 暴击升星时,暴击出的额外碎片
)

// NPC图鉴基础信息
type NPCLibraryInfo struct {
	ID         uint64 `json:"id"`           // 数据库索引
	HostUserID uint64 `json:"host_user_id"` // 主人用户ID
	NPCRoleID  uint64 `json:"npc_role_id"`  // NPC角色ID
	NPCUserID  uint64 `json:"npc_user_id"`  // NPC用户ID
	NPCCount   int32  `json:"npc_count"`    // NPC个数
	NPCLevel   int32  `json:"npc_level"`    // 图鉴级别
	CreateTime uint64 `json:"create_time"`  // 创建时间(初次解锁时间)
}

// 扩展图鉴信息
type NPCLibraryInfoEx struct {
	NPCLibraryInfo
	NextLevelNeedCount uint64 `json:"next_level_need_count"` // 下一级需要数量
	NextLevelCost      int32  `json:"next_level_cost"`       // 下一级消耗
	NextLevelCostType  int32  `json:"next_level_cost_type"`  // 消耗类型
	ComposeFlag        int32  `json:"compose_flag"`          // 合成标记
	NextCostStrength   int32  `json:"next_cost_strength"`    // 下一级消耗强度
	NextRequiredLevel  int32  `json:"next_required_level"`   // 下一级要求等级
	HP                 int32  `json:"hp"`                    // 血量
	FightVal           int32  `json:"fight_val"`             // 战斗力
	Prestige           int32  `json:"prestige"`              // 声望
	ComposeNeedCnt     int32  `json:"compose_need_cnt"`      // 合成需要数量
	Order              int32  `json:"order"`                 // 排序
	Quality            int32  `json:"quality"`               // 品质
	GroupOrder         int32  `json:"group_order"`           // 组排序
}

// 图鉴变化类型
type LibraryChangeType int32

const (
	LibraryChangeTypeAdd    LibraryChangeType = 1 // 新增
	LibraryChangeTypeUpdate LibraryChangeType = 2 // 更新
	LibraryChangeTypeDelete LibraryChangeType = 3 // 删除
)

// 图鉴变化通知
type NPCLibraryNotify struct {
	UserID     uint64              `json:"user_id"`
	Libraries  []*NPCLibraryInfoEx `json:"libraries"`
	ChangeType LibraryChangeType   `json:"change_type"`
}

// 获取Redis Key
func getNPCLibraryKey(userID uint64) string {
	return fmt.Sprintf("npc_library:%d", userID)
}

// 获取分片表名
func getNPCLibraryTableName(userID uint64) string {
	return fmt.Sprintf("npc_library_%d", (userID>>8)%64)
}

// 加载图鉴信息
func LoadNPCLibraryInfo(ctx context.Context, userID uint64) (info *NPCLibraryInfoEx, err error) {
	logger := fklog.ContextAppLogger(ctx)
	info = &NPCLibraryInfoEx{}
	if err = info.load(ctx, userID); err != nil {
		logger.CtxError(ctx, "LoadNPCLibraryInfo err",
			zap.Uint64("userID", userID), zap.Error(err))
		return nil, err
	}
	return
}

// 加载批量图鉴信息
func LoadBatchNPCLibraryInfo(ctx context.Context, userIDs []uint64) (infos []*NPCLibraryInfoEx, err error) {
	logger := fklog.ContextAppLogger(ctx)
	infos = make([]*NPCLibraryInfoEx, len(userIDs))
	if len(userIDs) > 0 {
		//  实现批量加载逻辑
		for i, userID := range userIDs {
			info, err := LoadNPCLibraryInfo(ctx, userID)
			if err != nil {
				logger.CtxError(ctx, "LoadBatchNPCLibraryInfo err",
					zap.Uint64("userID", userID), zap.Error(err))
				continue
			}
			infos[i] = info
		}
	}
	return
}

// 加载图鉴信息
func (nli *NPCLibraryInfoEx) load(ctx context.Context, userID uint64) (err error) {
	return io.LoadSvrData(ctx, getNPCLibraryKey(userID), nli)
}

// 保存图鉴信息
func (nli *NPCLibraryInfoEx) Save(ctx context.Context, userID uint64) (err error) {
	return io.SaveSvrData(ctx, getNPCLibraryKey(userID), nli)
}

// 删除图鉴信息
func (nli *NPCLibraryInfoEx) Delete(ctx context.Context, userID uint64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	// 实现Redis删除逻辑
	logger.CtxInfo(ctx, "DeleteNPCLibraryInfo",
		zap.Uint64("userID", userID))
	return
}

// 转换为JSON字符串
func (nli *NPCLibraryInfoEx) ToJSON() (string, error) {
	// TODO: 实现JSON序列化
	return "", nil
}

// 从JSON字符串填充数据
func (nli *NPCLibraryInfoEx) FromJSON(jsonStr string) error {
	// TODO: 实现JSON反序列化
	return nil
}

// 创建新的图鉴信息
func NewNPCLibraryInfo(userID uint64, npcRoleID uint64, npcUserID uint64, count int32, fromType LibFromType) *NPCLibraryInfoEx {
	return &NPCLibraryInfoEx{
		NPCLibraryInfo: NPCLibraryInfo{
			HostUserID: userID,
			NPCRoleID:  npcRoleID,
			NPCUserID:  npcUserID,
			NPCCount:   count,
			NPCLevel:   1,
			CreateTime: uint64(time.Now().Unix()),
		},
		NextLevelNeedCount: 0,
		NextLevelCost:      0,
		NextLevelCostType:  0,
		ComposeFlag:        0,
		NextCostStrength:   0,
		NextRequiredLevel:  0,
		HP:                 0,
		FightVal:           0,
		Prestige:           0,
		ComposeNeedCnt:     0,
		Order:              0,
		Quality:            0,
		GroupOrder:         0,
	}
}

// 检查图鉴是否存在
func (nli *NPCLibraryInfoEx) IsEmpty() bool {
	return nli.ID == 0
}

// 更新图鉴数量
func (nli *NPCLibraryInfoEx) UpdateCount(count int32) {
	nli.NPCCount += count
}

// 升级图鉴
func (nli *NPCLibraryInfoEx) Upgrade() {
	nli.NPCLevel++
	// 根据配置更新其他属性
}
