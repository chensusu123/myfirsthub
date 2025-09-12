package mailservice

import (
	"context"
	"maze_game_server/model/mailmodel"
)

type MailService interface {
	GetMailListByLabel(ctx context.Context, userId uint64, label, start, end int32) (mailList []*mailmodel.MailInfo, err error)
	SendMail(ctx context.Context, title, context, senderName string, label int32, reciverId uint64, Attachments []*mailmodel.Attachment, expireTime int64) (err error)
	ReadMail(ctx context.Context, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, err error)
	ReadAllMail(ctx context.Context, userId uint64, label int32) (mailList []*mailmodel.MailInfo, err error)
	GetMailAttachment(ctx context.Context, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetAllMailAttachment(ctx context.Context, userId uint64, label int32) (mailList []*mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetMailAttachmentAfter(ctx context.Context, userId uint64, mailList []*mailmodel.MailInfo) (err error)
	DelMail(ctx context.Context, userId, mailId uint64, label int32) (err error)
	DelAllMail(ctx context.Context, userId uint64, label int32) (err error)
	PushMailToReciver(ctx context.Context, userId uint64, packetType uint16, v interface{})
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
