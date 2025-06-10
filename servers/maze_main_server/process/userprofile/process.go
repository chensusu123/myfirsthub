// @Author pangchenyang 2025/6/9 21:32:00
// @Desc: 
package userprofile

import "maze_game_server/lib/nano/component"

type Profile struct {
	component.Base
}

func NewUserProfile() *Profile {
	// 初始化本地缓存
	initCache()
	return &Profile{}
}
