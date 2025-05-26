package mazeuserinfo

import (
	"errors"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/redis/mazeuserlevelredis"
)

const (
	USER_INFO_MASK_LEVEL int64 = 1 << iota
	USER_INFO_MASK_EXP
	USER_INFO_MASK_TOTAL_EXP
	USER_INFO_MASK_USER_TYPE
	USER_INFO_MASK_BARRIER
	USER_INFO_MASK_PASS_BARRIER
	USER_INFO_MASK_EQUIP_POINT
	USER_INFO_MASK_ENERGY
	USER_INFO_MASK_ENERGY_LAST_TIME
)

type UserInfo struct {
	UserId         uint64
	Level          int64 //等级
	Exp            int64 //当前经验
	TotalExp       int64 //总经验
	UserType       int32 //用户类型
	Barrier        int32 //当前关卡id
	PassBarrier    int32 //最高通关关卡id
	EquipPoint     int32 //装备积分
	Energy         int32 //当前体力值
	EnergyLastTime int64 //上次更新体力的时间戳
	mask           int64 //发生变化的数值掩码
}

// 设置等级
func (u *UserInfo) SetLevel(level int64) {
	u.Level = level
	u.mask |= USER_INFO_MASK_LEVEL
}

// 设置经验
func (u *UserInfo) SetExp(exp int64) {
	u.Exp = exp
	u.mask |= USER_INFO_MASK_EXP
}

// 设置经验总值
func (u *UserInfo) SetTotalExp(totalExp int64) {
	u.TotalExp = totalExp
	u.mask |= USER_INFO_MASK_TOTAL_EXP
}

// 加经验
func (u *UserInfo) AddExp(addExp int64) (err error) {

	for {
		levelCfg := GMazeLevelV8Cfg.Get(int32(u.Level))
		if levelCfg == nil {
			return errors.New("cant find level cfg")
		}

		if u.TotalExp+addExp < levelCfg.All_exp {
			u.Level = int64(levelCfg.Order)
			u.mask |= USER_INFO_MASK_LEVEL
			u.Exp = u.TotalExp + addExp - levelCfg.All_exp + levelCfg.Next_level_need_exp
			u.mask |= USER_INFO_MASK_EXP
			break
		}

		u.Level++
	}

	u.TotalExp += addExp
	u.mask |= USER_INFO_MASK_TOTAL_EXP
	return
}

// 根据设置的经验总值更新等级经验
func (u *UserInfo) CalExp() (err error) {
	curLevel := int32(1)
	for {
		levelCfg := GMazeLevelV8Cfg.Get(int32(curLevel))
		if levelCfg == nil {
			return errors.New("cant find level cfg")
		}

		if u.TotalExp < levelCfg.All_exp {
			u.Level = int64(levelCfg.Order)
			u.mask |= USER_INFO_MASK_LEVEL

			u.Exp = u.TotalExp - levelCfg.All_exp + levelCfg.Next_level_need_exp
			u.mask |= USER_INFO_MASK_EXP
			break
		}

		curLevel++
	}
	return
}

// 设置用户类型
func (u *UserInfo) SetUserType(userType int32) {
	u.UserType = userType
	u.mask |= USER_INFO_MASK_USER_TYPE
}

// 设置当前关卡id
func (u *UserInfo) SetBarrier(barrierId int32) {
	u.Barrier = barrierId
	u.mask |= USER_INFO_MASK_BARRIER
}

// 设置通关关卡id
func (u *UserInfo) SetPassBarrier(barrierId int32) {
	u.PassBarrier = barrierId
	u.mask |= USER_INFO_MASK_PASS_BARRIER
}

// 设置装备积分
func (u *UserInfo) SetEquipPoint(equipPoint int32) {
	u.EquipPoint = equipPoint
	u.mask |= USER_INFO_MASK_EQUIP_POINT
}

// 设置体力值
func (u *UserInfo) SetEnergy(energy int32) {
	u.Energy = energy
	u.mask |= USER_INFO_MASK_ENERGY
}

// 设置体力更新时间
func (u *UserInfo) SetEnergyLastTime(lastTime int64) {
	u.EnergyLastTime = lastTime
	u.mask |= USER_INFO_MASK_ENERGY_LAST_TIME
}

func GetUserInfoV2(logger fklog.FKLogI, userId uint64) (userInfo *UserInfo, err error) {
	userMap, err := mazeuserlevelredis.GetUserInfo(logger, userId)
	if err != nil {
		return
	}

	if userMap == nil {
		userMap = map[string]int64{}
	}

	userInfo = &UserInfo{
		UserId:         userId,
		Level:          userMap[mazeuserlevelredis.USER_LEVEL],
		Exp:            userMap[mazeuserlevelredis.USER_EXP],
		TotalExp:       userMap[mazeuserlevelredis.USER_TOTAL_EXP],
		UserType:       int32(userMap[mazeuserlevelredis.USER_TYPE]),
		Barrier:        int32(userMap[mazeuserlevelredis.USER_BARRIER]),
		PassBarrier:    int32(userMap[mazeuserlevelredis.USER_PASS_BARRIER]),
		EquipPoint:     int32(userMap[mazeuserlevelredis.USER_EQUIP_POINT]),
		Energy:         int32(userMap[mazeuserlevelredis.USER_ENERGY]),
		EnergyLastTime: userMap[mazeuserlevelredis.USER_ENERGY_LAST_TIME],
	}

	return
}

func SetUserInfoV2(logger fklog.FKLogI, userId uint64, userInfo *UserInfo) (err error) {

	if userInfo == nil {
		logger.ErrorWF("SetUserInfoV2 userInfo is nil")
		err = errors.New("userInfo nil")
		return
	}
	userMap := make(map[string]int64)
	if userInfo.mask&USER_INFO_MASK_LEVEL > 0 {
		userMap[mazeuserlevelredis.USER_LEVEL] = userInfo.Level
	}
	if userInfo.mask&USER_INFO_MASK_EXP > 0 {
		userMap[mazeuserlevelredis.USER_EXP] = userInfo.Exp
	}
	if userInfo.mask&USER_INFO_MASK_TOTAL_EXP > 0 {
		userMap[mazeuserlevelredis.USER_TOTAL_EXP] = userInfo.TotalExp
	}
	if userInfo.mask&USER_INFO_MASK_USER_TYPE > 0 {
		userMap[mazeuserlevelredis.USER_TYPE] = int64(userInfo.UserType)
	}
	if userInfo.mask&USER_INFO_MASK_BARRIER > 0 {
		userMap[mazeuserlevelredis.USER_BARRIER] = int64(userInfo.Barrier)
	}
	if userInfo.mask&USER_INFO_MASK_PASS_BARRIER > 0 {
		userMap[mazeuserlevelredis.USER_PASS_BARRIER] = int64(userInfo.PassBarrier)
	}
	// if userInfo.mask&USER_INFO_MASK_EQUIP_POINT > 0 {
	// 	userMap[mazeuserlevelredis.USER_EQUIP_POINT] = int64(userInfo.EquipPoint)
	// }
	if userInfo.mask&USER_INFO_MASK_ENERGY > 0 {
		userMap[mazeuserlevelredis.USER_ENERGY] = int64(userInfo.Energy)
	}
	if userInfo.mask&USER_INFO_MASK_ENERGY_LAST_TIME > 0 {
		userMap[mazeuserlevelredis.USER_ENERGY_LAST_TIME] = int64(userInfo.EnergyLastTime)
	}

	if userInfo.mask == 0 {
		return
	}

	err = mazeuserlevelredis.SetUserInfo(logger, userId, userMap)
	if err != nil {
		return
	}

	userInfo.mask = 0
	return
}
