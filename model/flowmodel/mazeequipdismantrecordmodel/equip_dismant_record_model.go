package mazeequipdismantrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
)

const (
	// OpTypeFromDismantleList    int32 = 0 // 分解列表分解
	OpTypeFromBagEquipInfo      int32 = 1 // 背包中装备分解
	OpTypeFromTempBagEquipInfo  int32 = 2 // 临时背包中装备分解
	OpTypeFromTempBagTimeout    int32 = 3 // 临时背包过期分解
	OpTypeFromTempBagTimeoutFix int32 = 4 // 临时背包过期分解
)

const MazeEquipDismantRecordTableName = "maze_equip_dismant_record"

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 装备分解流水
type MazeGameEquipDismantleRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`         // 用户id
	EquipGuids string `json:"equip_guids" gorm:"column:equip_guids"` // 装备guid列表
	TradeNum   string `json:"trade_num" gorm:"column:trade_num"`     // 交易单号
	Award      string `json:"award" gorm:"column:award"`             // 分解获得的材料
	OpType     int32  `json:"op_type" gorm:"column:op_type"`         // 分解的操作来源
	IsFail     int32  `json:"is_fail" gorm:"column:is_fail"`         // 操作是否失败 0-成功 1-失败
}

func NewMazeGameEquipDismantleRecord(useID uint64, equipGuids string, tradeNum uint64, award string, opType int32, isFail int32) *MazeGameEquipDismantleRecord {
	res := &MazeGameEquipDismantleRecord{
		UserId:     useID,
		EquipGuids: equipGuids,
		TradeNum:   strconv.FormatUint(tradeNum, 10),
		Award:      award,
		OpType:     opType,
		IsFail:     isFail,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipDismantRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
