package npclibrary

import (
	"context"
	"maze_game_server/model/npclibrarymodel"
	"maze_game_server/pb/common/MessageType"
	"maze_game_server/pb/common/NPCLibrary"
	"maze_game_server/services/npclibraryservice"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// 图鉴处理器
type NPCLibraryHandler struct {
	service npclibraryservice.NPCLibraryService
}

// 创建图鉴处理器实例
func NewNPCLibraryHandler() *NPCLibraryHandler {
	return &NPCLibraryHandler{
		service: npclibraryservice.GlobalNPCLibraryService,
	}
}

// 处理添加图鉴请求
func (h *NPCLibraryHandler) HandleAddNPCLibrary(ctx context.Context, req *NPCLibrary.AddNPCLibraryRQ) (*NPCLibrary.AddNPCLibraryRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if req.GetNpcRoleId() == 0 {
		logger.CtxError(ctx, "HandleAddNPCLibrary invalid npcRoleId", zap.Uint64("npcRoleId", req.GetNpcRoleId()))
		return &NPCLibrary.AddNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(400),
				ErrMsg:  []byte("invalid npc role id"),
			},
		}, nil
	}

	// 获取用户ID（从请求头或上下文）
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleAddNPCLibrary invalid userID")
		return &NPCLibrary.AddNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层添加图鉴
	library, err := h.service.AddNPCLibrary(ctx, userID, req.GetNpcRoleId(), req.GetNpcUserId(), req.GetCount(), npclibrarymodel.LibFromType(req.GetFromType()))
	if err != nil {
		logger.CtxError(ctx, "HandleAddNPCLibrary service failed", zap.Error(err))
		return &NPCLibrary.AddNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.AddNPCLibraryRS{
		Success: proto.Bool(true),
		Library: h.convertToPbLibrary(library),
	}

	logger.CtxInfo(ctx, "HandleAddNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleId", req.GetNpcRoleId()),
		zap.Int32("count", req.GetCount()))

	return response, nil
}

// 处理查询图鉴请求
func (h *NPCLibraryHandler) HandleQueryNPCLibrary(ctx context.Context, req *NPCLibrary.QueryNPCLibraryRQ) (*NPCLibrary.QueryNPCLibraryRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取用户ID
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleQueryNPCLibrary invalid userID")
		return &NPCLibrary.QueryNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层查询图鉴
	result, err := h.service.QueryNPCLibrary(ctx, userID, req.GetCursor(), req.GetLimit())
	if err != nil {
		logger.CtxError(ctx, "HandleQueryNPCLibrary service failed", zap.Error(err))
		return &NPCLibrary.QueryNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.QueryNPCLibraryRS{
		Libraries:  h.convertToPbLibraries(result.Libraries),
		NextCursor: proto.String(result.NextCursor),
		HasMore:    proto.Bool(result.HasMore),
		TotalCount: proto.Int32(result.TotalCount),
	}

	logger.CtxInfo(ctx, "HandleQueryNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.String("cursor", req.GetCursor()),
		zap.Int32("limit", req.GetLimit()),
		zap.Int("resultCount", len(result.Libraries)))

	return response, nil
}

// 处理升级图鉴请求
func (h *NPCLibraryHandler) HandleUpgradeNPCLibrary(ctx context.Context, req *NPCLibrary.UpgradeNPCLibraryRQ) (*NPCLibrary.UpgradeNPCLibraryRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if req.GetNpcRoleId() == 0 {
		logger.CtxError(ctx, "HandleUpgradeNPCLibrary invalid npcRoleId", zap.Uint64("npcRoleId", req.GetNpcRoleId()))
		return &NPCLibrary.UpgradeNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(400),
				ErrMsg:  []byte("invalid npc role id"),
			},
		}, nil
	}

	// 获取用户ID
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleUpgradeNPCLibrary invalid userID")
		return &NPCLibrary.UpgradeNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层升级图鉴
	result, err := h.service.UpgradeNPCLibrary(ctx, userID, req.GetNpcRoleId())
	if err != nil {
		logger.CtxError(ctx, "HandleUpgradeNPCLibrary service failed", zap.Error(err))
		return &NPCLibrary.UpgradeNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.UpgradeNPCLibraryRS{
		Success:  proto.Bool(result.Success),
		NewLevel: proto.Int32(result.NewLevel),
		Cost:     proto.Int32(result.Cost),
	}

	logger.CtxInfo(ctx, "HandleUpgradeNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleId", req.GetNpcRoleId()),
		zap.Bool("success", result.Success),
		zap.Int32("newLevel", result.NewLevel))

	return response, nil
}

// 处理检查图鉴请求
func (h *NPCLibraryHandler) HandleCheckNPCLibrary(ctx context.Context, req *NPCLibrary.CheckNPCLibraryRQ) (*NPCLibrary.CheckNPCLibraryRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if len(req.GetNpcRoleIds()) == 0 {
		logger.CtxError(ctx, "HandleCheckNPCLibrary empty npcRoleIds")
		return &NPCLibrary.CheckNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(400),
				ErrMsg:  []byte("empty npc role ids"),
			},
		}, nil
	}

	// 获取用户ID
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleCheckNPCLibrary invalid userID")
		return &NPCLibrary.CheckNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层检查图鉴
	result, err := h.service.CheckNPCLibrary(ctx, userID, req.GetNpcRoleIds())
	if err != nil {
		logger.CtxError(ctx, "HandleCheckNPCLibrary service failed", zap.Error(err))
		return &NPCLibrary.CheckNPCLibraryRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.CheckNPCLibraryRS{
		Results: result.Results,
	}

	logger.CtxInfo(ctx, "HandleCheckNPCLibrary success",
		zap.Uint64("userID", userID),
		zap.Int("requestCount", len(req.GetNpcRoleIds())),
		zap.Int("resultCount", len(result.Results)))

	return response, nil
}

// 处理获取下一级图鉴信息请求
func (h *NPCLibraryHandler) HandleGetNextLevelInfo(ctx context.Context, req *NPCLibrary.GetNextLevelInfoRQ) (*NPCLibrary.GetNextLevelInfoRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if req.GetNpcRoleId() == 0 {
		logger.CtxError(ctx, "HandleGetNextLevelInfo invalid npcRoleId", zap.Uint64("npcRoleId", req.GetNpcRoleId()))
		return &NPCLibrary.GetNextLevelInfoRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(400),
				ErrMsg:  []byte("invalid npc role id"),
			},
		}, nil
	}

	// 获取用户ID
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleGetNextLevelInfo invalid userID")
		return &NPCLibrary.GetNextLevelInfoRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层获取下一级信息
	nextInfo, err := h.service.GetNextLevelInfo(ctx, userID, req.GetNpcRoleId())
	if err != nil {
		logger.CtxError(ctx, "HandleGetNextLevelInfo service failed", zap.Error(err))
		return &NPCLibrary.GetNextLevelInfoRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.GetNextLevelInfoRS{
		NextInfo: h.convertToPbLibrary(nextInfo),
	}

	logger.CtxInfo(ctx, "HandleGetNextLevelInfo success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleId", req.GetNpcRoleId()),
		zap.Int32("nextLevel", nextInfo.NPCLevel))

	return response, nil
}

// 处理强制合成NPC请求
func (h *NPCLibraryHandler) HandleForceComposeNPC(ctx context.Context, req *NPCLibrary.ForceComposeNPCRQ) (*NPCLibrary.ForceComposeNPCRS, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 参数验证
	if req.GetNpcRoleId() == 0 {
		logger.CtxError(ctx, "HandleForceComposeNPC invalid npcRoleId", zap.Uint64("npcRoleId", req.GetNpcRoleId()))
		return &NPCLibrary.ForceComposeNPCRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(400),
				ErrMsg:  []byte("invalid npc role id"),
			},
		}, nil
	}

	// 获取用户ID
	userID := h.getUserIDFromContext(ctx)
	if userID == 0 {
		logger.CtxError(ctx, "HandleForceComposeNPC invalid userID")
		return &NPCLibrary.ForceComposeNPCRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(401),
				ErrMsg:  []byte("invalid user id"),
			},
		}, nil
	}

	// 调用服务层强制合成
	err := h.service.ForceComposeNPC(ctx, userID, req.GetNpcRoleId())
	if err != nil {
		logger.CtxError(ctx, "HandleForceComposeNPC service failed", zap.Error(err))
		return &NPCLibrary.ForceComposeNPCRS{
			ErrInfo: &MessageType.ErrorInfo{
				ErrCode: proto.Int64(500),
				ErrMsg:  []byte("internal server error"),
			},
		}, nil
	}

	// 构建响应
	response := &NPCLibrary.ForceComposeNPCRS{
		Success: proto.Bool(true),
	}

	logger.CtxInfo(ctx, "HandleForceComposeNPC success",
		zap.Uint64("userID", userID),
		zap.Uint64("npcRoleId", req.GetNpcRoleId()))

	return response, nil
}

// 从上下文获取用户ID
func (h *NPCLibraryHandler) getUserIDFromContext(ctx context.Context) uint64 {
	// TODO: 实现从上下文获取用户ID的逻辑
	// 这里需要根据实际的用户认证机制来实现
	return 0
}

// 转换图鉴信息为Protocol Buffer
func (h *NPCLibraryHandler) convertToPbLibrary(library *npclibrarymodel.NPCLibraryInfoEx) *NPCLibrary.NPCLibraryInfo {
	if library == nil {
		return nil
	}

	return &NPCLibrary.NPCLibraryInfo{
		Id:                 proto.Uint64(library.ID),
		HostUserId:         proto.Uint64(library.HostUserID),
		NpcRoleId:          proto.Uint64(library.NPCRoleID),
		NpcUserId:          proto.Uint64(library.NPCUserID),
		NpcCount:           proto.Int32(library.NPCCount),
		NpcLevel:           proto.Int32(library.NPCLevel),
		CreateTime:         proto.Uint64(library.CreateTime),
		NextLevelNeedCount: proto.Uint64(library.NextLevelNeedCount),
		NextLevelCost:      proto.Int32(library.NextLevelCost),
		NextLevelCostType:  proto.Int32(library.NextLevelCostType),
		ComposeFlag:        proto.Int32(library.ComposeFlag),
		NextCostStrength:   proto.Int32(library.NextCostStrength),
		NextRequiredLevel:  proto.Int32(library.NextRequiredLevel),
		Hp:                 proto.Int32(library.HP),
		FightVal:           proto.Int32(library.FightVal),
		Prestige:           proto.Int32(library.Prestige),
		ComposeNeedCnt:     proto.Int32(library.ComposeNeedCnt),
		Order:              proto.Int32(library.Order),
		Quality:            proto.Int32(library.Quality),
		GroupOrder:         proto.Int32(library.GroupOrder),
	}
}

// 批量转换图鉴信息为Protocol Buffer
func (h *NPCLibraryHandler) convertToPbLibraries(libraries []*npclibrarymodel.NPCLibraryInfoEx) []*NPCLibrary.NPCLibraryInfo {
	if len(libraries) == 0 {
		return nil
	}

	pbLibraries := make([]*NPCLibrary.NPCLibraryInfo, 0, len(libraries))
	for _, library := range libraries {
		if pbLibrary := h.convertToPbLibrary(library); pbLibrary != nil {
			pbLibraries = append(pbLibraries, pbLibrary)
		}
	}

	return pbLibraries
}
