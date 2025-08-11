package barrierenergyservice

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/usecase/online"
	"sync"
	"time"
)

// 审核版本用于自动恢复体力
var (
	mu      sync.Mutex
	userMap = make(map[uint64]*time.Timer)
)

// 启动用户自动恢复体力
func (s service) startUserRecoverEnergy(logger fklog.FKLogI, userId uint64, nextUpdateTime int64) {
	mu.Lock()
	defer mu.Unlock()

	// 若已有 timer，跳过
	if timer, ok := userMap[userId]; ok {
		//timer.Stop()
		_ = timer
		return
	}

	// 设置体力自动恢复时间
	recoverTime := nextUpdateTime - time.Now().Unix()
	timer := time.AfterFunc(time.Duration(recoverTime)*time.Second, func() {
		s.safeTimer(logger, userId)
	})
	userMap[userId] = timer
	logger.InfoWF("startUserRecoverEnergy success", zap.Any("userId", userId), zap.Int64("recoverTime", recoverTime), zap.Any("timer", timer))
}

func (s service) safeTimer(logger fklog.FKLogI, userID uint64) {
	defer func() {
		if r := recover(); r != nil {
			logger.ErrorWF("handleRecoverUserEnergy panic.", zap.Any("r", r))
		}
	}()
	isOnline := online.IsOnline(userID)
	if !isOnline {
		s.stopUserRecoverTimer(logger, userID)
		return
	}
	s.handleRecoverUserEnergy(logger, userID)
}

// 自动恢复体力
func (s service) handleRecoverUserEnergy(logger fklog.FKLogI, userID uint64) {

	defer func() {
		// 重新设置timer
		mu.Lock()
		defer mu.Unlock()
		if old, ok := userMap[userID]; ok {
			old.Stop()
		}

		nextTriggerTime := GetEnergyRecoverCfg()
		timer := time.AfterFunc(time.Duration(nextTriggerTime)*time.Second, func() {
			s.safeTimer(logger, userID)
		})
		userMap[userID] = timer
		logger.InfoWF("handleRecoverUserEnergy add timer", zap.Any("userID", userID), zap.Int64("nextTriggerTime", nextTriggerTime), zap.Any("timer", timer))
	}()

	curEnergy, nextTime, err := s.calEnergy(logger, userID)
	if err != nil {
		logger.InfoWF("handleRecoverUserEnergy calEnergy failed", zap.Any("userID", userID), zap.Error(err))
		return
	}

	err = s.SendEnergyChgPack(logger, userID, curEnergy, nextTime)
	if err != nil {
		err = nil
		logger.ErrorWF("handleRecoverUserEnergy SendEnergyChgPack fail", zap.Error(err), zap.Any("userID", userID), zap.Int32("curEnergy", curEnergy), zap.Int64("nextTime", nextTime))
		//return
	}
}

func (s service) stopUserRecoverTimer(logger fklog.FKLogI, userID uint64) {
	mu.Lock()
	defer mu.Unlock()
	if timer, ok := userMap[userID]; ok {
		timer.Stop()
		delete(userMap, userID)
		logger.InfoWF("StopUserRecoverTimer success", zap.Uint64("userId", userID), zap.Any("timer", timer))
	}
}
