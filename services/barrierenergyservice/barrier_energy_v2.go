package barrierenergyservice

//
//import (
//	"errors"
//	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
//	"go.uber.org/zap"
//	"maze_game_server/config/GMazeConfigV8Cfg"
//	"maze_game_server/model/barrierenergymodel"
//	"maze_game_server/usecase/online"
//	"sync"
//	"time"
//)
//
//func (s service) GetBarrierEnergy(logger fklog.FKLogI, userId uint64) (*barrierenergymodel.BarrierEnergyModel, error) {
//	energyInfo, err := barrierenergymodel.NewBarrierEnergyModel(logger, userId)
//	if err != nil {
//		logger.ErrorWF("GetBarrierEnergy NewBarrierEnergyModel failed", zap.Uint64("userId", userId), zap.Error(err))
//		return nil, err
//	}
//
//	if energyInfo == nil {
//		initEnergy := getEnergyInitCfg()
//		energyInfo = &barrierenergymodel.BarrierEnergyModel{
//			Energy:         initEnergy,
//			EnergyLastTime: time.Now().Unix(),
//		}
//		err = energyInfo.Save(logger, userId)
//		if err != nil {
//			logger.ErrorWF("GetBarrierEnergy init BarrierEnergyModel Save failed", zap.Uint64("userId", userId), zap.Error(err))
//			return nil, err
//		}
//	} else {
//		delay := getRecoverEnergyCfg()
//		diffTime := time.Now().Unix() - energyInfo.EnergyLastTime
//		times := diffTime / delay
//		if times > 0 {
//			energyInfo.Energy += int32(times) * 1
//			energyInfo.EnergyLastTime += delay * times
//			if energyInfo.Energy >= getEnergyMaxCfg() {
//				energyInfo.Energy = getEnergyMaxCfg()
//			}
//			err = energyInfo.Save(logger, userId)
//			if err != nil {
//				logger.ErrorWF("GetBarrierEnergy BarrierEnergyModel Save failed", zap.Uint64("userId", userId), zap.Error(err))
//				return nil, err
//			}
//		}
//	}
//
//	startUserRecoverEnergy(logger, userId, energyInfo.EnergyLastTime)
//	return energyInfo, err
//}
//
//// 增加体力，结果大于最大值时体力等于最大值
//func (s service) AddEnergy(logger fklog.FKLogI, userId uint64, addVal int32) (*barrierenergymodel.BarrierEnergyModel, error) {
//	energyInfo, err := s.GetBarrierEnergy(logger, userId)
//	if err != nil {
//		return nil, err
//	}
//	oldVal := energyInfo.Energy
//	energyInfo.Energy += addVal
//	maxEnergy := getEnergyMaxCfg()
//	if energyInfo.Energy >= maxEnergy {
//		energyInfo.Energy = maxEnergy
//	}
//	err = energyInfo.Save(logger, userId)
//	if err != nil {
//		logger.ErrorWF("AddEnergy SetEnergyModel failed", zap.Uint64("userId", userId), zap.Error(err), zap.Int32("addVal", addVal), zap.Any("energyInfo", energyInfo))
//		return nil, err
//	}
//	logger.InfoWF("AddEnergy success", zap.Uint64("userId", userId), zap.Int32("addVal", addVal), zap.Int32("oldVal", oldVal), zap.Any("energyInfo", energyInfo))
//
//	// 通知客户端
//	SendEnergyChgPack(logger, userId, energyInfo.Energy, energyInfo.EnergyLastTime+getRecoverEnergyCfg()-time.Now().Unix())
//
//	return energyInfo, err
//}
//
//// 减少体力
//func (s service) SubEnergy(logger fklog.FKLogI, userId uint64, subVal int32) (int32, error) {
//	energyInfo, err := s.GetBarrierEnergy(logger, userId)
//	if err != nil {
//		return 0, err
//	}
//	if energyInfo.Energy < subVal {
//		logger.InfoWF("SubEnergy energy is less than v", zap.Uint64("userID", userId), zap.Int32("subVal", subVal), zap.Any("energyInfo", energyInfo))
//		return 0, errors.New("energy not enough")
//	}
//
//	if energyInfo.Energy == getEnergyMaxCfg() {
//		// 如果之前体力是最大值，说明需要从这个时间点开始恢复体力
//		energyInfo.EnergyLastTime = time.Now().Unix()
//	}
//	oldVal := energyInfo.Energy
//	energyInfo.Energy = energyInfo.Energy - subVal
//
//	err = energyInfo.Save(logger, userId)
//	if err != nil {
//		logger.ErrorWF("SubEnergy SetEnergyModel failed", zap.Uint64("userId", userId), zap.Error(err), zap.Int32("subVal", subVal), zap.Any("energyInfo", energyInfo))
//		return 0, err
//	}
//	logger.InfoWF("SubEnergy success", zap.Uint64("userId", userId), zap.Int32("subVal", subVal), zap.Int32("oldVal", oldVal), zap.Any("energyInfo", energyInfo))
//
//	// 通知客户端
//	SendEnergyChgPack(logger, userId, energyInfo.Energy, energyInfo.EnergyLastTime+getRecoverEnergyCfg()-time.Now().Unix())
//
//	return energyInfo.Energy, err
//}
//
//const EnergyItemId = 44100001
//
//// 自动恢复体力
//var (
//	mu      sync.Mutex
//	userMap = make(map[uint64]*time.Timer)
//)
//
//// 启动用户自动恢复体力
//func startUserRecoverEnergy(logger fklog.FKLogI, userId uint64, lastRecoverTime int64) {
//	mu.Lock()
//	defer mu.Unlock()
//
//	// 若已有 timer，跳过
//	if timer, ok := userMap[userId]; ok {
//		//timer.Stop()
//		_ = timer
//		return
//	}
//	nextRecoverTime := lastRecoverTime + getRecoverEnergyCfg() - time.Now().Unix()
//	// 设置体力自动恢复时间
//	timer := time.AfterFunc(time.Duration(nextRecoverTime)*time.Second, func() {
//		safeTimer(logger, userId)
//	})
//
//	userMap[userId] = timer
//}
//
//func safeTimer(logger fklog.FKLogI, userID uint64) {
//	defer func() {
//		if r := recover(); r != nil {
//			logger.ErrorWF("handleRecoverUserEnergy panic.", zap.Any("r", r))
//		}
//	}()
//	handleRecoverUserEnergy(logger, userID)
//}
//
//// 自动恢复体力
//func handleRecoverUserEnergy(logger fklog.FKLogI, userID uint64) {
//	var energyInfo *barrierenergymodel.BarrierEnergyModel
//	var err error
//	defer func() {
//		// 重新设置timer
//		mu.Lock()
//		defer mu.Unlock()
//		if old, ok := userMap[userID]; ok {
//			old.Stop()
//		}
//		var nextTriggerTime int64
//		if energyInfo == nil {
//			// 说明失败了，设置个短暂时间重试
//			nextTriggerTime = 10
//		} else {
//			if energyInfo.Energy == getEnergyMaxCfg() {
//				// 体力已满就设置timer为获取体力的最长时间
//				nextTriggerTime = getRecoverEnergyCfg()
//			} else {
//				nextTriggerTime = energyInfo.EnergyLastTime + getRecoverEnergyCfg() - time.Now().Unix()
//			}
//		}
//		timer := time.AfterFunc(time.Duration(nextTriggerTime)*time.Second, func() {
//			safeTimer(logger, userID)
//		})
//		userMap[userID] = timer
//	}()
//	energyInfo, err = barrierenergymodel.NewBarrierEnergyModel(logger, userID)
//	if err != nil {
//		logger.ErrorWF("handleRecoverUserEnergy GetEnergyModel failed", zap.Error(err), zap.Uint64("userId", userID))
//		return
//	}
//
//	maxEnergy := getEnergyMaxCfg()
//	if energyInfo.Energy >= maxEnergy {
//		// 用户体力已经满了不需要恢复
//		logger.InfoWF("handleRecoverUserEnergy user energy full", zap.Uint64("userId", userID), zap.Any("energyInfo", energyInfo))
//		return
//	}
//	delay := getRecoverEnergyCfg()
//	if time.Now().Unix() < energyInfo.EnergyLastTime+delay {
//		// 如果当前时间小于下次恢复时间不需要恢复
//		logger.InfoWF("handleRecoverUserEnergy current time small next recover time", zap.Uint64("userId", userID), zap.Any("energyInfo", energyInfo), zap.Any("delay", delay))
//		return
//	}
//	// 恢复1点
//	energyInfo.Energy = energyInfo.Energy + 1
//	energyInfo.EnergyLastTime = energyInfo.EnergyLastTime + delay
//	err = energyInfo.Save(logger, userID)
//	if err != nil {
//		logger.ErrorWF("handleRecoverUserEnergy SetEnergyModel failed", zap.Uint64("userId", userID), zap.Error(err))
//		return
//	}
//	// 通知客户端体力恢复
//	SendEnergyChgPack(logger, userID, energyInfo.Energy, energyInfo.EnergyLastTime+getRecoverEnergyCfg()-time.Now().Unix())
//}
//
//func SendEnergyChgPack(logger fklog.FKLogI, userId uint64, curEnergy int32, nextRecoverTime int64) {
//	maxEnergy := GetEnergyMaxCfg()
//	if curEnergy == maxEnergy {
//		nextRecoverTime = -1
//	}
//	energyPack := &DollMazeBarrier.MazeEnergyChgID{
//		CurEnergy:        proto.Int32(curEnergy),
//		MaxEnergy:        proto.Int32(maxEnergy),
//		NextRecoveryTime: proto.Int64(nextRecoverTime),
//	}
//	err := online.Push(logger, userId, 16668, energyPack)
//	//err := commonmustarriveredis.SendArrivePacketFix(userId, 16668, energyPack)
//	if err != nil {
//		logger.ErrorWF("SendEnergyChgPack send client failed", zap.Uint64("userID", userId), zap.Error(err), zap.Any("energyPack", energyPack))
//	} else {
//		logger.InfoWF("SendEnergyChgPack send client success", zap.Uint64("userID", userId), zap.Any("energyPack", energyPack))
//	}
//	return
//}
//
//// 体力上限
//func getEnergyMaxCfg() int32 {
//	row := GMazeConfigV8Cfg.Get(100001)
//	if row != nil {
//		return int32(row.Value_int)
//	}
//	return 100
//}
//
//// 迷宫初始体力
//func getEnergyInitCfg() int32 {
//	row := GMazeConfigV8Cfg.Get(100002)
//	if row != nil {
//		return int32(row.Value_int)
//	}
//	return 100
//}
//
//// 迷宫1点体力回复间隔时间（秒）
//func getRecoverEnergyCfg() int64 {
//	row := GMazeConfigV8Cfg.Get(100003)
//	if row != nil {
//		return row.Value_int
//	}
//	return 30
//}
//
//// 迷宫体力瓶道具id：对应的体力数量
//func getEnergyItemCfg() int32 {
//	row := GMazeConfigV8Cfg.Get(100004)
//	if row != nil {
//		v := row.Value_map[EnergyItemId]
//		if v == 0 {
//			return 50
//		}
//		return int32(v)
//	}
//	return 50
//}
