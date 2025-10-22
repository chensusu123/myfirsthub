package mailservice

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"

	"maze_game_server/common/constdef"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/redis/mail_send_queue"
	"maze_game_server/model/mailmodel"
)

// MailBatchService 群发邮件服务
type MailBatchService struct {
	mailService MailService
}

// NewMailBatchService 创建群发邮件服务
func NewMailBatchService() *MailBatchService {
	return &MailBatchService{
		mailService: GlobalMailService,
	}
}

// BatchMailRequest 群发邮件请求
type BatchMailRequest struct {
	Title       string                    `json:"title"`         // 邮件标题
	Content     string                    `json:"content"`       // 邮件内容
	SenderName  string                    `json:"sender_name"`   // 发送者名称
	MailType    int32                     `json:"mail_type"`     // 邮件类型
	MailSubType int32                     `json:"mail_sub_type"` // 邮件子类型
	Attachments []*mailmodel.Attachment   `json:"attachments"`   // 附件
	ExpireTime  int64                     `json:"expire_time"`   // 过期时间
	ReceiverIds []uint64                  `json:"receiver_ids"`  // 接收者ID列表
	PanelInfo   *structsdef.MailPanelInfo `json:"panel_info"`    // 富文本面板信息
	UseQueue    bool                      `json:"use_queue"`     // 是否使用队列异步发送
	BatchSize   int                       `json:"batch_size"`    // 批处理大小
	MaxWorkers  int                       `json:"max_workers"`   // 最大工作协程数
}

// BatchMailResult 群发邮件结果
type BatchMailResult struct {
	TotalCount   int      `json:"total_count"`   // 总数量
	SuccessCount int      `json:"success_count"` // 成功数量
	FailedCount  int      `json:"failed_count"`  // 失败数量
	FailedUsers  []uint64 `json:"failed_users"`  // 失败用户列表
	ProcessTime  int64    `json:"process_time"`  // 处理时间(毫秒)
	UseQueue     bool     `json:"use_queue"`     // 是否使用队列
}

// BatchSendMail 群发邮件
func (s *MailBatchService) BatchSendMail(ctx context.Context, req *BatchMailRequest) (*BatchMailResult, error) {
	logger := fklog.ContextAppLogger(ctx)
	startTime := time.Now()

	// 设置默认值
	if req.BatchSize <= 0 {
		req.BatchSize = 100 // 默认批处理大小
	}
	if req.MaxWorkers <= 0 {
		req.MaxWorkers = 10 // 默认最大工作协程数
	}

	result := &BatchMailResult{
		TotalCount:  len(req.ReceiverIds),
		FailedUsers: make([]uint64, 0),
		UseQueue:    req.UseQueue,
	}

	logger.CtxInfo(ctx, "BatchSendMail started",
		zap.Int("totalCount", result.TotalCount),
		zap.Bool("useQueue", req.UseQueue),
		zap.Int("batchSize", req.BatchSize),
		zap.Int("maxWorkers", req.MaxWorkers))

	if req.UseQueue {
		// 使用队列异步发送
		err := s.batchSendWithQueue(ctx, req, result)
		if err != nil {
			logger.CtxError(ctx, "BatchSendMail batchSendWithQueue failed", zap.Error(err))
			return result, err
		}
	} else {
		// 同步发送
		err := s.batchSendSync(ctx, req, result)
		if err != nil {
			logger.CtxError(ctx, "BatchSendMail batchSendSync failed", zap.Error(err))
			return result, err
		}
	}

	result.ProcessTime = time.Since(startTime).Milliseconds()

	logger.CtxInfo(ctx, "BatchSendMail completed",
		zap.Int("totalCount", result.TotalCount),
		zap.Int("successCount", result.SuccessCount),
		zap.Int("failedCount", result.FailedCount),
		zap.Int64("processTime", result.ProcessTime))

	return result, nil
}

// batchSendWithQueue 使用队列异步发送
func (s *MailBatchService) batchSendWithQueue(ctx context.Context, req *BatchMailRequest, result *BatchMailResult) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建富文本面板信息JSON
	panelInfoJSON := ""
	if req.PanelInfo != nil {
		panelInfoBytes, err := json.Marshal(req.PanelInfo)
		if err != nil {
			logger.CtxError(ctx, "batchSendWithQueue marshal panel info failed", zap.Error(err))
		} else {
			panelInfoJSON = string(panelInfoBytes)
		}
	}

	// 分批处理用户ID
	userMails := make(map[uint64]*mailmodel.MailInfo)

	for _, userId := range req.ReceiverIds {
		mailInfo := &mailmodel.MailInfo{
			ID:              uint64(time.Now().UnixNano() + int64(userId)), // 确保唯一性
			Title:           req.Title,
			Content:         req.Content,
			Label:           int32(constdef.MailLabelSystem),
			MailType:        req.MailType,
			MailSubType:     req.MailSubType,
			MailContentType: int32(constdef.CONTENT_TEXT_AND_AWARD),
			Sender:          req.SenderName,
			ReciverID:       userId,
			IsRead:          false,
			IsGetAttach:     false,
			SendTime:        time.Now().Unix(),
			ExpireTime:      req.ExpireTime,
			Attachments:     req.Attachments,
			PanelInfo:       panelInfoJSON,
			PagePopUp:       true,
			AllSupportPop:   true,
		}

		userMails[userId] = mailInfo
	}

	// 推送到Redis队列
	err := mail_send_queue.PushBatchMail(ctx, logger, userMails)
	if err != nil {
		logger.CtxError(ctx, "batchSendWithQueue PushBatchMail failed", zap.Error(err))
		return err
	}

	// 队列模式下，假设全部成功（实际处理由队列消费者完成）
	result.SuccessCount = result.TotalCount
	result.FailedCount = 0

	return nil
}

// batchSendSync 同步发送
func (s *MailBatchService) batchSendSync(ctx context.Context, req *BatchMailRequest, result *BatchMailResult) error {
	logger := fklog.ContextAppLogger(ctx)

	// 构建富文本面板信息JSON
	panelInfoJSON := ""
	if req.PanelInfo != nil {
		panelInfoBytes, err := json.Marshal(req.PanelInfo)
		if err != nil {
			logger.CtxError(ctx, "batchSendSync marshal panel info failed", zap.Error(err))
		} else {
			panelInfoJSON = string(panelInfoBytes)
		}
	}

	// 使用协程池处理
	userChan := make(chan uint64, req.BatchSize)
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// 启动工作协程
	for i := 0; i < req.MaxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for userId := range userChan {
				// 创建邮件信息
				mailInfo := &mailmodel.MailInfo{
					ID:              uint64(time.Now().UnixNano() + int64(userId) + int64(workerID)),
					Title:           req.Title,
					Content:         req.Content,
					Label:           int32(constdef.MailLabelSystem),
					MailType:        req.MailType,
					MailSubType:     req.MailSubType,
					MailContentType: int32(constdef.CONTENT_TEXT_AND_AWARD),
					Sender:          req.SenderName,
					ReciverID:       userId,
					IsRead:          false,
					IsGetAttach:     false,
					SendTime:        time.Now().Unix(),
					ExpireTime:      req.ExpireTime,
					Attachments:     req.Attachments,
					PanelInfo:       panelInfoJSON,
					PagePopUp:       true,
					AllSupportPop:   true,
				}

				// 发送邮件
				err := s.mailService.SendMail(logger, req.Title, req.Content, req.SenderName,
					int32(constdef.MailLabelSystem), userId, req.Attachments, req.ExpireTime)

				mutex.Lock()
				if err != nil {
					logger.CtxError(ctx, "batchSendSync SendMail failed",
						zap.Uint64("userId", userId),
						zap.Int("workerID", workerID),
						zap.Error(err))
					result.FailedCount++
					result.FailedUsers = append(result.FailedUsers, userId)
				} else {
					result.SuccessCount++
				}
				mutex.Unlock()
			}
		}(i)
	}

	// 发送用户ID到通道
	go func() {
		defer close(userChan)
		for _, userId := range req.ReceiverIds {
			userChan <- userId
		}
	}()

	// 等待所有协程完成
	wg.Wait()

	return nil
}

// BatchSendSystemMail 群发系统邮件
func (s *MailBatchService) BatchSendSystemMail(ctx context.Context, title, content, senderName string, receiverIds []uint64, useQueue bool) (*BatchMailResult, error) {
	req := &BatchMailRequest{
		Title:       title,
		Content:     content,
		SenderName:  senderName,
		MailType:    int32(constdef.SYSTEM),
		MailSubType: int32(constdef.AWARD),
		Attachments: nil,
		ExpireTime:  0, // 永不过期
		ReceiverIds: receiverIds,
		PanelInfo:   nil,
		UseQueue:    useQueue,
		BatchSize:   100,
		MaxWorkers:  10,
	}

	return s.BatchSendMail(ctx, req)
}

// BatchSendRewardMail 群发奖励邮件
func (s *MailBatchService) BatchSendRewardMail(ctx context.Context, title, content, senderName string, attachments []*mailmodel.Attachment, receiverIds []uint64, useQueue bool) (*BatchMailResult, error) {
	req := &BatchMailRequest{
		Title:       title,
		Content:     content,
		SenderName:  senderName,
		MailType:    int32(constdef.SYSTEM),
		MailSubType: int32(constdef.AWARD),
		Attachments: attachments,
		ExpireTime:  0, // 永不过期
		ReceiverIds: receiverIds,
		PanelInfo:   nil,
		UseQueue:    useQueue,
		BatchSize:   100,
		MaxWorkers:  10,
	}

	return s.BatchSendMail(ctx, req)
}

// BatchSendRichTextMail 群发富文本邮件
func (s *MailBatchService) BatchSendRichTextMail(ctx context.Context, title, content, senderName string, panelInfo *structsdef.MailPanelInfo, receiverIds []uint64, useQueue bool) (*BatchMailResult, error) {
	req := &BatchMailRequest{
		Title:       title,
		Content:     content,
		SenderName:  senderName,
		MailType:    int32(constdef.SYSTEM),
		MailSubType: int32(constdef.AWARD),
		Attachments: nil,
		ExpireTime:  0, // 永不过期
		ReceiverIds: receiverIds,
		PanelInfo:   panelInfo,
		UseQueue:    useQueue,
		BatchSize:   100,
		MaxWorkers:  10,
	}

	return s.BatchSendMail(ctx, req)
}

// 全局群发邮件服务实例
var GlobalMailBatchService *MailBatchService

func init() {
	GlobalMailBatchService = NewMailBatchService()
}
