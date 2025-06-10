package mazerebornkafka

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/dispatcher"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户复活流水
type MazeRebornRecord struct {
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`           // 用户id
	Barrier     int32  `json:"barrier" gorm:"column:barrier"`           // 关卡id
	RebornCount int64  `json:"reborn_count" gorm:"column:reborn_count"` // 复活次数 第n次复活
	RebornCost  string `json:"reborn_cost" gorm:"column:reborn_cost"`   // 复活的消耗
	GroupID     uint32 `json:"group_id" gorm:"column:group_id"`         // 组id
	CreateTime  int64  `json:"create_time" gorm:"column:create_time"`   // 操作时间
	ServerId    int32  `json:"server_id" gorm:"column:server_id"`
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazerebornkafka", 1001107, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeRebornRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeRebornRecord)) {
	d.Watch(fn)
}

func PushMazeRebornRecord(agent fklog.FKLogI, record *MazeRebornRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	d.Push(agent, record)
	agent.InfoWF("PushMazeRebornRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	return nil
}
