package mailservice

import (
	"context"
	"encoding/json"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"
)

// MailBroadcastService 邮件广播服务
type MailBroadcastService struct {
	mailService MailService
}

// NewMailBroadcastService 创建邮件广播服务
func NewMailBroadcastService() *MailBroadcastService {
	return &MailBroadcastService{
		mailService: GlobalMailService,
	}
}

// BroadcastMailRequest 广播邮件请求
type BroadcastMailRequest struct {
	Title             string                    `json:"title"`                // 邮件标题
	Content           string                    `json:"content"`              // 邮件内容
	SenderName        string                    `json:"sender_name"`          // 发送者名称
	MailType          int32                     `json:"mail_type"`            // 邮件类型
	MailSubType       int32                     `json:"mail_sub_type"`        // 邮件子类型
	Attachments       []*mailmodel.Attachment   `json:"attachments"`          // 附件
	ExpireTime        int64                     `json:"expire_time"`          // 过期时间
	BroadcastFlag     int32                     `json:"broadcast_flag"`       // 广播标志
	BroadcastId       uint64                    `json:"broadcast_id"`         // 广播ID
	PanelInfo         *structsdef.MailPanelInfo `json:"panel_info"`           // 富文本面板信息
	UserLimitType     int32                     `json:"user_limit_type"`      // 用户限制类型
	UserLimitValueMin int32                     `json:"user_limit_value_min"` // 用户限制最小值
	UserLimitValueMax int32                     `json:"user_limit_value_max"` // 用户限制最大值
}

// BroadcastResult 广播结果
type BroadcastResult struct {
	SuccessCount int      `json:"success_count"` // 成功数量
	FailedCount  int      `json:"failed_count"`  // 失败数量
	FailedUsers  []uint64 `json:"failed_users"`  // 失败用户列表
}

// BroadcastMail 广播邮件
func (s *MailBroadcastService) BroadcastMail(ctx context.Context, req *BroadcastMailRequest) (*BroadcastResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 根据广播标志获取目标用户列表
	userIds, err := s.getBroadcastUsers(ctx, req.BroadcastFlag, req.BroadcastId, req.UserLimitType, req.UserLimitValueMin, req.UserLimitValueMax)
	if err != nil {
		logger.CtxError(ctx, "BroadcastMail getBroadcastUsers failed", zap.Error(err))
		return nil, err
	}

	if len(userIds) == 0 {
		logger.CtxWarn(ctx, "BroadcastMail no users found")
		return &BroadcastResult{}, nil
	}

	// 构建富文本面板信息
	panelInfoJSON := ""
	if req.PanelInfo != nil {
		panelInfoBytes, err := json.Marshal(req.PanelInfo)
		if err != nil {
			logger.CtxError(ctx, "BroadcastMail marshal panel info failed", zap.Error(err))
		} else {
			panelInfoJSON = string(panelInfoBytes)
		}
	}

	// 批量发送邮件
	result := &BroadcastResult{
		FailedUsers: make([]uint64, 0),
	}

	for _, userId := range userIds {
		// 创建邮件信息
		mailInfo := &mailmodel.MailInfo{
			ID:                uint64(ctx.Value("mail_id").(int64)), // 从上下文获取邮件ID
			Title:             req.Title,
			Content:           req.Content,
			Label:             int32(constdef.MailLabelSystem),
			MailType:          req.MailType,
			MailSubType:       req.MailSubType,
			MailContentType:   int32(constdef.CONTENT_TEXT_AND_AWARD),
			Sender:            req.SenderName,
			ReciverID:         userId,
			IsRead:            false,
			IsGetAttach:       false,
			SendTime:          ctx.Value("send_time").(int64),
			ExpireTime:        req.ExpireTime,
			Attachments:       req.Attachments,
			BroadcastFlag:     req.BroadcastFlag,
			BroadcastId:       req.BroadcastId,
			PanelInfo:         panelInfoJSON,
			UserLimitType:     req.UserLimitType,
			UserLimitValueMin: req.UserLimitValueMin,
			UserLimitValueMax: req.UserLimitValueMax,
			PagePopUp:         true, // 广播邮件需要弹框
			AllSupportPop:     true,
		}

		// 发送邮件
		err := s.mailService.SendMail(ctx, req.Title, req.Content, req.SenderName,
			int32(constdef.MailLabelSystem), userId, req.Attachments, req.ExpireTime)

		if err != nil {
			logger.CtxError(ctx, "BroadcastMail SendMail failed",
				zap.Uint64("userId", userId),
				zap.Error(err))
			result.FailedCount++
			result.FailedUsers = append(result.FailedUsers, userId)
		} else {
			result.SuccessCount++
		}
	}

	logger.CtxInfo(ctx, "BroadcastMail completed",
		zap.Int("successCount", result.SuccessCount),
		zap.Int("failedCount", result.FailedCount),
		zap.Int32("broadcastFlag", req.BroadcastFlag),
		zap.Uint64("broadcastId", req.BroadcastId))

	return result, nil
}

// getBroadcastUsers 根据广播标志获取目标用户列表
func (s *MailBroadcastService) getBroadcastUsers(ctx context.Context, broadcastFlag int32, broadcastId uint64, userLimitType int32, userLimitValueMin int32, userLimitValueMax int32) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)

	switch broadcastFlag {
	case int32(constdef.MAIL_NOT_BROADCAST):
		return []uint64{}, nil

	case int32(constdef.MAIL_FAMILY_BROADCAST):
		return s.getFamilyUsers(ctx, broadcastId, userLimitType, userLimitValueMin, userLimitValueMax)

	case int32(constdef.MAIL_LEAGUE_BROADCAST):
		return s.getLeagueUsers(ctx, broadcastId, userLimitType, userLimitValueMin, userLimitValueMax)

	case int32(constdef.MAIL_MAP_BROADCAST):
		return s.getMapUsers(ctx, broadcastId, userLimitType, userLimitValueMin, userLimitValueMax)

	default:
		logger.CtxWarn(ctx, "getBroadcastUsers unknown broadcast flag", zap.Int32("broadcastFlag", broadcastFlag))
		return []uint64{}, nil
	}
}

// getFamilyUsers 获取家族用户列表
func (s *MailBroadcastService) getFamilyUsers(ctx context.Context, familyId uint64, userLimitType int32, userLimitValueMin int32, userLimitValueMax int32) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 这里需要调用家族服务获取家族成员列表
	// 暂时返回空列表，实际实现时需要调用相应的服务
	logger.CtxInfo(ctx, "getFamilyUsers called", zap.Uint64("familyId", familyId))

	// TODO: 实现获取家族用户列表的逻辑
	// 可以通过家族服务或数据库查询获取家族成员

	return []uint64{}, nil
}

// getLeagueUsers 获取联盟用户列表
func (s *MailBroadcastService) getLeagueUsers(ctx context.Context, leagueId uint64, userLimitType int32, userLimitValueMin int32, userLimitValueMax int32) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 这里需要调用联盟服务获取联盟成员列表
	// 暂时返回空列表，实际实现时需要调用相应的服务
	logger.CtxInfo(ctx, "getLeagueUsers called", zap.Uint64("leagueId", leagueId))

	// TODO: 实现获取联盟用户列表的逻辑
	// 可以通过联盟服务或数据库查询获取联盟成员

	return []uint64{}, nil
}

// getMapUsers 获取地图用户列表
func (s *MailBroadcastService) getMapUsers(ctx context.Context, mapId uint64, userLimitType int32, userLimitValueMin int32, userLimitValueMax int32) ([]uint64, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 这里需要调用地图服务获取地图上的用户列表
	// 暂时返回空列表，实际实现时需要调用相应的服务
	logger.CtxInfo(ctx, "getMapUsers called", zap.Uint64("mapId", mapId))

	// TODO: 实现获取地图用户列表的逻辑
	// 可以通过地图服务或在线用户服务获取地图上的用户

	return []uint64{}, nil
}

// BroadcastSystemMail 广播系统邮件
func (s *MailBroadcastService) BroadcastSystemMail(ctx context.Context, title, content, senderName string, broadcastFlag int32, broadcastId uint64) (*BroadcastResult, error) {
	req := &BroadcastMailRequest{
		Title:             title,
		Content:           content,
		SenderName:        senderName,
		MailType:          int32(constdef.SYSTEM),
		MailSubType:       int32(constdef.AWARD),
		Attachments:       nil,
		ExpireTime:        0, // 永不过期
		BroadcastFlag:     broadcastFlag,
		BroadcastId:       broadcastId,
		PanelInfo:         nil,
		UserLimitType:     0,
		UserLimitValueMin: 0,
		UserLimitValueMax: 0,
	}

	return s.BroadcastMail(ctx, req)
}

// BroadcastRewardMail 广播奖励邮件
func (s *MailBroadcastService) BroadcastRewardMail(ctx context.Context, title, content, senderName string, attachments []*mailmodel.Attachment, broadcastFlag int32, broadcastId uint64) (*BroadcastResult, error) {
	req := &BroadcastMailRequest{
		Title:             title,
		Content:           content,
		SenderName:        senderName,
		MailType:          int32(constdef.SYSTEM),
		MailSubType:       int32(constdef.AWARD),
		Attachments:       attachments,
		ExpireTime:        0, // 永不过期
		BroadcastFlag:     broadcastFlag,
		BroadcastId:       broadcastId,
		PanelInfo:         nil,
		UserLimitType:     0,
		UserLimitValueMin: 0,
		UserLimitValueMax: 0,
	}

	return s.BroadcastMail(ctx, req)
}

// 全局邮件广播服务实例
var GlobalMailBroadcastService *MailBroadcastService

func init() {
	GlobalMailBroadcastService = NewMailBroadcastService()
}
