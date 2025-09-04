package mazesweeprecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeSweepRecordTableName = "maze_sweep_record"

// 用户迷宫扫荡纪录
type MazeBarrierSweepRecord struct {
	KafkaCommon
	UserId  uint64 `json:"user_id" gorm:"column:user_id"` // 用户id
	Barrier int32  `json:"barrier" gorm:"column:barrier"` // 关卡id
	Awards  string `json:"awards" gorm:"column:awards"`   // 本次获得的奖励
}

func NewMazeBarrierSweepRecord(userID uint64, barrier int32, awards string) *MazeBarrierSweepRecord {
	res := &MazeBarrierSweepRecord{
		UserId:  userID,
		Barrier: barrier,
		Awards:  awards,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeSweepRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
