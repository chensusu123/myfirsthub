package mazebarrieruserkafka

import (
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
)

const (
	GameRetSucc  = 1
	GameRetDeath = 2
	GameRetSweep = 3
)

var d = dispatcher.NewDispatcher[*MazeBarrierUserGameRecord]()

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`         // 用户id
	Barrier    int32  `json:"barrier" gorm:"column:barrier"`         // 关卡id
	GameRet    int32  `json:"game_ret" gorm:"column:game_ret"`       // 用户闯关结果 1-通关成功 2-死亡失败 3-扫荡
	Awards     string `json:"awards" gorm:"column:awards"`           // 本次获得的奖励
	GroupID    uint32 `json:"group_id" gorm:"column:group_id"`       // 组id
	CreateTime int64  `json:"create_time" gorm:"column:create_time"` // 操作时间
	ServerId   int32  `json:"server_id" gorm:"column:server_id"`
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazebarrieruserkafka", 1001105, gKafka)
}

func PushMazeBarrierUserRecord(agent fklog.FKLogI, record *MazeBarrierUserGameRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	// agent.InfoWF("PushMazeBarrierUserRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	d.Push(agent, record)
	return nil
}

func Watch(fn func(logger fklog.FKLogI, msg *MazeBarrierUserGameRecord)) {
	d.Watch(fn)
}
