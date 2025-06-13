// @Author pangchenyang 2025/6/9 21:33:00
// @Desc: 
package userprofile

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkprometheus"
	"maze_game_server/common/errors"
	"go.uber.org/zap"
	"time"
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
func (p *Profile) OnQueryUserProfile_10481_10482(ctx fklog.FKLogI, shardingID int64, rqMsg proto.Message, rsMsg proto.Message, opData string) (err error) {
	defer fkprometheus.DebugPMT("OnQueryUserProfile")()
	req := rqMsg.(*UserProfile.QueryUserProfileRQ)
	res := rsMsg.(*UserProfile.QueryUserProfileRS)
	res.ErrInfo = errors.NO_ERROR
	ctx.WarnWF("OnQueryUserProfile with", zap.Any("rq", req))

	addStartTime := time.Now()
	defer func() {
		costTime := time.Since(addStartTime).Seconds()
		ctx.WarnWF("OnQueryUserProfile end ", zap.Any("req", req), zap.Any("res", res),
			zap.String("errMsg", string(res.GetErrInfo().GetErrMsg())),
			zap.Float64("costTime", costTime))
	}()

	// 检查rq
	if len(req.GetUserId()) == 0 || shardingID <= 0 {
		res.ErrInfo = errors.COMMON_ERROR_TIPS.Wrap("无效参数")
		return
	}
	// 并发查询用户资料
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	for _, v := range req.GetUserId() {
		userID := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			ret, err := queryUserProfile(ctx, uint64(shardingID), userID)
			if err != nil {
				ctx.ErrorWF("OnQueryUserProfile queryUserProfile failed", zap.Error(err))
				return
			}
			mu.Lock()
			res.UserProfile = append(res.UserProfile, ret)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return
}

func queryUserProfile(logger fklog.FKLogI, shardingID, userID uint64) (ret *UserProfile.UserProfile, err error) {
	if userID <= 0 {
		err = fmt.Errorf("invalid userID")
		logger.ErrorWF("queryUserProfile invalid userID", zap.Uint64("userID", userID))
		return
	}

	if val, ok := cache.Get(userID); ok {
		ret = val.(*UserProfile.UserProfile)
		logger.DebugWF("OnQueryUserProfile get from cache", zap.Uint64("userID", userID),
			zap.Any("ret", ret))
		return
	}

	ret, err = userprofileredis.GetCache(userID)
	if err != nil {
		logger.ErrorWF("queryUserProfile GetCache error", zap.Error(err))
		return
	}
	if ret != nil {
		cache.Add(userID, ret)
		logger.DebugWF("OnQueryUserProfile update cache", zap.Uint64("userID", userID), zap.Any("ret", ret))
		return
	}

	// todo 缓存穿透？
	dbProfile, err := userprofilemysql.GetByUserID(userID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.ErrorWF("OnQueryUserProfile GetByUserID error", zap.Error(err))
			return
		} else {
			// 处理errNil
			userprofilemysql.GetDB().Error = nil
		}
	}

	if dbProfile != nil && dbProfile.UserID > 0 {
		ret = convertToProfile(dbProfile)

		cache.Add(userID, ret)
		if err = userprofileredis.SetCache(userID, ret); err != nil {
			logger.ErrorWF("OnQueryUserProfile SetCache error", zap.Error(err))
		}
		return
	}

	if shardingID != userID {
		logger.WarnWF("OnQueryUserProfile not register user", zap.Uint64("userID", userID), zap.Uint64("shardingID", shardingID))
		return
	}

	return initUserProfile(logger, userID)
}

// initUserProfile 初始化用户资料 保证数据最终一致
func initUserProfile(logger fklog.FKLogI, userID uint64) (ret *UserProfile.UserProfile, err error) {
	logger.DebugWF("initUserProfile start", zap.Uint64("userID", userID))
	defer func() {
		logger.DebugWF("initUserProfile end", zap.Uint64("userID", userID), zap.Any("ret", ret), zap.Error(err))
	}()
	locked, err := userprofilelock.Lock(userID)
	if err != nil {
		logger.ErrorWF("initUserProfile get lock failed", zap.Uint64("userID", userID), zap.Error(err))
		return
	}
	if !locked {
		err = fmt.Errorf("服务器忙")
		logger.WarnWF("initUserProfile lock already exists", zap.Uint64("userID", userID))
		return
	}
	defer func() {
		_ = userprofilelock.Unlock(userID)
	}()

	if val, ok := cache.Get(userID); ok {
		logger.DebugWF("initUserProfile get from cache", zap.Uint64("userID", userID),
			zap.Any("ret", ret))
		ret = val.(*UserProfile.UserProfile)
		return
	}

	mysqlDb := userprofilemysql.GetDB()
	if mysqlDb == nil || mysqlDb.Error != nil {
		err = fmt.Errorf("initUserProfile userprofilemysql.GetDB nil, user:%v", userID)
		logger.ErrorWF("initUserProfile get db failed", zap.Uint64("userID", userID))
		return
	}
	var profile *UserProfile.UserProfile
	if err = mysqlDb.Transaction(func(tx *gorm.DB) error {
		profile = generateInitProfile(logger, userID)

		dbProfile := &userprofilemysql.UserProfile{
			UserID:    userID,
			NickName:  profile.GetNickName(),
			IconUrl:   profile.GetIconUrl(),
			Sex:       uint8(profile.GetSex()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err = tx.Create(dbProfile).Error; err != nil {
			logger.ErrorWF("initUserProfile create db profile failed", zap.Uint64("userID", userID), zap.Error(err))
			return err
		}

		cache.Add(userID, profile)
		if err = userprofileredis.SetCache(userID, profile); err != nil {
			tx.Rollback()
			cache.Remove(userID)
			logger.ErrorWF("initUserProfile set cache failed", zap.Uint64("userID", userID), zap.Error(err))
			return err
		}
		return nil
	}); err != nil {
		logger.ErrorWF("initUserProfile tx failed", zap.Uint64("userID", userID), zap.Error(err))
		cache.Remove(userID)
		err = userprofileredis.DeleteCache(userID)
		if err != nil {
			logger.ErrorWF("initUserProfile DeleteCache failed", zap.Uint64("userID", userID), zap.Error(err))
		}
		return
	}

	ret = profile
	return
}

// generateInitProfile 生成初始用户资料
func generateInitProfile(logger fklog.FKLogI, userId uint64) *UserProfile.UserProfile {
	// todo 应从配置读取
	maleAvatars := []string{"1001", "1002", "1003"}
	femaleAvatars := []string{"2001", "2002", "2003"}
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
		IconUrl:  proto.String(avatar),
		Sex:      proto.Int32(sex),
	}
	logger.DebugWF("generateInitProfile success", zap.Any("initUser", initUser))
	return initUser
}
