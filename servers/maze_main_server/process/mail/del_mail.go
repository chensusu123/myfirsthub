// @Author: ZhaoXiming 2025/3/21 16:34
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/mailutil"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/MailBoxBodyRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/broadcastmailstatdb"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/mailboxsetdb"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/mailunreaddb"
	"go.uber.org/zap"
	"math"
	"time"
)

var MaxMailLoopCount = 40 // 40* mailunreaddb.GetMailIndexesLimitCount 信封个数

func OnTotalMailDelRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.TotalMailDelRQ)
	res := rsMsg.(*MazeMailCli.TotalMailDelRS)
	res.ErrInfo = errors.NO_ERROR
	userID := shardingID
	logger := ctx.FKLogI
	res.Header = req.Header
	classType := req.GetClassType()
	res.ClassType = req.ClassType

	beginOffset := uint64(math.MaxUint64)
	ctx.InfoWF("OnTotalMailDelRQ begin", zap.Uint64("beginOffset", beginOffset), zap.Any("req", req))
	defer fkprometheus.DebugPMT("OnTotalMailDelRQ")()

	needDelMailIDS := make([]uint64, 0, 30)
	now := time.Now().Unix()
	offset := beginOffset
	for loopIndex := 0; loopIndex < MaxMailLoopCount; loopIndex++ {
		indexes, _ := mailboxsetdb.GetMailIndexes(logger, userID, classType, offset, mailunreaddb.GetMailIndexesLimitCount)
		if len(indexes) <= 0 {
			//结束了
			ctx.InfoWF("OnTotalMailDelRQ indexes empty", zap.Uint64("offset", offset), zap.Int32("classType", classType), zap.Any("loopIndex", loopIndex))
			break
		}
		searchIndexes := make([]uint64, 0)
		for _, index := range indexes {
			searchIndexes = append(searchIndexes, index.GetId())
			ctx.InfoWF("OnTotalMailDelRQ  index", zap.Any("index", index))
		}

		mails, _ := MailBoxBodyRedis.GetMailList(logger, searchIndexes)

		for i, mail := range mails {
			ctx.InfoWF("OnTotalMailDelRQ show mail", zap.Any("mail.Base", mail.GetBase()))
			if mail.GetBase() == nil {
				ctx.InfoWF("OnTotalMailDelRQ invalid need del 1", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
				needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				continue
			}
			if mail.Base.GetMailId() == 0 {
				if i < len(searchIndexes) {
					ctx.InfoWF("OnTotalMailDelRQ invalid need del 2", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
					needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				}
			} else {
				if mail.GetBase().GetExpires() > 0 && now > mail.GetBase().GetExpires() {
					ctx.InfoWF("OnTotalMailDelRQ Expires need del", zap.Any("mail.Base", mail.GetBase()), zap.Uint64("delMailID", searchIndexes[i]))
					needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
				} else {
					broadcastFlag := mail.Base.GetBroadcastFlag()
					var hasRead, hasOperator int32

					if broadcastFlag != 0 {
						hasRead, hasOperator, _, _ = broadcastmailstatdb.GetUserBroadcastMailStat(ctx, userID, searchIndexes[i])
					} else {
						hasRead = mail.GetMailHasRead()
						hasOperator = mail.GetHasOperator()
					}

					if hasRead != 1 {
						continue
					}
					if mailutil.IsMailTextOrReport(mail) {
						needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
					} else {
						if hasOperator == 1 {
							needDelMailIDS = append(needDelMailIDS, searchIndexes[i])
						}
					}
				}
			}
		}

		if len(indexes) < mailunreaddb.GetMailIndexesLimitCount {
			ctx.InfoWF("OnTotalMailDelRQ is end empty", zap.Any("offset", offset), zap.Any("loopIndex", loopIndex))
			break
		}

		offset = indexes[len(indexes)-1].GetSequence()
		if offset == 0 {
			break
		} else {
			offset--
		}
	}

	ctx.InfoWF("OnTotalMailDelRQ getMailList end",
		zap.Any("beginOffset", beginOffset),
		zap.Any("offset", offset),
		zap.Any("needDelMailIDS", needDelMailIDS),
	)

	if len(needDelMailIDS) <= 0 {
		return
	}
	syncDelMail(logger, userID, classType, needDelMailIDS, "TotalMailDelRQ")
	res.DelMailList = needDelMailIDS

	num, err := mailunreaddb.UnreadMailNum(logger, userID, classType)
	if err != nil {
		logger.ErrorWF("OnTotalMailDelRQ get unread mail num", zap.Error(err))
	} else {
		res.UnreadNum = proto.Int32(num)
	}

	// 获取最新的过期时间
	latestExpireTime, err := mailunreaddb.GetUnreaMailExpireTime(logger, userID, classType)
	if err != nil {
		logger.ErrorWF("OnTotalMailDelRQ load data error", zap.Error(err))
	} else if latestExpireTime > 0 {
		res.LatestExpireTime = proto.Int64(latestExpireTime)
	}

	return nil
}
