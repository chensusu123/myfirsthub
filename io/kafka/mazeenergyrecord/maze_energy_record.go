/*
 * @Author: majian
 * @Date: 2025-03-27 11:20:40
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-27 16:06:37
 */
package mazeenergyrecord

import (
	"time"

	"maze_game_server/io/dispatcher"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	TimerRecovery int32 = 100 // 时间恢复
	InitEnergy    int32 = 101 // 初始化体力
)

// 迷宫体力变化流水
type MazeEnergyChgRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	OldVal     int32  `json:"old_val"`     // 旧值
	ChgVal     int32  `json:"chg_val"`     // 变化值
	NewVal     int32  `json:"new_val"`     // 新值
	LastTime   int64  `json:"last_time"`   // 上次恢复时间
	OpType     int32  `json:"op_type"`     // 操作类型
	TradeNo    uint64 `json:"trade_no"`    // 交易号
	CostItem   string `json:"cost_item"`   // 花费
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// // 3876 db_maze_energy_chg_log 迷宫体力变化流水
	// // 1001106 topic-maze-energy-chg-log 迷宫体力变化流水
	// fkconfig.RegisterNameNode("mazeenergyrecord", 1001106, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeEnergyChgRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeEnergyChgRecord)) {
	d.Watch(fn)
}

func SendMazeEnergyChgRecord(agent fklog.FKLogI, record *MazeEnergyChgRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	// record.GroupID = fkconfig.EnvVal.GroupID
	// data, err := json.Marshal(record)
	// if err != nil {
	// return err
	// }
	d.Push(agent, record)
	agent.InfoWF("SendMazeEnergyChgRecord data", zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, data)
	return nil
}
