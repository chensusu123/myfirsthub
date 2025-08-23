package mailservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/common/errors"
	"maze_game_server/common/tradeno"
	"maze_game_server/model/mailmodel"
	"maze_game_server/usecase/online"
	"sort"
	"time"
)

func (s service) GetMailListByLabel(logger fklog.FKLogI, userId uint64, label, start, end int32) (mailList []*mailmodel.MailInfo, err error) {
	model, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMailList NewMailModel fail", zap.Error(err))
		return
	}

	delList := make([]uint64, 0)
	mailList = make([]*mailmodel.MailInfo, 0)
	for _, info := range model.MailMap[label] {
		if checkMailExpire(logger, info) {
			delList = append(delList, info.ID)
			continue
		}
		mailList = append(mailList, info)
	}

	sort.Slice(mailList, func(i, j int) bool {
		return mailList[i].SendTime > mailList[j].SendTime
	})

	mailList = mailList[start:end]

	//获取的时候检测过期邮件
	if len(delList) > 0 {
		for _, mailId := range delList {
			delete(model.MailMap[label], mailId)
		}
		err = model.Save(logger, userId)
		if err != nil {
			logger.ErrorWF("GetMailList del mail Save fail", zap.Error(err))
			return mailList, err
		}
	}

	return
}

func (s service) SendMail(logger fklog.FKLogI, title, context, senderName string, label int32, reciverId uint64, attachments []*mailmodel.Attachment, expireTime int64) (err error) {
	if expireTime > 0 && expireTime < time.Now().Unix() {
		logger.InfoWF("SendMail expire time is less than now")
		return errors.New("expire time is less than now")
	}

	if len(title) == 0 || title == "" {
		logger.InfoWF("SendMail title is empty")
		return errors.New("SendMail title is empty")
	}

	if len(context) == 0 || context == "" {
		logger.InfoWF("SendMail context is empty")
		return errors.New("SendMail context is empty")
	}

	if len(senderName) == 0 || senderName == "" {
		logger.InfoWF("SendMail senderName is empty")
		senderName = "系统"
	}

	nMail := newMail(logger, title, context, senderName, label, reciverId, attachments, expireTime)
	logger.InfoWF("SendMail newMail success", zap.Any("nMail", nMail))

	mailModel, err := mailmodel.NewMailModel(logger, reciverId)
	if err != nil {
		logger.ErrorWF("SendMail NewMailModel fail", zap.Error(err))
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
		logger.ErrorWF("SendMail mail already exists")
		return errors.New("SendMail mail already exists")
	}
	mailModel.MailMap[label] = infoMap

	err = mailModel.Save(logger, reciverId)
	if err != nil {
		logger.ErrorWF("SendMail Save fail", zap.Error(err))
		return err
	}

	return
}

func (s service) ReadMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, err error) {
	if mailId == 0 {
		logger.ErrorWF("ReadMail mailId is empty")
		return nil, errors.New("mailId is null")
	}

	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("ReadMail NewMailModel fail", zap.Error(err))
		return
	}

	mailInfo = getMailById(logger, userId, mailId, label)
	if mailInfo == nil {
		logger.InfoWF("ReadMail getMailById is nil")
		return nil, errors.New("mail not found")
	}

	if mailInfo.IsRead {
		logger.InfoWF("ReadMail status already read", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil
	}
	mailInfo.IsRead = true

	mailModel.MailMap[label][mailId] = mailInfo
	err = mailModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("ReadMail Save fail", zap.Error(err))
		return
	}
	return
}

func (s service) ReadAllMail(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, err error) {
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("ReadAllMail NewMailModel fail", zap.Error(err))
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

	err = mailModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("ReadAllMail Save fail", zap.Error(err))
		return
	}

	return
}

func (s service) GetMailAttachment(logger fklog.FKLogI, userId, mailId uint64, label int32) (mailInfo *mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error) {
	mailInfo = getMailById(logger, userId, mailId, label)
	if mailInfo == nil {
		logger.InfoWF("GetMailAttachment getMailById is nil")
		return mailInfo, nil, errors.New("mail not found")
	}

	if len(mailInfo.Attachments) == 0 {
		logger.InfoWF("GetMailAttachment Attachments is nil")
		return mailInfo, nil, errors.New("没有可领取的附件")
	}

	if checkMailExpire(logger, mailInfo) {
		logger.InfoWF("GetAllMailAttachment already expire", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil, errors.New("<UNK>")
	}

	if mailInfo.IsGetAttach {
		logger.InfoWF("GetMailAttachment already claimed", zap.Any("mailInfo", mailInfo))
		return mailInfo, nil, errors.New("mail not found")
	}

	attachments = append(attachments, mailInfo.Attachments...)

	mailInfo.IsGetAttach = true
	mailInfo.IsRead = true
	return
}

func (s service) GetAllMailAttachment(logger fklog.FKLogI, userId uint64, label int32) (mailList []*mailmodel.MailInfo, attachments []*mailmodel.Attachment, err error) {
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetAllMailAttachment NewMailModel fail", zap.Error(err))
		return
	}

	mailList = make([]*mailmodel.MailInfo, 0)

	for _, info := range mailModel.MailMap[label] {

		if checkMailExpire(logger, info) {
			logger.InfoWF("GetAllMailAttachment already expire", zap.Any("mailInfo", info))
			continue
		}
		if len(info.Attachments) == 0 {
			logger.InfoWF("GetAllMailAttachment Attachments is nil", zap.Any("mailInfo", info))
			continue
		}

		if info.IsGetAttach {
			logger.InfoWF("GetAllMailAttachment already GetAttach", zap.Any("info", info))
			continue
		}

		attachments = append(attachments, info.Attachments...)
		info.IsGetAttach = true
		info.IsRead = true
		mailList = append(mailList, info)
	}

	return
}

func (s service) DelMail(logger fklog.FKLogI, userId, mailId uint64, label int32) (err error) {
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("DelMail NewMailModel fail", zap.Error(err))
		return
	}

	mailInfo := getMailById(logger, userId, mailId, label)
	if mailInfo == nil {
		logger.InfoWF("DelMail getMailById is nil")
		return errors.New("mail not found")
	}

	if !mailInfo.IsGetAttach {
		logger.InfoWF("DelMail mail not GetAttach", zap.Any("mailInfo", mailInfo))
		return errors.New("邮件附件未领取")
	}

	delete(mailModel.MailMap[label], mailId)

	err = mailModel.Save(logger, userId)
	if err != nil {
		return err
	}
	return
}

func (s service) DelAllMail(logger fklog.FKLogI, userId uint64, label int32) (err error) {
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("DelAllMail NewMailModel fail", zap.Error(err))
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

	err = mailModel.Save(logger, userId)
	if err != nil {
		return err
	}

	return
}

func (s service) PushMailToReciver(logger fklog.FKLogI, userId uint64, packetType uint16, v interface{}) {
	isOnline := online.IsOnline(userId)
	if !isOnline {
		logger.InfoWF("PushMailToReciver reciver is not online", zap.Uint64("userId", userId))
		return
	}

	err := online.Push(logger, userId, packetType, v)
	if err != nil {
		logger.ErrorWF("PushMailToReciver fail", zap.Error(err))
	}
}

func (s service) GetMailAttachmentAfter(logger fklog.FKLogI, userId uint64, mailList []*mailmodel.MailInfo) (err error) {
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMailAttachmentAfter NewMailModel fail", zap.Error(err))
		return
	}

	for _, info := range mailList {
		mailModel.MailMap[info.Label][info.ID].IsRead = true
		mailModel.MailMap[info.Label][info.ID].IsGetAttach = true
		//mailModel.MailMap[info.ID].Status = int32(MailStatusReadClaimed)
	}

	err = mailModel.Save(logger, userId)
	if err != nil {
		logger.ErrorWF("GetMailAttachmentAfter Save fail", zap.Error(err))
		return
	}
	return nil
}

func newMail(logger fklog.FKLogI, title, context, senderName string, label int32, reciverID uint64, attachments []*mailmodel.Attachment, expireTime int64) *mailmodel.MailInfo {
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

func getMailById(logger fklog.FKLogI, userId, mailId uint64, label int32) *mailmodel.MailInfo {
	if mailId == 0 {
		return nil
	}
	mailModel, err := mailmodel.NewMailModel(logger, userId)
	if err != nil {
		logger.ErrorWF("getMailById NewMailModel fail", zap.Error(err))
		return nil
	}
	mailMap, ok := mailModel.MailMap[label]
	if !ok {
		logger.ErrorWF("getMailById label mail is empty", zap.Int32("label", label))
		return nil
	}
	info, ok := mailMap[mailId]
	if !ok {
		logger.ErrorWF("getMailById mail not found", zap.Int32("label", label), zap.Uint64("mailId", mailId))
		return nil
	}
	return info
}

func checkMailExpire(logger fklog.FKLogI, mailInfo *mailmodel.MailInfo) bool {
	if mailInfo.ExpireTime == 0 {
		return false
	}
	if mailInfo.ExpireTime <= time.Now().Unix() {
		logger.InfoWF("checkMailExpire mailInfo.ExpireTime <= time.Now().Unix()", zap.Any("mailInfo", mailInfo))
		return true
	}
	return false
}
