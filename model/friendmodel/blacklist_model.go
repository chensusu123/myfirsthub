package friendmodel

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/redis/friendredis"
	"maze_game_server/lib/serialize"
)

// 黑名单信息
type BlacklistInfo struct {
	UserId   uint64 `json:"user_id,omitempty"`
	CreateAt int64  `json:"create_at,omitempty"`
}
type BlacklistModel struct {
	Blacklist []*BlacklistInfo `json:"blacklist,omitempty"`
}

func NewBlacklistModel(logger fklog.FKLogI, userID uint64) (*BlacklistModel, error) {
	sendModel := &BlacklistModel{}
	if err := sendModel.load(logger, userID); err != nil {
		return nil, err
	}
	return sendModel, nil
}

func (f *BlacklistModel) load(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := friendredis.GetBlacklist(logger, userId)
	if err != nil {
		logger.ErrorWF("load blacklist redis err", zap.Error(err))
		return nil
	}
	if len(bytes) == 0 {
		f.Blacklist = make([]*BlacklistInfo, 0)
		return nil
	}
	err = serialize.Unmarshal(bytes, f)
	if err != nil {
		logger.ErrorWF("load blacklist unmarshal err", zap.Error(err), zap.Any("bytes", string(bytes)))
		return err
	}

	return nil
}

func (f *BlacklistModel) Save(logger fklog.FKLogI, userId uint64) (err error) {
	bytes, err := serialize.Marshal(f)
	if err != nil {
		logger.ErrorWF("save blacklist marshal err", zap.Error(err), zap.Any("friends", f))
		return err
	}
	if err = friendredis.SetBlacklist(logger, userId, bytes); err != nil {
		logger.ErrorWF("save blacklist err", zap.Error(err), zap.Any("friends", f))
		return err
	}
	logger.InfoWF("save blacklist success", zap.Any("friends", f))
	return nil
}

// 删除所有发送的好友申请
func (f *BlacklistModel) Del(logger fklog.FKLogI, userId uint64) (err error) {
	if err = friendredis.DelBlacklist(logger, userId); err != nil {
		logger.ErrorWF("del blacklist err", zap.Error(err))
		return err
	}
	logger.InfoWF("del blacklist success")
	return nil
}
