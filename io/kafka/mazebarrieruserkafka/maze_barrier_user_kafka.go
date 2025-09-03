package mazebarrieruserkafka

import (
	"context"
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazebarrieruserrecordmodel"
	"maze_game_server/model/flowmodel/mazesweeprecordmodel"
	"maze_game_server/services/flowservice"
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
	UserId         uint64 `json:"user_id" gorm:"column:user_id"`            // 用户id
	Barrier        int32  `json:"barrier" gorm:"column:barrier"`            // 关卡id
	GameRet        int32  `json:"game_ret" gorm:"column:game_ret"`          // 用户闯关结果 1-通关成功 2-死亡失败 3-扫荡
	Awards         string `json:"awards" gorm:"column:awards"`              // 本次获得的奖励
	GroupID        uint32 `json:"group_id" gorm:"column:group_id"`          // 组id
	CreateTime     int64  `json:"create_time" gorm:"column:create_time"`    // 操作时间
	KillMonsterNum int64  `json:"kill_monster_num" gorm:"kill_monster_num"` // 杀怪数量
	DeathReason    uint32 `json:"death_reason" gorm:"death_reason"`         // 死亡原因
	UserData       string
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazebarrieruserkafka", 1001105, gKafka)
}

// 流水和通知均使用
func PushMazeBarrierUserRecord(ctx context.Context, record *MazeBarrierUserGameRecord) error {
	if record.GameRet == 3 {
		flowData := mazesweeprecordmodel.NewMazeBarrierSweepRecord(record.UserId, record.Barrier, record.Awards)
		flowservice.GflowService.SendFlowData(ctx, flowData)
	}
	record.CreateTime = time.Now().UnixNano() / 1e6
	flowData := mazebarrieruserrecordmodel.NewMazeBarrierUserGameRecord(record.UserId, record.Barrier, record.GameRet, record.Awards,
		record.KillMonsterNum, record.DeathReason, record.UserData)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	// agent.InfoWF("PushMazeBarrierUserRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	d.Push(ctx, record)
	return nil
}

func Watch(fn func(ctx context.Context, msg *MazeBarrierUserGameRecord)) {
	d.Watch(fn)
}
