package mazeuserenergylogmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeUserEnergyRecordTableName = "maze_user_energy_log"

const (
	TimerRecovery int32 = 100 // 时间恢复
	InitEnergy    int32 = 101 // 初始化体力
	ItemEnergy    int32 = 102 // 体力瓶
	EnterBarrier  int32 = 103 // 进入关卡
	SweepBarrier  int32 = 104 // 扫荡关卡
	GMAdd         int32 = 105 // GM添加
	Reset         int32 = 106 // 重置
)

// 迷宫体力变化流水
type MazeEnergyChgRecord struct {
	KafkaCommon
	UserId   uint64 `json:"user_id"`   // 用户id
	OldVal   int32  `json:"old_val" `  // 旧值
	NewVal   int32  `json:"new_val"`   // 新值
	LastTime int64  `json:"last_time"` // 上次恢复时间
	OpType   int32  `json:"op_type"`   // 操作类型
}

func NewMazeEnergyChgRecord(userID uint64, oldVal int32, newVal int32, lastTime int64, opType int32) *MazeEnergyChgRecord {
	res := &MazeEnergyChgRecord{
		UserId:   userID,
		OldVal:   oldVal,
		NewVal:   newVal,
		LastTime: lastTime,
		OpType:   opType,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserEnergyRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
