// @Author: ZhaoXiming 2025/3/21 16:52
// @Desc: 发信封

package mail

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/cgkargs"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/mailutil"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mail_module"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/FamilyAllocUserRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/LeagueFamilyRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/common/new_map_db/WorldLeagueRedis"
	"gitlab.ifreetalk.com/plate/io_interface/redis_interface/mail/broadcastmailrecordredis"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	"gitlab.ifreetalk.com/plate/protodef/PhoneNotify"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

func MailSendProcess(ctx context.Context, logger fklog.FKLogI, index int, data []byte) (err error) {
	ctx = context.Background()
	defer fkprometheus.DebugPMT("MailQueueSend")()

	m := &MazeMail.SendMailRQ{}
	err = proto.Unmarshal(data, m)
	if err != nil {
		logger.ErrorWF("MailSendProcess proto unmarshal err", zap.Error(err))
		return
	}

	logger.WarnWF("MailSendProcess start",
		zap.Uint64("mailID", m.MailInfo.Base.GetMailId()),
		//zap.Uint64("fromUID", m.MailInfo.GetFromUser().GetUserId()),
		//zap.Uint64("toUID", m.GetToUser()),
		zap.Uint64s("list", m.UserList),
		zap.Any("obj", m))

	mailInfo := m.GetMailInfo()
	if mailInfo == nil {
		logger.ErrorWF("MailSendProcess  mail_Info==nil", zap.Any("req", m))
		return nil
	}
	baseInfo := mailInfo.GetBase()
	if baseInfo == nil {
		logger.ErrorWF("MailSendProcess mail_base==nil", zap.Any("req", m))
		return nil
	}

	if mailutil.IsMailHasAward(mailInfo) && baseInfo.GetExtType() <= 0 {
		// 为可领取附件类型 没有来源
		logger.ErrorWF("MailSendProcess award mail extType nil", zap.Any("req", m))
		return
	}

	_ = processSendMail(ctx, logger, m)
	return nil
}

// 处理信封发送
func processSendMail(ctx context.Context, logger fklog.FKLogI, m *MazeMail.SendMailRQ) error {
	baseInfo := m.MailInfo.GetBase()
	if baseInfo == nil {
		logger.ErrorWF("MailSend Do: mail_base==nil", zap.Any("req", m))
		return nil
	}

	bType := baseInfo.GetBroadcastFlag()

	if bType == 0 {
		toUser := m.MailInfo.GetToUser()
		if toUser == 0 {
			logger.ErrorWF("MailSend Do: userid ==0", zap.Any("req", m))
			return nil
		}
		sendOneMail(ctx, logger, toUser, m.MailInfo)
		return nil
	}
	return sendBroadcastMail(ctx, logger, m.MailInfo, m.UserList)
}

func sendOneMail(ctx context.Context, logger fklog.FKLogI, toUser uint64, mailInfo *MazeMail.MailInfo) (err error) {

	for i := 0; i < 3; i++ {
		_, err = mail_module.SendMail2User(ctx, logger, toUser, mailInfo)
		if err == nil || err == errors.ERR_MAIL_ALREADY_EXIST {
			logger.InfoWF("MailSend Do:start mail send op", zap.Any("mail", mailInfo), zap.Uint64("toUID", toUser))
			err = nil
			return nil
		} else if err == errors.ARGS_NOT_MATCH {
			logger.ErrorWF("MailSend Do: send param error", zap.Uint64("userId", toUser), zap.Any("info", mailInfo), zap.Error(err))
			err = nil
			return nil
		}

	}

	//TODO 异常流程先不要
	//if err != nil {
	//	// 有异常
	//	data, _ := json.Marshal(&MailBoxSvr.SvrSendMailRQ{ToUser: proto.Uint64(toUser), MailInfo: mailInfo})
	//	logger.ErrorWF("MailSend Do: send error", zap.ByteString("req", data), zap.Error(err))
	//
	//	if err == errors.DB_SAVE_ERROR {
	//		// 信封索引保存失败
	//		mail_module.PushMailSend(ctx, logger, toUser, mailInfo, mailRecord.SaveErr)
	//	} else {
	//		// 其余中断错误
	//		mail_module.PushMailSend(ctx, logger, toUser, mailInfo, mailRecord.OtherErr)
	//	}
	//
	//	if _, isExist := cgkargs.ReSendOpTypeMap[mailInfo.GetBase().GetExtType()]; isExist {
	//		// 需要重发 保存备份数据
	//		err = mailfailbackupdb.SetBackup(logger, toUser, mailInfo.GetBase().GetMailId(), mailfailbackupdb.Mail,
	//			mailInfo, mailInfo.GetBase().GetExpires()-time.Now().Unix()+86400,
	//		)
	//		if err != nil {
	//			logger.ErrorWF("MailSend Do: SetBackup err", zap.Any("mailInfo", mailInfo), zap.Error(err))
	//		}
	//	}
	//}
	return nil
}

// 广播信封，这里的信封不处理发送信封回调，AfterMailSend
func sendBroadcastMail(ctx context.Context, logger fklog.FKLogI, mail *MazeMail.MailInfo, userList []uint64) (err error) {
	if mail.Base.GetMailId() == 0 {
		logger.ErrorWF("sendBroadcastMail MailSend Do: mailId ==0", zap.Any("req", mail))
		return nil
	}

	bid := mail.Base.GetBroadcastId()
	if bid == 0 {
		logger.ErrorWF("sendBroadcastMail BroadcastId not filled", zap.Any("mail", mail))
		return
	}

	bType := mail.Base.GetBroadcastFlag()

	if bType == 4 {
		if len(userList) == 0 {
			err = fmt.Errorf("input user 0")
			return
		}

	} else {
		userList, err = getUserList(logger, bType, bid)
		if err != nil {
			return err
		}
	}

	err = mail_module.BroadcastMail(ctx, logger, mail)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	_broadUserCount := cgkargs.BroadUserCount
	ch := make(chan uint64, _broadUserCount)
	wg.Add(int(_broadUserCount))
	count := int64(0)
	for i := int64(0); i < _broadUserCount; i++ {
		go func(l fklog.FKLogI, sendMail *MazeMail.MailInfo, execIndex int64) {
			defer wg.Done()
			for userID := range ch {
				atomic.AddInt64(&count, 1)
				broadSendMail(l, execIndex, userID, sendMail)
			}
		}(logger, mail, i)
	}

	startTime := time.Now()
	logger.InfoWF("sendBroadcastMail start", zap.Uint64("mailId", mail.Base.GetMailId()))
	for _, uid := range userList {
		ch <- uid
	}
	close(ch)
	wg.Wait()
	logger.InfoWF("sendBroadcastMail end", zap.Any("mailId", mail.Base.GetMailId()), zap.Int64("sendUser", count), zap.Int64("totalUser", int64(len(userList))), zap.Duration("costTime", time.Since(startTime)))
	broadcastmailrecordredis.ExpireRecord(logger, mail.Base.GetMailId(), 2*86400)
	return nil
}

// 并发发送信封
func broadSendMail(logger fklog.FKLogI, idx int64, userID uint64, mail *MazeMail.MailInfo) {
	if userID <= 0 {
		return
	}
	ctx := context.Background()
	mail_module.SendMail2User(ctx, logger, userID, mail)

	//TODO 限制等级
	//if mail.GetUserLimitType() != 1 {
	//	logger.DebugWF("broadSendMail user debug", zap.Uint64("userID", userID), zap.Int64("execIndex", idx))
	//	mail_module.SendMail2User(ctx, logger, userID, mail)
	//	return
	//}
	//
	////根据岛等级限制发送用户
	//level, err := mailutil.GetLevelByVersion(logger, userID)
	//logger.DebugWF("broadSendMail level", zap.Any("level", level), zap.Uint64("userID", userID),
	//	zap.Int64("execIndex", idx),
	//	zap.Int32("userLimitValMax", mail.GetUserLimitValueMax()),
	//	zap.Int32("userLimitValMin", mail.GetUserLimitValueMin()))
	//if err != nil {
	//	logger.InfoWF("broadSendMail get user level error", zap.Uint64("uid", userID), zap.Int32("level", level), zap.Error(err))
	//}
	//
	//if mail.GetUserLimitValueMax() > 0 {
	//	logger.DebugWF("broadSendMail level max",
	//		zap.Any("GetUserLimitValueMin", mail.GetUserLimitValueMin()), zap.Any("GetUserLimitValueMax", mail.GetUserLimitValueMax()))
	//	if level >= mail.GetUserLimitValueMin() && level <= mail.GetUserLimitValueMax() {
	//		mail_module.SendMail2User(ctx, logger, userID, mail)
	//	}
	//} else if level >= mail.GetUserLimitValueMin() {
	//	mail_module.SendMail2User(ctx, logger, userID, mail)
	//}
	return
}

// 家族用户列表
func getFamilyUsers(logger fklog.FKLogI, familyId uint64) (list []uint64, err error) {
	return FamilyAllocUserRedis.BatchGetAllFamilyUsers(logger, familyId, 300)
}

// 获取联盟用户列表
func getLeagueUsers(logger fklog.FKLogI, leagueId uint64) (list []uint64, err error) {
	// 获取联盟下的家族列表
	list = make([]uint64, 0, 2000)
	ids, err := LeagueFamilyRedis.GetAllLeagueFamilyIDs(logger, leagueId)
	if err != nil {
		return
	}
	logger.DebugWF("getLeagueUsers get league family id", zap.Uint64("league", leagueId), zap.Strings("family", ids))
	wg1 := &sync.WaitGroup{}
	_getLeagueUserCount := cgkargs.GetLeagueUserCount
	ch1 := make(chan string, _getLeagueUserCount)
	lock := sync.Mutex{}
	startTime := time.Now()

	wg1.Add(int(_getLeagueUserCount))
	for i := int64(1); i <= _getLeagueUserCount; i++ {
		go func(l fklog.FKLogI) {
			defer wg1.Done()
			for _familyID := range ch1 {
				familyId := fkutil.ToUint64(_familyID)
				if familyId <= 0 {
					continue
				}
				uids, err := getFamilyUsers(l, familyId)
				if err != nil {
					logger.ErrorWF("getLeagueUsers getFamilyUsers error", zap.Error(err),
						zap.Uint64("familyID", familyId), zap.String("familyIDStr", _familyID))
					continue
				}
				if len(uids) <= 0 {
					continue
				}
				lock.Lock()
				list = append(list, uids...)
				lock.Unlock()
			}
		}(logger)
	}
	for _, fid := range ids {
		ch1 <- fid
	}

	close(ch1)
	wg1.Wait()
	logger.InfoWF("getLeagueUsers get users end", zap.Uint64("leagueId", leagueId), zap.Int64("leagueUserCount", int64(len(list))), zap.Duration("costTime", time.Since(startTime)))
	return
}

// 获取地图用户id
func getMapUsers(logger fklog.FKLogI, mapId uint64) (list []uint64, err error) {
	leagues, err := WorldLeagueRedis.GetWorldLeagueInfo(logger, mapId)
	if err != nil {
		return
	}
	logger.DebugWF("get map league id", zap.Uint64("mapId", mapId), zap.Any("league", leagues))
	list = make([]uint64, 0, 10000)
	for lid := range leagues {
		uids, err := getLeagueUsers(logger, lid)
		if err != nil {
			continue
		}
		list = append(list, uids...)
	}
	return
}

// 根据广播标志获取用户列表
func getUserList(logger fklog.FKLogI, btype int32, id uint64) ([]uint64, error) {
	switch btype {
	case int32(PhoneNotify.PhoneBroadcastFlag_PHONE_FAMILY_BROADCAST):
		return getFamilyUsers(logger, id)
	case int32(PhoneNotify.PhoneBroadcastFlag_PHONE_LEAGUE_BROADCAST):
		return getLeagueUsers(logger, id)
	case int32(PhoneNotify.PhoneBroadcastFlag_PHONE_MAP_BROADCAST):
		return getMapUsers(logger, id)
	}
	return nil, nil
}
