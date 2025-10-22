package mailservice

import (
	"context"
	"maze_game_server/common/structsdef"
	"maze_game_server/model/mailmodel"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type MailService interface {
	// 基础邮件功能
	GetMailListByLabel(logger fklog.FKLogI, userId uint64, label, start, end int32) (mailList []*mailmodel.MailInfo, err error)
	SendMail(logger fklog.FKLogI, title, context, senderName string, label int32, reciverId uint64, Attachments []*mailmodel.Attachment, expireTime int64) (err error)
	ReadMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, err error)
	ReadAllMail(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, err error)
	GetMailAttachment(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetAllMailAttachment(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error)
	GetMailAttachmentAfter(logger fklog.FKLogI, userId uint64, mailList []*mailmodel.MailInfo) (err error)
	GetAllMailList(logger fklog.FKLogI, userId uint64) (mailMap map[int32][]*mailmodel.MailInfo, err error)
	DelMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (err error)
	DelAllMail(logger fklog.FKLogI, userId uint64, label int32) (err error)
	PushMailToReciver(logger fklog.FKLogI, userId uint64, packetType uint16, v interface{})

	// 群发邮件功能
	BatchSendMail(ctx context.Context, req *BatchMailRequest) (result *BatchMailResult, err error)
	BatchSendSystemMail(ctx context.Context, title, content, senderName string, receiverIds []uint64, useQueue bool) (result *BatchMailResult, err error)
	BatchSendRewardMail(ctx context.Context, title, content, senderName string, attachments []*mailmodel.Attachment, receiverIds []uint64, useQueue bool) (result *BatchMailResult, err error)
	BatchSendRichTextMail(ctx context.Context, title, content, senderName string, panelInfo *structsdef.MailPanelInfo, receiverIds []uint64, useQueue bool) (result *BatchMailResult, err error)

	// 删除限制功能
	CheckDeleteRestrictions(ctx context.Context, userId, mailId uint64, label int32) (result *DeleteRestrictionResult, err error)
	ForceDeleteMail(ctx context.Context, userId, mailId uint64, label int32, adminId uint64) (err error)
	BatchCheckDeleteRestrictions(ctx context.Context, userId uint64, mailIds []uint64, label int32) (results map[uint64]*DeleteRestrictionResult, err error)
	GetDeletableMails(ctx context.Context, userId uint64, label int32) (mails []*mailmodel.MailInfo, err error)

	// 附件领取功能
	ClaimMailAttachment(ctx context.Context, userId, mailId uint64, label int32) (result *AttachmentClaimResult, err error)
	ClaimAllMailAttachments(ctx context.Context, userId uint64, label int32) (result *AttachmentClaimResult, err error)
	GetClaimableAttachments(ctx context.Context, userId uint64, label int32) (mails []*mailmodel.MailInfo, err error)

	// 背包集成功能
	AddAttachmentsToBag(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) (result *BagIntegrationResult, err error)
	GetBagSpaceInfo(ctx context.Context, userId uint64) (currentSpace, remainingSpace int32, err error)
	CheckBagSpaceForAttachments(ctx context.Context, userId uint64, attachments []*mailmodel.Attachment) (hasEnoughSpace bool, remainingSpace int32, err error)
	BatchAddAttachmentsToBag(ctx context.Context, userId uint64, mailAttachments map[uint64][]*mailmodel.Attachment) (result *BagIntegrationResult, err error)

	// 高级邮件功能
	SendRichTextMail(ctx context.Context, title, content, senderName string, panelInfo *structsdef.MailPanelInfo, receiverId uint64, attachments []*mailmodel.Attachment, expireTime int64) (err error)
	SendVoteMail(ctx context.Context, req *VoteMailRequest) (err error)
	SendBattleReportMail(ctx context.Context, req *BattleReportMailRequest) (err error)
	SendGiftBoxMail(ctx context.Context, req *GiftBoxMailRequest) (err error)
	SendLinkMail(ctx context.Context, req *LinkMailRequest) (err error)
	BroadcastMail(ctx context.Context, req *BroadcastMailRequest) (result *BroadcastResult, err error)
	SendMailWithConfig(ctx context.Context, configId int32, userId uint64, wildcardData map[string]interface{}) (err error)
	GetMailMetrics(ctx context.Context) (metrics map[string]interface{}, err error)
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
