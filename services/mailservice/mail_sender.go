package mailservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// IMailSender 邮件发送接口
type IMailSender interface {
	// ErrorMailSender 错误邮件发送 - 服务内部数据或逻辑错误，必须处理
	ErrorMailSender(ctx context.Context, mailID int32, errorMsg string) error

	// AlertMailSender 警告邮件发送 - 服务内部数据或逻辑错误，但不影响业务继续执行
	AlertMailSender(ctx context.Context, mailID int32, alertMsg string) error

	// NotifyMailSender 通知邮件发送 - 服务运行状态正常上报
	NotifyMailSender(ctx context.Context, mailID int32, notifyMsg string) error
}

// MailSender 邮件发送器实现
type MailSender struct {
	mailService MailService
}

// NewMailSender 创建邮件发送器
func NewMailSender() IMailSender {
	return &MailSender{
		mailService: GlobalMailService,
	}
}

// ErrorMailSender 发送错误邮件
func (s *MailSender) ErrorMailSender(ctx context.Context, mailID int32, errorMsg string) error {

	// 构建错误邮件内容
	title := fmt.Sprintf("系统错误报告 [%d]", mailID)
	content := fmt.Sprintf("系统发生错误，请及时处理：\n%s", errorMsg)

	// 创建富文本内容
	panelInfo := s.buildErrorMailPanel(mailID, errorMsg)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 创建邮件信息
	mailInfo := &mailmodel.MailInfo{
		ID:              uint64(time.Now().UnixNano()),
		Title:           title,
		Content:         content,
		Label:           int32(constdef.MailLabelSystem),
		MailType:        int32(constdef.SYSTEM),
		MailSubType:     int32(constdef.AWARD),
		MailContentType: int32(constdef.CONTENT_TEXT),
		Sender:          "系统监控",
		ReciverID:       0, // 系统邮件，发送给管理员
		IsRead:          false,
		IsGetAttach:     false,
		SendTime:        time.Now().Unix(),
		ExpireTime:      time.Now().Add(24 * time.Hour).Unix(), // 24小时后过期
		PanelInfo:       string(panelInfoJSON),
		PagePopUp:       true, // 错误邮件需要弹框提醒
		AllSupportPop:   true,
	}

	// 获取logger
	logger := fklog.ContextAppLogger(ctx)

	// 发送邮件
	err := s.mailService.SendMail(logger, title, content, "系统监控",
		int32(constdef.MailLabelSystem), 0, nil, mailInfo.ExpireTime)

	if err != nil {
		logger.CtxError(ctx, "ErrorMailSender failed",
			zap.Int32("mailID", mailID),
			zap.String("errorMsg", errorMsg),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "ErrorMailSender success",
		zap.Int32("mailID", mailID),
		zap.String("errorMsg", errorMsg))

	return nil
}

// AlertMailSender 发送警告邮件
func (s *MailSender) AlertMailSender(ctx context.Context, mailID int32, alertMsg string) error {

	// 构建警告邮件内容
	title := fmt.Sprintf("系统警告 [%d]", mailID)
	content := fmt.Sprintf("系统出现警告，请注意：\n%s", alertMsg)

	// 创建富文本内容
	panelInfo := s.buildAlertMailPanel(mailID, alertMsg)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 创建邮件信息
	mailInfo := &mailmodel.MailInfo{
		ID:              uint64(time.Now().UnixNano()),
		Title:           title,
		Content:         content,
		Label:           int32(constdef.MailLabelSystem),
		MailType:        int32(constdef.SYSTEM),
		MailSubType:     int32(constdef.AWARD),
		MailContentType: int32(constdef.CONTENT_TEXT),
		Sender:          "系统监控",
		ReciverID:       0, // 系统邮件，发送给管理员
		IsRead:          false,
		IsGetAttach:     false,
		SendTime:        time.Now().Unix(),
		ExpireTime:      time.Now().Add(12 * time.Hour).Unix(), // 12小时后过期
		PanelInfo:       string(panelInfoJSON),
		PagePopUp:       false, // 警告邮件不需要弹框
		AllSupportPop:   false,
	}

	// 获取logger
	logger := fklog.ContextAppLogger(ctx)

	// 发送邮件
	err := s.mailService.SendMail(logger, title, content, "系统监控",
		int32(constdef.MailLabelSystem), 0, nil, mailInfo.ExpireTime)

	if err != nil {
		logger.CtxError(ctx, "AlertMailSender failed",
			zap.Int32("mailID", mailID),
			zap.String("alertMsg", alertMsg),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "AlertMailSender success",
		zap.Int32("mailID", mailID),
		zap.String("alertMsg", alertMsg))

	return nil
}

// NotifyMailSender 发送通知邮件
func (s *MailSender) NotifyMailSender(ctx context.Context, mailID int32, notifyMsg string) error {

	// 构建通知邮件内容
	title := fmt.Sprintf("系统通知 [%d]", mailID)
	content := fmt.Sprintf("系统运行状态正常：\n%s", notifyMsg)

	// 创建富文本内容
	panelInfo := s.buildNotifyMailPanel(mailID, notifyMsg)
	panelInfoJSON, _ := json.Marshal(panelInfo)

	// 创建邮件信息
	mailInfo := &mailmodel.MailInfo{
		ID:              uint64(time.Now().UnixNano()),
		Title:           title,
		Content:         content,
		Label:           int32(constdef.MailLabelSystem),
		MailType:        int32(constdef.SYSTEM),
		MailSubType:     int32(constdef.AWARD),
		MailContentType: int32(constdef.CONTENT_TEXT),
		Sender:          "系统监控",
		ReciverID:       0, // 系统邮件，发送给管理员
		IsRead:          false,
		IsGetAttach:     false,
		SendTime:        time.Now().Unix(),
		ExpireTime:      time.Now().Add(6 * time.Hour).Unix(), // 6小时后过期
		PanelInfo:       string(panelInfoJSON),
		PagePopUp:       false, // 通知邮件不需要弹框
		AllSupportPop:   false,
	}

	// 获取logger
	logger := fklog.ContextAppLogger(ctx)

	// 发送邮件
	err := s.mailService.SendMail(logger, title, content, "系统监控",
		int32(constdef.MailLabelSystem), 0, nil, mailInfo.ExpireTime)

	if err != nil {
		logger.CtxError(ctx, "NotifyMailSender failed",
			zap.Int32("mailID", mailID),
			zap.String("notifyMsg", notifyMsg),
			zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "NotifyMailSender success",
		zap.Int32("mailID", mailID),
		zap.String("notifyMsg", notifyMsg))

	return nil
}

// buildErrorMailPanel 构建错误邮件面板
func (s *MailSender) buildErrorMailPanel(mailID int32, errorMsg string) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: fmt.Sprintf("系统错误报告 [%d]", mailID),
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的管理员:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: fmt.Sprintf("系统发生错误，请及时处理：\n%s", errorMsg),
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统监控",
			},
		},
		User: []*structsdef.MailUserInfo{
			{
				Key:  "{u1}",
				Type: 0,
				Style: &structsdef.MailStyleInfo{
					Color:    "#FF0000", // 红色表示错误
					Font:     "sh_sdf",
					FontS:    40,
					FontW:    1,
					ShowIcon: false,
				},
			},
		},
	}
}

// buildAlertMailPanel 构建警告邮件面板
func (s *MailSender) buildAlertMailPanel(mailID int32, alertMsg string) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: fmt.Sprintf("系统警告 [%d]", mailID),
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的管理员:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: fmt.Sprintf("系统出现警告，请注意：\n%s", alertMsg),
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统监控",
			},
		},
		User: []*structsdef.MailUserInfo{
			{
				Key:  "{u1}",
				Type: 0,
				Style: &structsdef.MailStyleInfo{
					Color:    "#FFA500", // 橙色表示警告
					Font:     "sh_sdf",
					FontS:    40,
					FontW:    1,
					ShowIcon: false,
				},
			},
		},
	}
}

// buildNotifyMailPanel 构建通知邮件面板
func (s *MailSender) buildNotifyMailPanel(mailID int32, notifyMsg string) *structsdef.MailPanelInfo {
	return &structsdef.MailPanelInfo{
		Content: fmt.Sprintf("系统通知 [%d]", mailID),
		Contents: []*structsdef.MailPanelContent{
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_RECEIVER,
				Content: "亲爱的管理员:",
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_CONTENT,
				Content: fmt.Sprintf("系统运行状态正常：\n%s", notifyMsg),
			},
			{
				Type:    constdef.MAIL_PANEL_CONTENT_MAIL_FROM,
				Content: "发件人: 系统监控",
			},
		},
		User: []*structsdef.MailUserInfo{
			{
				Key:  "{u1}",
				Type: 0,
				Style: &structsdef.MailStyleInfo{
					Color:    "#00FF00", // 绿色表示正常
					Font:     "sh_sdf",
					FontS:    40,
					FontW:    1,
					ShowIcon: false,
				},
			},
		},
	}
}

// 全局邮件发送器实例
var GlobalMailSender IMailSender

func init() {
	GlobalMailSender = NewMailSender()
}
