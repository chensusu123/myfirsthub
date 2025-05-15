package mazemoneykafka

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/dispatcher"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户货币变化流水
type MazeMoneyRecord struct {
	UserId        uint64 `json:"user_id"`         // 用户id
	OldMoneyId    int32  `json:"old_money_id"`    // 旧货币id
	OldMoneyCount int64  `json:"old_money_count"` // 旧货币数量
	NewMoneyId    int32  `json:"new_money_id"`    // 新货币id
	NewMoneyCount int64  `json:"new_money_count"` // 新货币数量
	TradeNo       int64  `json:"trade_no"`        // 交易号
	ChgReason     int32  `json:"chg_reason"`      // 变化原因
	GroupID       uint32 `json:"group_id"`        // 组id
	CreateTime    int64  `json:"create_time"`     // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazemoneykafka", 1001100, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeMoneyRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeMoneyRecord)) {
	d.Watch(fn)
}

func PushMazeMoneyRecord(agent fklog.FKLogI, record *MazeMoneyRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	d.Push(agent, record)
	agent.InfoWF("PushMazeMoneyRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	return nil
}
