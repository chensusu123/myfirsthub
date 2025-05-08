// @Author: ZhaoXiming 2025/3/21 16:32
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/cgkargs"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/mailutil"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/MailBoxBodyRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/broadcastmailstatdb"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/mailunreaddb"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"gitlab.ifreetalk.com/plate/protodef/MazeMailCli"
	"go.uber.org/zap"
	"time"
)

func OnGetMailInfoRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.GetMailInfoRQ)
	res := rsMsg.(*MazeMailCli.GetMailInfoRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	userId := shardingID

	ctx.InfoWF("OnGetMailInfoRQ req", zap.Any("with", req))
	defer fkprometheus.DebugPMT("OnGetMailInfoRQ")()

	mailID := req.GetReceiptItems().GetId()
	mail, err := MailBoxBodyRedis.GetMail(ctx, mailID)
	if err == redis.ErrNil {
		ctx.InfoWF("OnGetMailInfoRQ get mail info empty ", zap.Any("req", req))
		return nil
	}
	if err != nil {
		ctx.ErrorWF("OnGetMailInfoRQ get mail info fail ", zap.Any("req", req), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return nil
	}
	if mail == nil {
		ctx.ErrorWF("OnGetMailInfoRQ get mail nil", zap.Any("req", req), zap.Error(err))
		res.ErrInfo = errors.ERR_MAIL_EXPIRE.ToInfo()
		return
	}
	classType := mail.Base.GetClassType()

	isBroadcast, hasRead, hasOperator, _, err := getMailFlag(ctx, userId, mail)
	if err != nil {
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return nil
	}

	// 设置邮件是否已经读取过
	ctx.InfoWF("OnGetMailInfoRQ detail", zap.Any("hasRead", hasRead), zap.Any("mailContentType", mail.GetBase().GetMailContentType()))
	if hasRead == 0 {
		mail.MailHasRead = proto.Int32(1)
		readGuid := make([]uint64, 0, 1)

		if mailutil.IsMailTextOrReport(mail) {
			mail.HasOperator = proto.Int32(1)
			readGuid = append(readGuid, mailID)
			if isBroadcast == 0 {
				now := time.Now().Unix()
				mail.Base.Expires = proto.Int64(now + cgkargs.AlreadyReadMailExpireTime)
			}
		} else {
			//已读已操作
			if mail.GetHasOperator() == 1 {
				readGuid = append(readGuid, mailID)
			}
		}

		if isBroadcast == 0 {
			_ = MailBoxBodyRedis.AddMail(ctx, mail)
		} else {
			// 广播
			err = broadcastmailstatdb.SetUseBroadcastMailState(ctx, userId, mailID, 1, mail.GetHasOperator(), nil, mail.Base.GetExpires())
			if err != nil {
				ctx.ErrorWF("OnGetMailInfoRQ set broadcast mail stat fail", zap.Uint64("uid", userId), zap.Uint64("mailId", mailID), zap.Error(err))
			}
		}
		if len(readGuid) > 0 {
			_ = mailunreaddb.DelUnreadMail(ctx, userId, classType, readGuid)
			var num int32
			num, err = mailunreaddb.UnreadMailNum(ctx, userId, classType)
			if err != nil {
				ctx.ErrorWF("OnGetMailInfoRQ get unread mail num", zap.Error(err), zap.Any("classType", classType))
			} else {
				res.UnreadNum = proto.Int32(num)
			}
			// 获取最新的过期时间
			latestExpireTime, loadErr := mailunreaddb.GetUnreaMailExpireTime(ctx, userId, classType)
			if loadErr != nil {
				ctx.WarnWF("OnGetMailInfoRQ load data error", zap.Error(loadErr), zap.Any("classType", classType))
			} else if latestExpireTime > 0 {
				res.LatestExpireTime = proto.Int64(latestExpireTime)
			}
		}
	} else {
		if isBroadcast != 0 {
			mail.MailHasRead = proto.Int32(1)
			mail.HasOperator = proto.Int32(hasOperator)

		}
	}

	res.MailInfo = mail
	return nil
}

func getMailFlag(logger fklog.FKLogI, userId uint64, mail *MazeMail.MailInfo) (isBroadcast, hasRead, hasOperator int32, vote *MazeMail.UserVoteAnswer, err error) {
	isBroadcast = mail.Base.GetBroadcastFlag()
	if isBroadcast != 0 {
		//var voteData []byte
		hasRead, hasOperator, _, err = broadcastmailstatdb.GetUserBroadcastMailStat(logger, userId, mail.Base.GetMailId())
		if err != nil {
			logger.ErrorWF("getMailFlag check broadcast mail stat", zap.Any("mail", mail), zap.Error(err))
		}
		//if len(voteData) > 0 && mail.Vote != nil {
		//	vote = &MazeMail.UserVoteAnswer{}
		//	err = proto.Unmarshal(voteData, vote)
		//	if err != nil {
		//		logger.ErrorWF("getMailFlag unmarshal broadcast mail vote info fail", zap.Uint64("uid", userId), zap.Any("mail", mail), zap.Error(err))
		//	}
		//}
		return
	}

	hasRead = mail.GetMailHasRead()
	hasOperator = mail.GetHasOperator()
	//vote = mail.Vote.GetUserAnswer()
	return
}
