package dollequipdismantlekafka

import (
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

const (
	// OpTypeFromDismantleList    int32 = 0 // 分解列表分解
	OpTypeFromBagEquipInfo      int32 = 1 // 背包中装备分解
	OpTypeFromTempBagEquipInfo  int32 = 2 // 临时背包中装备分解
	OpTypeFromTempBagTimeout    int32 = 3 // 临时背包过期分解
	OpTypeFromTempBagTimeoutFix int32 = 4 // 临时背包过期分解
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 装备分解流水
type MazeGameEquipDismantleRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`         // 用户id
	EquipGuids string `json:"equip_guids" gorm:"column:equip_guids"` // 装备guid列表
	TradeNum   uint64 `json:"trade_num" gorm:"column:trade_num"`     // 交易单号
	Award      string `json:"award" gorm:"column:award"`             // 分解获得的材料
	OpType     int32  `json:"op_type" gorm:"column:op_type"`         // 分解的操作来源
	IsFail     int32  `json:"is_fail" gorm:"column:is_fail"`         // 操作是否失败 0-成功 1-失败
	GroupID    uint32 `json:"group_id" gorm:"column:group_id"`       // 组id
	CreateTime int64  `json:"create_time" gorm:"column:create_time"` // 操作时间
	ServerId   int32  `json:"server_id" gorm:"column:server_id"`
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// // 1001089 topic-maze-equip-dismantle-record-log 迷宫游戏装备分解流水
	// fkconfig.RegisterNameNode("dollequipdismantlekafka", 1001089, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeGameEquipDismantleRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeGameEquipDismantleRecord)) {
	d.Watch(fn)
}

func PushDollEquipDismantleRecord(agent fklog.FKLogI, record *MazeGameEquipDismantleRecord) error {
	record.CreateTime = time.Now().Unix()
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// 	return err
	// }
	// return gKafka.SendWithUserID(record.UserId, cnt)
	d.Push(agent, record)
	agent.InfoWF("PushDollEquipDismantleRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	return nil
}
