package mazebarrierareakafka

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkafka"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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

var gKafka = &fkafka.KafkaProducer{}

func init() {
	fkconfig.RegisterNameNode("mazebarrierareakafka", 1001099, gKafka)
}

func PushMazeBarrierAreaRecord(agent fklog.FKLogI, record *MazeBarrierAreaRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	cnt, err := json.Marshal(record)
	if err != nil {
		return err
	}
	agent.InfoWF("PushMazeBarrierAreaRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	return gKafka.SendWithUserID(record.UserId, cnt)
}
