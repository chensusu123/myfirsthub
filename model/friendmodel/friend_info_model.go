package friendmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/friendredis"
	"maze_game_server/lib/serialize"
)

// 好友信息
type FriendInfo struct {
	UserId   uint64 `json:"user_id,omitempty"`
	CreateAt int64  `json:"create_at,omitempty"`
}
type FriendModel struct {
	FriendList []*FriendInfo `json:"friend_list,omitempty"`
}

func NewFriendModel(logger fklog.FKLogI, userID uint64) (*FriendModel, error) {
	friendModel := &FriendModel{}
	if err := friendModel.load(logger, userID); err != nil {
		return nil, err
	}
	return friendModel, nil
}

func (f *FriendModel) load(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := friendredis.GetFriends(logger, userId)
	if err != nil {
		logger.ErrorWF("load redis err", zap.Error(err))
		return nil
	}
	if len(bytes) == 0 {
		f.FriendList = make([]*FriendInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, f)
	if err != nil {
		logger.ErrorWF("load unmarshal err", zap.Error(err), zap.Any("bytes", string(bytes)))
		return err
	}

	return nil
}

func (f *FriendModel) Save(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := serialize.Marshal(f)
	if err != nil {
		logger.ErrorWF("SetFriends marshal failed", zap.Error(err), zap.Any("friends", f.FriendList))
		return err
	}
	if err = friendredis.SetFriends(logger, userId, bytes); err != nil {
		logger.ErrorWF("SetFriends failed", zap.Error(err), zap.Any("friends", f.FriendList))
		return err
	}
	logger.InfoWF("SetFriends success", zap.Any("friends", f.FriendList))
	return nil
}

func (f *FriendModel) Del(logger fklog.FKLogI, userId uint64) (err error) {
	if err = friendredis.DelFriends(logger, userId); err != nil {
		logger.ErrorWF("DelFriends err", zap.Error(err))
		return err
	}
	logger.InfoWF("DelFriends success")
	return nil
}
