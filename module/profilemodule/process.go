// @Author pangchenyang 2025/6/19 14:50:00
// @Desc: 
package profilemodule

import (
	"maze_game_server/pb/common/UserProfile"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

type ProfileModule interface {
	AlterUserProfile(logger fklog.FKLogI, data *UserProfile.UserProfile) (*UserProfile.UserProfile, error)
	QueryUserProfile(logger fklog.FKLogI, userID uint64, viewIDs []uint64) ([]*UserProfile.UserProfile, error)
	QueryAvatarToken(logger fklog.FKLogI) (string, error)
}

// todo
func NewProfileModule() ProfileModule {
	return NewUserProfileModule()
}

type UserProfileModule struct {
}

func NewUserProfileModule() *UserProfileModule {
	return &UserProfileModule{}
}

func InitModule() {
	// 初始化缓存
	initUserProfileCache()
}

// notifyProfileChange 通知用户资料变更
func notifyProfileChange(ctx fklog.FKLogI, userId uint64) {
	// TODO: 实现通知逻辑，可以通过消息队列或其他方式通知相关服务
}
