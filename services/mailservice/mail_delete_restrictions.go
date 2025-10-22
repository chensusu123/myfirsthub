package mailservice

import (
	"context"
	"fmt"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/model/mailmodel"
)

// MailDeleteRestrictions 邮件删除限制服务
type MailDeleteRestrictions struct {
	mailService MailService
}

// NewMailDeleteRestrictions 创建邮件删除限制服务
func NewMailDeleteRestrictions() *MailDeleteRestrictions {
	return &MailDeleteRestrictions{
		mailService: GlobalMailService,
	}
}

// DeleteRestrictionType 删除限制类型
type DeleteRestrictionType int32

const (
	DELETE_RESTRICTION_NONE           DeleteRestrictionType = 0 // 无限制
	DELETE_RESTRICTION_UNREAD         DeleteRestrictionType = 1 // 未读邮件不能删除
	DELETE_RESTRICTION_HAS_ATTACHMENT DeleteRestrictionType = 2 // 有附件未领取不能删除
	DELETE_RESTRICTION_IMPORTANT      DeleteRestrictionType = 3 // 重要邮件不能删除
	DELETE_RESTRICTION_SYSTEM         DeleteRestrictionType = 4 // 系统邮件不能删除
	DELETE_RESTRICTION_EXPIRED        DeleteRestrictionType = 5 // 过期邮件不能删除
)

// DeleteRestrictionResult 删除限制检查结果
type DeleteRestrictionResult struct {
	CanDelete    bool                    `json:"can_delete"`   // 是否可以删除
	Restrictions []DeleteRestrictionType `json:"restrictions"` // 限制类型列表
	Reason       string                  `json:"reason"`       // 限制原因
	MailInfo     *mailmodel.MailInfo     `json:"mail_info"`    // 邮件信息
}

// CheckDeleteRestrictions 检查删除限制
func (s *MailDeleteRestrictions) CheckDeleteRestrictions(ctx context.Context, userId, mailId uint64, label int32) (*DeleteRestrictionResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取邮件信息
	mailInfo, err := s.mailService.ReadMail(logger, userId, mailId, label)
	if err != nil {
		logger.CtxError(ctx, "CheckDeleteRestrictions ReadMail failed",
			zap.Uint64("userId", userId),
			zap.Uint64("mailId", mailId),
			zap.Error(err))
		return nil, err
	}

	if mailInfo == nil {
		return &DeleteRestrictionResult{
			CanDelete:    false,
			Restrictions: []DeleteRestrictionType{DELETE_RESTRICTION_NONE},
			Reason:       "邮件不存在",
			MailInfo:     nil,
		}, nil
	}

	result := &DeleteRestrictionResult{
		CanDelete:    true,
		Restrictions: make([]DeleteRestrictionType, 0),
		MailInfo:     mailInfo,
	}

	// 检查各种删除限制
	s.checkUnreadRestriction(mailInfo, result)
	s.checkAttachmentRestriction(mailInfo, result)
	s.checkImportantRestriction(mailInfo, result)
	s.checkSystemRestriction(mailInfo, result)
	s.checkExpiredRestriction(mailInfo, result)

	// 如果有任何限制，则不能删除
	if len(result.Restrictions) > 0 {
		result.CanDelete = false
		result.Reason = s.buildRestrictionReason(result.Restrictions)
	}

	logger.CtxInfo(ctx, "CheckDeleteRestrictions completed",
		zap.Uint64("userId", userId),
		zap.Uint64("mailId", mailId),
		zap.Bool("canDelete", result.CanDelete),
		zap.Int("restrictionCount", len(result.Restrictions)))

	return result, nil
}

// checkUnreadRestriction 检查未读限制
func (s *MailDeleteRestrictions) checkUnreadRestriction(mailInfo *mailmodel.MailInfo, result *DeleteRestrictionResult) {
	if !mailInfo.IsRead {
		result.Restrictions = append(result.Restrictions, DELETE_RESTRICTION_UNREAD)
	}
}

// checkAttachmentRestriction 检查附件限制
func (s *MailDeleteRestrictions) checkAttachmentRestriction(mailInfo *mailmodel.MailInfo, result *DeleteRestrictionResult) {
	if len(mailInfo.Attachments) > 0 && !mailInfo.IsGetAttach {
		result.Restrictions = append(result.Restrictions, DELETE_RESTRICTION_HAS_ATTACHMENT)
	}
}

// checkImportantRestriction 检查重要邮件限制
func (s *MailDeleteRestrictions) checkImportantRestriction(mailInfo *mailmodel.MailInfo, result *DeleteRestrictionResult) {
	// 根据邮件类型判断是否为重要邮件
	importantTypes := []int32{
		int32(constdef.SYSTEM),
		int32(constdef.BATTLE),
		int32(constdef.TEAM),
	}

	for _, importantType := range importantTypes {
		if mailInfo.MailType == importantType {
			result.Restrictions = append(result.Restrictions, DELETE_RESTRICTION_IMPORTANT)
			break
		}
	}
}

// checkSystemRestriction 检查系统邮件限制
func (s *MailDeleteRestrictions) checkSystemRestriction(mailInfo *mailmodel.MailInfo, result *DeleteRestrictionResult) {
	if mailInfo.MailType == int32(constdef.SYSTEM) && mailInfo.Sender == "系统" {
		result.Restrictions = append(result.Restrictions, DELETE_RESTRICTION_SYSTEM)
	}
}

// checkExpiredRestriction 检查过期限制
func (s *MailDeleteRestrictions) checkExpiredRestriction(mailInfo *mailmodel.MailInfo, result *DeleteRestrictionResult) {
	// 如果邮件已过期，不能删除（需要先清理过期邮件）
	if mailInfo.ExpireTime > 0 && mailInfo.ExpireTime < time.Now().Unix() {
		result.Restrictions = append(result.Restrictions, DELETE_RESTRICTION_EXPIRED)
	}
}

// buildRestrictionReason 构建限制原因
func (s *MailDeleteRestrictions) buildRestrictionReason(restrictions []DeleteRestrictionType) string {
	reasons := make([]string, 0)

	for _, restriction := range restrictions {
		switch restriction {
		case DELETE_RESTRICTION_UNREAD:
			reasons = append(reasons, "未读邮件不能删除")
		case DELETE_RESTRICTION_HAS_ATTACHMENT:
			reasons = append(reasons, "有未领取附件不能删除")
		case DELETE_RESTRICTION_IMPORTANT:
			reasons = append(reasons, "重要邮件不能删除")
		case DELETE_RESTRICTION_SYSTEM:
			reasons = append(reasons, "系统邮件不能删除")
		case DELETE_RESTRICTION_EXPIRED:
			reasons = append(reasons, "过期邮件不能删除")
		}
	}

	return fmt.Sprintf("删除限制: %s", fmt.Sprintf("%v", reasons))
}

// ForceDeleteMail 强制删除邮件（管理员权限）
func (s *MailDeleteRestrictions) ForceDeleteMail(ctx context.Context, userId, mailId uint64, label int32, adminId uint64) error {
	logger := fklog.ContextAppLogger(ctx)

	// 记录强制删除操作
	logger.CtxWarn(ctx, "ForceDeleteMail executed",
		zap.Uint64("userId", userId),
		zap.Uint64("mailId", mailId),
		zap.Uint64("adminId", adminId))

	// 直接删除邮件
	err := s.mailService.DelMail(logger, userId, mailId, label)
	if err != nil {
		logger.CtxError(ctx, "ForceDeleteMail DelMail failed",
			zap.Uint64("userId", userId),
			zap.Uint64("mailId", mailId),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "ForceDeleteMail success",
		zap.Uint64("userId", userId),
		zap.Uint64("mailId", mailId),
		zap.Uint64("adminId", adminId))

	return nil
}

// BatchCheckDeleteRestrictions 批量检查删除限制
func (s *MailDeleteRestrictions) BatchCheckDeleteRestrictions(ctx context.Context, userId uint64, mailIds []uint64, label int32) (map[uint64]*DeleteRestrictionResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	results := make(map[uint64]*DeleteRestrictionResult)

	for _, mailId := range mailIds {
		result, err := s.CheckDeleteRestrictions(ctx, userId, mailId, label)
		if err != nil {
			logger.CtxError(ctx, "BatchCheckDeleteRestrictions CheckDeleteRestrictions failed",
				zap.Uint64("userId", userId),
				zap.Uint64("mailId", mailId),
				zap.Error(err))
			continue
		}

		results[mailId] = result
	}

	logger.CtxInfo(ctx, "BatchCheckDeleteRestrictions completed",
		zap.Uint64("userId", userId),
		zap.Int("totalCount", len(mailIds)),
		zap.Int("resultCount", len(results)))

	return results, nil
}

// GetDeletableMails 获取可删除的邮件列表
func (s *MailDeleteRestrictions) GetDeletableMails(ctx context.Context, userId uint64, label int32) ([]*mailmodel.MailInfo, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取用户的所有邮件
	mailList, err := s.mailService.GetMailListByLabel(logger, userId, label, 0, 1000)
	if err != nil {
		logger.CtxError(ctx, "GetDeletableMails GetMailListByLabel failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return nil, err
	}

	deletableMails := make([]*mailmodel.MailInfo, 0)

	for _, mailInfo := range mailList {
		// 检查删除限制
		result, err := s.CheckDeleteRestrictions(ctx, userId, mailInfo.ID, label)
		if err != nil {
			logger.CtxError(ctx, "GetDeletableMails CheckDeleteRestrictions failed",
				zap.Uint64("userId", userId),
				zap.Uint64("mailId", mailInfo.ID),
				zap.Error(err))
			continue
		}

		if result.CanDelete {
			deletableMails = append(deletableMails, mailInfo)
		}
	}

	logger.CtxInfo(ctx, "GetDeletableMails completed",
		zap.Uint64("userId", userId),
		zap.Int("totalCount", len(mailList)),
		zap.Int("deletableCount", len(deletableMails)))

	return deletableMails, nil
}

// 全局邮件删除限制服务实例
var GlobalMailDeleteRestrictions *MailDeleteRestrictions

func init() {
	GlobalMailDeleteRestrictions = NewMailDeleteRestrictions()
}
