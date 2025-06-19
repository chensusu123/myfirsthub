// @Author pangchenyang 2025/6/19 15:24:00
// @Desc: 
package profilemodule

import (
	"maze_game_server/pb/common/UserProfile"
	"sync"
	"go.uber.org/zap"
	"fmt"
	"maze_game_server/io/redis/userprofileredis"
	"maze_game_server/io/redis/userprofilelock"
	"maze_game_server/io/mysql/flowrecord"
	"time"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"google.golang.org/protobuf/proto"
)

func (m *UserProfileModule) QueryUserProfile(logger fklog.FKLogI, userID uint64, viewIDs []uint64) ([]*UserProfile.UserProfile, error) {
	var profiles []*UserProfile.UserProfile
	// 并发查询用户资料
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	for _, v := range viewIDs {
		viewID := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			ret, err := queryUserProfile(logger, userID, viewID)
			if err != nil {
				logger.ErrorWF("OnQueryUserProfile queryUserProfile failed", zap.Error(err))
				return
			}
			mu.Lock()
			profiles = append(profiles, ret)
			mu.Unlock()
		}()
	}
	wg.Wait()
	return profiles, nil
}

func queryUserProfile(logger fklog.FKLogI, userID, viewID uint64) (ret *UserProfile.UserProfile, err error) {
	if viewID <= 0 {
		err = fmt.Errorf("invalid viewID")
		logger.ErrorWF("queryUserProfile invalid userID", zap.Uint64("viewID", viewID))
		return
	}

	if val, ok := cache.Get(viewID); ok {
		ret = val.(*UserProfile.UserProfile)
		logger.DebugWF("OnQueryUserProfile get from cache", zap.Uint64("viewID", viewID),
			zap.Any("ret", ret))
		return
	}

	ret, err = userprofileredis.GetProfile(viewID)
	if err != nil {
		logger.ErrorWF("queryUserProfile GetCache error", zap.Uint64("viewID", viewID), zap.Error(err))
		return
	}
	if err == redis.Nil {
		if userID != viewID {
			err = fmt.Errorf("无效用户:%v", viewID)
			logger.WarnWF("OnQueryUserProfile query unregister user", zap.Uint64("viewID", viewID), zap.Uint64("userID", userID))
			return
		} else {
			ret, err = initUserProfile(logger, viewID)
			if err != nil {
				err = fmt.Errorf("初始化用户资料失败:%v", err)
				logger.ErrorWF("queryUserProfile initUserProfile error", zap.Uint64("userID", userID), zap.Error(err))
				return
			}
		}
	} else {
		// 更新缓存
		cache.Add(viewID, ret)
	}
	return
}

// initUserProfile 初始化用户资料 加锁保证数据最终一致
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

	if val, ok := cache.Get(userID); ok {
		logger.DebugWF("initUserProfile get from cache", zap.Uint64("userID", userID),
			zap.Any("ret", ret))
		ret = val.(*UserProfile.UserProfile)
		return
	}

	newProfile := generateInitProfile(logger, userID)

	var (
		chgDesc, newVal, oldVal string
	)

	// 修改用户资料
	ret = alterProfile(newProfile, nil, &chgDesc, &newVal, &oldVal)

	if err = userprofileredis.SetProfile(userID, ret); err != nil {
		logger.ErrorWF("initUserProfile set cache failed", zap.Uint64("userID", userID), zap.Error(err))
		return
	}
	cache.Add(userID, ret)
	// 解锁
	_ = userprofilelock.Unlock(userID)
	// 修改资料流水
	flowrecord.SaveAlterProfileRecord(logger, flowrecord.AlterProfileRecord{
		UserId:     userID,
		NewVal:     newVal,
		OldVal:     oldVal,
		ChgDesc:    chgDesc,
		CreateTime: time.Now().UnixMilli(), // ms
	})
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
		Avatar:   proto.String(avatar),
		Sex:      proto.Int32(sex),
	}
	logger.DebugWF("generateInitProfile success", zap.Any("initUser", initUser))
	return initUser
}
