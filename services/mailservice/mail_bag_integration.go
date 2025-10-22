package mailservice

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/model/mailmodel"
	// "maze_game_server/services/itemservice"
)

// MailBagIntegration 邮件背包集成服务
type MailBagIntegration struct {
	mailService MailService
	// itemService itemservice.ItemService
}

// NewMailBagIntegration 创建邮件背包集成服务
func NewMailBagIntegration() *MailBagIntegration {
	return &MailBagIntegration{
		mailService: GlobalMailService,
		// itemService: itemservice.GlobalItemService,
	}
}

// BagIntegrationResult 背包集成结果
type BagIntegrationResult struct {
	Success      bool                    `json:"success"`        // 是否成功
	AddedItems   []*mailmodel.Attachment `json:"added_items"`    // 成功添加的物品
	FailedItems  []*mailmodel.Attachment `json:"failed_items"`   // 添加失败的物品
	BagFullItems []*mailmodel.Attachment `json:"bag_full_items"` // 背包已满的物品
	BagSpace     int32                   `json:"bag_space"`      // 背包剩余空间
	TotalValue   int64                   `json:"total_value"`    // 总价值
	Reason       string                  `json:"reason"`         // 失败原因
}

// AddAttachmentsToBag 添加附件到背包
func (s *MailBagIntegration) AddAttachmentsToBag(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) (*BagIntegrationResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	result := &BagIntegrationResult{
		Success:      true,
		AddedItems:   make([]*mailmodel.Attachment, 0),
		FailedItems:  make([]*mailmodel.Attachment, 0),
		BagFullItems: make([]*mailmodel.Attachment, 0),
		BagSpace:     0,
		TotalValue:   0,
		Reason:       "",
	}

	// 检查背包空间
	bagSpace, err := s.checkBagSpace(ctx, userId, attachments)
	if err != nil {
		logger.CtxError(ctx, "AddAttachmentsToBag checkBagSpace failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return nil, err
	}

	result.BagSpace = bagSpace

	// 如果背包空间不足，标记为背包已满
	if bagSpace < int32(len(attachments)) {
		result.BagFullItems = attachments
		result.Success = false
		result.Reason = "背包空间不足"
		return result, nil
	}

	// 逐个添加物品到背包
	for _, attachment := range attachments {
		success, reason := s.addSingleItemToBag(ctx, userId, attachment)
		if success {
			result.AddedItems = append(result.AddedItems, attachment)
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

	// 设置失败原因
	if len(result.FailedItems) > 0 || len(result.BagFullItems) > 0 {
		result.Success = false
		if len(result.BagFullItems) > 0 {
			result.Reason = "背包空间不足"
		} else if len(result.FailedItems) > 0 {
			result.Reason = "部分物品添加失败"
		}
	}

	logger.CtxInfo(ctx, "AddAttachmentsToBag completed",
		zap.Uint64("userId", userId),
		zap.Bool("success", result.Success),
		zap.Int("addedCount", len(result.AddedItems)),
		zap.Int("failedCount", len(result.FailedItems)),
		zap.Int("bagFullCount", len(result.BagFullItems)),
		zap.Int32("bagSpace", result.BagSpace),
		zap.Int64("totalValue", result.TotalValue))

	return result, nil
}

// checkBagSpace 检查背包空间
func (s *MailBagIntegration) checkBagSpace(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) (int32, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取背包当前空间
	currentSpace, err := s.itemService.GetBagSpace(logger, userId)
	if err != nil {
		logger.CtxError(ctx, "checkBagSpace GetBagSpace failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return 0, err
	}

	// 计算需要的空间
	neededSpace := int32(len(attachments))

	// 返回剩余空间
	remainingSpace := currentSpace - neededSpace
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	logger.CtxInfo(ctx, "checkBagSpace completed",
		zap.Uint64("userId", userId),
		zap.Int32("currentSpace", currentSpace),
		zap.Int32("neededSpace", neededSpace),
		zap.Int32("remainingSpace", remainingSpace))

	return remainingSpace, nil
}

// addSingleItemToBag 添加单个物品到背包
func (s *MailBagIntegration) addSingleItemToBag(ctx context.Context, userId uint64, attachment *mailmodel.Attachment) (bool, string) {
	logger := fklog.ContextAppLogger(ctx)

	// 验证物品信息
	if attachment.ItemID <= 0 || attachment.Count <= 0 {
		logger.CtxWarn(ctx, "addSingleItemToBag invalid attachment",
			zap.Uint64("userId", userId),
			zap.Int32("itemId", attachment.ItemID),
			zap.Int64("count", attachment.Count))
		return false, "invalid_item"
	}

	// 调用物品服务添加物品
	err := s.itemService.AddItem(logger, userId, attachment.ItemID, attachment.Count, attachment.Extra)
	if err != nil {
		logger.CtxError(ctx, "addSingleItemToBag AddItem failed",
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

	logger.CtxInfo(ctx, "addSingleItemToBag success",
		zap.Uint64("userId", userId),
		zap.Int32("itemId", attachment.ItemID),
		zap.Int64("count", attachment.Count))

	return true, ""
}

// calculateItemValue 计算物品价值
func (s *MailBagIntegration) calculateItemValue(attachment *mailmodel.Attachment) int64 {
	// 根据物品类型计算价值
	switch attachment.ItemID {
	case 46700001: // 装备碎片
		return attachment.Count * 100
	case 46200001: // 普通物品
		return attachment.Count * 10
	case 46100001: // 货币
		return attachment.Count * 1
	default:
		return attachment.Count * 5
	}
}

// GetBagSpaceInfo 获取背包空间信息
func (s *MailBagIntegration) GetBagSpaceInfo(ctx context.Context, userId uint64) (int32, int32, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取背包当前空间
	currentSpace, err := s.itemService.GetBagSpace(logger, userId)
	if err != nil {
		logger.CtxError(ctx, "GetBagSpaceInfo GetBagSpace failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return 0, 0, err
	}

	// 获取背包最大空间
	maxSpace, err := s.itemService.GetMaxBagSpace(logger, userId)
	if err != nil {
		logger.CtxError(ctx, "GetBagSpaceInfo GetMaxBagSpace failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return 0, 0, err
	}

	remainingSpace := maxSpace - currentSpace
	if remainingSpace < 0 {
		remainingSpace = 0
	}

	logger.CtxInfo(ctx, "GetBagSpaceInfo completed",
		zap.Uint64("userId", userId),
		zap.Int32("currentSpace", currentSpace),
		zap.Int32("maxSpace", maxSpace),
		zap.Int32("remainingSpace", remainingSpace))

	return currentSpace, remainingSpace, nil
}

// CheckBagSpaceForAttachments 检查背包空间是否足够容纳附件
func (s *MailBagIntegration) CheckBagSpaceForAttachments(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) (bool, int32, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取背包空间信息
	currentSpace, remainingSpace, err := s.GetBagSpaceInfo(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "CheckBagSpaceForAttachments GetBagSpaceInfo failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return false, 0, err
	}

	// 计算需要的空间
	neededSpace := int32(len(attachments))

	// 检查空间是否足够
	hasEnoughSpace := remainingSpace >= neededSpace

	logger.CtxInfo(ctx, "CheckBagSpaceForAttachments completed",
		zap.Uint64("userId", userId),
		zap.Int32("currentSpace", currentSpace),
		zap.Int32("remainingSpace", remainingSpace),
		zap.Int32("neededSpace", neededSpace),
		zap.Bool("hasEnoughSpace", hasEnoughSpace))

	return hasEnoughSpace, remainingSpace, nil
}

// BatchAddAttachmentsToBag 批量添加附件到背包
func (s *MailBagIntegration) BatchAddAttachmentsToBag(ctx context.Context, userId uint64, mailAttachments map[uint64][]*mailmodel.Attachment) (*BagIntegrationResult, error) {
	logger := fklog.ContextAppLogger(ctx)

	result := &BagIntegrationResult{
		Success:      true,
		AddedItems:   make([]*mailmodel.Attachment, 0),
		FailedItems:  make([]*mailmodel.Attachment, 0),
		BagFullItems: make([]*mailmodel.Attachment, 0),
		BagSpace:     0,
		TotalValue:   0,
		Reason:       "",
	}

	// 收集所有附件
	allAttachments := make([]*mailmodel.Attachment, 0)
	for _, attachments := range mailAttachments {
		allAttachments = append(allAttachments, attachments...)
	}

	// 检查背包空间
	hasEnoughSpace, remainingSpace, err := s.CheckBagSpaceForAttachments(ctx, userId, allAttachments)
	if err != nil {
		logger.CtxError(ctx, "BatchAddAttachmentsToBag CheckBagSpaceForAttachments failed",
			zap.Uint64("userId", userId),
			zap.Error(err))
		return nil, err
	}

	result.BagSpace = remainingSpace

	if !hasEnoughSpace {
		result.BagFullItems = allAttachments
		result.Success = false
		result.Reason = "背包空间不足"
		return result, nil
	}

	// 逐个添加物品
	for _, attachment := range allAttachments {
		success, reason := s.addSingleItemToBag(ctx, userId, attachment)
		if success {
			result.AddedItems = append(result.AddedItems, attachment)
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

	// 设置失败原因
	if len(result.FailedItems) > 0 || len(result.BagFullItems) > 0 {
		result.Success = false
		if len(result.BagFullItems) > 0 {
			result.Reason = "背包空间不足"
		} else if len(result.FailedItems) > 0 {
			result.Reason = "部分物品添加失败"
		}
	}

	logger.CtxInfo(ctx, "BatchAddAttachmentsToBag completed",
		zap.Uint64("userId", userId),
		zap.Bool("success", result.Success),
		zap.Int("addedCount", len(result.AddedItems)),
		zap.Int("failedCount", len(result.FailedItems)),
		zap.Int("bagFullCount", len(result.BagFullItems)),
		zap.Int32("bagSpace", result.BagSpace),
		zap.Int64("totalValue", result.TotalValue))

	return result, nil
}

// 全局邮件背包集成服务实例
var GlobalMailBagIntegration *MailBagIntegration

func init() {
	GlobalMailBagIntegration = NewMailBagIntegration()
}
