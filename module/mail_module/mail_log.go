package mail_module

import (
	"context"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
)

func PushMailSend(ctx context.Context, logger fklog.FKLogI, uid uint64, mail *MazeMail.MailInfo, errState int32) {
	// TODO
	//log := &mailRecord.MailSendLog{}
	//log.Time = time.Now().UnixNano() / 1000000
	//log.Expire = mail.Base.GetExpires()
	//log.ContentType = mail.Base.GetMailContentType()
	//log.MailId = mail.Base.GetMailId()
	//log.ClassType = mail.Base.GetClassType()
	//log.AwardType = mail.Base.GetExtType()
	//log.SubType = mail.Base.GetSubType()
	//log.ErrState = errState
	//
	//// 领奖类型,需要判断贸易(随机宝箱)
	//if mail.TradeBox != nil {
	//	annex := &MazeMail.AnnexItem{Items: mail.TradeBox.Award}
	//	log.AwardList = append(log.AwardList, annex)
	//} else {
	//	log.AwardList = mail.Annex
	//}
	//
	//if mail.EnemyAward.GetHasAward() && len(mail.EnemyAward.GetAward()) > 0 {
	//	annex := &MazeMail.AnnexItem{Items: mail.EnemyAward.Award}
	//	log.AwardList = append(log.AwardList, annex)
	//}
	//
	//if mail.FriendshipAward != nil {
	//	bAnnex := mail.FriendshipAward.BaseAward
	//	mAnnex := mail.FriendshipAward.MaterialAward
	//	if mAnnex != nil {
	//		bAnnex = append(bAnnex, mAnnex...)
	//	}
	//	log.AwardList = bAnnex
	//}
	//log.MailBody = ""
	//log.MailTitle = mail.GetTitle()
	//fromServer := mail.Base.GetFromServer()
	//if fromServer != nil {
	//	log.FromServer = string(fromServer)
	//} else {
	//	log.FromServer = ""
	//}
	//log.Uid = uid
	//log.OpType = 1
	//
	//mailRecord.PushMailFlowLog(logger, log)
}

// 领奖流水
func PushMailAwardLog(logger fklog.FKLogI, uid uint64, mail *MazeMail.MailInfo, items []*MazeCommon.MazeItem) {
	// TODO
	//log := &mailRecord.MailSendLog{}
	//log.Time = time.Now().UnixNano() / 1000000
	//log.Expire = mail.Base.GetExpires()
	//log.ContentType = mail.Base.GetMailContentType()
	//log.MailId = mail.Base.GetMailId()
	//log.ClassType = mail.Base.GetClassType()
	//log.AwardType = mail.Base.GetExtType()
	//log.SubType = mail.Base.GetSubType()
	//log.AwardList = []*MazeMail.AnnexItem{
	//	&MazeMail.AnnexItem{
	//		HasOpen: proto.Int32(1),
	//		Items:   items,
	//	}}
	//log.MailBody = ""
	//log.MailTitle = mail.GetTitle()
	//fromServer := mail.Base.GetFromServer()
	//if fromServer != nil {
	//	log.FromServer = string(fromServer)
	//} else {
	//	log.FromServer = ""
	//}
	//log.Uid = uid
	//log.OpType = 2
	//
	//mailRecord.PushMailFlowLog(logger, log)
}
