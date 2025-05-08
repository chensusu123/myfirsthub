package mail_module

import (
	"context"
	"gitlab.ifreetalk.com/maze/maze_mail_server/common/tradeno"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/MailBoxBodyRedis"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/MailBoxUserRedis"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/broadcastmailrecordredis"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/mailboxsetdb"
	"gitlab.ifreetalk.com/maze/maze_mail_server/db/mailunreaddb"
	"gitlab.ifreetalk.com/maze/maze_mail_server/kafka/mailRecord"
	"gitlab.ifreetalk.com/maze/maze_mail_server/module/checkexpire"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/io_interface/frontcache/UserRuntimeFC"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/MustArriveRedis2"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"gitlab.ifreetalk.com/plate/protodef/MazeMailCli"
	"time"

	"go.uber.org/zap"
)

const (
	MailID = 000 // TODO 信封ID包
)

//发送邮件
func SendMail2User(ctx context.Context, logger fklog.FKLogI, toUserID uint64, mail *MazeMail.MailInfo) (uint64, error) {
	baseInfo := mail.Base
	if baseInfo == nil {
		logger.ErrorWF("SendMail2User mail base not fill", zap.Uint64("toUserID", toUserID), zap.Any("mailInfo", mail))
		return 0, errors.ARGS_NOT_MATCH
	}
	broadcastFlag := baseInfo.GetBroadcastFlag()
	classType := baseInfo.GetClassType()

	contentType := baseInfo.GetMailContentType()
	if baseInfo.GetFromServer() == nil || contentType < 1 {
		logger.ErrorWF("SendMail2User from_server or content_type not fill", zap.Uint64("toUserID", toUserID), zap.Any("mail", mail))
		return 0, errors.ARGS_NOT_MATCH
	}

	if mail.GetTitle() == "" {
		logger.ErrorWF("SendMail2User mail title not fill", zap.Uint64("toUserID", toUserID), zap.Any("mailInfo", mail))
		return 0, errors.ARGS_NOT_MATCH
	}

	guid := baseInfo.GetMailId()
	if guid == 0 {
		if broadcastFlag != 0 {
			return 0, errors.ARGS_NOT_MATCH
		}

		// 发送非广播信封，不带信封id 需要做处理
		logger.WarnWF("SendMail2User no mailID", zap.Uint64("toUserID", toUserID), zap.Any("mail", mail))
		guid = tradeno.GetTradeNoMaker().MakeTradeNo()
		if guid == 0 {
			logger.ErrorWF("SendMail2User call GenerateID failed", zap.Uint64("toUserID", toUserID), zap.Any("mailInfo", mail.GetBase()))
			return 0, errors.ERR_MAIL_SEND
		}
		baseInfo.MailId = proto.Uint64(guid)
	} else {
		//判断信封是否存在
		isExist, err := mailboxsetdb.IsMailBoxExist(logger, toUserID, classType, guid)
		if err != nil {
			logger.ErrorWF("SendMail2User tradeNum check failed", zap.Error(err), zap.Uint64("toUserID", toUserID), zap.Any("mail", mail))
			return 0, errors.ERR_MAIL_SEND
		}

		if isExist {
			logger.WarnWF("SendMail2User tradeNum already exist", zap.Uint64("uid", toUserID), zap.Any("mail", mail))
			// TODO: 可能没有流水
			return 0, errors.ERR_MAIL_ALREADY_EXIST
		}
	}

	if baseInfo.GetExpires() <= 0 {
		baseInfo.Expires = proto.Int64(time.Now().Unix() + 86400*7)
	}

	// 内网过期时间改为3天
	if checkexpire.IsInner() {
		baseInfo.Expires = proto.Int64(checkexpire.FixExpire(baseInfo.GetExpires()))
	}

	//needLimitCheck, err := reachLimit(logger, toUserID, mail)
	//if err != nil {
	//	logger.WarnWF("SendMail2User reachLimit error", zap.Error(err), zap.Uint64("toUserID", toUserID), zap.Any("mail", mail))
	//	return 0, err
	//}

	seq, err := MailBoxUserRedis.GetMailBoxSequence(logger, toUserID, classType)
	if err != nil {
		logger.ErrorWF("SendMail2User get mailbox sequence err", zap.Uint64("toUserID", toUserID), zap.Int32("classType", classType), zap.Error(err))
		return 0, errors.COMMON_ERROR_TIPS
	}

	baseInfo.Sequence = proto.Uint64(seq)
	if len(mail.Annex) > 0 {
		mail.HasAnnex = proto.Int32(1)
	}
	baseInfo.Time = proto.Int64(time.Now().Unix())
	// 感谢信使用from_user
	if mail.Base.GetClassType() == 0 {
		mail.FromUser = nil
	}
	if broadcastFlag == 0 {
		// 广播信息在函数外调用保存
		err = MailBoxBodyRedis.AddMail(logger, mail)
		if err != nil {
			logger.ErrorWF("SendMail2User add mail body", zap.Uint64("uid", toUserID), zap.Any("mail", mail), zap.Error(err))
			return 0, errors.COMMON_ERROR_TIPS
		}
	}

	err = mailboxsetdb.AddMailBox(logger, toUserID, classType, guid, seq)
	if err != nil {
		logger.ErrorWF("SendMail2User add mail index", zap.Uint64("uid", toUserID), zap.Any("mail", mail), zap.Error(err))
		err = mailboxsetdb.DelMailBox(logger, toUserID, classType, []uint64{guid})
		if err != nil {
			logger.ErrorWF("SendMail2User remove fail mail index", zap.Uint64("uid", toUserID), zap.Any("mail", mail), zap.Error(err))
		}
		return 0, errors.DB_SAVE_ERROR
	}
	if mail.GetHasOperator() == 0 && mail.GetMailHasRead() == 0 {
		_ = mailunreaddb.AddUnreadMail(logger, toUserID, classType, guid, baseInfo.GetExpires())
	}
	num, err := mailunreaddb.UnreadMailNum(logger, toUserID, classType)

	sendPack := &MazeMailCli.MailNotifyID{Mail: mail}
	sendPack.ClassType = proto.Int32(classType)
	if err != nil {
		logger.ErrorWF("SendMail2User get unread num", zap.Error(err), zap.Any("mail", mail))
	} else {
		sendPack.UnreadNum = proto.Int32(num)
	}
	// 获取最新的过期时间 //TODO 考虑废弃
	latestExpireTime, loadErr := mailunreaddb.GetUnreaMailExpireTime(logger, toUserID, classType)
	if loadErr != nil {
		logger.WarnWF("SendMail2User GetUnreaMailExpireTime error", zap.Error(loadErr), zap.Any("mail", mail))
	} else if latestExpireTime > 0 {
		sendPack.LatestExpireTime = proto.Int64(latestExpireTime)
	}

	//// 盟运翻倍信息
	//if mailSpecial := GMailSpecaialCfg.GetMailSpecaialConfig(mail.Base.GetExtType()); mailSpecial != nil {
	//	mail.Base.SpecialConfig = &MazeMail.SpecialAwardConfig{
	//		CostItem:   proto.Int32(mailSpecial.Cost_item),
	//		EffectItem: mailSpecial.Effect_item,
	//		ExtraRatio: proto.Int32(mailSpecial.Extra_ratio)}
	//
	//	ttl := int32(86401)
	//	windowId := &PopupWindow.SeaPopupWindowID{
	//		Type: proto.Uint32(2),
	//		ShipBerthAward: &PopupWindow.ShipBerthStatusData{
	//			Type:   proto.Uint32(62),
	//			UserId: proto.Uint64(toUserID),
	//			MailId: mail.Base.MailId,
	//		},
	//	}
	//
	//	if len(mail.Annex) > 0 {
	//		windowId.ShipBerthAward.Awards = mail.Annex[0].Items
	//	}
	//
	//	// 盟运发送奖励电视包，原来由王阔发的
	//	MustArriveRedis2.SendArrivePacketWithTtl(logger, toUserID, 13108, ttl, windowId)
	//}

	//esId, err := UserRuntimeFC.GetUserEsId(logger, toUserID, 30)
	//if err == nil {
	//	if esId != 0 {
	//		MustArriveRedis2.SendArrivePacketDebug(logger, toUserID, MailID, sendPack)
	//	} else {
	//		//离线推送
	//		OfflineNotify(ctx, logger, toUserID, baseInfo.GetExtType(), mail.GetOfflinePushText())
	//	}
	//} else {
	//	MustArriveRedis2.SendArrivePacketDebug(logger, toUserID, MailID, sendPack)
	//	OfflineNotify(ctx, logger, toUserID, baseInfo.GetExtType(), mail.GetOfflinePushText())
	//}

	//notifyEquipStone(logger, toUserID, mail)

	PushMailSend(ctx, logger, toUserID, mail, mailRecord.NoErr)
	//if needLimitCheck {
	//	maillimitredis.PushMailTime(logger, toUserID, baseInfo.GetExtType(), baseInfo.GetTime())
	//}

	if broadcastFlag != 0 {
		_ = broadcastmailrecordredis.RecordUser(logger, guid, toUserID)
	}

	// 信封ID包
	esId, err := UserRuntimeFC.GetUserEsId(logger, toUserID, 30)
	if err != nil {
		_ = MustArriveRedis2.SendArrivePacketDebug(logger, toUserID, MailID, sendPack)
	} else {
		if esId != 0 {
			_ = MustArriveRedis2.SendArrivePacketDebug(logger, toUserID, MailID, sendPack)
		}
	}

	return guid, nil
}

////更新邮件附件
//func UpdateMailAnnex(ctx context.Context, logger fklog.FKLogI, userId uint64, mail *MazeMail.MailInfo, items []*MazeCommon.MazeItem) (bool, error) {
//	logger.InfoWF("updateMailInfo begin", zap.Any("mail", mail), zap.Uint64("uid", userId))
//
//	completed := true
//	for _, annex := range mail.Annex {
//		if annex.GetHasOpen() != 1 {
//			completed = false
//			break
//		}
//	}
//
//	if completed {
//		mail.HasOperator = proto.Int32(1)
//		mail.MailHasRead = proto.Int32(1)
//	}
//	classType := mail.Base.GetClassType()
//
//	sendPack := &MazeMailCli.MailNotifyID{Mail: mail}
//
//	if completed {
//		//处理未读
//		mailunreaddb.DelUnreadMail(logger, userId, classType, []uint64{mail.Base.GetMailId()})
//
//		num, err := mailunreaddb.UnreadMailNum(logger, userId, classType)
//		if err != nil {
//			logger.ErrorWF("get unread mail num", zap.Error(err))
//		} else {
//			sendPack.UnreadNum = proto.Int32(num)
//		}
//		expireTime, loadErr := mailunreaddb.GetUnreaMailExpireTime(logger, userId, classType)
//		if loadErr == nil && expireTime > 0 {
//			sendPack.LatestExpireTime = proto.Int64(expireTime)
//		}
//	}
//
//	PushMailAwardLog(logger, userId, mail, items)
//
//	err := MailBoxBodyRedis.AddMail(logger, mail)
//	if err != nil {
//		logger.ErrorWF("add mail fail", zap.Any("mail", mail), zap.Error(err))
//		return false, err
//	}
//
//	MustArriveRedis2.SendArrivePacket(userId, packet_type.DEF_PROTO_MAIL_BOX_CLI_MAIL_NOTIFY_ID, sendPack)
//	logger.InfoWF("UpdateMailAnnex end", zap.Any("toUserID", userId), zap.Any("mail", mail.GetBase()))
//	return completed, nil
//}

//func notifyEquipStone(logger fklog.FKLogI, userId uint64, mail *MazeMail.MailInfo) {
//	var chgType int32
//	switch mail.GetBase().GetMailContentType() {
//	case int32(MazeMail.MailContentType_CONTENT_AWARD),
//		int32(MazeMail.MailContentType_CONTENT_FAMILY_TEAM),
//		int32(MazeMail.MailContentType_CONTENT_GIFT_PACK),
//		int32(MazeMail.MailContentType_CONTENT_AWARD_AND_INVITE_FRIENT),
//		int32(MazeMail.MailContentType_CONTENT_TRAIN_BOX):
//		chgType = normalAward(mail.Annex, mail.EnemyAward)
//	case int32(MazeMail.MailContentType_CONTENT_TRADE_BOX):
//		if mail.TradeBox.GetIsDiamond() != 0 {
//			// 文本类型
//			return
//		}
//		chgType = tradeBox(mail.TradeBox, mail.EnemyAward)
//	default:
//		return
//	}
//
//	logger.DebugWF("equip stone notify", zap.Uint64("uid", userId), zap.Int32("info", chgType))
//
//	if chgType == 0 {
//		return
//	}
//
//	info := &MailEquipStoneNotifyKafka.MailStoneMsg{
//		UserID:  userId,
//		ChgType: chgType,
//		Time:    time.Now().Unix(),
//	}
//	MailEquipStoneNotifyKafka.PushMailFlowLog(logger, info)
//}

//func tradeBox(box *MazeMail.TradeBoxInfo, enemy *MazeMail.EnemyAward) int32 {
//	var infos []*MazeCommon.MazeItem
//	infos = append(infos, box.Award...)
//	if enemy.GetHasAward() && len(enemy.GetAward()) > 0 {
//		infos = append(infos, enemy.GetAward()...)
//	}
//	if len(box.VipAward) > 0 {
//		infos = append(infos, box.VipAward...)
//	}
//	return hasShiLiZhi(infos)
//}

//func normalAward(annexes []*MazeMail.AnnexItem, enemy *MazeMail.EnemyAward) int32 {
//	var infos []*MazeCommon.MazeItem
//	for _, annex := range annexes {
//		if annex.GetHasOpen() == 0 {
//			infos = append(infos, annex.Items...)
//		}
//	}
//	if enemy.GetHasAward() && len(enemy.GetAward()) > 0 {
//		infos = append(infos, enemy.GetAward()...)
//	}
//	return hasShiLiZhi(infos)
//}
//
//func hasShiLiZhi(items []*MazeCommon.MazeItem) int32 {
//	var cType int32
//	for _, item := range items {
//		if item.GetItemId() == 30000001 {
//			// 势力值
//			cType = 2
//		}
//	}
//	return cType
//}
