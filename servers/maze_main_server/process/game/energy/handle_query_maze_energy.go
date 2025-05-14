/*
 * @Author: majian
 * @Date: 2025-03-21 19:14:33
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 14:01:38
 */
package energy

import (
	"time"

	"github.com/lonng/nano/session"
	"gitlab.ifreetalk.com/plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/plate/freetk/common/errors"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkprometheus"
	"gitlab.ifreetalk.com/plate/protodef/MazeEnergy"

	"gitlab.ifreetalk.com/maze/maze_game_server/excel/mazeconfigv8"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/kafka/mazeenergyrecord"
	"gitlab.ifreetalk.com/maze/maze_game_server/module/mazeuserinfo"

	"go.uber.org/zap"
)

func (e *Energy) OnQueryMazeEnergyRQ(s *session.Session, req *MazeEnergy.QueryMazeEnergyRQ) (err error) {
	defer fkprometheus.DebugPMT("OnQueryMazeEnergyRQ")()

	logger := fklog.AppLogger().Clone("energy")
	res := &MazeEnergy.QueryMazeEnergyRS{}

	res.ErrInfo = errors.NO_ERROR
	res.Header = req.Header

	userID := uint64(s.UID())

	defer func() {
		logger.InfoWF("OnQueryMazeEnergyRQ end", zap.Any("res", res))
	}()

	logger.InfoWF("OnQueryMazeEnergyRQ with", zap.Any("req", req))

	uInfo, err := mazeuserinfo.GetUserInfoV2(logger, userID)
	if err != nil {
		logger.ErrorWF("OnQueryMazeEnergyRQ GetUserInfoV2 fail", zap.Error(err))
		res.ErrInfo = errors.MODULE_ERROR.ToInfo()
		return
	}
	var updateFlag int32 //是否需要更新
	now := time.Now().Unix()
	maxVal := mazeconfigv8.GetEnergyMax()     // 体力最大值
	cost, val := mazeconfigv8.GetEnergyRate() // 每n秒回复多少体力
	var nextUpdateTime int64                  // 下次更新时间
	record := BeginRecord(userID, mazeenergyrecord.TimerRecovery, uInfo)
	var chgVal int32 // 变化值

	if uInfo.EnergyLastTime == 0 { // 首次初始化
		uInfo.SetEnergyLastTime(now)
		initVal := mazeconfigv8.GetEnergyInitVal()
		chgVal = initVal - uInfo.Energy
		uInfo.SetEnergy(initVal)
		updateFlag = 1
		nextUpdateTime = now + int64(cost)
		record.OpType = mazeenergyrecord.InitEnergy

	} else {
		curVal := uInfo.Energy
		if curVal < mazeconfigv8.GetEnergyMax() { // 未恢复满
			cycleNum := (now - uInfo.EnergyLastTime) / int64(cost)        // 周期数
			addVal := cycleNum * int64(val)                               // 周期数*每周期增加的体力
			lastUpdateTime := uInfo.EnergyLastTime + cycleNum*int64(cost) // 计算上次更新时间
			nextUpdateTime = lastUpdateTime + int64(cost)
			if addVal > 0 {
				chgVal = int32(addVal)
				curVal += int32(addVal)
				if curVal >= maxVal { // 如果恢复到满值,上次恢复时间设置为当前时间
					curVal = maxVal
					lastUpdateTime = now
					nextUpdateTime = now + int64(cost)
				}
				uInfo.SetEnergyLastTime(lastUpdateTime)
				uInfo.SetEnergy(curVal)
				updateFlag = 2
			}
		} else {
			nextUpdateTime = now + int64(cost)
		}
	}
	if updateFlag > 0 {
		err = mazeuserinfo.SetUserInfoV2(logger, userID, uInfo)
		if err != nil {
			logger.ErrorWF("OnQueryMazeEnergyRQ SetUserInfoV2 fail", zap.Error(err))
			res.ErrInfo = errors.DB_SAVE_ERROR.ToInfo()
			return
		}
		EndRecord(logger, record, chgVal, uInfo)
	}
	res.EnergyInfo = &MazeEnergy.EnergyInfo{
		CurVal:           proto.Int32(uInfo.Energy),
		MaxVal:           proto.Int32(maxVal),
		NextRecoveryTime: proto.Int64(nextUpdateTime)}
	return
}

func BeginRecord(userId uint64, opType int32, energyInfo *mazeuserinfo.UserInfo) *mazeenergyrecord.MazeEnergyChgRecord {
	rd := &mazeenergyrecord.MazeEnergyChgRecord{}
	rd.UserId = userId
	rd.OldVal = energyInfo.Energy
	rd.OpType = opType
	return rd
}

func EndRecord(logger fklog.FKLogI, record *mazeenergyrecord.MazeEnergyChgRecord, chgval int32, energyInfo *mazeuserinfo.UserInfo) error {
	record.NewVal = energyInfo.Energy
	record.ChgVal = chgval
	record.LastTime = energyInfo.EnergyLastTime
	return mazeenergyrecord.SendMazeEnergyChgRecord(logger, record)
}
