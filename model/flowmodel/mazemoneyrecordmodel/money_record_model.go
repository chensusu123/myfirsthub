package mazemoneyrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeMoneyRecordTableName = "maze_money_record"

// 用户货币变化流水
type MazeMoneyRecord struct {
	KafkaCommon
	UserId        uint64 `json:"user_id" gorm:"column:user_id"`                 // 用户id
	OldMoneyId    int32  `json:"old_money_id" gorm:"column:old_money_id"`       // 旧货币id
	OldMoneyCount int64  `json:"old_money_count" gorm:"column:old_money_count"` // 旧货币数量
	NewMoneyId    int32  `json:"new_money_id" gorm:"column:new_money_id"`       // 新货币id
	NewMoneyCount int64  `json:"new_money_count" gorm:"column:new_money_count"` // 新货币数量
	TradeNo       string `json:"trade_no" gorm:"column:trade_no"`               // 交易号
	ChgReason     int32  `json:"chg_reason" gorm:"column:chg_reason"`           // 变化原因
}

func NewMazeMoneyRecord(userID uint64, oldMoneyID int32, oldMoneyCount int64, newMoneyID int32, newMoneyCount int64, tradeNo uint64, chgReason int32) *MazeMoneyRecord {
	res := &MazeMoneyRecord{
		UserId:        userID,
		OldMoneyId:    oldMoneyID,
		OldMoneyCount: oldMoneyCount,
		NewMoneyId:    newMoneyID,
		NewMoneyCount: newMoneyCount,
		TradeNo:       strconv.FormatUint(tradeNo, 10),
		ChgReason:     chgReason,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeMoneyRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
