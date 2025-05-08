// @Author: ZhaoXiming 2025/3/21 16:22
// @Desc:

package mail

import (
	"gitlab.ifreetalk.com/maze/maze_game_server/common/cache/simCache"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/cgkargs"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mail_module"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"gitlab.ifreetalk.com/plate/protodef/MazeMailCli"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/MailBoxBodyRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/broadcastmailstatdb"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/mailboxsetdb"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/mailunreaddb"
	"gitlab.ifreetalk.com/plate/protodef/Common"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"go.uber.org/zap"
	"time"
)

var awardCache = simCache.NewCache()

func OnAnnexOpenRQ(ctx fknet.TCPContext, shardingID uint64, rqMsg proto.Message, rsMsg proto.Message) (err error) {
	req := rqMsg.(*MazeMailCli.AnnexOpenRQ)
	res := rsMsg.(*MazeMailCli.AnnexOpenRS)
	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header
	res.Mail = req.Mail

	userID := shardingID

	mailID := req.GetMail().GetId()
	ctx.InfoWF("OnAnnexOpenRQ with", zap.Any("req ", req))
	defer fkprometheus.DebugPMT("OnAnnexOpenRQ")()

	// 限制重复请求
	if awardCache.SetCache(userID, struct{}{}) {
		res.ErrInfo = errors.NewCodeError(62091, "").ToInfo()
		return
	}
	defer awardCache.DelCache(userID)

	mail, err := MailBoxBodyRedis.GetMail(ctx, mailID)
	if err != nil || mail == nil || mail.Base == nil || mail.Base.GetMailId() == 0 {
		ctx.WarnWF("OnAnnexOpenRQ call GetMail failed", zap.Any("Mail", req.GetMail()))
		res.ErrInfo = errors.ERR_MAIL_NOT_FOUND.ToInfo()
		return nil
	}
	ctx.InfoWF("OnAnnexOpenRQ req", zap.Any("mail", mail))

	classType := mail.Base.GetClassType()

	exist, err := mailunreaddb.IsMailBoxOpen(ctx, userID, classType, mailID)
	if err != nil {
		ctx.WarnWF("OnAnnexOpenRQ check user mail exist fail", zap.Uint64("uid", userID), zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return nil
	}
	var notInUnread bool
	if !exist {
		exist, err := mailboxsetdb.IsMailBoxExist(ctx, userID, classType, mailID)
		if err != nil {
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			return nil
		}
		if !exist {
			// 不在个人信封列表
			ctx.WarnWF("OnAnnexOpenRQ unread mail set not exist ", zap.Uint64("uid", userID), zap.Int32("type", classType), zap.Uint64("id", mailID))
			res.ErrInfo = errors.ERR_ANNEX_HAS_OPEN.ToInfo()
			return nil
		}
		notInUnread = true
	}

	res.MailInfo = mail
	now := time.Now().Unix()
	var hasOperator int32
	var isBroadcast bool
	if mail.Base.GetBroadcastFlag() != 0 {
		isBroadcast = true
		// 广播消息
		//var voteData []byte
		_, hasOperator, _, err = broadcastmailstatdb.GetUserBroadcastMailStat(ctx, userID, mailID)
		if err != nil {
			ctx.WarnWF("OnAnnexOpenRQ load broadcast mail stat fail", zap.Any("Mail", req.GetMail()), zap.Error(err))
			res.ErrInfo = errors.MODULE_ERROR.ToInfo()
			//return errors.MODULE_ERROR
			return nil
		}
		//if mail.Vote != nil && len(voteData) > 0 {
		//	vote := &MazeMail.UserVoteAnswer{}
		//	err = proto.Unmarshal(voteData, vote)
		//	if err != nil {
		//		ctx.ErrorWF("OnAnnexOpenRQ unmarshal broadcast mail vote info fail", zap.Uint64("uid", userID), zap.Any("mail", mail), zap.Error(err))
		//	}
		//	mail.Vote.UserAnswer = vote
		//}
	} else {
		hasOperator = mail.GetHasOperator()
	}

	if hasOperator == 1 {
		ctx.WarnWF("OnAnnexOpenRQ already open", zap.Uint64("uid", userID), zap.Any("Mail", req.GetMail()))
		mail.HasOperator = proto.Int32(1)
		mail.MailHasRead = proto.Int32(1)
		res.MailInfo = mail
		res.ErrInfo = errors.ERR_ANNEX_HAS_OPEN.ToInfo()
		return nil
	}

	if notInUnread {
		ctx.ErrorWF("OnAnnexOpenRQ mail not in unread set openAnn ", zap.Uint64("uid", userID), zap.Int32("type", classType), zap.Uint64("id", mailID))
	}

	if mail.GetBase().GetExpires() < now {
		ctx.WarnWF("OnAnnexOpenRQ Expires", zap.Any("Mail", req.GetMail()))
		res.ErrInfo = errors.ERR_MAIL_EXPIRE.ToInfo()
		return nil
	}

	mail.HasOperator = proto.Int32(1)
	mail.MailHasRead = proto.Int32(1)
	res.MailInfo = mail

	// 计算奖励 设置领取状态
	items, openAnnex := filterLuckGiftPack(ctx, mail.Annex)

	if isBroadcast {
		err = broadcastmailstatdb.SetUseBroadcastMailState(ctx, userID, mailID, 1, 1, nil, mail.Base.GetExpires())
	} else {
		// 非广播信息保存
		// 已读、已领取信封设置过期时间是1天
		mail.Base.Expires = proto.Int64(now + cgkargs.AlreadyReadMailExpireTime)
		for _, annex := range mail.Annex {
			if annex.GetHasOpen() == 0 {
				annex.HasOpen = proto.Int32(1)
			}
		}
		err = MailBoxBodyRedis.AddMail(ctx, mail)
	}
	if err != nil {
		ctx.ErrorWF("OnAnnexOpenRQ openAnnex add mail", zap.Any("mail", mail), zap.Error(err))
		return err
	}

	if len(items) > 0 { // 只根据上面计算的有无奖励 决定是否领奖
		res.Annex = openAnnex
		receiveAward(ctx, userID, mail, items, req.Header, res)
	} else {
		res.ErrInfo = errors.ERR_MAIL_NOT_FOUND_ANNEX.ToInfo()
		return nil
	}

	//处理未读
	_ = mailunreaddb.DelUnreadMail(ctx, userID, classType, []uint64{mailID})
	num, err := mailunreaddb.UnreadMailNum(ctx, userID, classType)
	if err != nil {
		ctx.ErrorWF("OnAnnexOpenRQ get unread mail num", zap.Error(err))
	} else {
		res.UnreadNum = proto.Int32(num)
	}

	// 获取最新的过期时间
	latestExpireTime, err := mailunreaddb.GetUnreaMailExpireTime(ctx, userID, classType)
	if err != nil {
		ctx.WarnWF("OnAnnexOpenRQ load data error", zap.Error(err))
	} else if latestExpireTime > 0 {
		res.LatestExpireTime = proto.Int64(latestExpireTime)
	}

	ctx.InfoWF("OnAnnexOpenRQ success", zap.Any("Mail", res.Mail), zap.Any("Annex", res.Annex))
	return nil
}

func receiveAward(ctx fknet.TCPContext, userID uint64, mail *MazeMail.MailInfo, items []*MazeCommon.MazeItem, header *Common.PacketHeader, res *MazeMailCli.AnnexOpenRS) {
	opType := mail.Base.GetExtType()
	mailID := mail.Base.GetMailId()
	isBroadcast := false
	if mail.Base.GetBroadcastFlag() != 0 {
		isBroadcast = true
	}

	tradeNo := mailID
	if isBroadcast {
		tradeNo = tradeno.GetTradeNoMaker().MakeTradeNo()
	}

	errInfo, err := itemrpc.GatherItems(ctx, userID, header, opType, tradeNo, items...)
	if err != nil {
		ctx.ErrorWF("receiveAward DeductItems err",
			zap.Uint64("tradeNo", tradeNo),
			zap.Error(err),
		)
		res.ErrInfo = errors.NewCommonCodeError("add item err")
		return
	}
	if errInfo != nil {
		ctx.WarnWF("receiveAward DeductItems invalid", zap.Uint64("tradeNo", tradeNo), zap.Any("items", items),
			zap.Any("errorInfo", errInfo), zap.Error(err),
		)
		res.ErrInfo = errInfo
		return
	}

	//流水
	mail_module.PushMailAwardLog(ctx, userID, mail, items)

	return
}

// 聚合加道具的奖励信息,客户端票动画展示获得东西使用  // 考虑废弃
//func oneMailAggregationResults(err *MessageType.ErrorInfo, addRes *ItemSvr.AddItemRS) (annex *MazeMail.AnnexItem) {
//	if addRes == nil {
//		return
//	}
//
//	// 1、加成功并且没有道具转化
//	// 2、加道具失败，展示全部信息
//	if err != nil || len(addRes.ConvertItems) <= 0 {
//		return nil
//	}
//
//	// 加道具成功并且发生了道具转化
//	openAnn := &MazeMail.AnnexItem{}
//	openAnn.HasOpen = proto.Int32(1)
//	if len(addRes.SuccItems) > 0 {
//		openAnn.Items = append(openAnn.Items, addRes.SuccItems...)
//	}
//	if len(addRes.Items) > 0 {
//		for _, item := range addRes.Items {
//			for _, citem := range item.ChgItems {
//				openAnn.Items = append(openAnn.Items, &MazeCommon.MazeItem{
//					ItemId: citem.ItemId,
//					Count:  citem.Count,
//				})
//			}
//		}
//	}
//	return openAnn
//}

//// 取除rmb、派豆之外的奖励
//func getOtherAward(items []*MazeCommon.MazeItem) (res []*MazeCommon.MazeItem) {
//	for k, item := range items {
//		if item.GetItemId() == 1200001 || item.GetItemId() == 2600000 {
//			continue
//		}
//		res = append(res, items[k])
//	}
//	return
//}
//
//// 是否含有rmb奖励
//func hasRMBAward(items []*MazeCommon.MazeItem) (has bool, moneyCount int64) {
//	for _, item := range items {
//		if item.GetItemId() == 1200001 {
//			return true, item.GetCount()
//		}
//	}
//	return false, 0
//}

// 礼包数量，普通奖励，附件信息
func filterLuckGiftPack(logger fklog.FKLogI, annexs []*MazeMail.AnnexItem) ([]*MazeCommon.MazeItem, []*MazeMail.AnnexItem) {
	var items []*MazeCommon.MazeItem
	openAnnex := make([]*MazeMail.AnnexItem, 0, len(annexs))
	for _, annex := range annexs {
		if annex.GetHasOpen() == 0 {
			annex.HasOpen = proto.Int32(1)

			for _, item := range annex.GetItems() {
				items = append(items, item)
			}
			openAnnex = append(openAnnex, annex)
		}
	}
	return items, openAnnex
}

//var failSetAwardFlag = sync.Map{}
//
//func isMailAwardFailFlag(logger fklog.FKLogI, userId, guid uint64) bool {
//	info, ok := failSetAwardFlag.Load(fmt.Sprintf("%d:%d", userId, guid))
//	if !ok {
//		return false
//	}
//	mail, ok := info.(*MazeMail.MailInfo)
//	if !ok {
//		logger.ErrorWF("load inner fail award flag fail", zap.Uint64("userId", userId), zap.Any("mail", info))
//		return true
//	}
//	resetAwardMailFlag(logger, userId, mail)
//	return true
//}
//
//func setMailAwardFailFlag(userId, guid uint64, mail *MazeMail.MailInfo) {
//	failSetAwardFlag.Store(fmt.Sprintf("%d:%d", userId, guid), mail)
//}
//
//func resetAwardMailFlag(logger fklog.FKLogI, userId uint64, mail *MazeMail.MailInfo) {
//	mailId := mail.GetBase().GetMailId()
//	var err error
//	if mail.GetBase().GetBroadcastFlag() != 0 {
//		// 广播信封,设置已读、已操作状态
//		//var voteInfo proto.Message
//		err = broadcastmailstatdb.SetUseBroadcastMailState(logger, userId, mailId, 1, 1, nil, mail.Base.GetExpires())
//	} else {
//		// 非广播信息保存
//		err = MailBoxBodyRedis.AddMail(logger, mail)
//	}
//	if err != nil {
//		return
//	}
//	failSetAwardFlag.Delete(fmt.Sprintf("%d:%d", userId, mailId))
//	logger.WarnWF("process inner fail award flag success", zap.Uint64("userId", userId), zap.Uint64("mailId", mailId))
//}
