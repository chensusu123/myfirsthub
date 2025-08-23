package mazefoerecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeFoeRecordTableName = "maze_foe_record"

// 用户打怪变化流水
type MazeFoeRecord struct {
	KafkaCommon
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`           // 用户id
	Barrier     int32  `json:"barrier" gorm:"column:barrier"`           // 关卡id
	Area        int32  `json:"area" gorm:"column:area"`                 // 区域id
	Level       int32  `json:"level" gorm:"column:level"`               // 用户等级
	AwardList   string `json:"award_list" gorm:"column:award_list"`     // 打怪获得的奖励 1银子 2装备积分 3经验
	Equips      string `json:"equips" gorm:"column:equips"`             // 打怪获得的装备
	EquipPoints int32  `json:"equip_points" gorm:"column:equip_points"` // 本次打怪后的当前装备积分
	MasterId    int64  `json:"master_id" gorm:"column:master_id"`       // 怪物id
}

func NewDollMazeFoeRecord(userID uint64, barrier int32, area int32, level int32, awardList, equips string, equipPoints int32, masterID int64) *MazeFoeRecord {
	res := &MazeFoeRecord{
		UserId:      userID,
		Barrier:     barrier,
		Area:        area,
		Level:       level,
		AwardList:   awardList,
		Equips:      equips,
		EquipPoints: equipPoints,
		MasterId:    masterID,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeFoeRecordTableName), ".")
	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
