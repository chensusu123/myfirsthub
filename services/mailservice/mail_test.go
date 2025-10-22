package mailservice

import (
	"context"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"maze_game_server/common/constdef"
)

// TestGetAllMailList 测试GetAllMailList方法
func TestGetAllMailList(t *testing.T) {
	// 创建测试上下文和日志器
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	// 测试用户ID
	userId := uint64(40000005)

	// 调用GetAllMailList方法
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		t.Errorf("GetAllMailList failed: %v", err)
		return
	}

	// 验证结果
	if mailMap == nil {
		t.Error("GetAllMailList returned nil mailMap")
		return
	}

	// 打印结果
	t.Logf("GetAllMailList success for userId: %d", userId)
	t.Logf("Total labels: %d", len(mailMap))

	for label, mailList := range mailMap {
		t.Logf("Label %d has %d mails", label, len(mailList))
		for i, mail := range mailList {
			if i < 3 { // 只打印前3个邮件
				t.Logf("  Mail %d: ID=%d, Title=%s, IsRead=%v, IsGetAttach=%v",
					i, mail.ID, mail.Title, mail.IsRead, mail.IsGetAttach)
			}
		}
	}
}

// TestGetAllMailListWithEmptyUser 测试空用户的情况
func TestGetAllMailListWithEmptyUser(t *testing.T) {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	// 使用一个不存在的用户ID
	userId := uint64(99999999)

	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		t.Errorf("GetAllMailList failed: %v", err)
		return
	}

	// 新用户应该返回空的邮件映射
	if mailMap == nil {
		t.Error("GetAllMailList returned nil mailMap for new user")
		return
	}

	t.Logf("GetAllMailList for new user returned %d labels", len(mailMap))
}

// TestGetAllMailListIntegration 集成测试
func TestGetAllMailListIntegration(t *testing.T) {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	userId := uint64(40000005)

	// 1. 先发送一些测试邮件
	testMails := []struct {
		title   string
		content string
		label   int32
	}{
		{"测试邮件1", "这是第一封测试邮件", int32(constdef.MailLabelSystem)},
		{"测试邮件2", "这是第二封测试邮件", int32(constdef.MailLabelSystem)},
		{"测试邮件3", "这是第三封测试邮件", int32(constdef.MailLabelFamily)},
	}

	for _, testMail := range testMails {
		err := GlobalMailService.SendMail(logger, testMail.title, testMail.content, "测试系统",
			testMail.label, userId, nil, 0)
		if err != nil {
			t.Logf("SendMail failed for %s: %v", testMail.title, err)
		}
	}

	// 2. 获取所有邮件列表
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		t.Errorf("GetAllMailList failed: %v", err)
		return
	}

	// 3. 验证结果
	t.Logf("Integration test - GetAllMailList returned %d labels", len(mailMap))

	for label, mailList := range mailMap {
		t.Logf("Label %d: %d mails", label, len(mailList))
	}
}
