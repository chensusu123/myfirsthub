package mazeuserlevelkafka

import (
	"context"
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazeuserlevelrecordmodel"
	"maze_game_server/services/flowservice"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 用户等级变化流水
type MazeUserLevelRecord struct {
	KafkaCommon
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`             // 用户id
	OldLevel    int32  `json:"old_level" gorm:"column:old_level"`         // 旧等级
	OldTotalExp int64  `json:"old_total_exp" gorm:"column:old_total_exp"` // 旧经验总值
	NewLevel    int32  `json:"new_level" gorm:"column:new_level"`         // 新等级
	NewTotalExp int32  `json:"new_total_exp" gorm:"column:new_total_exp"` // 新经验总值
	GroupID     uint32 `json:"group_id" gorm:"column:group_id"`           // 组id
	CreateTime  int64  `json:"create_time" gorm:"column:create_time"`     // 操作时间 毫秒
}

// var gKafka = &fkafka.KafkaProducer{}
var d = dispatcher.NewDispatcher[*MazeUserLevelRecord]()

func init() {
	// fkconfig.RegisterNameNode("mazeuserlevelkafka", 1001084, gKafka)
}

// 流水和通知均使用
func PushMazeLevelRecord(ctx context.Context, record *MazeUserLevelRecord) error {
	agent := fklog.ContextAppLogger(ctx)
	record.CreateTime = time.Now().UnixNano() / 1e6
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }

	flowData := mazeuserlevelrecordmodel.NewMazeUserLevelRecord(record.UserId, record.OldLevel, record.OldTotalExp, record.NewLevel, record.NewTotalExp)
	flowservice.GflowService.SendFlowData(ctx, flowData)

	agent.InfoWF("PushMazeLevelRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))

	// ctx := context.TODO()
	// err = collect.HandleMazeLevelMsg(ctx, agent, 0, nil, cnt)
	// if err != nil {
	// 	agent.ErrorWF("PushMazeLevelRecord HandleMazeLevelMsg", zap.Any("cnt", cnt), zap.Error(err))
	// 	return err
	// }
	// err = equip.HandleMazeLvChg(ctx, agent, 0, nil, cnt)
	// if err != nil {
	// 	agent.ErrorWF("PushMazeLevelRecord HandleMazeLvChg", zap.Any("cnt", cnt), zap.Error(err))
	// 	return err
	// }
	d.Push(ctx, record)
	return nil
}

func Watch(fn func(ctx context.Context, record *MazeUserLevelRecord)) {
	d.Watch(fn)
}
