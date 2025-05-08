/*
@Author: xiaobo
@Date: 2023/11/1 14:28
@Description:
*/

package mailutil

import (
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
)

//func IsVoteMail(mail *MazeMail.MailInfo) bool {
//	if mail.GetBase().GetMailContentType() == int32(MazeMail.MailContentType_CONTENT_VOTE) {
//		return true
//	}
//	return false
//}

// IsMailHasAward 是否有奖励
func IsMailHasAward(mail *MazeMail.MailInfo) bool {
	switch mail.GetBase().GetMailContentType() {
	case int32(MazeMail.MailContentType_CONTENT_AWARD):
		return true
	default:
		return false
	}
}

// IsMailTextOrReport 是否是文本类型
func IsMailTextOrReport(mail *MazeMail.MailInfo) bool {
	switch mail.GetBase().GetMailContentType() {
	case int32(MazeMail.MailContentType_CONTENT_TEXT),
		int32(MazeMail.MailContentType_CONTENT_TEXT_AND_AWARD):
		return true
	default:
		return false
	}
}

//// IsEffectItem 检查道具是否翻倍
//func IsEffectItem(items []int32, itemId int32) bool {
//	for _, item := range items {
//		if item == itemId {
//			return true
//		}
//	}
//	return false
//}

//// GetLevelByVersion 根据版本取不同等级 岛等级、人偶等级
//func GetLevelByVersion(logger fklog.FKLogI, uid uint64) (level int32, err error) {
//	userVersion := BreedVersionFC.GetUserVersion(logger, uid)
//	if BreedVersionFC.IsDollVersionEx(logger, userVersion) {
//		level, err = dollLevelFcWithRedis.GetDollLevelFC(logger, uid, 50)
//		if err != nil {
//			logger.ErrorWF("GetLevelByUserVersion GetDollLevelFc err", zap.Uint64("uid", uid), zap.Error(err))
//			return
//		}
//	} else {
//		level, err = IslandLevelFcWithRedis.GetUserLandSizeFC(logger, uid, 50)
//		if err != nil {
//			logger.ErrorWF("GetLevelByUserVersion GetUserLandSizeFC err", zap.Uint64("uid", uid), zap.Error(err))
//			return
//		}
//	}
//	return
//}
