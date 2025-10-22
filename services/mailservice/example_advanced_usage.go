package mailservice

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"
)

// ExampleAdvancedUsage 展示高级邮件功能的使用
func ExampleAdvancedUsage() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	// 1. 群发邮件示例
	exampleBatchSendMail(ctx, logger)

	// 2. 删除邮件限制示例
	exampleDeleteRestrictions(ctx, logger)

	// 3. 附件领取示例
	exampleAttachmentClaim(ctx, logger)

	// 4. 背包集成示例
	exampleBagIntegration(ctx, logger)
}

// exampleBatchSendMail 群发邮件示例
func exampleBatchSendMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 群发邮件示例 ===")

	// 准备接收者列表
	receiverIds := []uint64{40000001, 40000002, 40000003, 40000004, 40000005}

	// 1. 群发系统邮件（同步）
	result, err := GlobalMailBatchService.BatchSendSystemMail(ctx,
		"系统维护通知",
		"系统将于今晚进行维护，请提前做好准备",
		"系统管理员",
		receiverIds,
		false) // 同步发送

	if err != nil {
		logger.ErrorWF("群发系统邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("群发系统邮件成功",
			zap.Int("totalCount", result.TotalCount),
			zap.Int("successCount", result.SuccessCount),
			zap.Int("failedCount", result.FailedCount),
			zap.Int64("processTime", result.ProcessTime))
	}

	// 2. 群发奖励邮件（异步队列）
	attachments := []*mailmodel.Attachment{
		{
			ItemID: 46700001,
			Count:  10,
			Extra:  "equip",
		},
		{
			ItemID: 46200001,
			Count:  5,
			Extra:  "item",
		},
	}

	result, err = GlobalMailBatchService.BatchSendRewardMail(ctx,
		"节日奖励",
		"恭喜您获得节日奖励，请查收！",
		"系统",
		attachments,
		receiverIds,
		true) // 异步队列发送

	if err != nil {
		logger.ErrorWF("群发奖励邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("群发奖励邮件成功",
			zap.Int("totalCount", result.TotalCount),
			zap.Int("successCount", result.SuccessCount),
			zap.Int("failedCount", result.FailedCount),
			zap.Bool("useQueue", result.UseQueue))
	}

	// 3. 群发富文本邮件
	panelInfo := &structsdef.MailPanelInfo{
		Content: "富文本群发邮件",
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: "欢迎参与我们的活动！",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 活动系统",
			},
		},
	}

	result, err = GlobalMailBatchService.BatchSendRichTextMail(ctx,
		"活动邀请",
		"欢迎参与我们的活动！",
		"活动系统",
		panelInfo,
		receiverIds,
		false)

	if err != nil {
		logger.ErrorWF("群发富文本邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("群发富文本邮件成功",
			zap.Int("totalCount", result.TotalCount),
			zap.Int("successCount", result.SuccessCount),
			zap.Int("failedCount", result.FailedCount))
	}
}

// exampleDeleteRestrictions 删除邮件限制示例
func exampleDeleteRestrictions(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 删除邮件限制示例 ===")

	userId := uint64(40000001)
	mailId := uint64(12345)
	label := int32(constdef.MailLabelSystem)

	// 1. 检查单个邮件的删除限制
	result, err := GlobalMailDeleteRestrictions.CheckDeleteRestrictions(ctx, userId, mailId, label)
	if err != nil {
		logger.ErrorWF("检查删除限制失败", zap.Error(err))
	} else {
		logger.InfoWF("删除限制检查结果",
			zap.Bool("canDelete", result.CanDelete),
			zap.String("reason", result.Reason),
			zap.Int("restrictionCount", len(result.Restrictions)))
	}

	// 2. 批量检查删除限制
	mailIds := []uint64{12345, 12346, 12347}
	results, err := GlobalMailDeleteRestrictions.BatchCheckDeleteRestrictions(ctx, userId, mailIds, label)
	if err != nil {
		logger.ErrorWF("批量检查删除限制失败", zap.Error(err))
	} else {
		logger.InfoWF("批量删除限制检查结果",
			zap.Int("totalCount", len(mailIds)),
			zap.Int("resultCount", len(results)))

		for mailId, result := range results {
			logger.InfoWF("邮件删除限制",
				zap.Uint64("mailId", mailId),
				zap.Bool("canDelete", result.CanDelete),
				zap.String("reason", result.Reason))
		}
	}

	// 3. 获取可删除的邮件列表
	deletableMails, err := GlobalMailDeleteRestrictions.GetDeletableMails(ctx, userId, label)
	if err != nil {
		logger.ErrorWF("获取可删除邮件列表失败", zap.Error(err))
	} else {
		logger.InfoWF("可删除邮件列表",
			zap.Int("deletableCount", len(deletableMails)))
	}

	// 4. 强制删除邮件（管理员权限）
	adminId := uint64(10000001)
	err = GlobalMailDeleteRestrictions.ForceDeleteMail(ctx, userId, mailId, label, adminId)
	if err != nil {
		logger.ErrorWF("强制删除邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("强制删除邮件成功")
	}
}

// exampleAttachmentClaim 附件领取示例
func exampleAttachmentClaim(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 附件领取示例 ===")

	userId := uint64(40000001)
	mailId := uint64(12345)
	label := int32(constdef.MailLabelSystem)

	// 1. 领取单个邮件的附件
	result, err := GlobalMailAttachmentClaim.ClaimMailAttachment(ctx, userId, mailId, label)
	if err != nil {
		logger.ErrorWF("领取邮件附件失败", zap.Error(err))
	} else {
		logger.InfoWF("领取邮件附件结果",
			zap.Bool("success", result.Success),
			zap.String("reason", result.Reason),
			zap.Int("claimedCount", len(result.ClaimedItems)),
			zap.Int("failedCount", len(result.FailedItems)),
			zap.Int("bagFullCount", len(result.BagFullItems)),
			zap.Int64("totalValue", result.TotalValue))
	}

	// 2. 领取所有邮件的附件
	result, err = GlobalMailAttachmentClaim.ClaimAllMailAttachments(ctx, userId, label)
	if err != nil {
		logger.ErrorWF("领取所有邮件附件失败", zap.Error(err))
	} else {
		logger.InfoWF("领取所有邮件附件结果",
			zap.Bool("success", result.Success),
			zap.String("reason", result.Reason),
			zap.Int("claimedCount", len(result.ClaimedItems)),
			zap.Int("failedCount", len(result.FailedItems)),
			zap.Int("bagFullCount", len(result.BagFullItems)),
			zap.Int64("totalValue", result.TotalValue))
	}

	// 3. 获取可领取的附件列表
	claimableMails, err := GlobalMailAttachmentClaim.GetClaimableAttachments(ctx, userId, label)
	if err != nil {
		logger.ErrorWF("获取可领取附件列表失败", zap.Error(err))
	} else {
		logger.InfoWF("可领取附件列表",
			zap.Int("claimableCount", len(claimableMails)))
	}
}

// exampleBagIntegration 背包集成示例
func exampleBagIntegration(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 背包集成示例 ===")

	userId := uint64(40000001)

	// 1. 获取背包空间信息
	currentSpace, remainingSpace, err := GlobalMailBagIntegration.GetBagSpaceInfo(ctx, userId)
	if err != nil {
		logger.ErrorWF("获取背包空间信息失败", zap.Error(err))
	} else {
		logger.InfoWF("背包空间信息",
			zap.Int32("currentSpace", currentSpace),
			zap.Int32("remainingSpace", remainingSpace))
	}

	// 2. 准备附件列表
	attachments := []*mailmodel.Attachment{
		{
			ItemID: 46700001,
			Count:  10,
			Extra:  "equip",
		},
		{
			ItemID: 46200001,
			Count:  5,
			Extra:  "item",
		},
		{
			ItemID: 46100001,
			Count:  100,
			Extra:  "currency",
		},
	}

	// 3. 检查背包空间是否足够
	hasEnoughSpace, remainingSpace, err := GlobalMailBagIntegration.CheckBagSpaceForAttachments(ctx, userId, attachments)
	if err != nil {
		logger.ErrorWF("检查背包空间失败", zap.Error(err))
	} else {
		logger.InfoWF("背包空间检查结果",
			zap.Bool("hasEnoughSpace", hasEnoughSpace),
			zap.Int32("remainingSpace", remainingSpace))
	}

	// 4. 添加附件到背包
	result, err := GlobalMailBagIntegration.AddAttachmentsToBag(ctx, userId, attachments)
	if err != nil {
		logger.ErrorWF("添加附件到背包失败", zap.Error(err))
	} else {
		logger.InfoWF("添加附件到背包结果",
			zap.Bool("success", result.Success),
			zap.String("reason", result.Reason),
			zap.Int("addedCount", len(result.AddedItems)),
			zap.Int("failedCount", len(result.FailedItems)),
			zap.Int("bagFullCount", len(result.BagFullItems)),
			zap.Int32("bagSpace", result.BagSpace),
			zap.Int64("totalValue", result.TotalValue))
	}

	// 5. 批量添加附件到背包
	mailAttachments := map[uint64][]*mailmodel.Attachment{
		12345: {
			{
				ItemID: 46700001,
				Count:  5,
				Extra:  "equip",
			},
		},
		12346: {
			{
				ItemID: 46200001,
				Count:  3,
				Extra:  "item",
			},
		},
	}

	result, err = GlobalMailBagIntegration.BatchAddAttachmentsToBag(ctx, userId, mailAttachments)
	if err != nil {
		logger.ErrorWF("批量添加附件到背包失败", zap.Error(err))
	} else {
		logger.InfoWF("批量添加附件到背包结果",
			zap.Bool("success", result.Success),
			zap.String("reason", result.Reason),
			zap.Int("addedCount", len(result.AddedItems)),
			zap.Int("failedCount", len(result.FailedItems)),
			zap.Int("bagFullCount", len(result.BagFullItems)),
			zap.Int32("bagSpace", result.BagSpace),
			zap.Int64("totalValue", result.TotalValue))
	}
}
