// @Author pangchenyang 2025/6/9 21:33:00
// @Desc: 
package userprofile

import (
	"sync"
	lru "github.com/hashicorp/golang-lru"
	"maze_game_server/io/mysql/userprofilemysql"
	"maze_game_server/pb/common/UserProfile"
	"google.golang.org/protobuf/proto"
	"time"
)

var (
	cache *lru.Cache
	// 初始化本地缓存
	once sync.Once
)

func initCache() {
	once.Do(func() {
		var err error
		cache, err = lru.New(10000) // 缓存1万个用户资料
		if err != nil {
			panic(err)
		}
	})
}

func convertToProfile(dbProfile *userprofilemysql.UserProfile) *UserProfile.UserProfile {
	return &UserProfile.UserProfile{
		UserId:   proto.Uint64(dbProfile.UserID),
		NickName: proto.String(dbProfile.NickName),
		IconUrl:  proto.String(dbProfile.IconUrl),
		Sex:      proto.Int32(int32(dbProfile.Sex)),
	}
}

func convertToDbProfile(profile *UserProfile.UserProfile, creatTime, updateTime time.Time) *userprofilemysql.UserProfile {
	if profile == nil {
		return nil
	}
	return &userprofilemysql.UserProfile{
		UserID:    profile.GetUserId(),
		NickName:  profile.GetNickName(),
		IconUrl:   profile.GetIconUrl(),
		Sex:       uint8(profile.GetSex()),
		CreatedAt: creatTime,
		UpdatedAt: updateTime,
	}
}
