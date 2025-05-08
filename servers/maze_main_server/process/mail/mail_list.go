// @Author: ZhaoXiming 2025/3/21 16:21
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/io/redis/MailBoxUserRedis"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/mail/MazeMail"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/mail/MazeMailCli"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/broadcastmailstatdb"
	"go.uber.org/zap"
	"math"
	"time"
)

func OnMailListQueryRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.MailListQueryRQ)
	res := rsMsg.(*MazeMailCli.MailListQueryRS)
	res.Header = req.Header
	res.ErrInfo = errors.NO_ERROR

	defer fkprometheus.DebugPMT("OnMailListQueryRQ")()

	userID := shardingID
	maxReceiptItem := &MazeMail.ReceiptItem{}
	token := req.GetToken()
	classType := req.GetClassType()
	res.ClassType = req.ClassType

	nowToken, err := MailBoxUserRedis.ReadMailBoxSequence(ctx, userID, classType)
	if err != nil {
		ctx.WarnWF("OnMailListQueryRQ read mail token fail", zap.Uint64("uid", userID), zap.Any("classType", classType), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		//return errors.MODULE_ERROR
		return nil
	}

	ctx.InfoWF("OnMailListQueryRQ with", zap.Any("req", req))
	for _, receiptItem := range req.GetReceiptItems() {
		if maxReceiptItem.GetSequence() < receiptItem.GetSequence() {
			maxReceiptItem = receiptItem
		}
	}
	offset := maxReceiptItem.GetSequence()

	if offset > 0 {
		offset -= 1
	} else {
		//第一次获取
		if nowToken == token {
			ctx.DebugWF("OnMailListQueryRQ mail token same", zap.Uint64("uid", userID), zap.Uint64("token", nowToken))
			res.Token = req.Token
			return
		}
		offset = math.MaxUint64
	}
	ctx.InfoWF("OnMailListQueryRQ with maxReceiptItem", zap.Any("offset", offset), zap.Any("maxReceiptItem", maxReceiptItem))

	mailList, needDelMailIDS, isEnd, err := getMailList(ctx, userID, classType, offset)
	if err != nil {
		ctx.ErrorWF("OnMailListQueryRQ getMailList err", zap.Int32("classType", classType),
			zap.Uint64("nowToken", nowToken),
		)
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	res.Mails = mailList
	if isEnd {
		res.IsEnd = proto.Int32(1)
	}
	if len(needDelMailIDS) > 0 {
		delMail(ctx, userID, classType, needDelMailIDS, "OnMailListQueryRQ")
	}

	res.Token = proto.Uint64(nowToken)

	ctx.InfoWF("OnMailListQueryRQ success", zap.Any("IsEnd", res.GetIsEnd()), zap.Any("MailsLen", len(res.Mails)))
	return nil
}

func getMailList(logger fklog.FKLogI, userID uint64, classType int32, beginOffset uint64) (resList []*MazeMail.MailInfo, needDelMailIDS []uint64, isEnd bool, err error) {
	defer fkprometheus.DebugPMT("getMailList")()

	var (
		now           = time.Now().Unix()
		indexes       []*MazeMail.ReceiptItem
		mails         []*MazeMail.MailInfo
		searchIndexes []uint64
		offset        = beginOffset
	)
	for loopIndex := 0; loopIndex < 6; loopIndex++ {
		// 获取信封索引
		indexes, err = mailboxsetdb.GetMailIndexes(logger, userID, classType, offset, mailunreaddb.GetMailIndexesLimitCount)
		if err != nil {
			logger.ErrorWF("getMailList GetMailIndexes err", zap.Int32("classType", classType),
				zap.Uint64("offset", offset), zap.Error(err),
			)
			return
		}
		if len(indexes) <= 0 {
			//结束了
			isEnd = true
			logger.InfoWF("getMailList indexes empty", zap.Int32("classType", classType), zap.Uint64("offset", offset),
				zap.Int("loopIndex", loopIndex),
			)
			break
		}

		searchIndexes = make([]uint64, 0, len(indexes))
		for _, index := range indexes {
			searchIndexes = append(searchIndexes, index.GetId())
		}
		logger.InfoWF("getMailList all index", zap.Uint64s("index", searchIndexes))

		// 根据索引查信封信息
		mails, err = MailBoxBodyRedis.GetMailList(logger, searchIndexes)
		if err != nil {
			logger.ErrorWF("getMailList GetMailList err", zap.Uint64s("searchIndexes", searchIndexes), zap.Error(err))
			return
		}

		for i, mail := range mails {
			logger.InfoWF("getMailList show mail", zap.Any("mail.Base", mail.GetBase()))
			if mail.GetBase() == nil {
				logger.InfoWF("getMailList invalid need del 1", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
				needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				continue
			}
			if mail.Base.GetMailId() == 0 {
				if i < len(searchIndexes) {
					logger.InfoWF("getMailList invalid need del 2", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
					needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				}
			} else {
				if mail.GetBase().GetExpires() > 0 && now > mail.GetBase().GetExpires() {
					logger.InfoWF("getMailList Expires need del", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
					needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				} else {
					//获取列表，不返回body
					mail.Body = nil
					mail.SubBody = nil
					mail.BattleRecordBody = nil
					//mail.RandomBox = nil
					if mail.Base.GetBroadcastFlag() != 0 {
						//广播信封
						hasRead, hasOperator, _, getErr := broadcastmailstatdb.GetUserBroadcastMailStat(logger, userID, searchIndexes[i])
						if getErr == nil {
							mail.MailHasRead = proto.Int32(hasRead)
							mail.HasOperator = proto.Int32(hasOperator)
							//if mail.Vote != nil && vote != nil {
							//	answer := &MazeMail.UserVoteAnswer{}
							//	getErr = proto.Unmarshal(vote, answer)
							//	if getErr != nil {
							//		logger.ErrorWF("GetUserBroadcastMailStat vote unmarshal err", zap.Error(err))
							//	} else {
							//		mail.Vote.UserAnswer = answer
							//	}
							//}
						}
					}
					//if mailSpecial := GMailSpecaialCfg.GetMailSpecaialConfig(mail.Base.GetExtType()); mailSpecial != nil {
					//	mail.Base.SpecialConfig = &MazeMail.SpecialAwardConfig{
					//		CostItem:   proto.Int32(mailSpecial.Cost_item),
					//		EffectItem: mailSpecial.Effect_item,
					//		ExtraRatio: proto.Int32(mailSpecial.Extra_ratio)}
					//}
					resList = append(resList, mail)
				}
			}
		}

		if len(indexes) < mailunreaddb.GetMailIndexesLimitCount {
			logger.InfoWF("getMailList is end empty", zap.Any("offset", offset), zap.Any("loopIndex", loopIndex))
			isEnd = true
			break
		}

		offset = indexes[len(indexes)-1].GetSequence()
		if offset == 0 {
			break
		} else {
			offset--
		}
	}
	logger.InfoWF("getMailList end",
		zap.Uint64("beginOffset", beginOffset),
		zap.Uint64("offset", offset),
		zap.Bool("isEnd", isEnd),
		zap.Int("MailsLen", len(resList)),
		zap.Uint64s("needDelMailIDS", needDelMailIDS),
	)
	return
}

func delMail(logger fklog.FKLogI, userID uint64, classType int32, mailIDS []uint64, callFuncName string) {
	if len(mailIDS) < 1 {
		return
	}
	go func() {
		logger.InfoWF("delMail req", zap.String("callFuncName", callFuncName), zap.Any("mailIDS", mailIDS))
		err := mailboxsetdb.DelMailBox(logger, userID, classType, mailIDS)
		if err != nil {
			logger.ErrorWF("delMail end", zap.String("callFuncName", callFuncName), zap.Any("mailIDS", mailIDS), zap.Error(err))
		}
		err = mailunreaddb.DelUnreadMail(logger, userID, classType, mailIDS)
		if err != nil {
			logger.ErrorWF("delMail DelUnreadMail error", zap.String("callFuncName", callFuncName), zap.Uint64s("mailIDS", mailIDS), zap.Error(err))
		}
	}()
}
