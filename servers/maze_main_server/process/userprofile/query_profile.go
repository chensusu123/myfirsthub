// @Author pangchenyang 2025/6/9 21:33:00
// @Desc: 
package userprofile

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"maze_game_server/common/errors"
	"go.uber.org/zap"
	"time"
	"context"
	"google.golang.org/protobuf/proto"
	"maze_game_server/pb/common/UserProfile"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/io/mysql/userprofilemysql"
	"gorm.io/gorm"
	"fmt"
	"math/rand"
	"sync"
	"maze_game_server/io/redis/userprofilelock"
)

// OnQueryUserProfile 查询用户资料
func OnQueryUserProfile(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message, opData string) (err error) {
	defer fkprometheus.DebugPMT("OnQueryUserProfile")()
	userCtx := fkserver.NewUserContext(context.TODO(), uint64(shardingID), ctx)
	req := rqMsg.(*UserProfile.QueryUserProfileRQ)
	res := rsMsg.(*UserProfile.QueryUserProfileRS)
	res.ErrInfo = errors.NO_ERROR
	userCtx.WarnWF("OnQueryUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		userCtx.WarnWF("OnQueryUserProfile end ", zap.Any("req", req), zap.Any("res", res), zap.Float64("costTime", costTime))
	}()

	// 检查rq
	if len(req.GetUserId()) == 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 并发查询用户资料
	var (
		wg sync.WaitGroup
	)
	for _, v := range req.GetUserId() {
		userID := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			ret, err := queryUserProfile(userCtx, userID)
			if err != nil {
				userCtx.ErrorWF("OnQueryUserProfile queryUserProfile failed", zap.Error(err))
				return
			}
			res.UserProfile = append(res.UserProfile, ret)
		}()
	}
	return
}

func queryUserProfile(logger fklog.FKLogI, userID uint64) (ret *UserProfile.UserProfile, err error) {
	if userID <= 0 {
		err = fmt.Errorf("invalid userID")
		logger.ErrorWF("queryUserProfile invalid userID", zap.Uint64("userID", userID))
		return
	}
	// 1. 检查缓存
	if val, ok := cache.Get(userID); ok {
		ret = val.(*UserProfile.UserProfile)
		return
	}

	// 2. 读取redis
	ret, err = userprofileredis.GetCache(userID)
	if err != nil {
		logger.ErrorWF("queryUserProfile GetCache error", zap.Error(err))
		return
	}
	if ret != nil {
		// 更新本地缓存
		cache.Add(userID, ret)
		return
	}

	// 3. 查询数据库
	dbProfile, err := userprofilemysql.GetByUserID(userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.ErrorWF("OnQueryUserProfile GetByUserID error", zap.Error(err))
		return
	}

	if dbProfile != nil {
		ret = &UserProfile.UserProfile{
			UserId:   proto.Uint64(dbProfile.UserID),
			NickName: proto.String(dbProfile.Nickname),
			IconUrl:  proto.Int32(int32(dbProfile.IconUrl)),
			Sex:      proto.Int32(int32(dbProfile.Sex)),
		}

		// 更新缓存
		cache.Add(userID, ret)
		if err = userprofileredis.SetCache(userID, ret); err != nil {
			logger.ErrorWF("OnQueryUserProfile SetCache error", zap.Error(err))
		}
		return
	}

	// 4. 初始化用户资料
	return initUserProfile(logger, userID)
}

// initUserProfile 初始化用户资料
func initUserProfile(logger fklog.FKLogI, userID uint64) (ret *UserProfile.UserProfile, err error) {
	// 获取分布式锁
	locked, err := userprofilelock.Lock(userID)
	if err != nil {
		logger.ErrorWF("initUserProfile get lock failed", zap.Error(err))
		return
	}
	if !locked {
		err = fmt.Errorf("服务器忙")
		logger.WarnWF("initUserProfile lock already exists")
		return
	}
	defer func() {
		_ = userprofilelock.Unlock(userID)
	}()

	if val, ok := cache.Get(userID); ok {
		ret = val.(*UserProfile.UserProfile)
		return
	}

	profile := generateInitProfile(logger, userID)

	// 保存到数据库
	dbProfile := &userprofilemysql.UserProfile{
		UserID:   userID,
		Nickname: profile.GetNickName(),
		IconUrl:  uint32(profile.GetIconUrl()),
		Sex:      uint8(profile.GetSex()),
	}
	if err = userprofilemysql.Create(dbProfile); err != nil {
		logger.ErrorWF("initUserProfile create db profile failed", zap.Error(err))
		return
	}

	// 更新缓存
	cache.Add(userID, profile)
	if err := userprofileredis.SetCache(userID, profile); err != nil {
		logger.ErrorWF("initUserProfile set cache failed", zap.Error(err))
	}

	ret = profile
	return
}

// generateInitProfile 生成初始用户资料
func generateInitProfile(logger fklog.FKLogI, userId uint64) *UserProfile.UserProfile {
	// todo 应从配置读取
	maleAvatars := []uint32{1001, 1002, 1003}
	femaleAvatars := []uint32{2001, 2002, 2003}
	nicknamePrefix := []string{
		"Shadow", "Blood", "Dark", "Iron", "Storm",
		"Night", "Ghost", "Rage", "Frost", "Flame",
		"Thunder", "Silver", "Gold", "Steel", "Demon",
		"Angel", "Dragon", "Wolf", "Phoenix", "Void",
	}

	nicknameSuffix := []string{
		"blade", "fist", "heart", "fang", "wing",
		"hunter", "slayer", "reaver", "lord", "king",
		"queen", "mage", "ranger", "assassin", "warden",
		"bringer", "walker", "guardian", "watcher", "caller",
		"born", "seeker", "weaver", "shade", "born",
	}

	// 默认性别为男
	sex := int32(1)
	avatar := maleAvatars[0]

	rand.Seed(time.Now().UnixNano())
	if rand.Intn(2) == 1 { // 50%概率为女性
		sex = 2
		avatar = femaleAvatars[rand.Intn(len(femaleAvatars))]
	} else {
		avatar = maleAvatars[rand.Intn(len(maleAvatars))]
	}

	nickname := fmt.Sprintf("%s%s%d",
		nicknamePrefix[rand.Intn(len(nicknamePrefix))],
		nicknameSuffix[rand.Intn(len(nicknameSuffix))],
		userId%10000)

	initUser := &UserProfile.UserProfile{
		UserId:   proto.Uint64(userId),
		NickName: proto.String(nickname),
		IconUrl:  proto.Int32(int32(avatar)),
		Sex:      proto.Int32(sex),
	}
	logger.DebugWF("generateInitProfile success", zap.Any("initUser", initUser))
	return initUser
}
