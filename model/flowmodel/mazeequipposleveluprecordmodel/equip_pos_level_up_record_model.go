package mazeequipposleveluprecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 装备位强化流水
type EquipPosLevelUpRecord struct {
	KafkaCommon
	UserId       uint64 `json:"user_id" gorm:"column:user_id"`
	PosId        int32  `json:"pos_id" gorm:"column:pos_id"`
	OldPosLv     int32  `json:"old_pos_lv" gorm:"column:old_pos_lv"`           // 强化前装备位等级
	NewPosLv     int32  `json:"new_pos_lv" gorm:"column:new_pos_lv"`           // 强化后装备位等级
	OldPosSuitId int32  `json:"old_pos_suit_id" gorm:"column:old_pos_suit_id"` // 强化前装备位套装id
	NewPosSuitId int32  `json:"new_pos_suit_id" gorm:"column:new_pos_suit_id"` // 强化后装备位套装id
	TradeNo      uint64 `json:"trade_no" gorm:"column:trade_no"`               // 扣物品流水号
	CostItems    string `json:"cost_items" gorm:"column:cost_items"`           // 扣物品
	Result       int32  `json:"result" gorm:"column:result"`                   // 结果 0:成功 1:强化失败 2:存储武力值属性失败 3:存储非武力值属性失败
}

const MazeEquipPosLevelUpRecordTableName = "maze_equip_pos_level_up_record"

func NewEquipPosLevelUpRecord(userID uint64, posID int32, oldPosLv int32, newPosLv int32, oldPosSuitID int32,
	newPosSuitID int32, tradeNo uint64, costItems string, result int32) *EquipPosLevelUpRecord {
	res := &EquipPosLevelUpRecord{
		UserId:       userID,
		PosId:        posID,
		OldPosLv:     oldPosLv,
		NewPosLv:     newPosLv,
		OldPosSuitId: oldPosSuitID,
		NewPosSuitId: newPosSuitID,
		TradeNo:      tradeNo,
		CostItems:    costItems,
		Result:       result,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipPosLevelUpRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
