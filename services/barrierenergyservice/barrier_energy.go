package barrierenergyservice

import (
	"context"
	"errors"
	"time"

	"maze_game_server/common/constdef"
	"maze_game_server/config/GMazeConfigV8Cfg"
	"maze_game_server/io/kafka/mazeenergyrecord"
	"maze_game_server/module/mazeuserinfo"
	"maze_game_server/pb/common/MazeEnergy"
	"maze_game_server/usecase/online"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s service) GetBarrierEnergy(logger fklog.FKLogI, userId uint64) (curEnergy int32, nextTime int64, err error) {
	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("GetBarrierEnergy GetUserInfoV2 fail", zap.Error(err))
		return 0, 0, errors.New("userInfo not find")
	}

	curEnergy, nextUpdateTime, err := s.calEnergy(logger, userId)
	if err != nil {
		logger.ErrorWF("GetBarrierEnergy calEnergy fail", zap.Error(err))
		return uInfo.Energy, uInfo.EnergyLastTime, err
	}

	// 服务器添加定时器,补发ID包
	s.startUserRecoverEnergy(logger, userId, nextUpdateTime)

	logger.InfoWF("GetBarrierEnergy success", zap.Any("userId", userId), zap.Any("curEnergy", uInfo.Energy), zap.Any("nextUpdateTime", nextUpdateTime))
	return curEnergy, nextUpdateTime, err
}

// 增加体力，结果大于最大值时体力等于最大值
func (s service) AddEnergy(logger fklog.FKLogI, userId uint64, addVal int32) (curEnergy int32, nextTime int64, err error) {
	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("AddEnergy GetUserInfoV2 fail", zap.Error(err))
		return 0, 0, errors.New("userInfo not find")
	}
	maxVal := GetEnergyMax() // 体力最大值

	cur, nextTime, err := s.GetBarrierEnergy(logger, userId)
	if err != nil {
		logger.ErrorWF("AddEnergy GetBarrierEnergy fail", zap.Error(err))
		return uInfo.Energy, uInfo.EnergyLastTime, err
	}

	remain := cur + addVal
	if remain >= maxVal { // 如果加到满值，更新上次恢复时间
		remain = maxVal
		uInfo.SetEnergyLastTime(time.Now().Unix())
	}
	uInfo.SetEnergy(remain)
	err = mazeuserinfo.SetUserInfoV2(logger, userId, uInfo)
	if err != nil {
		logger.ErrorWF("AddEnergy SetUserInfoV2 fail", zap.Error(err), zap.Any("uInfo", uInfo))
		return uInfo.Energy, uInfo.EnergyLastTime, errors.New("userInfo not find")
	}
	curEnergy = remain
	logger.InfoWF("AddEnergy success", zap.Any("userId", userId), zap.Any("uInfo", uInfo), zap.Any("nextUpdateTime", nextTime), zap.Any("addVal", addVal), zap.Any("curEnergy", curEnergy))
	return curEnergy, nextTime, err
}

// 减少体力
func (s service) SubEnergy(logger fklog.FKLogI, userId uint64, subVal int32) (int32, error) {
	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("AddEnergy GetUserInfoV2 fail", zap.Error(err))
		return 0, errors.New("userInfo not find")
	}

	curEnergy, _, err := s.GetBarrierEnergy(logger, userId)
	if err != nil {
		logger.ErrorWF("SubEnergy GetBarrierEnergy fail", zap.Error(err))
		return uInfo.Energy, err
	}

	if curEnergy < subVal {
		logger.WarnWF("SubEnergy energy less", zap.Int32("has", uInfo.Energy), zap.Int32("need", subVal))
		return uInfo.Energy, errors.New("energy not enough")
	}

	remain := curEnergy - subVal
	if uInfo.Energy >= s.GetEnergyMaxValue() { // 如果满值时扣除，更新上次恢复时间
		uInfo.SetEnergyLastTime(time.Now().Unix())
	}
	uInfo.SetEnergy(remain)
	err = mazeuserinfo.SetUserInfoV2(logger, userId, uInfo)
	if err != nil {
		logger.ErrorWF("SubEnergy SetUserInfoV2 fail", zap.Error(err), zap.Any("uInfo", uInfo))
		return uInfo.Energy, errors.New("userInfo not find")
	}

	err = s.SendEnergyChgPack(logger, userId, remain, uInfo.EnergyLastTime)
	if err != nil {
		logger.ErrorWF("SubEnergy SendEnergyChgPack fail", zap.Error(err), zap.Any("uInfo", uInfo))
		err = nil
		// return uInfo.Energy, err
	}
	logger.InfoWF("SubEnergy success", zap.Any("userId", userId), zap.Any("uInfo", uInfo), zap.Any("curEnergy", remain), zap.Int32("subVal", subVal))
	return uInfo.Energy, err
}

func (s service) GetEnergyMaxValue() int32 {
	return GetEnergyMax()
}

func (s service) GetEnergyRecoverCfg() int64 {
	return GetEnergyRecoverCfg()
}

func (s service) GetEnergyItemCfg() (int32, int32) {
	return GetEnergyItemCfg()
}

// 推送体力变化ID包
func (s service) SendEnergyChgPack(logger fklog.FKLogI, userId uint64, curEnergy int32, nextRecoverTime int64) error {
	maxEnergy := GetEnergyMax() // 体力最大值
	if curEnergy == maxEnergy {
		nextRecoverTime = -1
	}
	energyPack := &MazeEnergy.EnergyChangeID{
		EnergyInfo: &MazeEnergy.EnergyInfo{
			CurVal:           proto.Int32(curEnergy),
			MaxVal:           proto.Int32(maxEnergy),
			NextRecoveryTime: proto.Int64(nextRecoverTime),
		},
	}
	err := online.ClusterPush(context.TODO(), userId, 10610, energyPack)
	if err != nil {
		logger.ErrorWF("SendEnergyChgPack send client failed", zap.Uint64("userID", userId), zap.Error(err), zap.Any("energyPack", energyPack))
	} else {
		logger.InfoWF("SendEnergyChgPack send client success", zap.Uint64("userID", userId), zap.Any("energyPack", energyPack))
	}
	return err
}

// 体力初始值
func GetEnergyInitVal() int32 {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId302)
	if row != nil {
		return int32(row.Value_int)
	}
	return 200
}

// 体力上限
func GetEnergyMax() int32 {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId301)
	if row != nil {
		return int32(row.Value_int)
	}
	return 200
}

// 体力恢复效率
func GetEnergyRate() (costTime, recoverVal int32) {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId303)
	if row != nil {
		for k, v := range row.Value_map {
			return k, int32(v)
		}
	}
	return 90, 1
}

// 迷宫体力回复间隔时间（秒）
func GetEnergyRecoverCfg() int64 {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId303)
	if row != nil {
		for k := range row.Value_map {
			return int64(k)
		}
	}
	return 30
}

// 迷宫体力瓶道具id：对应的体力数量
func GetEnergyItemCfg() (id, count int32) {
	row := GMazeConfigV8Cfg.GetMazeConfigV8Config(constdef.MazeCfgId304)
	if row != nil {
		for k, v := range row.Value_map {
			return k, int32(v)
		}
	}
	return 49000001, 60
}

// 计算体力
func (s service) calEnergy(logger fklog.FKLogI, userId uint64) (curEnergy int32, nextTime int64, err error) {
	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("GetBarrierEnergy GetUserInfoV2 fail", zap.Error(err))
		return 0, 0, errors.New("userInfo not find")
	}

	oldEnergy := uInfo.Energy
	var updateFlag int32 // 是否需要更新
	now := time.Now().Unix()
	maxVal := GetEnergyMax()     // 体力最大值
	cost, val := GetEnergyRate() // 每n秒回复多少体力
	var nextUpdateTime int64     // 下次更新时间
	// var chgVal int32             // 变化值

	if uInfo.EnergyLastTime == 0 { // 首次初始化
		uInfo.SetEnergyLastTime(now)
		initVal := GetEnergyInitVal()
		// chgVal = initVal - uInfo.Energy
		uInfo.SetEnergy(initVal)
		updateFlag = 1
		nextUpdateTime = now + int64(cost)
		s.PushEnergyRecord(logger, userId, oldEnergy, uInfo.Energy, mazeenergyrecord.InitEnergy, nextUpdateTime)

	} else {
		curVal := uInfo.Energy
		if curVal < GetEnergyMax() { // 未恢复满
			cycleNum := (now - uInfo.EnergyLastTime) / int64(cost)        // 周期数
			addVal := cycleNum * int64(val)                               // 周期数*每周期增加的体力
			lastUpdateTime := uInfo.EnergyLastTime + cycleNum*int64(cost) // 计算上次更新时间
			nextUpdateTime = lastUpdateTime + int64(cost)
			if addVal > 0 {
				// chgVal = int32(addVal)
				curVal += int32(addVal)
				if curVal >= maxVal { // 如果恢复到满值,上次恢复时间设置为当前时间
					curVal = maxVal
					lastUpdateTime = now
					nextUpdateTime = now + int64(cost)
				}
				uInfo.SetEnergyLastTime(lastUpdateTime)
				uInfo.SetEnergy(curVal)
				updateFlag = 2
				s.PushEnergyRecord(logger, userId, oldEnergy, uInfo.Energy, mazeenergyrecord.TimerRecovery, nextUpdateTime)
			}
		} else {
			nextUpdateTime = now + int64(cost)
		}
	}
	if updateFlag > 0 {
		err = mazeuserinfo.SetUserInfoV2(logger, userId, uInfo)
		if err != nil {
			logger.ErrorWF("OnQueryMazeEnergyRQ SetUserInfoV2 fail", zap.Error(err))
			return uInfo.Energy, nextUpdateTime, errors.New("保存数据错误")
		}
		// nextTime := uInfo.EnergyLastTime + GetEnergyRecoverCfg() - time.Now().Unix()
		err := s.SendEnergyChgPack(logger, userId, uInfo.Energy, nextUpdateTime)
		if err != nil {
			err = nil
			logger.ErrorWF("OnQueryMazeEnergyRQ SendEnergyChgPack fail", zap.Error(err))
			// return uInfo.Energy, nextUpdateTime, err
		}
	}
	return uInfo.Energy, nextUpdateTime, nil
}

func (s service) ResetEnergy(logger fklog.FKLogI, userId uint64) (err error) {
	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userId)
	if err != nil {
		logger.ErrorWF("ResetEnergy GetUserInfoV2 fail", zap.Error(err))
		return errors.New("userInfo not find")
	}
	oldEnergy := uInfo.Energy

	now := time.Now().Unix()
	cost, _ := GetEnergyRate() // 每n秒回复多少体力

	uInfo.SetEnergyLastTime(now)
	initVal := GetEnergyInitVal()
	uInfo.SetEnergy(initVal)
	nextUpdateTime := now + int64(cost)

	err = mazeuserinfo.SetUserInfoV2(logger, userId, uInfo)
	if err != nil {
		logger.ErrorWF("ResetEnergy SetUserInfoV2 fail", zap.Error(err))
		return errors.New("保存数据错误")
	}
	s.SendEnergyChgPack(logger, userId, uInfo.Energy, nextUpdateTime)

	s.PushEnergyRecord(logger, userId, oldEnergy, uInfo.Energy, mazeenergyrecord.Reset, nextUpdateTime)

	return nil
}

func (s service) PushEnergyRecord(logger fklog.FKLogI, userId uint64, oldEnergy, newEnergy, opType int32, lastTime int64) {
	if oldEnergy != newEnergy {
		record := &mazeenergyrecord.MazeEnergyChgRecord{
			UserId:   userId,
			OldVal:   oldEnergy,
			NewVal:   newEnergy,
			LastTime: lastTime,
			OpType:   opType,
		}
		err := mazeenergyrecord.PushMazeEnergyChgRecord(logger, record)
		if err != nil {
			logger.ErrorWF("OnUseMazeEnergyItemRQ PushMazeEnergyChgRecord fail", zap.Error(err))
		}
	}
}
