package mazemoneykafka

import (
	"context"
	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazemoneyrecordmodel"
	"maze_game_server/services/flowservice"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 用户货币变化流水
type MazeMoneyRecord struct {
	KafkaCommon
	UserId        uint64 `json:"user_id" gorm:"column:user_id"`                 // 用户id
	OldMoneyId    int32  `json:"old_money_id" gorm:"column:old_money_id"`       // 旧货币id
	OldMoneyCount int64  `json:"old_money_count" gorm:"column:old_money_count"` // 旧货币数量
	NewMoneyId    int32  `json:"new_money_id" gorm:"column:new_money_id"`       // 新货币id
	NewMoneyCount int64  `json:"new_money_count" gorm:"column:new_money_count"` // 新货币数量
	TradeNo       uint64 `json:"trade_no" gorm:"column:trade_no"`               // 交易号
	ChgReason     int32  `json:"chg_reason" gorm:"column:chg_reason"`           // 变化原因
	GroupID       uint32 `json:"group_id" gorm:"column:group_id"`               // 组id
	CreateTime    int64  `json:"create_time" gorm:"column:create_time"`         // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazemoneykafka", 1001100, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeMoneyRecord]()

func Watch(fn func(ctx context.Context, msg *MazeMoneyRecord)) {
	d.Watch(fn)
}

// 流水打点使用
func PushMazeMoneyRecord(ctx context.Context, record *MazeMoneyRecord) error {
	flowData := mazemoneyrecordmodel.NewMazeMoneyRecord(record.UserId, record.OldMoneyId, record.OldMoneyCount, record.NewMoneyId, record.NewMoneyCount, record.TradeNo, record.ChgReason)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	// d.Push(agent, record)
	// agent.InfoWF("PushMazeMoneyRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	return nil
}
