package mail_module

import (
	"context"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/MailBoxBodyRedis"

	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/TradeNumberCheckRedis"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"go.uber.org/zap"
)

// 广播信封
func BroadcastMail(ctx context.Context, logger fklog.FKLogI, mail *MazeMail.MailInfo) error {
	baseInfo := mail.Base
	if baseInfo == nil {
		logger.ErrorWF("mail base not fill", zap.Any("mailInfo", mail))
		return errors.ARGS_NOT_MATCH
	}

	contentType := baseInfo.GetMailContentType()
	if baseInfo.GetFromServer() == nil || contentType < 1 {
		logger.ErrorWF("sendMail2User from_server or content_type not fill", zap.Any("mail", mail))
		return errors.ARGS_NOT_MATCH
	}

	if mail.GetTitle() == "" {
		logger.ErrorWF("mail title not fill", zap.Any("mailInfo", mail))
		return errors.ARGS_NOT_MATCH
	}

	bid := baseInfo.GetBroadcastFlag()
	guid := baseInfo.GetMailId()
	//判断信封是否存在
	isExist, err := TradeNumberCheckRedis.IsMailExist(ctx, logger, uint64(bid), guid)
	if err != nil {
		logger.ErrorWF("tradeNum check failed", zap.Error(err), zap.Any("mail", mail))
		return errors.ERR_MAIL_SEND
	}

	if isExist {
		logger.WarnWF("tradeNum already exist", zap.Any("req", mail))
		return errors.ERR_MAIL_ALREADY_EXIST
	}

	// 广播信息再函数外调用保存
	err = MailBoxBodyRedis.AddMail(logger, mail)
	if err != nil {
		logger.ErrorWF("add mail body", zap.Any("mail", mail), zap.Error(err))
		return errors.COMMON_ERROR_TIPS
	}
	return nil
}
