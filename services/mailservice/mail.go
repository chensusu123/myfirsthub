package mailservice

import (
	"context"
	"sort"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"

	"maze_game_server/common/errors"
	"maze_game_server/common/tradeno"
	"maze_game_server/model/mailmodel"
	"maze_game_server/usecase/online"

	"go.uber.org/zap"
)

func (s service) GetMailListByLabel(ctx context.Context, userId uint64, label, start, end int32) (mailList []*mailmodel.MailInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	model, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetMailListByLabel NewMailModel fail", zap.Error(err))
		return
	}

	delList := make([]uint64, 0)
	mailList = make([]*mailmodel.MailInfo, 0)
	for _, info := range model.MailMap[label] {
		if checkMailExpire(ctx, info) {
			delList = append(delList, info.ID)
			continue
		}
		mailList = append(mailList, info)
	}

	sort.Slice(mailList, func(i, j int) bool {
		return mailList[i].SendTime > mailList[j].SendTime
	})

	mailList = mailList[start:end]

	// 获取的时候检测过期邮件
	if len(delList) > 0 {
		for _, mailId := range delList {
			delete(model.MailMap[label], mailId)
		}
		err = model.Save(ctx, userId)
		if err != nil {
			logger.CtxError(ctx, "GetMailList del mail Save fail", zap.Error(err))
			return mailList, err
		}
	}

	return
}

func (s service) SendMail(ctx context.Context, title, context, senderName string, label int32, reciverId uint64, attachments []*mailmodel.Attachment, expireTime int64) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	if expireTime > 0 && expireTime < time.Now().Unix() {
		logger.CtxInfo(ctx, "SendMail expire time is less than now")
		return errors.New("expire time is less than now")
	}

	if len(title) == 0 || title == "" {
		logger.CtxInfo(ctx, "SendMail title is empty")
		return errors.New("SendMail title is empty")
	}

	if len(context) == 0 || context == "" {
		logger.CtxInfo(ctx, "SendMail context is empty")
		return errors.New("SendMail context is empty")
	}

	if len(senderName) == 0 || senderName == "" {
		logger.CtxInfo(ctx, "SendMail senderName is empty")
		senderName = "系统"
	}

	nMail := newMail(ctx, title, context, senderName, label, reciverId, attachments, expireTime)
	logger.CtxInfo(ctx, "SendMail newMail success", zap.Any("nMail", nMail))

	mailModel, err := mailmodel.NewMailModel(ctx, reciverId)
	if err != nil {
		logger.CtxError(ctx, "SendMail NewMailModel fail", zap.Error(err))
		return
	}

	infoMap, ok := mailModel.MailMap[label]
	if !ok {
		infoMap = make(map[uint64]*mailmodel.MailInfo)
	}
	_, ok = infoMap[nMail.ID]
	if !ok {
		infoMap[nMail.ID] = nMail
	} else {
		logger.CtxInfo(ctx, "SendMail mail already exists")
		return errors.New("SendMail mail already exists")
	}
	mailModel.MailMap[label] = infoMap

	err = mailModel.Save(ctx, reciverId)
	if err != nil {
		logger.CtxError(ctx, "SendMail Save fail", zap.Error(err))
		return err
	}

	return
}

func (s service) ReadMail(ctx context.Context, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	if mailId == 0 {
		logger.CtxInfo(ctx, "ReadMail mailId is empty")
		return nil, errors.New("mailId is null")
	}

	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReadMail NewMailModel fail", zap.Error(err))
		return
	}

	mailInfo = getMailById(ctx, userId, mailId, label)
	if mailInfo == nil {
		logger.CtxInfo(ctx, "ReadMail getMailById is nil")
		return nil, errors.New("mail not found")
	}

	if mailInfo.IsRead {
		logger.CtxInfo(ctx, "ReadMail status already read", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil
	}
	mailInfo.IsRead = true

	mailModel.MailMap[label][mailId] = mailInfo
	err = mailModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReadMail Save fail", zap.Error(err))
		return
	}
	return
}

func (s service) ReadAllMail(ctx context.Context, userId uint64, label int32) (mailList []*mailmodel.MailInfo, err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReadAllMail NewMailModel fail", zap.Error(err))
		return
	}

	mailList = make([]*mailmodel.MailInfo, 0)
	for _, info := range mailModel.MailMap[label] {
		if info.IsRead {
			continue
		}
		info.IsRead = true
		mailList = append(mailList, info)
	}

	err = mailModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "ReadAllMail Save fail", zap.Error(err))
		return
	}

	return
}

func (s service) GetMailAttachment(ctx context.Context, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailInfo = getMailById(ctx, userId, mailId, label)
	if mailInfo == nil {
		logger.CtxInfo(ctx, "GetMailAttachment getMailById is nil")
		return mailInfo, nil, errors.New("mail not found")
	}

	if len(mailInfo.Attachments) == 0 {
		logger.CtxInfo(ctx, "GetMailAttachment Attachments is nil")
		return mailInfo, nil, errors.New("没有可领取的附件")
	}

	if checkMailExpire(ctx, mailInfo) {
		logger.CtxInfo(ctx, "GetAllMailAttachment already expire", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil, errors.New("<UNK>")
	}

	if mailInfo.IsGetAttach {
		logger.CtxInfo(ctx, "GetMailAttachment already claimed", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil, errors.New("mail not found")
	}

	attachments = append(attachments, mailInfo.Attachments...)

	mailInfo.IsGetAttach = true
	mailInfo.IsRead = true
	return
}

func (s service) GetAllMailAttachment(ctx context.Context, userId uint64, label int32) (mailList []*mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetAllMailAttachment NewMailModel fail", zap.Error(err))
		return
	}

	mailList = make([]*mailmodel.MailInfo, 0)

	for _, info := range mailModel.MailMap[label] {

		if checkMailExpire(ctx, info) {
			logger.CtxInfo(ctx, "GetAllMailAttachment already expire", zap.Any("mailInfo", info))
			continue
		}
		if len(info.Attachments) == 0 {
			logger.CtxInfo(ctx, "GetAllMailAttachment Attachments is nil", zap.Any("mailInfo", info))
			continue
		}

		if info.IsGetAttach {
			logger.CtxInfo(ctx, "GetAllMailAttachment already GetAttach", zap.Any("info", info))
			continue
		}

		attachments = append(attachments, info.Attachments...)
		info.IsGetAttach = true
		info.IsRead = true
		mailList = append(mailList, info)
	}

	return
}

func (s service) DelMail(ctx context.Context, userId, mailId uint64, label int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "DelMail NewMailModel fail", zap.Error(err))
		return
	}

	mailInfo := getMailById(ctx, userId, mailId, label)
	if mailInfo == nil {
		logger.CtxInfo(ctx, "DelMail getMailById is nil")
		return errors.New("mail not found")
	}

	if !mailInfo.IsGetAttach {
		logger.CtxInfo(ctx, "DelMail mail not GetAttach", zap.Any("mailInfo", mailInfo))
		return errors.New("邮件附件未领取")
	}

	delete(mailModel.MailMap[label], mailId)

	err = mailModel.Save(ctx, userId)
	if err != nil {
		return err
	}
	return
}

func (s service) DelAllMail(ctx context.Context, userId uint64, label int32) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "DelAllMail NewMailModel fail", zap.Error(err))
		return
	}

	delList := make([]uint64, 0)
	for _, info := range mailModel.MailMap[label] {
		if !info.IsGetAttach || !info.IsRead {
			continue
		}
		delList = append(delList, info.ID)
	}

	for _, mailId := range delList {
		delete(mailModel.MailMap[label], mailId)
	}

	err = mailModel.Save(ctx, userId)
	if err != nil {
		return err
	}

	return
}

func (s service) PushMailToReciver(ctx context.Context, userId uint64, packetType uint16, v interface{}) {
	logger := fklog.ContextAppLogger(ctx)
	isOnline := online.IsOnline(userId)
	if !isOnline {
		logger.CtxInfo(ctx, "PushMailToReciver reciver is not online", zap.Uint64("userId", userId))
		return
	}

	err := online.ClusterPush(ctx, userId, packetType, v)
	if err != nil {
		logger.CtxError(ctx, "PushMailToReciver fail", zap.Error(err))
	}
}

func (s service) GetMailAttachmentAfter(ctx context.Context, userId uint64, mailList []*mailmodel.MailInfo) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetMailAttachmentAfter NewMailModel fail", zap.Error(err))
		return
	}

	for _, info := range mailList {
		mailModel.MailMap[info.Label][info.ID].IsRead = true
		mailModel.MailMap[info.Label][info.ID].IsGetAttach = true
		// mailModel.MailMap[info.ID].Status = int32(MailStatusReadClaimed)
	}

	err = mailModel.Save(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "GetMailAttachmentAfter Save fail", zap.Error(err))
		return
	}
	return nil
}

func newMail(ctx context.Context, title, context, senderName string, label int32, reciverID uint64, attachments []*mailmodel.Attachment, expireTime int64) *mailmodel.MailInfo {
	if len(attachments) == 0 || attachments == nil {
		attachments = make([]*mailmodel.Attachment, 0)
	}
	return &mailmodel.MailInfo{
		ID:          tradeno.GetTradeNum(),
		Title:       title,
		Content:     context,
		Label:       label,
		Attachments: attachments,
		Sender:      senderName,
		ReciverID:   reciverID,
		IsRead:      false,
		IsGetAttach: false,
		SendTime:    time.Now().Unix(),
		ExpireTime:  expireTime,
	}
}

func getMailById(ctx context.Context, userId, mailId uint64, label int32) *mailmodel.MailInfo {
	logger := fklog.ContextAppLogger(ctx)
	if mailId == 0 {
		return nil
	}
	mailModel, err := mailmodel.NewMailModel(ctx, userId)
	if err != nil {
		logger.CtxError(ctx, "getMailById NewMailModel fail", zap.Error(err))
		return nil
	}
	mailMap, ok := mailModel.MailMap[label]
	if !ok {
		logger.CtxInfo(ctx, "getMailById label mail is empty", zap.Int32("label", label))
		return nil
	}
	info, ok := mailMap[mailId]
	if !ok {
		logger.CtxInfo(ctx, "getMailById mail not found", zap.Int32("label", label), zap.Uint64("mailId", mailId))
		return nil
	}
	return info
}

func checkMailExpire(ctx context.Context, mailInfo *mailmodel.MailInfo) bool {
	logger := fklog.ContextAppLogger(ctx)
	if mailInfo.ExpireTime == 0 {
		return false
	}
	if mailInfo.ExpireTime <= time.Now().Unix() {
		logger.CtxInfo(ctx, "checkMailExpire mailInfo.ExpireTime <= time.Now().Unix()", zap.Any("mailInfo", mailInfo))
		return true
	}
	return false
}
