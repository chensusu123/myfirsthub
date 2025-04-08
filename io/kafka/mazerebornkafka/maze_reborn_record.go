package mazerebornkafka

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkafka"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户复活流水
type MazeRebornRecord struct {
	UserId      uint64 `json:"user_id"`      // 用户id
	Barrier     int32  `json:"barrier"`      // 关卡id
	RebornCount int64  `json:"reborn_count"` // 复活次数 第n次复活
	RebornCost  string `json:"reborn_cost"`  // 复活的消耗
	GroupID     uint32 `json:"group_id"`     // 组id
	CreateTime  int64  `json:"create_time"`  // 操作时间
}

var gKafka = &fkafka.KafkaProducer{}

func init() {
	fkconfig.RegisterNameNode("mazerebornkafka", 1001107, gKafka)
}

func PushMazeRebornRecord(agent fklog.FKLogI, record *MazeRebornRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	cnt, err := json.Marshal(record)
	if err != nil {
		return err
	}
	agent.InfoWF("PushMazeRebornRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	return gKafka.SendWithUserID(record.UserId, cnt)
}
