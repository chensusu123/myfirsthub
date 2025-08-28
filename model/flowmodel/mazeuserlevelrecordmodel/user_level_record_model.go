package mazeuserlevelrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeUserLevelRecordTableName = "maze_user_level_record"

// 用户等级变化流水
type MazeUserLevelRecord struct {
	KafkaCommon
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`             // 用户id
	OldLevel    int32  `json:"old_level" gorm:"column:old_level"`         // 旧等级
	OldTotalExp int64  `json:"old_total_exp" gorm:"column:old_total_exp"` // 旧经验总值
	NewLevel    int32  `json:"new_level" gorm:"column:new_level"`         // 新等级
	NewTotalExp int32  `json:"new_total_exp" gorm:"column:new_total_exp"` // 新经验总值
}

func NewMazeUserLevelRecord(userID uint64, oldLevel int32, oldTotalExp int64, newLevel int32, newTotalExp int32) *MazeUserLevelRecord {
	res := &MazeUserLevelRecord{
		UserId:      userID,
		OldLevel:    oldLevel,
		OldTotalExp: oldTotalExp,
		NewLevel:    newLevel,
		NewTotalExp: newTotalExp,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeUserLevelRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
