package mazeequipbagrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strconv"
	"strings"
)

const (
	MazeAddEquip         int32 = 1 // 添加装备
	MazeDelEquip         int32 = 2 // 出售装备
	MazeInstanceEquip    int32 = 3 // 实例化装备
	MazeDelInstanceEquip int32 = 4 // 删除实例化装备
	MazeDressEquip       int32 = 5 // 装备穿戴
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeEquipBagRecordTableName = "maze_equip_bag_record"

// 装备背包流水
type MazeGameEquipBagRecord struct {
	KafkaCommon
	UserId        uint64 `json:"user_id" gorm:"column:user_id"`                 //用户id
	ChgType       int32  `json:"chg_type" gorm:"column:chg_type"`               //变化原因 1 添加 2 删除 3 更新 4 锁定 5 解锁 6 实例化装备 7 删除实例化装备
	TradeNum      string `json:"trade_num" gorm:"column:trade_num"`             //交易单号
	AddEquipGuids string `json:"add_equip_guids" gorm:"column:add_equip_guids"` //新增装备guid列表
	DelEquipGuids string `json:"del_equip_guids" gorm:"column:del_equip_guids"` //删除装备guid列表
	OpType        int32  `json:"op_type" gorm:"column:op_type"`                 // 业务类型 挂机/锻造/购买
	IsFail        int32  `json:"is_fail" gorm:"column:is_fail"`                 //操作是否失败 0-成功 1-失败
}

func NewMazeGameEquipBagRecord(userID uint64, chgType int32, tradeNum uint64, addEquipGuids, delEquipGuids string, opType int32, isFail int32) *MazeGameEquipBagRecord {
	res := &MazeGameEquipBagRecord{
		UserId:        userID,
		ChgType:       chgType,
		TradeNum:      strconv.FormatUint(tradeNum, 10),
		AddEquipGuids: addEquipGuids,
		DelEquipGuids: delEquipGuids,
		OpType:        opType,
		IsFail:        isFail,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipBagRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
