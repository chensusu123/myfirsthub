package mailunreaddb

import (
	"context"
	"strconv"
	"time"

	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis/redis"
	"gitlab.ifreetalk.com/plate/freetk/fkutil"
	"gitlab.ifreetalk.com/plate/protodef/MazeMail"
	_ "gitlab.ifreetalk.com/plate/protodef/MazeMail"

	"go.uber.org/zap"
)

func init() {
	// mailunreaddb 16095 UN_CGK_REDIS_TYPE_MAIL_BOX_USER_REDIS mail:seq:%d,mail:box:%d,mail:box:%d,annex:open:%d:%d,mail:readed:%d {addr:127.0.0.1:6379,index:0,}
	fkconfig.RegisterNameNode("mailunreaddb", 20154, redisConn) //TODO 用新类型
}

var redisConn = &fkredis.FkRedis{}

func getUnreadKey(userId uint64, classType int32) string {
	//return fmt.Sprintf("u:%d:mail:%d:unread", userId, classType)
	return redisConn.GetKey(userId, classType)
}

const GetMailIndexesLimitCount = 20

const openAnnexFlagExpireTime = 86400 * 30

//获取未读邮件数量
func UnreadMailNum(logger fklog.FKLogI, userID uint64, classType int32) (num int32, err error) {
	key := getUnreadKey(userID, classType)
	ctx := context.Background()
	reply, err := redis.Int(redisConn.Do(ctx, "ZCOUNT", key, time.Now().Unix(), "+inf"))
	if err != nil {
		logger.ErrorWF("call HaveUnreadMail do", zap.Error(err))
	}

	logger.DebugWF("call HaveUnreadMail ret", zap.Int("num", reply))
	num = int32(reply)
	return
}

//添加未读邮件
func AddUnreadMail(logger fklog.FKLogI, userID uint64, classType int32, guid uint64, expireTime int64) (err error) {
	key := getUnreadKey(userID, classType)
	ctx := context.Background()
	_, err = redis.Int(redisConn.Do(ctx, "ZADD", key, expireTime, guid))
	if err != nil {
		logger.ErrorWF("call AddUnreadMail do", zap.Error(err), zap.String("key", key), zap.Uint64("guid", guid), zap.Int64("expireTime", expireTime))
	}
	return
}

func DelUnreadMail(logger fklog.FKLogI, userID uint64, classType int32, guids []uint64) (err error) {
	_, err = delUnreadMail(logger, userID, classType, guids)
	return err
}

func DelUnreadMail2(logger fklog.FKLogI, userID uint64, classType int32, guids []uint64) (num int, err error) {
	return delUnreadMail(logger, userID, classType, guids)
}

//删除未读邮件
func delUnreadMail(logger fklog.FKLogI, userID uint64, classType int32, guids []uint64) (num int, err error) {
	if len(guids) < 1 {
		return
	}

	key := getUnreadKey(userID, classType)
	ctx := context.Background()

	logger.DebugWF("DelUnreadMail req", zap.Any("guids", guids))
	fields := make([]interface{}, 0, len(guids)+1)
	fields = append(fields, key)
	for _, mailID := range guids {
		fields = append(fields, mailID)
	}

	num, err = redis.Int(redisConn.Do(ctx, "ZREM", fields...))
	if err != nil {
		logger.ErrorWF("call DelUnreadMail do", zap.String("key", key), zap.Uint64s("guids", guids), zap.Error(err))
	}
	return
}

// 获取距离现在最近的未读信封的过期时间
func GetUnreaMailExpireTime(logger fklog.FKLogI, userID uint64, classType int32) (retTime int64, err error) {
	key := getUnreadKey(userID, classType)
	ctx := context.Background()

	info, loadErr := redis.ByteSlices(redisConn.Do(ctx, "ZRANGE", key, 0, 0, "WITHSCORES"))
	if loadErr == redis.ErrNil {
		return
	} else if loadErr != nil {
		err = loadErr
		logger.WarnWF("GetLatestMailInfo error", zap.Error(loadErr))
	}

	if len(info) == 2 {
		retTime = fkutil.ToInt64(string(info[1]))
	}
	logger.DebugWF("GetLatestMailInfo res", zap.Any("expireTime", retTime), zap.Any("nowTime", time.Now().Unix()),
		zap.Any("loadErr", loadErr))
	return
}

// 获取过期时间最小的邮件id和时间
// 先找有没有过期的，如果没有过期的,则过期时间是集合的第一个元素。如果过期了，取过期之后的下一个元素
func GetLatestUnreadMail(logger fklog.FKLogI, userID uint64, classType int32) (mailList []uint64, retTime int64, err error) {
	key := getUnreadKey(userID, classType)
	ctx := context.Background()
	nowTime := time.Now().Unix()

	// 获取第一个
	info, loadErr := redis.ByteSlices(redisConn.Do(ctx, "ZRANGEBYSCORE", key, "-inf", nowTime, "WITHSCORES"))
	if loadErr == redis.ErrNil {
		logger.WarnWF("GetLatestUnreadMail empty")
		return
	}

	idx := 0
	// 添加过期信封id
	if len(info) >= 2 {
		for i := 0; i < len(info); i += 2 {
			mailList = append(mailList, fkutil.ToUint64(string(info[i])))
			idx++
		}
	}

	// 当前时间大于第一个时间，表名有过期未读消息
	info2, loadErr2 := redis.ByteSlices(redisConn.Do(ctx, "ZRANGE", key, idx, idx, "WITHSCORES"))
	if loadErr2 != nil {
		err = loadErr2
		logger.WarnWF("GetLatestMailInfo error", zap.Error(loadErr2))
	}

	if len(info2) == 2 {
		retTime = fkutil.ToInt64(string(info2[1]))
	}
	logger.DebugWF("GetLatestMailInfo res", zap.Any("mailID", mailList), zap.Any("expireTime", retTime), zap.Any("nowTime", time.Now().Unix()),
		zap.Any("loadErr", loadErr), zap.Any("loadErr2", loadErr2))
	return
}

//从未读消息里读取列表
func GetMailIndexesFromUnreadList(logger fklog.FKLogI, userID uint64, classType int32, maxSeq uint64, limit int32) (mailIndexes []*MazeMail.ReceiptItem, err error) {
	key := getUnreadKey(userID, classType)
	ctx := context.Background()
	reply, err := redisConn.Do(ctx, "ZREVRANGEBYSCORE", key, maxSeq, "-inf", "WITHSCORES", "LIMIT", 0, limit)
	if err != nil {
		logger.ErrorWF("call GetMailIndexes do", zap.Error(err))
	}

	res, err := redis.ByteSlices(reply, err)
	if err != nil {
		return nil, err
	}

	ret := make([]*MazeMail.ReceiptItem, 0, len(res)/2)

	for i := 0; i < len(res); i += 2 {
		id, err := strconv.ParseUint(string(res[i]), 10, 64)
		if err != nil {
			continue
		}
		seq, err := strconv.ParseUint(string(res[i+1]), 10, 64)
		if err != nil {
			continue
		}
		ret = append(ret, &MazeMail.ReceiptItem{Id: proto.Uint64(id), Sequence: proto.Uint64(seq)})
	}
	return ret, nil
}

//判断信封是有已经存在
func IsMailBoxOpen(logger fklog.FKLogI, userID uint64, classType int32, mailID uint64) (bool, error) {
	ctx := context.Background()
	key := getUnreadKey(userID, classType)
	_, err := redis.Int(redisConn.Do(ctx, "ZSCORE", key, mailID))
	if err == redis.ErrNil {
		logger.DebugWF("IsMailBoxOpen not exist", zap.String("key", key), zap.Uint64("mailID", mailID))
		return false, nil
	}

	if err != nil {
		logger.ErrorWF("IsMailBoxOpen fail", zap.String("key", key), zap.Error(err))
		return false, err
	}
	logger.InfoWF("success", zap.String("key", key), zap.Uint64("mailID", mailID))
	return true, err
}
