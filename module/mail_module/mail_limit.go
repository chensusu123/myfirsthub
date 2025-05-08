package mail_module

//import (
//	"errors"
//	"gitlab.ifreetalk.com/maze/maze_mail_server/common/mailutil"
//	"gitlab.ifreetalk.com/maze/maze_mail_server/db/maillimitfaulttolerantredis"
//	"gitlab.ifreetalk.com/maze/maze_mail_server/db/maillimitredis"
//	"sort"
//
//	"gitlab.ifreetalk.com/plate/excel/auto/GMailLimitConfigCfg"
//	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
//	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
//	"gitlab.ifreetalk.com/plate/protodef/MailBoxSvr"
//	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
//	"go.uber.org/zap"
//)
//
//// 东八区时间
//var DongBaQu int64 = 8 * 3600
//
//var ErrMailLimit = errors.New("mail reach limit")
//var ErrCheckMailLimit = errors.New("check mail limit fail")
//
//// 检查信封是否达到次数上限了
//func reachLimit(logger fklog.FKLogI, userId uint64, mail *MazeMail.MailInfo) (bool, error) {
//	awardMail := mailutil.IsMailHasAward(mail)
//	if !awardMail {
//		return false, nil
//	}
//
//	opType := mail.Base.GetExtType()
//	if opType == 0 {
//		return false, nil
//	}
//	createTime := mail.Base.GetTime()
//	cfg := GMailLimitConfigCfg.GetMailLimitConfigConfig(opType)
//	if cfg == nil {
//		logger.DebugWF("get mail limit config empty", zap.Int32("opType", opType))
//		return false, nil
//	}
//
//	qietian := cfg.Is_clear_at_0
//	duration := cfg.Time
//	limitNum := cfg.Num
//
//	sendTime, err := maillimitredis.GetMailSendTime(userId, opType, limitNum)
//	if err != nil {
//		logger.ErrorWF("get mail send Time fail", zap.Uint64("uid", userId), zap.Int32("opType", opType), zap.Error(err))
//		return true, ErrCheckMailLimit
//	}
//	// 最早一次发送时间
//	var lastSend int64
//	if len(sendTime) > 0 {
//		lastSend = sendTime[len(sendTime)-1]
//	}
//
//	beginTime := createTime - int64(duration)
//
//	if qietian == 1 {
//		todyBegin := createTime - (createTime+DongBaQu)%86400
//		if todyBegin > beginTime {
//			beginTime = todyBegin
//		}
//	}
//
//	sendNum := len(sendTime)
//	index := sort.Search(sendNum, func(index int) bool {
//		if sendTime[index] < beginTime {
//			return true
//		}
//		return false
//	})
//	logger.DebugWF("mail_limit_check", zap.Uint64("userId", userId), zap.Int32("opType", opType), zap.Int64("lastSend", lastSend), zap.Int64("beginTime", beginTime), zap.Int("index", index))
//
//	if index != sendNum {
//		sendTime = sendTime[:index]
//	}
//
//	if int32(len(sendTime)) >= limitNum {
//		logger.WarnWF("limit fail", zap.Uint64("userId", userId), zap.Int32("opType", opType), zap.Int32("cfg", limitNum), zap.Int("nowNum", len(sendTime)))
//		maillimitfaulttolerantredis.Push(logger, &MailBoxSvr.SvrSendMailRQ{ToUser: proto.Uint64(userId), MailInfo: mail})
//		return true, ErrMailLimit
//	}
//
//	// todo: 如果发生并发，可能会多删除记录
//	if lastSend != 0 {
//		if createTime-lastSend > int64(duration) {
//			if index != 0 {
//				index--
//			}
//			maillimitredis.TrimMailSendTime(logger, userId, opType, int32(index))
//		}
//	}
//	return true, nil
//}
