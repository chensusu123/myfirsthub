package npclibraryservice

import (
	"context"
	"fmt"
	"maze_game_server/model/npclibrarymodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// 图鉴服务接口
type NPCLibraryService interface {
	// 添加图鉴
	AddNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64, npcUserID uint64, count int32, fromType npclibrarymodel.LibFromType) (*npclibrarymodel.NPCLibraryInfoEx, error)

	// 查询图鉴
	QueryNPCLibrary(ctx context.Context, userID uint64, cursor string, limit int32) (*npclibrarymodel.NPCLibraryQueryResult, error)

	// 升级图鉴
	UpgradeNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryUpgradeResult, error)

	// 检查图鉴
	CheckNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) (*npclibrarymodel.NPCLibraryCheckResult, error)

	// 获取下一级图鉴信息
	GetNextLevelInfo(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryInfoEx, error)

	// 强制合成NPC
	ForceComposeNPC(ctx context.Context, userID uint64, npcRoleID uint64) error

	// 删除图鉴
	DeleteNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error

	// 批量获取图鉴
	BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error)
}

// 图鉴服务实现
type npcLibraryService struct {
	model npclibrarymodel.NPCLibraryModel
}

// 创建图鉴服务实例
func NewNPCLibraryService() NPCLibraryService {
	return &npcLibraryService{
		model: npclibrarymodel.GlobalNPCLibraryModel,
	}
}

// 添加图鉴
func (nls *npcLibraryService) AddNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64, npcUserID uint64, count int32, fromType npclibrarymodel.LibFromType) (*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || npcRoleID == 0 {
		logger.CtxError(ctx, "AddNPCLibrary invalid parameters",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID))
		return nil, fmt.Errorf("invalid parameters")
	}

	// 构建添加参数
	params := &npclibrarymodel.NPCLibraryAddParams{
		UserID:    userID,
		NPCRoleID: npcRoleID,
		NPCUserID: npcUserID,
		Count:     count,
		FromType:  fromType,
	}

	// 调用模型层添加图鉴
	result, err := nls.model.AddNPCLibrary(ctx, params)
	if err != nil {
		logger.CtxError(ctx, "AddNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID),
			zap.Error(err))
		return nil, err
	}

	// 发送图鉴变化通知
	go nls.sendLibraryChangeNotify(ctx, userID, []*npclibrarymodel.NPCLibraryInfoEx{result}, npclibrarymodel.LibraryChangeTypeAdd)

	logger.CtxInfo(ctx, "AddNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID),
		zap.Int32("count", count),
		zap.Int32("fromType", int32(fromType)))

	return result, nil
}

// 查询图鉴
func (nls *npcLibraryService) QueryNPCLibrary(ctx context.Context, userID uint64, cursor string, limit int32) (*npclibrarymodel.NPCLibraryQueryResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 {
		logger.CtxError(ctx, "QueryNPCLibrary invalid userID", zap.Uint64("userID", userID))
		return nil, fmt.Errorf("invalid userID")
	}

	// 设置默认限制
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// 构建查询参数
	params := &npclibrarymodel.NPCLibraryQueryParams{
		UserID: userID,
		Cursor: cursor,
		Limit:  limit,
	}

	// 调用模型层查询图鉴
	result, err := nls.model.QueryNPCLibrary(ctx, params)
	if err != nil {
		logger.CtxError(ctx, "QueryNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "QueryNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.String("cursor", cursor),
		zap.Int32("limit", limit),
		zap.Int("resultCount", len(result.Libraries)))

	return result, nil
}

// 升级图鉴
func (nls *npcLibraryService) UpgradeNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryUpgradeResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || npcRoleID == 0 {
		logger.CtxError(ctx, "UpgradeNPCLibrary invalid parameters",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID))
		return nil, fmt.Errorf("invalid parameters")
	}

	// 构建升级参数
	params := &npclibrarymodel.NPCLibraryUpgradeParams{
		UserID:    userID,
		NPCRoleID: npcRoleID,
	}

	// 调用模型层升级图鉴
	result, err := nls.model.UpgradeNPCLibrary(ctx, params)
	if err != nil {
		logger.CtxError(ctx, "UpgradeNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID),
			zap.Error(err))
		return nil, err
	}

	// 发送图鉴变化通知
	if result.Success {
		go nls.sendLibraryChangeNotify(ctx, userID, nil, npclibrarymodel.LibraryChangeTypeUpdate)
	}

	logger.CtxInfo(ctx, "UpgradeNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID),
		zap.Bool("success", result.Success),
		zap.Int32("newLevel", result.NewLevel))

	return result, nil
}

// 检查图鉴
func (nls *npcLibraryService) CheckNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) (*npclibrarymodel.NPCLibraryCheckResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || len(npcRoleIDs) == 0 {
		logger.CtxError(ctx, "CheckNPCLibrary invalid parameters",
			zap.Uint64("userID", userID),
			zap.Int("npcRoleIDsCount", len(npcRoleIDs)))
		return nil, fmt.Errorf("invalid parameters")
	}

	// 构建检查参数
	params := &npclibrarymodel.NPCLibraryCheckParams{
		UserID:     userID,
		NPCRoleIDs: npcRoleIDs,
	}

	// 调用模型层检查图鉴
	result, err := nls.model.CheckNPCLibrary(ctx, params)
	if err != nil {
		logger.CtxError(ctx, "CheckNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "CheckNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Int("requestCount", len(npcRoleIDs)),
		zap.Int("resultCount", len(result.Results)))

	return result, nil
}

// 获取下一级图鉴信息
func (nls *npcLibraryService) GetNextLevelInfo(ctx context.Context, userID uint64, npcRoleID uint64) (*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || npcRoleID == 0 {
		logger.CtxError(ctx, "GetNextLevelInfo invalid parameters",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID))
		return nil, fmt.Errorf("invalid parameters")
	}

	// 调用模型层获取下一级信息
	result, err := nls.model.GetNextLevelInfo(ctx, userID, npcRoleID)
	if err != nil {
		logger.CtxError(ctx, "GetNextLevelInfo model failed",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID),
			zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "GetNextLevelInfo success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID),
		zap.Int32("nextLevel", result.NPCLevel))

	return result, nil
}

// 强制合成NPC
func (nls *npcLibraryService) ForceComposeNPC(ctx context.Context, userID uint64, npcRoleID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || npcRoleID == 0 {
		logger.CtxError(ctx, "ForceComposeNPC invalid parameters",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID))
		return fmt.Errorf("invalid parameters")
	}

	// 调用模型层强制合成
	err := nls.model.ForceComposeNPC(ctx, userID, npcRoleID)
	if err != nil {
		logger.CtxError(ctx, "ForceComposeNPC model failed",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID),
			zap.Error(err))
		return err
	}

	// 发送图鉴变化通知
	go nls.sendLibraryChangeNotify(ctx, userID, nil, npclibrarymodel.LibraryChangeTypeUpdate)

	logger.CtxInfo(ctx, "ForceComposeNPC success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 删除图鉴
func (nls *npcLibraryService) DeleteNPCLibrary(ctx context.Context, userID uint64, npcRoleID uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 || npcRoleID == 0 {
		logger.CtxError(ctx, "DeleteNPCLibrary invalid parameters",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID))
		return fmt.Errorf("invalid parameters")
	}

	// 调用模型层删除图鉴
	err := nls.model.DeleteNPCLibrary(ctx, userID, npcRoleID)
	if err != nil {
		logger.CtxError(ctx, "DeleteNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Uint64("npcRoleID", npcRoleID),
			zap.Error(err))
		return err
	}

	// 发送图鉴变化通知
	go nls.sendLibraryChangeNotify(ctx, userID, nil, npclibrarymodel.LibraryChangeTypeDelete)

	logger.CtxInfo(ctx, "DeleteNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleID", npcRoleID))

	return nil
}

// 批量获取图鉴
func (nls *npcLibraryService) BatchGetNPCLibrary(ctx context.Context, userID uint64, npcRoleIDs []uint64) ([]*npclibrarymodel.NPCLibraryInfoEx, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if userID == 0 {
		logger.CtxError(ctx, "BatchGetNPCLibrary invalid userID", zap.Uint64("userID", userID))
		return nil, fmt.Errorf("invalid userID")
	}

	// 调用模型层批量获取
	result, err := nls.model.BatchGetNPCLibrary(ctx, userID, npcRoleIDs)
	if err != nil {
		logger.CtxError(ctx, "BatchGetNPCLibrary model failed",
			zap.Uint64("userID", userID),
			zap.Error(err))
		return nil, err
	}

	logger.CtxInfo(ctx, "BatchGetNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Int("requestCount", len(npcRoleIDs)),
		zap.Int("resultCount", len(result)))

	return result, nil
}

// 发送图鉴变化通知
func (nls *npcLibraryService) sendLibraryChangeNotify(ctx context.Context, userID uint64, libraries []*npclibrarymodel.NPCLibraryInfoEx, changeType npclibrarymodel.LibraryChangeType) {
	logger := fklog.ContextAppLogger(ctx)

	// 构建通知消息
	notify := &npclibrarymodel.NPCLibraryNotify{
		UserID:     userID,
		Libraries:  libraries,
		ChangeType: changeType,
	}

	// 实现通知发送逻辑
	// 这里可以通过消息队列或直接推送给客户端
	_ = notify // 暂时忽略未使用警告

	logger.CtxInfo(ctx, "sendLibraryChangeNotify",
		zap.Uint64("userID", userID),
		zap.Int("librariesCount", len(libraries)),
		zap.Int32("changeType", int32(changeType)))
}

// 全局图鉴服务实例
var GlobalNPCLibraryService NPCLibraryService

func init() {
	GlobalNPCLibraryService = NewNPCLibraryService()
}
