package friendredis

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	globalredis "maze_game_server/io/redis"
)

const (
	// 好友key
	KeyFriends = "friend:list:%d"
	// 待处理好友请求key
	KeyReceiveFriendRequest = "friend:receive:%d"
	//已发送好友请求key
	KeySendFriendRequest = "friend:send:%d"
	// 黑名单key
	KeyBlacklist = "blacklist:%d"
)

func getKeyFriends(userId uint64) string {
	return fmt.Sprintf(KeyFriends, userId)
}
func getKeyReceiveFriendRequest(userId uint64) string {
	return fmt.Sprintf(KeyReceiveFriendRequest, userId)
}

func getKeySendFriendRequest(userId uint64) string {
	return fmt.Sprintf(KeySendFriendRequest, userId)
}

func getKeyBlacklist(userId uint64) string {
	return fmt.Sprintf(KeyBlacklist, userId)
}

// 设置已发送好友请求
func SetSendFriendRequest(logger fklog.FKLogI, userId uint64, sends []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("SetSendFriendRequest get redis err", zap.Error(err), zap.Any("sends", sends))
		return err
	}
	key := db.MakeSectionKey(getKeySendFriendRequest(userId))
	if err = db.Set(context.TODO(), key, sends, 0).Err(); err != nil {
		logger.ErrorWF("SetSendFriendRequest set err", zap.Error(err), zap.String("key", key), zap.Any("sends", sends))
		return err
	}
	return nil
}

// 获取已发送的好友申请
func GetSendFriendRequest(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetSendFriendRequest get redis err", zap.Error(err))
		return nil, err
	}
	key := db.MakeSectionKey(getKeySendFriendRequest(userId))
	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ErrorWF("GetSendFriendRequest get err", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

// 删除所有已发送的好友申请
func DelSendFriendRequest(logger fklog.FKLogI, userId uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("DelSendFriendRequest get redis err", zap.Error(err))
		return err
	}
	key := db.MakeSectionKey(getKeySendFriendRequest(userId))
	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelSendFriendRequest del err", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// 设置收到的好友请求
func SetReceiveFriendRequest(logger fklog.FKLogI, userId uint64, receives []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("SetReceiveFriendRequest get redis err", zap.Error(err), zap.Any("receives", string(receives)))
		return err
	}
	key := db.MakeSectionKey(getKeyReceiveFriendRequest(userId))
	if err = db.Set(context.TODO(), key, receives, 0).Err(); err != nil {
		logger.ErrorWF("SetReceiveFriendRequest set err", zap.String("key", key), zap.Error(err), zap.Any("receives", string(receives)))
		return err
	}
	return nil
}

// 获取收到的好友请求
func GetReceiveFriendRequest(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetReceiveFriendRequest get redis err", zap.Error(err))
		return nil, err
	}
	key := db.MakeSectionKey(getKeyReceiveFriendRequest(userId))
	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ErrorWF("GetReceiveFriendRequest get err", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

// 删除所有收到的好友请求
func DelReceiveFriendRequest(logger fklog.FKLogI, userId uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("DelReceiveFriendRequest get redis err", zap.Error(err))
		return err
	}
	key := db.MakeSectionKey(getKeyReceiveFriendRequest(userId))
	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelReceiveFriendRequest del err", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// 获取好友
func GetFriends(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetFriends get redis err", zap.Error(err))
		return nil, err
	}
	key := db.MakeSectionKey(getKeyFriends(userId))
	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ErrorWF("GetFriends get err", zap.String("key", key), zap.Error(err))
		return nil, err
	}
	return bytes, nil
}

// 设置好友
func SetFriends(logger fklog.FKLogI, userId uint64, friends []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("SetFriends get redis err", zap.Error(err), zap.Any("friends", friends))
		return err
	}
	key := db.MakeSectionKey(getKeyFriends(userId))
	if err = db.Set(context.TODO(), key, friends, 0).Err(); err != nil {
		logger.ErrorWF("SetFriends set err", zap.String("key", key), zap.Error(err), zap.Any("bytes", friends))
		return err
	}
	return nil
}

// 删除所有好友
func DelFriends(logger fklog.FKLogI, userId uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("DelFriends get redis err", zap.Error(err))
		return err
	}
	key := db.MakeSectionKey(getKeyFriends(userId))
	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelFriends del err", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}

// 设置黑名单
func SetBlacklist(logger fklog.FKLogI, userId uint64, blacklist []byte) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("SetBlacklist get redis err", zap.Error(err), zap.Any("blacklist", blacklist))
		return err
	}
	key := db.MakeSectionKey(getKeyBlacklist(userId))
	if err = db.Set(context.TODO(), key, blacklist, 0).Err(); err != nil {
		logger.ErrorWF("SetBlacklist set err", zap.String("key", key), zap.Error(err), zap.Any("bytes", blacklist))
		return err
	}
	return nil
}

// 获取黑名单
func GetBlacklist(logger fklog.FKLogI, userId uint64) ([]byte, error) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("GetBlacklist get redis err", zap.Error(err), zap.Uint64("userId", userId))
		return nil, err
	}
	key := db.MakeSectionKey(getKeyBlacklist(userId))
	bytes, err := db.Get(context.TODO(), key).Bytes()
	if err != nil {
		logger.ErrorWF("GetBlacklist get err", zap.String("key", key), zap.Error(err), zap.Uint64("userId", userId))
		return nil, err
	}
	return bytes, nil
}

// 删除所有好友
func DelBlacklist(logger fklog.FKLogI, userId uint64) error {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		logger.ErrorWF("DelBlacklist get redis err", zap.Error(err))
		return err
	}
	key := db.MakeSectionKey(getKeyBlacklist(userId))
	err = db.Del(context.TODO(), key).Err()
	if err != nil {
		logger.ErrorWF("DelBlacklist del err", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
