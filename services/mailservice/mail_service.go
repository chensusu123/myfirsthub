package mailservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/model/mailmodel"
)

type MailService interface {
	GetMailListByLabel(logger fklog.FKLogI, userId uint64, label, start, end int32) (mailList []*mailmodel.MailInfo, err error)
	SendMail(logger fklog.FKLogI, title, context, senderName string, label int32, reciverId uint64, Attachments []*mailmodel.Attachment, expireTime int64) (err error)
	ReadMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, err error)
	ReadAllMail(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, err error)
	GetMailAttachment(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetAllMailAttachment(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetMailAttachmentAfter(logger fklog.FKLogI, userId uint64, mailList []*mailmodel.MailInfo) (err error)
	DelMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (err error)
	DelAllMail(logger fklog.FKLogI, userId uint64, label int32) (err error)
	PushMailToReciver(logger fklog.FKLogI, userId uint64, packetType uint16, v interface{})
}

var GlobalMailService MailService

func init() {
	GlobalMailService = newMailService()
}

type service struct {
}

func newMailService() MailService {
	return &service{}
}

// 邮件状态枚举
type MailStatus int

const (
	MailStatusUnread      MailStatus = iota // 未读
	MailStatusRead                          // 已读未领取附件
	MailStatusReadClaimed                   // 已读已领取附件
	MailStatusDeleted                       // 已删除
)

type MailLabel int

const (
	MailLabelSystem   MailLabel = iota //系统
	MailLabelFamily                    //家族
	MailLabelAlliance                  //联盟
)

const mailShowCount = 30 //邮件初始展示数量
