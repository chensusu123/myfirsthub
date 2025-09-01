package userinfomodel

import (
	"context"
	"errors"
	"maze_game_server/config/GMazeLevelV8Cfg"
	"maze_game_server/io/redis/mazeuserlevelredis"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var (
	ErrUserLevelConfigNotFound = errors.New("user level config not found")
)

const (
	UserLevel          = "level"
	UserExp            = "exp"
	UserTotalExp       = "total_exp"
	UserType           = "user_type"
	UserBarrier        = "barrier"
	UserPassBarrier    = "pass_barrier"
	UserFixedBarrier   = "fixed_barrier"
	UserEquipPoint     = "equip_point"
	UserEnergy         = "energy"
	UserEnergyLastTime = "energy_last_time"
)

type UserInfoModel struct {
	UserID         uint64 // 用户ID
	Level          int64  `json:"level,omitempty"`            //等级
	Exp            int64  `json:"exp,omitempty"`              //当前经验
	TotalExp       int64  `json:"total_exp,omitempty"`        //总经验
	UserType       int32  `json:"user_type,omitempty"`        //用户类型
	Barrier        int32  `json:"barrier,omitempty"`          //当前关卡id
	PassBarrier    int32  `json:"pass_barrier,omitempty"`     //最高通关关卡id
	FixedBarrier   int32  `json:"fixed_barrier,omitempty"`    // 锁定关卡id
	EquipPoint     int32  `json:"equip_point,omitempty"`      //装备积分
	Energy         int32  `json:"energy,omitempty"`           //当前体力值
	EnergyLastTime int64  `json:"energy_last_time,omitempty"` //上次更新体力的时间戳
}

// NewUserInfoModel
func NewUserInfoModel(ctx context.Context, userID uint64) (u *UserInfoModel, err error) {
	u = &UserInfoModel{
		UserID: userID,
	}
	if err = u.load(ctx); err != nil {
		return nil, err
	}
	return u, nil
}

// load
func (u *UserInfoModel) load(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	userMap, err := mazeuserlevelredis.GetUserInfo(ctx, u.UserID)
	if err != nil {
		logger.CtxError(ctx, "load GetUserInfo fail", zap.Error(err), zap.Uint64("userID", u.UserID))
		return
	}
	u.Level = userMap[UserLevel]
	u.Exp = userMap[UserExp]
	u.TotalExp = userMap[UserTotalExp]
	u.UserType = int32(userMap[UserType])
	u.Barrier = int32(userMap[UserBarrier])
	u.PassBarrier = int32(userMap[UserPassBarrier])
	u.FixedBarrier = int32(userMap[UserFixedBarrier])
	u.EquipPoint = int32(userMap[UserEquipPoint])
	u.Energy = int32(userMap[UserEnergy])
	u.EnergyLastTime = userMap[UserEnergyLastTime]
	return
}

// AddExp 加经验
func (u *UserInfoModel) AddExp(ctx context.Context, addExp int64) (err error) {

	for {
		levelCfg := GMazeLevelV8Cfg.GetWithCtx(ctx, int32(u.Level))
		if levelCfg == nil {
			return ErrUserLevelConfigNotFound
		}

		if u.TotalExp+addExp < levelCfg.All_exp {
			u.Level = int64(levelCfg.Order)
			u.Exp = u.TotalExp + addExp - levelCfg.All_exp + levelCfg.Next_level_need_exp
			break
		}

		u.Level++
	}

	u.TotalExp += addExp
	return
}

// CalExp 根据设置的经验总值更新等级经验
func (u *UserInfoModel) CalExp(ctx context.Context) (err error) {
	curLevel := int32(1)
	for {
		levelCfg := GMazeLevelV8Cfg.GetWithCtx(ctx, int32(curLevel))
		if levelCfg == nil {
			return ErrUserLevelConfigNotFound
		}

		if u.TotalExp < levelCfg.All_exp {
			u.Level = int64(levelCfg.Order)
			u.Exp = u.TotalExp - levelCfg.All_exp + levelCfg.Next_level_need_exp
			break
		}

		curLevel++
	}
	return
}

// Save
func (u *UserInfoModel) Save(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	userMap := map[string]int64{
		UserLevel:          u.Level,
		UserExp:            u.Exp,
		UserTotalExp:       u.TotalExp,
		UserType:           int64(u.UserType),
		UserBarrier:        int64(u.Barrier),
		UserPassBarrier:    int64(u.PassBarrier),
		UserFixedBarrier:   int64(u.FixedBarrier),
		UserEquipPoint:     int64(u.EquipPoint),
		UserEnergy:         int64(u.Energy),
		UserEnergyLastTime: u.EnergyLastTime,
	}
	err = mazeuserlevelredis.SetUserInfo(ctx, u.UserID, userMap)
	if err != nil {
		logger.CtxError(ctx, "load SetUserInfo fail", zap.Error(err), zap.Uint64("userID", u.UserID), zap.Any("userMap", userMap))
	}
	return
}
