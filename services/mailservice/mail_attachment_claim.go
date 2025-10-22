package mailservice

import (
	"context"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/model/mailmodel"
	// "maze_game_server/services/itemservice"
)

// MailAttachmentClaim 邮件附件领取服务
type MailAttachmentClaim struct {
	mailService MailService
	// itemService itemservice.ItemService
}

// NewMailAttachmentClaim 创建邮件附件领取服务
func NewMailAttachmentClaim() *MailAttachmentClaim {
	return &MailAttachmentClaim{
		mailService: GlobalMailService,
		// itemService: itemservice.GlobalItemService,
	}
}

// AttachmentClaimResult 附件领取结果
type AttachmentClaimResult struct {
	Success      bool                    `json:"success"`        // 是否成功
	ClaimedItems []*mailmodel.Attachment `json:"claimed_items"`  // 已领取的物品
	FailedItems  []*mailmodel.Attachment `json:"failed_items"`   // 领取失败的物品
	BagFullItems []*mailmodel.Attachment `json:"bag_full_items"` // 背包已满的物品
	TotalValue   int64                   `json:"total_value"`    // 总价值
	Reason       string                  `json:"reason"`         // 失败原因
}

// ClaimMailAttachment 领取邮件附件
func (s *MailAttachmentClaim) ClaimMailAttachment(ctx context.Context, userId, mailId uint64, label int32) (*AttachmentClaimResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取邮件信息
	mailInfo, err := s.mailService.ReadMail(logger, userId, mailId, label)
	if err != nil {
		logger.CtxError(ctx, "ClaimMailAttachment ReadMail failed",
			zap.Uint64("userId", userId),
			zap.Uint64("mailId", mailId),
			zap.Error(err))
		return nil, err
	}

	if mailInfo == nil {
		return &AttachmentClaimResult{
			Success:      false,
			ClaimedItems: make([]*mailmodel.Attachment, 0),
			FailedItems:  make([]*mailmodel.Attachment, 0),
			BagFullItems: make([]*mailmodel.Attachment, 0),
			TotalValue:   0,
			Reason:       "邮件不存在",
		}, nil
	}

	// 检查邮件状态
	if mailInfo.IsGetAttach {
		return &AttachmentClaimResult{
			Success:      false,
			ClaimedItems: make([]*mailmodel.Attachment, 0),
			FailedItems:  make([]*mailmodel.Attachment, 0),
			BagFullItems: make([]*mailmodel.Attachment, 0),
			TotalValue:   0,
			Reason:       "附件已领取",
		}, nil
	}

	// 检查邮件是否过期
	if mailInfo.ExpireTime > 0 && mailInfo.ExpireTime < time.Now().Unix() {
		return &AttachmentClaimResult{
			Success:      false,
			ClaimedItems: make([]*mailmodel.Attachment, 0),
			FailedItems:  make([]*mailmodel.Attachment, 0),
			BagFullItems: make([]*mailmodel.Attachment, 0),
			TotalValue:   0,
			Reason:       "邮件已过期",
		}, nil
	}

	// 检查是否有附件
	if len(mailInfo.Attachments) == 0 {
		return &AttachmentClaimResult{
			Success:      false,
			ClaimedItems: make([]*mailmodel.Attachment, 0),
			FailedItems:  make([]*mailmodel.Attachment, 0),
			BagFullItems: make([]*mailmodel.Attachment, 0),
			TotalValue:   0,
			Reason:       "邮件无附件",
		}, nil
	}

	// 领取附件
	result := s.claimAttachments(ctx, userId, mailInfo.Attachments)

	// 如果领取成功，更新邮件状态
	if result.Success {
		// 标记邮件为已领取附件
		mailInfo.IsGetAttach = true
		mailInfo.IsRead = true // 领取附件时自动标记为已读

		// 保存邮件状态
		err = s.saveMailStatus(ctx, mailInfo)
		if err != nil {
			logger.CtxError(ctx, "ClaimMailAttachment saveMailStatus failed",
				zap.Uint64("userId", userId),
				zap.Uint64("mailId", mailId),
				zap.Error(err))
			// 即使保存失败，附件已经领取成功
		}
	}

	logger.CtxInfo(ctx, "ClaimMailAttachment completed",
		zap.Uint64("userId", userId),
		zap.Uint64("mailId", mailId),
		zap.Bool("success", result.Success),
		zap.Int("claimedCount", len(result.ClaimedItems)),
		zap.Int("failedCount", len(result.FailedItems)),
		zap.Int64("totalValue", result.TotalValue))

	return result, nil
}

// claimAttachments 领取附件
func (s *MailAttachmentClaim) claimAttachments(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) *AttachmentClaimResult {
	logger := fklog.ContextAppLogger(ctx)

	result := &AttachmentClaimResult{
		Success:      true,
		ClaimedItems: make([]*mailmodel.Attachment, 0),
		FailedItems:  make([]*mailmodel.Attachment, 0),
		BagFullItems: make([]*mailmodel.Attachment, 0),
		TotalValue:   0,
		Reason:       "",
	}

	for _, attachment := range attachments {
		// 检查物品是否有效
		if attachment.ItemID <= 0 || attachment.Count <= 0 {
			logger.CtxWarn(ctx, "ClaimMailAttachment invalid attachment",
				zap.Uint64("userId", userId),
				zap.Int32("itemId", attachment.ItemID),
				zap.Int64("count", attachment.Count))
			result.FailedItems = append(result.FailedItems, attachment)
			continue
		}

		// 尝试添加到背包
		success, reason := s.addItemToBag(ctx, userId, attachment)
		if success {
			result.ClaimedItems = append(result.ClaimedItems, attachment)
			result.TotalValue += s.calculateItemValue(attachment)
		} else {
			if reason == "bag_full" {
				result.BagFullItems = append(result.BagFullItems, attachment)
			} else {
				result.FailedItems = append(result.FailedItems, attachment)
			}
			result.Success = false
		}
	}

	// 如果有任何失败，设置失败原因
	if len(result.FailedItems) > 0 || len(result.BagFullItems) > 0 {
		result.Success = false
		if len(result.BagFullItems) > 0 {
			result.Reason = "背包空间不足"
		} else if len(result.FailedItems) > 0 {
			result.Reason = "部分物品领取失败"
		}
	}

	return result
}

// addItemToBag 添加物品到背包
func (s *MailAttachmentClaim) addItemToBag(ctx context.Context, userId uint64, attachment *mailmodel.Attachment) (bool, string) {
	logger := fklog.ContextAppLogger(ctx)

	// 调用物品服务添加物品到背包
	err := s.itemService.AddItem(logger, userId, attachment.ItemID, attachment.Count, attachment.Extra)
	if err != nil {
		logger.CtxError(ctx, "addItemToBag AddItem failed",
			zap.Uint64("userId", userId),
			zap.Int32("itemId", attachment.ItemID),
			zap.Int64("count", attachment.Count),
			zap.Error(err))

		// 根据错误类型判断失败原因
		if err.Error() == "bag_full" {
			return false, "bag_full"
		}
		return false, "add_failed"
	}

	logger.CtxInfo(ctx, "addItemToBag success",
		zap.Uint64("userId", userId),
		zap.Int32("itemId", attachment.ItemID),
		zap.Int64("count", attachment.Count))

	return true, ""
}

// calculateItemValue 计算物品价值
func (s *MailAttachmentClaim) calculateItemValue(attachment *mailmodel.Attachment) int64 {
	// 这里可以根据物品ID查询物品配置来计算价值
	// 暂时使用简单的计算方式
	return attachment.Count * int64(attachment.ItemID%1000)
}

// saveMailStatus 保存邮件状态
func (s *MailAttachmentClaim) saveMailStatus(ctx context.Context, mailInfo *mailmodel.MailInfo) error {
	logger := fklog.ContextAppLogger(ctx)

	// 这里需要调用邮件服务保存邮件状态
	// 由于现有接口限制，这里暂时记录日志
	logger.CtxInfo(ctx, "saveMailStatus",
		zap.Uint64("mailId", mailInfo.ID),
		zap.Bool("isGetAttach", mailInfo.IsGetAttach),
		zap.Bool("isRead", mailInfo.IsRead))

	// TODO: 实现保存邮件状态的逻辑
	// 可以通过邮件服务的内部方法或者直接操作数据库

	return nil
}

// ClaimAllMailAttachments 领取所有邮件附件
func (s *MailAttachmentClaim) ClaimAllMailAttachments(ctx context.Context, userId uint64, label int32) (*AttachmentClaimResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取用户的所有邮件
	mailList, err := s.mailService.GetMailListByLabel(logger, userId, label, 0, 1000)
	if err != nil {
		logger.CtxError(ctx, "ClaimAllMailAttachments GetMailListByLabel failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return nil, err
	}

	result := &AttachmentClaimResult{
		Success:      true,
		ClaimedItems: make([]*mailmodel.Attachment, 0),
		FailedItems:  make([]*mailmodel.Attachment, 0),
		BagFullItems: make([]*mailmodel.Attachment, 0),
		TotalValue:   0,
		Reason:       "",
	}

	// 遍历所有邮件，领取附件
	for _, mailInfo := range mailList {
		// 跳过已领取附件的邮件
		if mailInfo.IsGetAttach {
			continue
		}

		// 跳过过期邮件
		if mailInfo.ExpireTime > 0 && mailInfo.ExpireTime < time.Now().Unix() {
			continue
		}

		// 跳过无附件的邮件
		if len(mailInfo.Attachments) == 0 {
			continue
		}

		// 领取单个邮件的附件
		singleResult := s.claimAttachments(ctx, userId, mailInfo.Attachments)

		// 合并结果
		result.ClaimedItems = append(result.ClaimedItems, singleResult.ClaimedItems...)
		result.FailedItems = append(result.FailedItems, singleResult.FailedItems...)
		result.BagFullItems = append(result.BagFullItems, singleResult.BagFullItems...)
		result.TotalValue += singleResult.TotalValue

		// 如果单个邮件领取成功，更新邮件状态
		if singleResult.Success {
			mailInfo.IsGetAttach = true
			mailInfo.IsRead = true
			s.saveMailStatus(ctx, mailInfo)
		}
	}

	// 检查整体结果
	if len(result.FailedItems) > 0 || len(result.BagFullItems) > 0 {
		result.Success = false
		if len(result.BagFullItems) > 0 {
			result.Reason = "背包空间不足"
		} else if len(result.FailedItems) > 0 {
			result.Reason = "部分物品领取失败"
		}
	}

	logger.CtxInfo(ctx, "ClaimAllMailAttachments completed",
		zap.Uint64("userId", userId),
		zap.Bool("success", result.Success),
		zap.Int("claimedCount", len(result.ClaimedItems)),
		zap.Int("failedCount", len(result.FailedItems)),
		zap.Int("bagFullCount", len(result.BagFullItems)),
		zap.Int64("totalValue", result.TotalValue))

	return result, nil
}

// GetClaimableAttachments 获取可领取的附件列表
func (s *MailAttachmentClaim) GetClaimableAttachments(ctx context.Context, userId uint64, label int32) ([]*mailmodel.MailInfo, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取用户的所有邮件
	mailList, err := s.mailService.GetMailListByLabel(logger, userId, label, 0, 1000)
	if err != nil {
		logger.CtxError(ctx, "GetClaimableAttachments GetMailListByLabel failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return nil, err
	}

	claimableMails := make([]*mailmodel.MailInfo, 0)

	for _, mailInfo := range mailList {
		// 检查是否可领取
		if s.isClaimable(mailInfo) {
			claimableMails = append(claimableMails, mailInfo)
		}
	}

	logger.CtxInfo(ctx, "GetClaimableAttachments completed",
		zap.Uint64("userId", userId),
		zap.Int("totalCount", len(mailList)),
		zap.Int("claimableCount", len(claimableMails)))

	return claimableMails, nil
}

// isClaimable 检查邮件是否可领取
func (s *MailAttachmentClaim) isClaimable(mailInfo *mailmodel.MailInfo) bool {
	// 已领取附件
	if mailInfo.IsGetAttach {
		return false
	}

	// 无附件
	if len(mailInfo.Attachments) == 0 {
		return false
	}

	// 已过期
	if mailInfo.ExpireTime > 0 && mailInfo.ExpireTime < time.Now().Unix() {
		return false
	}

	return true
}

// 全局邮件附件领取服务实例
var GlobalMailAttachmentClaim *MailAttachmentClaim

func init() {
	GlobalMailAttachmentClaim = NewMailAttachmentClaim()
}
