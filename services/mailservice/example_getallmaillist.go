package mailservice

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/model/mailmodel"
)

// ExampleGetAllMailList 展示如何使用GetAllMailList方法
func ExampleGetAllMailList() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	// 用户ID
	userId := uint64(40000005)

	// 调用GetAllMailList方法
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailList failed", zap.Error(err))
		return
	}

	// 处理结果
	logger.InfoWF("GetAllMailList success",
		zap.Uint64("userId", userId),
		zap.Int("totalLabels", len(mailMap)))

	// 遍历所有标签的邮件
	for label, mailList := range mailMap {
		logger.InfoWF("Processing mail label",
			zap.Int32("label", label),
			zap.Int("mailCount", len(mailList)))

		// 处理每个标签的邮件
		for i, mail := range mailList {
			logger.InfoWF("Mail info",
				zap.Int("index", i),
				zap.Uint64("mailId", mail.ID),
				zap.String("title", mail.Title),
				zap.String("content", mail.Content),
				zap.Bool("isRead", mail.IsRead),
				zap.Bool("isGetAttach", mail.IsGetAttach),
				zap.Int64("sendTime", mail.SendTime),
				zap.Int64("expireTime", mail.ExpireTime))
		}
	}
}

// ExampleGetAllMailListWithFilter 展示如何过滤和处理邮件
func ExampleGetAllMailListWithFilter() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	userId := uint64(40000005)

	// 获取所有邮件
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailList failed", zap.Error(err))
		return
	}

	// 统计信息
	totalMails := 0
	unreadMails := 0
	unclaimedAttachments := 0

	// 按标签处理邮件
	for label, mailList := range mailMap {
		labelStats := struct {
			total     int
			unread    int
			unclaimed int
		}{}

		for _, mail := range mailList {
			labelStats.total++
			totalMails++

			if !mail.IsRead {
				labelStats.unread++
				unreadMails++
			}

			if len(mail.Attachments) > 0 && !mail.IsGetAttach {
				labelStats.unclaimed++
				unclaimedAttachments++
			}
		}

		logger.InfoWF("Label statistics",
			zap.Int32("label", label),
			zap.Int("total", labelStats.total),
			zap.Int("unread", labelStats.unread),
			zap.Int("unclaimed", labelStats.unclaimed))
	}

	// 总体统计
	logger.InfoWF("Overall statistics",
		zap.Int("totalMails", totalMails),
		zap.Int("unreadMails", unreadMails),
		zap.Int("unclaimedAttachments", unclaimedAttachments))
}

// ExampleGetAllMailListByLabel 展示如何获取特定标签的邮件
func ExampleGetAllMailListByLabel() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	userId := uint64(40000005)

	// 获取所有邮件
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailList failed", zap.Error(err))
		return
	}

	// 获取系统邮件
	systemMails, exists := mailMap[int32(constdef.MailLabelSystem)]
	if exists {
		logger.InfoWF("System mails",
			zap.Int("count", len(systemMails)))

		for _, mail := range systemMails {
			logger.InfoWF("System mail",
				zap.Uint64("mailId", mail.ID),
				zap.String("title", mail.Title),
				zap.Bool("isRead", mail.IsRead))
		}
	} else {
		logger.InfoWF("No system mails found")
	}

	// 获取家族邮件
	familyMails, exists := mailMap[int32(constdef.MailLabelFamily)]
	if exists {
		logger.InfoWF("Family mails",
			zap.Int("count", len(familyMails)))
	} else {
		logger.InfoWF("No family mails found")
	}

	// 获取联盟邮件
	allianceMails, exists := mailMap[int32(constdef.MailLabelAlliance)]
	if exists {
		logger.InfoWF("Alliance mails",
			zap.Int("count", len(allianceMails)))
	} else {
		logger.InfoWF("No alliance mails found")
	}
}

// ExampleGetAllMailListWithPagination 展示如何实现分页
func ExampleGetAllMailListWithPagination() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	userId := uint64(40000005)
	pageSize := 10
	page := 1

	// 获取所有邮件
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailList failed", zap.Error(err))
		return
	}

	// 合并所有标签的邮件
	allMails := make([]*mailmodel.MailInfo, 0)
	for _, mailList := range mailMap {
		allMails = append(allMails, mailList...)
	}

	// 计算分页
	totalMails := len(allMails)
	totalPages := (totalMails + pageSize - 1) / pageSize

	logger.InfoWF("Pagination info",
		zap.Int("totalMails", totalMails),
		zap.Int("pageSize", pageSize),
		zap.Int("totalPages", totalPages),
		zap.Int("currentPage", page))

	// 计算当前页的邮件
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalMails {
		end = totalMails
	}

	if start >= totalMails {
		logger.InfoWF("No mails for this page")
		return
	}

	// 获取当前页的邮件
	currentPageMails := allMails[start:end]

	logger.InfoWF("Current page mails",
		zap.Int("count", len(currentPageMails)),
		zap.Int("start", start),
		zap.Int("end", end))

	for i, mail := range currentPageMails {
		logger.InfoWF("Mail on current page",
			zap.Int("index", i),
			zap.Uint64("mailId", mail.ID),
			zap.String("title", mail.Title),
			zap.Int32("label", mail.Label))
	}
}

// ExampleGetAllMailListWithSearch 展示如何搜索邮件
func ExampleGetAllMailListWithSearch() {
	ctx := context.Background()
	logger := fklog.ContextAppLogger(ctx)

	userId := uint64(40000005)
	searchKeyword := "测试" // 搜索关键词

	// 获取所有邮件
	mailMap, err := GlobalMailService.GetAllMailList(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailList failed", zap.Error(err))
		return
	}

	// 搜索邮件
	searchResults := make([]*mailmodel.MailInfo, 0)

	for label, mailList := range mailMap {
		for _, mail := range mailList {
			// 在标题和内容中搜索关键词
			if contains(mail.Title, searchKeyword) || contains(mail.Content, searchKeyword) {
				searchResults = append(searchResults, mail)
				logger.InfoWF("Found matching mail",
					zap.Uint64("mailId", mail.ID),
					zap.String("title", mail.Title),
					zap.Int32("label", label))
			}
		}
	}

	logger.InfoWF("Search results",
		zap.String("keyword", searchKeyword),
		zap.Int("resultCount", len(searchResults)))
}

// contains 检查字符串是否包含子字符串（简单实现）
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

// containsSubstring 检查字符串是否包含子字符串
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
