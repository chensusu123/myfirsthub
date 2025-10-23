package mail_send_queue

import (
	"context"
	"encoding/json"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkredis"
	"go.uber.org/zap"

	"maze_game_server/model/mailmodel"
)

func init() {
	fkconfig.RegisterNode(16252, redisQueue)
}

const redisKey = "ppworld:mail:send:queue"

var redisQueue = &fkredis.QueueRedis{MsgKey: redisKey}

// MailQueueMessage 邮件队列消息结构
type MailQueueMessage struct {
	ToUser   uint64              `json:"to_user"`
	MailInfo *mailmodel.MailInfo `json:"mail_info"`
	Priority int32               `json:"priority"` // 优先级，数字越小优先级越高
}

// PushNewInfo 推送邮件信息到Redis队列
func PushNewInfo(ctx context.Context, logger fklog.FKLogI, uid uint64, mailInfo *mailmodel.MailInfo) error {
	return PushNewInfoWithPriority(ctx, logger, uid, mailInfo, 0)
}

// PushNewInfoWithPriority 推送邮件信息到Redis队列（带优先级）
func PushNewInfoWithPriority(ctx context.Context, logger fklog.FKLogI, uid uint64, mailInfo *mailmodel.MailInfo, priority int32) error {
	message := &MailQueueMessage{
		ToUser:   uid,
		MailInfo: mailInfo,
		Priority: priority,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logger.CtxError(ctx, "PushNewInfoWithPriority Marshal error", zap.Error(err))
		return err
	}

	_, err = redisQueue.Do(ctx, "LPUSH", redisQueue.GetKey(), data)
	if err != nil {
		logger.CtxError(ctx, "PushNewInfoWithPriority push mail error", zap.Error(err))
		return err
	}

	logger.CtxInfo(ctx, "PushNewInfoWithPriority success",
		zap.Uint64("uid", uid),
		zap.Uint64("mailId", mailInfo.ID),
		zap.Int32("priority", priority))

	return nil
}

// PushBatchMail 批量推送邮件到队列
func PushBatchMail(ctx context.Context, logger fklog.FKLogI, userMails map[uint64]*mailmodel.MailInfo) error {
	if len(userMails) == 0 {
		return nil
	}

	// 批量推送

	for uid, mailInfo := range userMails {
		message := &MailQueueMessage{
			ToUser:   uid,
			MailInfo: mailInfo,
			Priority: 0,
		}

		data, err := json.Marshal(message)
		if err != nil {
			logger.CtxError(ctx, "PushBatchMail Marshal error",
				zap.Uint64("uid", uid),
				zap.Error(err))
			continue
		}

		_, err = redisQueue.Do(ctx, "LPUSH", redisQueue.GetKey(), data)
		if err != nil {
			logger.CtxError(ctx, "PushBatchMail LPUSH error", zap.Error(err))
			return err
		}
	}

	logger.CtxInfo(ctx, "PushBatchMail success", zap.Int("count", len(userMails)))
	return nil
}

// GetQueueLength 获取队列长度
func GetQueueLength(ctx context.Context) (int64, error) {
	result, err := redisQueue.Do(ctx, "LLEN", redisQueue.GetKey())
	if err != nil {
		return 0, err
	}
	return result.(int64), nil
}

// ClearQueue 清空队列
func ClearQueue(ctx context.Context) error {
	_, err := redisQueue.Do(ctx, "DEL", redisQueue.GetKey())
	return err
}
