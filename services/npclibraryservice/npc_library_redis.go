package npclibraryservice

import (
	"context"
	"encoding/json"
	"fmt"
	"maze_game_server/model/npclibrarymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 图鉴Redis操作接口
type NPCLibraryRedis interface {
	// 获取图鉴信息
	GetNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryInfoEx, error)

	// 设置图鉴信息
	SetNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64, info *npclibrarymodel.NPCLibraryInfoEx) error

	// 删除图鉴信息
	DelNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error

	// 批量获取图鉴信息
	BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error)

	// 获取用户所有图鉴
	GetUserAllNPCLibrary(ctx context.Context, userID uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error)

	// 设置用户图鉴列表
	SetUserNPCLibraryList(ctx context.Context, userID uint64, libraries []*npclibrarymodel.NPCLibraryInfoEx) error
}

// 图鉴Redis操作实现
type npcLibraryRedis struct {
	// TODO: 添加Redis客户端
}

// 创建图鉴Redis操作实例
func NewNPCLibraryRedis() NPCLibraryRedis {
	return &npcLibraryRedis{}
}

// 获取Redis Key
func (nlr *npcLibraryRedis) getNPCLibraryKey(userID uint64, npcRoleID uint64) string {
	return fmt.Sprintf("npc_library:%d:%d", userID, npcRoleID)
}

// 获取用户图鉴列表Key
func (nlr *npcLibraryRedis) getUserNPCLibraryListKey(userID uint64) string {
	return fmt.Sprintf("npc_library_list:%d", userID)
}

// 获取图鉴信息
func (nlr *npcLibraryRedis) GetNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	//  实现Redis获取逻辑
	// 1. 从Redis获取数据
	// 2. 反序列化为NPCLibraryInfoEx
	// 3. 返回结果

	logger.CtxInfo(ctx, "GetNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil, nil
}

// 设置图鉴信息
func (nlr *npcLibraryRedis) SetNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64, info *npclibrarymodel.NPCLibraryInfoEx) error {
	logger := fklog.ContextAppLogger(ctx)

	//  实现Redis设置逻辑
	// 1. 序列化NPCLibraryInfoEx
	// 2. 存储到Redis
	// 3. 设置过期时间

	logger.CtxInfo(ctx, "SetNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 删除图鉴信息
func (nlr *npcLibraryRedis) DelNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	// 实现Redis删除逻辑
	// 1. 删除单个图鉴
	// 2. 从用户图鉴列表中移除

	logger.CtxInfo(ctx, "DelNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 批量获取图鉴信息
func (nlr *npcLibraryRedis) BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	//  实现批量获取逻辑
	// 1. 构建批量查询Key
	// 2. 批量从Redis获取
	// 3. 反序列化结果

	results := make([]*npclibrarymodel.NPCLibraryInfoEx, 0, len(npcRoleIDs))

	for _, npcRoleID := range npcRoleIDs {
		info, err := nlr.GetNPCLibrary(ctx, userID, npcRoleID)
		if err != nil {
			logger.CtxError(ctx, "BatchGetNPCLibrary get failed",
				zap.Uint64("npcRoleID", npcRoleID),
				zap.Error(err))
			continue
		}
		if info != nil {
			results = append(results, info)
		}
	}

	logger.CtxInfo(ctx, "BatchGetNPCLibrary",
		zap.Uint64("userID", userID),
		zap.Int("requestCount", len(npcRoleIDs)),
		zap.Int("resultCount", len(results)))

	return results, nil
}

// 获取用户所有图鉴
func (nlr *npcLibraryRedis) GetUserAllNPCLibrary(ctx context.Context, userID uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 实现获取用户所有图鉴逻辑
	// 1. 从用户图鉴列表Key获取所有图鉴ID
	// 2. 批量获取图鉴详情
	// 3. 返回结果

	logger.CtxInfo(ctx, "GetUserAllNPCLibrary",
		zap.Uint64("userID", userID))

	return nil, nil
}

// 设置用户图鉴列表
func (nlr *npcLibraryRedis) SetUserNPCLibraryList(ctx context.Context, userID uint64, libraries []*npclibrarymodel.NPCLibraryInfoEx) error {
	logger := fklog.ContextAppLogger(ctx)

	// 实现设置用户图鉴列表逻辑
	// 1. 序列化图鉴列表
	// 2. 存储到Redis
	// 3. 设置过期时间

	logger.CtxInfo(ctx, "SetUserNPCLibraryList",
		zap.Uint64("userID", userID),
		zap.Int("librariesCount", len(libraries)))

	return nil
}

// 序列化图鉴信息
func (nlr *npcLibraryRedis) serializeNPCLibrary(info *npclibrarymodel.NPCLibraryInfoEx) (string, error) {
	data, err := json.Marshal(info)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 反序列化图鉴信息
func (nlr *npcLibraryRedis) deserializeNPCLibrary(data string) (*npclibrarymodel.NPCLibraryInfoEx, error) {
	var info npclibrarymodel.NPCLibraryInfoEx
	err := json.Unmarshal([]byte(data), &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// 全局图鉴Redis操作实例
var GlobalNPCLibraryRedis NPCLibraryRedis

func init() {
	GlobalNPCLibraryRedis = NewNPCLibraryRedis()
}


