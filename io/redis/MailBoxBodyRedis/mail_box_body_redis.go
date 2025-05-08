package MailBoxBodyRedis

import (
	"context"
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"

	"go.uber.org/zap"
)

func init() {
	// MailBoxBodyRedis 16094 UN_CGK_REDIS_TYPE_MAIL_BOX_BODY_REDIS mail:body:%d {addr:127.0.0.1:6379,index:0,}
	_ = fkconfig.RegisterNameNode("MailBoxBodyRedis", 16094, redisConn) //TODO 用新类型
}

var redisConn = &fkredis.FkRedis{}

func getKey(mailID uint64) string {
	return redisConn.GetKey(mailID)
}

// 获取邮件列表
func GetMailList(logger fklog.FKLogI, mailIDS []uint64) ([]*MazeMail.MailInfo, error) {
	if len(mailIDS) <= 0 {
		logger.ErrorWF("mailIDS is empty")
		return nil, nil
	}

	keys := make([]interface{}, 0, len(mailIDS))
	for _, mailID := range mailIDS {
		keys = append(keys, getKey(mailID))
	}
	r, err := redisConn.Do(context.Background(), "MGET", keys...)
	if err != nil {
		logger.ErrorWF("call do", zap.Error(err))
		return nil, err
	}

	list, err := redis.ByteSlices(r, err)
	if err == redis.ErrNil {
		logger.WarnWF("ByteSlices is nill", zap.Error(err))
		return nil, nil
	}

	if err != nil {
		logger.WarnWF("result ByteSlices error", zap.Error(err))
		return nil, err
	}

	ret := make([]*MazeMail.MailInfo, 0, len(list))
	for i, result := range list {
		tmp := &MazeMail.MailInfo{}
		err := proto.Unmarshal(result, tmp)
		if err != nil {
			logger.ErrorWF("not mail", zap.Any("mailID", mailIDS[i]), zap.Error(err))
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}

const mailNormalExpireTime = 86400 * 30

// 获取邮件列表
func AddMail(logger fklog.FKLogI, mail *MazeMail.MailInfo) error {
	key := getKey(mail.GetBase().GetMailId())
	body, _ := proto.Marshal(mail)
	now := time.Now().Unix()
	var expire int64 = mailNormalExpireTime
	if mail.GetBase().GetExpires() > 0 {
		expire = mail.GetBase().GetExpires() - now
	}
	if expire <= 0 {
		expire = mailNormalExpireTime
	}
	_, err := redisConn.Do(context.Background(), "set", key, body)
	if err != nil {
		logger.ErrorWF("MailBoxBodyRedis.AddMail",
			zap.Any("Base", mail.GetBase()),
			zap.Error(err))
		return err
	}
	logger.InfoWF("MailBoxBodyRedis.AddMail success",
		zap.Any("Base", mail),
		zap.Error(err))
	return err
}

// 获取邮件列表
func GetMail(logger fklog.FKLogI, mailID uint64) (mail *MazeMail.MailInfo, err error) {
	key := getKey(mailID)

	result, err := redis.Bytes(redisConn.Do(context.Background(), "GET", key))
	if err != nil {
		logger.WarnWF("load mail data error", zap.Error(err))
		return nil, err
	}

	mail = &MazeMail.MailInfo{}
	err = proto.Unmarshal(result, mail)
	if err != nil {
		logger.ErrorWF("not mail", zap.Any("mailID", mailID), zap.Error(err))
	}

	return mail, err
}

// 删除邮件列表
func DelMail(logger fklog.FKLogI, mailID uint64) (ok bool, err error) {
	key := getKey(mailID)

	ok, err = redis.Bool(redisConn.Do(context.Background(), "DEL", key))
	if err != nil {
		logger.WarnWF("del mail data error", zap.Error(err))
	}
	return ok, err
}

func BatchDelMail(_ fklog.FKLogI, mailIDs []uint64, batchCount int) (ok bool, err error) {
	if len(mailIDs) == 0 {
		return
	}

	keys := make([]interface{}, 0, len(mailIDs))
	for _, mailID := range mailIDs {
		keys = append(keys, getKey(mailID))
	}

	getNum := len(mailIDs) / batchCount
	if len(mailIDs)%batchCount != 0 {
		// 有余数再加1
		getNum++
	}

	var delKeys []interface{}
	for index := 1; index <= getNum; index++ {
		if index == getNum {
			// 最后一个
			delKeys = keys[(index-1)*batchCount:]
		} else {
			delKeys = keys[(index-1)*batchCount : index*batchCount]
		}

		_, err = redisConn.Do(context.Background(), "del", delKeys...)
		if err != nil {
			return
		}
	}
	return
}
