package mazebarrierareakafka

import (
	"context"
	"time"

	"maze_game_server/io/dispatcher"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户关卡区域变化流水
type MazeBarrierAreaRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	OldBarrier int32  `json:"old_barrier"` // 旧关卡id
	OldArea    int32  `json:"old_area"`    // 旧区域id
	NewBarrier int32  `json:"new_barrier"` // 新关卡id
	NewArea    int32  `json:"new_area"`    // 新区域id
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazebarrierareakafka", 1001099, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeBarrierAreaRecord]()

func Watch(fn func(ctx context.Context, msg *MazeBarrierAreaRecord)) {
	d.Watch(fn)
}

func PushMazeBarrierAreaRecord(ctx context.Context, record *MazeBarrierAreaRecord) error {
	agent := fklog.ContextAppLogger(ctx)
	record.CreateTime = time.Now().UnixNano() / 1e6
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	agent.InfoWF("PushMazeBarrierAreaRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	d.Push(ctx, record)
	return nil
}
