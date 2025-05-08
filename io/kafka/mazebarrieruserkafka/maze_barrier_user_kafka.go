package mazebarrieruserkafka

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkafka"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/buff"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/collect"
	"context"
)

const (
	GameRetSucc  = 1
	GameRetDeath = 2
	GameRetSweep = 3
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户迷宫闯关纪录
type MazeBarrierUserGameRecord struct {
	UserId     uint64 `json:"user_id"`     // 用户id
	Barrier    int32  `json:"barrier"`     // 关卡id
	GameRet    int32  `json:"game_ret"`    // 用户闯关结果 1-通关成功 2-死亡失败 3-扫荡
	Awards     string `json:"awards"`      // 本次获得的奖励
	GroupID    uint32 `json:"group_id"`    // 组id
	CreateTime int64  `json:"create_time"` // 操作时间
}

var gKafka = &fkafka.KafkaProducer{}

func init() {
	fkconfig.RegisterNameNode("mazebarrieruserkafka", 1001105, gKafka)
}

func PushMazeBarrierUserRecord(agent fklog.FKLogI, record *MazeBarrierUserGameRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	cnt, err := json.Marshal(record)
	if err != nil {
		return err
	}
	agent.InfoWF("PushMazeBarrierUserRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	ctx := context.TODO()
	err = buff.MazeBarrierNotifyProcess(ctx, agent, 0, nil, cnt)
	if err != nil {
		return err
	}
	err = collect.HandleMazeBarrierMsg(ctx, agent, 0, nil, cnt)
	if err != nil {
		return err
	}
	return nil
	// return gKafka.SendWithUserID(record.UserId, cnt)
}
