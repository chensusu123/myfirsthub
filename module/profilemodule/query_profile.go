// @Author pangchenyang 2025/6/19 15:24:00
// @Desc: 
package profilemodule

import (
	"maze_game_server/pb/common/UserProfile"
	"go.uber.org/zap"
	"maze_game_server/io/redis/userprofileredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

// MGetUserProfile 获取用户资料
func MGetUserProfile(logger fklog.FKLogI, userID uint64, viewIDs []uint64) ([]*UserProfile.UserProfile, error) {
	profiles, err := userprofileredis.BatchGetProfile(viewIDs)
	if err != nil {
		logger.ErrorWF("MGetUserProfile BatchGetProfile failed", zap.Uint64("userID", userID), zap.Error(err))
		return profiles, err
	}
	return profiles, nil
}

// // initUserProfile 初始化用户资料 加锁保证数据最终一致
// func initUserProfile(logger fklog.FKLogI, userID uint64) (ret *UserProfile.UserProfile, err error) {
// 	logger.DebugWF("initUserProfile start", zap.Uint64("userID", userID))
// 	defer func() {
// 		logger.DebugWF("initUserProfile end", zap.Uint64("userID", userID), zap.Any("ret", ret), zap.Error(err))
// 	}()
// 	locked, err := userprofilelock.Lock(userID)
// 	if err != nil {
// 		logger.ErrorWF("initUserProfile get lock failed", zap.Uint64("userID", userID), zap.Error(err))
// 		return
// 	}
// 	if !locked {
// 		err = fmt.Errorf("服务器忙")
// 		logger.WarnWF("initUserProfile lock already exists", zap.Uint64("userID", userID))
// 		return
// 	}
//
// 	if val, ok := cache.Get(userID); ok {
// 		logger.DebugWF("initUserProfile get from cache", zap.Uint64("userID", userID),
// 			zap.Any("ret", ret))
// 		ret = val.(*UserProfile.UserProfile)
// 		return
// 	}
//
// 	newProfile := generateInitProfile(logger, userID)
//
// 	var (
// 		chgDesc, newVal, oldVal string
// 	)
//
// 	// 修改用户资料
// 	ret = alterProfile(newProfile, nil, &chgDesc, &newVal, &oldVal)
//
// 	if err = userprofileredis.SetProfile(userID, ret); err != nil {
// 		logger.ErrorWF("initUserProfile set cache failed", zap.Uint64("userID", userID), zap.Error(err))
// 		return
// 	}
// 	cache.Add(userID, ret)
// 	// 解锁
// 	_ = userprofilelock.Unlock(userID)
// 	// 修改资料流水
// 	flowrecord.SaveAlterProfileRecord(logger, flowrecord.AlterProfileRecord{
// 		UserId:     userID,
// 		NewVal:     newVal,
// 		OldVal:     oldVal,
// 		ChgDesc:    chgDesc,
// 		CreateTime: time.Now().UnixMilli(), // ms
// 	})
// 	return
// }
//
// // generateInitProfile 生成初始用户资料
// func generateInitProfile(logger fklog.FKLogI, userId uint64) *UserProfile.UserProfile {
// 	// todo 应从配置读取
// 	maleAvatars := []string{"1001", "1002", "1003"}
// 	femaleAvatars := []string{"2001", "2002", "2003"}
// 	nicknamePrefix := []string{
// 		"Shadow", "Blood", "Dark", "Iron", "Storm",
// 		"Night", "Ghost", "Rage", "Frost", "Flame",
// 		"Thunder", "Silver", "Gold", "Steel", "Demon",
// 		"Angel", "Dragon", "Wolf", "Phoenix", "Void",
// 	}
//
// 	nicknameSuffix := []string{
// 		"blade", "fist", "heart", "fang", "wing",
// 		"hunter", "slayer", "reaver", "lord", "king",
// 		"queen", "mage", "ranger", "assassin", "warden",
// 		"bringer", "walker", "guardian", "watcher", "caller",
// 		"born", "seeker", "weaver", "shade", "born",
// 	}
//
// 	// 默认性别为男
// 	sex := int32(1)
// 	avatar := maleAvatars[0]
//
// 	rand.Seed(time.Now().UnixNano())
// 	if rand.Intn(2) == 1 { // 50%概率为女性
// 		sex = 2
// 		avatar = femaleAvatars[rand.Intn(len(femaleAvatars))]
// 	} else {
// 		avatar = maleAvatars[rand.Intn(len(maleAvatars))]
// 	}
//
// 	nickname := fmt.Sprintf("%s%s%d",
// 		nicknamePrefix[rand.Intn(len(nicknamePrefix))],
// 		nicknameSuffix[rand.Intn(len(nicknameSuffix))],
// 		userId%10000)
//
// 	initUser := &UserProfile.UserProfile{
// 		UserId:   proto.Uint64(userId),
// 		NickName: proto.String(nickname),
// 		Avatar:   proto.String(avatar),
// 		Sex:      proto.Int32(sex),
// 	}
// 	logger.DebugWF("generateInitProfile success", zap.Any("initUser", initUser))
// 	return initUser
// }
