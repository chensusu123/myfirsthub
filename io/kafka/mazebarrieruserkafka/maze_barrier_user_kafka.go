package mazebarrieruserkafka

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/io/dispatcher"
)

const (
	GameRetSucc  = 1
	GameRetDeath = 2
	GameRetSweep = 3
)

var d = dispatcher.NewDispatcher[*MazeBarrierUserGameRecord]()

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	Barrier    int32  `json:"barrier"`     // 关卡id
	GameRet    int32  `json:"game_ret"`    // 用户闯关结果 1-通关成功 2-死亡失败 3-扫荡
	Awards     string `json:"awards"`      // 本次获得的奖励
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazebarrieruserkafka", 1001105, gKafka)
}

func PushMazeBarrierUserRecord(agent fklog.FKLogI, record *MazeBarrierUserGameRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
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
