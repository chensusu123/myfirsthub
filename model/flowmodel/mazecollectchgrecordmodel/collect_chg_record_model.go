package mazecollectchgrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

const (
	MazeCollectInit    = 1 // 迷宫挂机初始化
	MazeCollectTimeOut = 2 // 定时收集
	MazeCollectReceive = 3 // 领取
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeCollectChgRecordTableName = "maze_collect_chg_record"

// 迷宫挂机变化流水
type MazeCollectChgRecord struct {
	KafkaCommon
	UserId        uint64 `json:"user_id" gorm:"column:user_id"`               // 用户id
	OpType        int32  `json:"op_type" gorm:"column:op_type"`               // 1.初始化 2.定时收集 3.领取
	StartTime     int64  `json:"start_time" gorm:"column:start_time"`         // 开始时间
	LastTime      int64  `json:"last_time" gorm:"column:last_time"`           // 上次收集结算时间
	NewLastTime   int64  `json:"new_last_time" gorm:"column:new_last_time"`   // 本次收集结算时间
	AvailableTime int64  `json:"available_time" gorm:"column:available_time"` // 可领取时间
	EndTime       int64  `json:"end_time" gorm:"column:end_time"`             // 结束时间
	PeriodTime    int32  `json:"period_time" gorm:"column:period_time"`       // 产出周期
	CollectTimes  int64  `json:"collect_times" gorm:"column:collect_times"`   // 道具产出周期数
	BarrierId     int32  `json:"barrier_id" gorm:"column:barrier_id"`         // 关卡id
	TradeNo       uint64 `json:"trade_no" gorm:"column:trade_no"`             // 加物品流水号
	AddItems      string `json:"add_items" gorm:"column:add_items"`           // 收集的道具/领取的道具
	RemainItems   string `json:"remain_items" gorm:"column:remain_items"`     // 累计产出道具/领取后遗留的道具
	RetCode       int64  `json:"ret_code" gorm:"column:ret_code"`             // 0:成功  其他失败
}

func NewMazeCollectChgRecord(userID uint64, opType int32, startTime int64, lastTime int64, newLastTime int64, availableTime int64, endTime int64,
	periodTime int32, collectTimes int64, barrierId int32, tradeNo uint64, addItems string, remainItems string, retCode int64) *MazeCollectChgRecord {
	res := &MazeCollectChgRecord{
		UserId:        userID,
		OpType:        opType,
		StartTime:     startTime,
		LastTime:      lastTime,
		NewLastTime:   newLastTime,
		AvailableTime: availableTime,
		EndTime:       endTime,
		PeriodTime:    periodTime,
		CollectTimes:  collectTimes,
		BarrierId:     barrierId,
		TradeNo:       tradeNo,
		AddItems:      addItems,
		RemainItems:   remainItems,
		RetCode:       retCode,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeCollectChgRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
