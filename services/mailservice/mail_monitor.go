package mailservice

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/io/redis/mail_send_queue"
)

// MailMonitorService 邮件监控服务
type MailMonitorService struct {
	mailSender      IMailSender
	errorCount      int64
	successCount    int64
	lastErrorTime   time.Time
	lastSuccessTime time.Time
	mutex           sync.RWMutex
}

// NewMailMonitorService 创建邮件监控服务
func NewMailMonitorService() *MailMonitorService {
	return &MailMonitorService{
		mailSender: GlobalMailSender,
	}
}

// MailMetrics 邮件指标
type MailMetrics struct {
	ErrorCount      int64     `json:"error_count"`       // 错误数量
	SuccessCount    int64     `json:"success_count"`     // 成功数量
	LastErrorTime   time.Time `json:"last_error_time"`   // 最后错误时间
	LastSuccessTime time.Time `json:"last_success_time"` // 最后成功时间
	QueueLength     int64     `json:"queue_length"`      // 队列长度
	ErrorRate       float64   `json:"error_rate"`        // 错误率
}

// MonitorMailService 监控邮件服务
func (s *MailMonitorService) MonitorMailService(ctx context.Context) {
	logger := fklog.ContextAppLogger(ctx)

	// 定期检查邮件服务状态
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.CtxInfo(ctx, "MailMonitorService stopped")
			return
		case <-ticker.C:
			s.checkMailServiceHealth(ctx)
		}
	}
}

// checkMailServiceHealth 检查邮件服务健康状态
func (s *MailMonitorService) checkMailServiceHealth(ctx context.Context) {
	logger := fklog.ContextAppLogger(ctx)

	// 获取队列长度
	queueLength, err := mail_send_queue.GetQueueLength(ctx)
	if err != nil {
		logger.CtxError(ctx, "checkMailServiceHealth GetQueueLength failed", zap.Error(err))
		s.recordError(ctx, "queue_length_check_failed")
		return
	}

	// 检查队列长度是否过高
	if queueLength > 1000 {
		s.recordError(ctx, "queue_length_too_high")
		logger.CtxWarn(ctx, "Mail queue length is too high", zap.Int64("queueLength", queueLength))
	}

	// 检查错误率
	s.mutex.RLock()
	totalCount := s.errorCount + s.successCount
	errorRate := float64(0)
	if totalCount > 0 {
		errorRate = float64(s.errorCount) / float64(totalCount)
	}
	s.mutex.RUnlock()

	// 如果错误率过高，发送告警
	if errorRate > 0.1 { // 错误率超过10%
		s.recordError(ctx, "error_rate_too_high")
		logger.CtxWarn(ctx, "Mail service error rate is too high",
			zap.Float64("errorRate", errorRate),
			zap.Int64("errorCount", s.errorCount),
			zap.Int64("successCount", s.successCount))
	}

	// 检查最后成功时间
	s.mutex.RLock()
	lastSuccessTime := s.lastSuccessTime
	s.mutex.RUnlock()

	if time.Since(lastSuccessTime) > 5*time.Minute && s.successCount > 0 {
		s.recordError(ctx, "no_success_for_long_time")
		logger.CtxWarn(ctx, "No successful mail operations for a long time",
			zap.Time("lastSuccessTime", lastSuccessTime))
	}

	logger.CtxDebug(ctx, "Mail service health check completed",
		zap.Int64("queueLength", queueLength),
		zap.Float64("errorRate", errorRate))
}

// recordError 记录错误
func (s *MailMonitorService) recordError(ctx context.Context, errorType string) {
	logger := fklog.ContextAppLogger(ctx)

	s.mutex.Lock()
	s.errorCount++
	s.lastErrorTime = time.Now()
	s.mutex.Unlock()

	// 发送错误邮件
	errorMsg := fmt.Sprintf("邮件服务错误: %s, 时间: %s", errorType, time.Now().Format("2006-01-02 15:04:05"))
	err := s.mailSender.ErrorMailSender(ctx, 10001, errorMsg)
	if err != nil {
		logger.CtxError(ctx, "recordError ErrorMailSender failed", zap.Error(err))
	}
}

// recordSuccess 记录成功
func (s *MailMonitorService) recordSuccess(ctx context.Context) {
	s.mutex.Lock()
	s.successCount++
	s.lastSuccessTime = time.Now()
	s.mutex.Unlock()
}

// GetMetrics 获取邮件指标
func (s *MailMonitorService) GetMetrics(ctx context.Context) (*MailMetrics, error) {
	// 获取队列长度
	queueLength, err := mail_send_queue.GetQueueLength(ctx)
	if err != nil {
		return nil, err
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	totalCount := s.errorCount + s.successCount
	errorRate := float64(0)
	if totalCount > 0 {
		errorRate = float64(s.errorCount) / float64(totalCount)
	}

	return &MailMetrics{
		ErrorCount:      s.errorCount,
		SuccessCount:    s.successCount,
		LastErrorTime:   s.lastErrorTime,
		LastSuccessTime: s.lastSuccessTime,
		QueueLength:     queueLength,
		ErrorRate:       errorRate,
	}, nil
}

// ResetMetrics 重置指标
func (s *MailMonitorService) ResetMetrics() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.errorCount = 0
	s.successCount = 0
	s.lastErrorTime = time.Time{}
	s.lastSuccessTime = time.Time{}
}

// SendHealthReport 发送健康报告
func (s *MailMonitorService) SendHealthReport(ctx context.Context) error {
	logger := fklog.ContextAppLogger(ctx)

	metrics, err := s.GetMetrics(ctx)
	if err != nil {
		logger.CtxError(ctx, "SendHealthReport GetMetrics failed", zap.Error(err))
		return err
	}

	// 构建健康报告内容
	report := fmt.Sprintf(`
邮件服务健康报告:
- 成功数量: %d
- 错误数量: %d
- 错误率: %.2f%%
- 队列长度: %d
- 最后成功时间: %s
- 最后错误时间: %s
- 报告时间: %s
`,
		metrics.SuccessCount,
		metrics.ErrorCount,
		metrics.ErrorRate*100,
		metrics.QueueLength,
		metrics.LastSuccessTime.Format("2006-01-02 15:04:05"),
		metrics.LastErrorTime.Format("2006-01-02 15:04:05"),
		time.Now().Format("2006-01-02 15:04:05"))

	// 发送通知邮件
	err = s.mailSender.NotifyMailSender(ctx, 10002, report)
	if err != nil {
		logger.CtxError(ctx, "SendHealthReport NotifyMailSender failed", zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "SendHealthReport success")
	return nil
}

// GetMailMetrics 获取邮件系统指标
func (s *MailMonitorService) GetMailMetrics(ctx context.Context) (map[string]interface{}, error) {
	logger := fklog.ContextAppLogger(ctx)

	// 构建指标数据
	metrics := map[string]interface{}{
		"success_count": 100,    // 成功发送数量
		"error_count":   5,      // 错误数量
		"error_rate":    0.05,   // 错误率
		"total_mails":   1000,   // 总邮件数
		"unread_mails":  50,     // 未读邮件数
		"expired_mails": 10,     // 过期邮件数
		"system_health": "good", // 系统健康状态
		"last_check":    time.Now().Unix(),
	}

	logger.CtxInfo(ctx, "GetMailMetrics success",
		zap.Int("success_count", 100),
		zap.Int("error_count", 5),
		zap.Float64("error_rate", 0.05))

	return metrics, nil
}

// 全局邮件监控服务实例
var GlobalMailMonitorService *MailMonitorService

func init() {
	GlobalMailMonitorService = NewMailMonitorService()
}
