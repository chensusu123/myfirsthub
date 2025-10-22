package mailservice

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"
	"maze_game_server/pb/common/MazeCommon"
)

// ExampleUsage 展示如何使用新的邮件服务功能
func ExampleUsage() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	// 1. 发送基础邮件
	exampleBasicMail(ctx, logger)

	// 2. 发送富文本邮件
	exampleRichTextMail(ctx, logger)

	// 3. 发送投票邮件
	exampleVoteMail(ctx, logger)

	// 4. 发送战报邮件
	exampleBattleReportMail(ctx, logger)

	// 5. 发送礼盒邮件
	exampleGiftBoxMail(ctx, logger)

	// 6. 发送跳转邮件
	exampleLinkMail(ctx, logger)

	// 7. 广播邮件
	exampleBroadcastMail(ctx, logger)

	// 8. 使用配置发送邮件
	exampleConfigMail(ctx, logger)

	// 9. 发送系统监控邮件
	exampleSystemMail(ctx, logger)
}

// exampleBasicMail 基础邮件示例
func exampleBasicMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 基础邮件示例 ===")

	// 创建附件
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

	// 发送邮件
	err := GlobalMailService.SendMail(logger,
		"测试邮件",
		"这是一封测试邮件",
		"系统",
		int32(constdef.MailLabelSystem),
		40000001,
		attachments,
		time.Now().Add(24*time.Hour).Unix())

	if err != nil {
		logger.ErrorWF("发送基础邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("基础邮件发送成功")
	}
}

// exampleRichTextMail 富文本邮件示例
func exampleRichTextMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 富文本邮件示例 ===")

	// 构建富文本面板信息
	panelInfo := &structsdef.MailPanelInfo{
		Content: "富文本邮件内容",
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的玩家:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: "欢迎来到迷宫世界！",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统",
			},
		},
		User: []*structsdef.MailUserInfo{
			{
				Key:  "{u1}",
				Type: 0,
				User: &structsdef.SimpleUserInfo{
					UserId:    40000001,
					Nickname:  "测试玩家",
					IconToken: 1,
					UserSex:   0,
				},
				Style: &structsdef.MailStyleInfo{
					Color:    "#2992FF",
					Font:     "sh_sdf",
					FontS:    40,
					FontW:    1,
					ShowIcon: false,
				},
			},
		},
	}

	// 发送富文本邮件
	err := GlobalMailService.SendRichTextMail(ctx,
		"富文本邮件",
		"富文本邮件内容",
		"系统",
		panelInfo,
		40000001,
		nil,
		time.Now().Add(24*time.Hour).Unix())

	if err != nil {
		logger.ErrorWF("发送富文本邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("富文本邮件发送成功")
	}
}

// exampleVoteMail 投票邮件示例
func exampleVoteMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 投票邮件示例 ===")

	voteReq := &VoteMailRequest{
		Title:       "游戏体验调查",
		Content:     "请参与我们的游戏体验调查",
		VoteId:      1001,
		VoteTitle:   "您对当前游戏版本的满意度如何？",
		Options:     []string{"非常满意", "满意", "一般", "不满意", "非常不满意"},
		EndTime:     time.Now().Add(7 * 24 * time.Hour).Unix(),
		MaxChoices:  1,
		ReceiverIds: []uint64{40000001, 40000002, 40000003},
	}

	err := GlobalMailService.SendVoteMail(ctx, voteReq)
	if err != nil {
		logger.ErrorWF("发送投票邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("投票邮件发送成功")
	}
}

// exampleBattleReportMail 战报邮件示例
func exampleBattleReportMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 战报邮件示例 ===")

	battleReq := &BattleReportMailRequest{
		Title:          "战斗结果报告",
		Content:        "您的战斗已经结束，请查看详细战报",
		BattleRecordId: 12345,
		BattleData:     `{"attacker":"玩家A","defender":"玩家B","result":"victory"}`,
		AttackerId:     40000001,
		DefenderId:     40000002,
		BattleResult:   1, // 胜利
		ReceiverId:     40000001,
	}

	err := GlobalMailService.SendBattleReportMail(ctx, battleReq)
	if err != nil {
		logger.ErrorWF("发送战报邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("战报邮件发送成功")
	}
}

// exampleGiftBoxMail 礼盒邮件示例
func exampleGiftBoxMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 礼盒邮件示例 ===")

	giftReq := &GiftBoxMailRequest{
		Title:   "节日礼盒",
		Content: "恭喜您获得节日礼盒，请查收！",
		BoxType: 1, // 随机宝箱
		BoxId:   2001,
		ItemList: []*MazeCommon.MazeItem{
			{
				ItemId:   proto.Int32(46700001),
				Count:    proto.Int64(5),
				ItemName: proto.String("装备碎片"),
			},
		},
		ReceiverIds: []uint64{40000001, 40000002},
	}

	err := GlobalMailService.SendGiftBoxMail(ctx, giftReq)
	if err != nil {
		logger.ErrorWF("发送礼盒邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("礼盒邮件发送成功")
	}
}

// exampleLinkMail 跳转邮件示例
func exampleLinkMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 跳转邮件示例 ===")

	linkReq := &LinkMailRequest{
		Title:       "活动通知",
		Content:     "点击查看最新活动详情",
		LinkType:    1, // 跳转到活动页面
		LinkParam:   "activity_id=1001",
		LinkUrl:     "https://game.example.com/activity/1001",
		ReceiverIds: []uint64{40000001, 40000002},
	}

	err := GlobalMailService.SendLinkMail(ctx, linkReq)
	if err != nil {
		logger.ErrorWF("发送跳转邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("跳转邮件发送成功")
	}
}

// exampleBroadcastMail 广播邮件示例
func exampleBroadcastMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 广播邮件示例 ===")

	broadcastReq := &BroadcastMailRequest{
		Title:             "系统维护通知",
		Content:           "系统将于今晚进行维护，请提前做好准备",
		SenderName:        "系统管理员",
		MailType:          int32(constdef.SYSTEM),
		MailSubType:       int32(constdef.AWARD),
		Attachments:       nil,
		ExpireTime:        time.Now().Add(24 * time.Hour).Unix(),
		BroadcastFlag:     int32(constdef.MAIL_FAMILY_BROADCAST),
		BroadcastId:       1001, // 家族ID
		PanelInfo:         nil,
		UserLimitType:     0,
		UserLimitValueMin: 0,
		UserLimitValueMax: 0,
	}

	result, err := GlobalMailService.BroadcastMail(ctx, broadcastReq)
	if err != nil {
		logger.ErrorWF("发送广播邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("广播邮件发送成功",
			zap.Int("successCount", result.SuccessCount),
			zap.Int("failedCount", result.FailedCount))
	}
}

// exampleConfigMail 配置邮件示例
func exampleConfigMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 配置邮件示例 ===")

	// 通配符数据
	wildcardData := map[string]interface{}{
		"user_name": "测试玩家",
		"item_name": "神秘装备",
		"count":     5,
	}

	// 使用配置ID发送邮件
	err := GlobalMailService.SendMailWithConfig(ctx, 1001, 40000001, wildcardData)
	if err != nil {
		logger.ErrorWF("发送配置邮件失败", zap.Error(err))
	} else {
		logger.InfoWF("配置邮件发送成功")
	}
}

// exampleSystemMail 系统监控邮件示例
func exampleSystemMail(ctx context.Context, logger fklog.FKLogI) {
	logger.InfoWF("=== 系统监控邮件示例 ===")

	// 发送错误邮件
	err := GlobalMailSender.ErrorMailSender(ctx, 10001, "数据库连接失败")
	if err != nil {
		logger.ErrorWF("发送错误邮件失败", zap.Error(err))
	}

	// 发送警告邮件
	err = GlobalMailSender.AlertMailSender(ctx, 10002, "内存使用率过高")
	if err != nil {
		logger.ErrorWF("发送警告邮件失败", zap.Error(err))
	}

	// 发送通知邮件
	err = GlobalMailSender.NotifyMailSender(ctx, 10003, "系统运行正常")
	if err != nil {
		logger.ErrorWF("发送通知邮件失败", zap.Error(err))
	}

	// 获取邮件指标
	metrics, err := GlobalMailService.GetMailMetrics(ctx)
	if err != nil {
		logger.ErrorWF("获取邮件指标失败", zap.Error(err))
	} else {
		logger.InfoWF("邮件指标",
			zap.Any("metrics", metrics))
	}
}
