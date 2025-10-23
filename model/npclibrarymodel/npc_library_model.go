package npclibrarymodel

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 图鉴查询结果
type NPCLibraryQueryResult struct {
	Libraries  []*NPCLibraryInfoEx `json:"libraries"`
	NextCursor string              `json:"next_cursor"`
	HasMore    bool                `json:"has_more"`
	TotalCount int32               `json:"total_count"`
}

// 图鉴检查结果
type NPCLibraryCheckResult struct {
	Results map[uint64]bool `json:"results"`
}

// 图鉴升级结果
type NPCLibraryUpgradeResult struct {
	Success  bool  `json:"success"`
	NewLevel int32 `json:"new_level"`
	Cost     int32 `json:"cost"`
}

// 图鉴分页查询参数
type NPCLibraryQueryParams struct {
	UserID    uint64 `json:"user_id"`
	Cursor    string `json:"cursor"`
	Limit     int32  `json:"limit"`
	NPCRoleID uint64 `json:"npc_role_id,omitempty"`
	Quality   int32  `json:"quality,omitempty"`
	GroupID   int32  `json:"group_id,omitempty"`
}

// 图鉴添加参数
type NPCLibraryAddParams struct {
	UserID    uint64      `json:"user_id"`
	NPCRoleID uint64      `json:"npc_role_id"`
	NPCUserID uint64      `json:"npc_user_id"`
	Count     int32       `json:"count"`
	FromType  LibFromType `json:"from_type"`
}

// 图鉴升级参数
type NPCLibraryUpgradeParams struct {
	UserID    uint64 `json:"user_id"`
	NPCRoleID uint64 `json:"npc_role_id"`
}

// 图鉴检查参数
type NPCLibraryCheckParams struct {
	UserID     uint64   `json:"user_id"`
	NPCRoleIDs []uint64 `json:"npc_role_ids"`
}

// 图鉴模型接口
type NPCLibraryModel interface {
	// 添加图鉴
	AddNPCLibrary(ctx context.Context, params *NPCLibraryAddParams) (*NPCLibraryInfoEx, error)

	// 查询图鉴
	QueryNPCLibrary(ctx context.Context, params *NPCLibraryQueryParams) (*NPCLibraryQueryResult, error)

	// 升级图鉴
	UpgradeNPCLibrary(ctx context.Context, params *NPCLibraryUpgradeParams) (*NPCLibraryUpgradeResult, error)

	// 检查图鉴
	CheckNPCLibrary(ctx context.Context, params *NPCLibraryCheckParams) (*NPCLibraryCheckResult, error)

	// 获取下一级图鉴信息
	GetNextLevelInfo(ctx context.Context, userID uint64, npcRoleID uint64) (*NPCLibraryInfoEx, error)

	// 强制合成NPC
	ForceComposeNPC(ctx context.Context, userID uint64, npcRoleID uint64) error

	// 删除图鉴
	DeleteNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error

	// 批量获取图鉴
	BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*NPCLibraryInfoEx, error)
}

// 图鉴模型实现
type npcLibraryModel struct{}

// 创建图鉴模型实例
func NewNPCLibraryModel() NPCLibraryModel {
	return &npcLibraryModel{}
}

// 添加图鉴
func (nlm *npcLibraryModel) AddNPCLibrary(ctx context.Context, params *NPCLibraryAddParams) (*NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 检查参数
	if params.UserID == 0 || params.NPCRoleID == 0 {
		return nil, fmt.Errorf("invalid parameters")
	}

	// 检查是否已存在
	existingInfo, err := nlm.getNPCLibraryByRoleID(ctx, params.UserID, params.NPCRoleID)
	if err != nil {
		logger.CtxError(ctx, "AddNPCLibrary get existing failed", zap.Error(err))
		return nil, err
	}

	if existingInfo != nil {
		// 更新现有图鉴
		existingInfo.UpdateCount(params.Count)
		if err := existingInfo.Save(ctx, params.UserID); err != nil {
			logger.CtxError(ctx, "AddNPCLibrary save failed", zap.Error(err))
			return nil, err
		}
		return existingInfo, nil
	}

	// 创建新图鉴
	newInfo := NewNPCLibraryInfo(params.UserID, params.NPCRoleID, params.NPCUserID, params.Count, params.FromType)
	if err := newInfo.Save(ctx, params.UserID); err != nil {
		logger.CtxError(ctx, "AddNPCLibrary save new failed", zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "AddNPCLibrary success",
		zap.Uint64("userID", params.UserID),
		zap.Uint64("npcRoleID", params.NPCRoleID),
		zap.Int32("count", params.Count))

	return newInfo, nil
}

// 查询图鉴
func (nlm *npcLibraryModel) QueryNPCLibrary(ctx context.Context, params *NPCLibraryQueryParams) (*NPCLibraryQueryResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 实现分页查询逻辑
	// 这里需要根据实际需求实现分页查询
	result := &NPCLibraryQueryResult{
		Libraries:  make([]*NPCLibraryInfoEx, 0),
		NextCursor: "",
		HasMore:    false,
		TotalCount: 0,
	}

	logger.CtxInfo(ctx, "QueryNPCLibrary",
		zap.Uint64("userID", params.UserID),
		zap.String("cursor", params.Cursor),
		zap.Int32("limit", params.Limit))

	return result, nil
}

// 升级图鉴
func (nlm *npcLibraryModel) UpgradeNPCLibrary(ctx context.Context, params *NPCLibraryUpgradeParams) (*NPCLibraryUpgradeResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取图鉴信息
	info, err := nlm.getNPCLibraryByRoleID(ctx, params.UserID, params.NPCRoleID)
	if err != nil {
		logger.CtxError(ctx, "UpgradeNPCLibrary get failed", zap.Error(err))
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("npc library not found")
	}

	// 检查升级条件
	// 扣除升级消耗

	// 执行升级
	oldLevel := info.NPCLevel
	info.Upgrade()

	// 保存更新
	if err := info.Save(ctx, params.UserID); err != nil {
		logger.CtxError(ctx, "UpgradeNPCLibrary save failed", zap.Error(err))
		return nil, err
	}

	result := &NPCLibraryUpgradeResult{
		Success:  true,
		NewLevel: info.NPCLevel,
		Cost:     0, //  计算实际消耗
	}

	logger.CtxInfo(ctx, "UpgradeNPCLibrary success",
		zap.Uint64("userID", params.UserID),
		zap.Uint64("npcRoleID", params.NPCRoleID),
		zap.Int32("oldLevel", oldLevel),
		zap.Int32("newLevel", info.NPCLevel))

	return result, nil
}

// 检查图鉴
func (nlm *npcLibraryModel) CheckNPCLibrary(ctx context.Context, params *NPCLibraryCheckParams) (*NPCLibraryCheckResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	result := &NPCLibraryCheckResult{
		Results: make(map[uint64]bool),
	}

	// 批量检查图鉴是否存在
	for _, npcRoleID := range params.NPCRoleIDs {
		info, err := nlm.getNPCLibraryByRoleID(ctx, params.UserID, npcRoleID)
		if err != nil {
			logger.CtxError(ctx, "CheckNPCLibrary get failed",
				zap.Uint64("npcRoleID", npcRoleID), zap.Error(err))
			result.Results[npcRoleID] = false
			continue
		}
		result.Results[npcRoleID] = info != nil && !info.IsEmpty()
	}

	logger.CtxInfo(ctx, "CheckNPCLibrary",
		zap.Uint64("userID", params.UserID),
		zap.Int("count", len(params.NPCRoleIDs)))

	return result, nil
}

// 获取下一级图鉴信息
func (nlm *npcLibraryModel) GetNextLevelInfo(ctx context.Context, userID uint64, npcRoleID uint64) (*NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取当前图鉴信息
	info, err := nlm.getNPCLibraryByRoleID(ctx, userID, npcRoleID)
	if err != nil {
		logger.CtxError(ctx, "GetNextLevelInfo get failed", zap.Error(err))
		return nil, err
	}

	if info == nil {
		return nil, fmt.Errorf("npc library not found")
	}

	// 根据配置计算下一级信息
	nextInfo := &NPCLibraryInfoEx{}
	*nextInfo = *info
	nextInfo.NPCLevel++
	nextInfo.NextLevelNeedCount = 0 // 从配置获取
	nextInfo.NextLevelCost = 0      // 从配置获取
	nextInfo.NextLevelCostType = 0  //  从配置获取

	logger.CtxInfo(ctx, "GetNextLevelInfo",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID),
		zap.Int32("currentLevel", info.NPCLevel),
		zap.Int32("nextLevel", nextInfo.NPCLevel))

	return nextInfo, nil
}

// 强制合成NPC
func (nlm *npcLibraryModel) ForceComposeNPC(ctx context.Context, userID uint64, npcRoleID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	//  实现强制合成逻辑
	// 1. 检查合成条件
	// 2. 扣除合成消耗
	// 3. 执行合成

	logger.CtxInfo(ctx, "ForceComposeNPC",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 删除图鉴
func (nlm *npcLibraryModel) DeleteNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	// 实现删除逻辑
	// 1. 从Redis删除
	// 2. 从数据库删除

	logger.CtxInfo(ctx, "DeleteNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 批量获取图鉴
func (nlm *npcLibraryModel) BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	results := make([]*NPCLibraryInfoEx, 0, len(npcRoleIDs))

	for _, npcRoleID := range npcRoleIDs {
		info, err := nlm.getNPCLibraryByRoleID(ctx, userID, npcRoleID)
		if err != nil {
			logger.CtxError(ctx, "BatchGetNPCLibrary get failed",
				zap.Uint64("npcRoleID", npcRoleID), zap.Error(err))
			continue
		}
		if info != nil && !info.IsEmpty() {
			results = append(results, info)
		}
	}

	logger.CtxInfo(ctx, "BatchGetNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Int("requestCount", len(npcRoleIDs)),
		zap.Int("resultCount", len(results)))

	return results, nil
}

// 根据角色ID获取图鉴信息
func (nlm *npcLibraryModel) getNPCLibraryByRoleID(ctx context.Context, userID uint64, npcRoleID uint64) (*NPCLibraryInfoEx, error) {
	//  实现根据角色ID获取图鉴信息的逻辑
	// 这里需要从Redis或数据库查询
	return nil, nil
}

// 全局图鉴模型实例
var GlobalNPCLibraryModel NPCLibraryModel

func init() {
	GlobalNPCLibraryModel = NewNPCLibraryModel()
}
