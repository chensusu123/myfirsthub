package mazeequipassemblerecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

const (
	DollEquipAssembleOpDress    int32 = 1    // 穿戴装备
	DollEquipAssembleOpReplace  int32 = 2    // 更换装备
	DollEquipAssembleOpDown     int32 = 3    // 卸下装备
	DollEquipAssembleOpIdentify int32 = 4    // 鉴定装备
	DollEquipAssembleOpInit     int32 = 5    // 初始穿戴
	DollEquipAssembleOpBag      int32 = 1000 // 背包操作最终类型= DollEquipAssembleOpBag+ 背包ENUM_EQUIP_BAG_OP_TYPE
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeEquipAssembleRecordTableName = "maze_equip_assemble_record"

type MazeGameEquipAssembleRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`           // 用户Id
	EquipPos   int32  `json:"equip_pos" gorm:"column:equip_pos"`       // 装备位ID
	OpType     int32  `json:"op_type" gorm:"column:op_type"`           // 穿戴装备/更换装备/卸下装备
	NewEquipId int32  `json:"new_equip_id" gorm:"column:new_equip_id"` // 穿戴装备配置ID
	NewGuid    uint64 `json:"new_guid" gorm:"column:new_guid"`         // 穿戴装备guid
	OldEquipId int32  `json:"old_equip_id" gorm:"column:old_equip_id"` // 卸下装备配置ID
	OldGuid    uint64 `json:"old_guid" gorm:"column:old_guid"`         // 卸下装备guid
	OldFElem   string `json:"old_f_elem" gorm:"column:old_f_elem"`     // 变化前激活信息
	NewFElem   string `json:"new_f_elem" gorm:"column:new_f_elem"`     // 变化后激活信息
	RetCode    int32  `json:"ret_code" gorm:"column:ret_code"`         // 0:成功  其他失败
	CodeMask   int32  `json:"code_mask" gorm:"column:code_mask"`       // 业务掩码
	TransID    uint64 `json:"trans_id" gorm:"column:trans_id"`         // 事务Id
}

func NewMazeGameEquipAssembleRecord(userID uint64, equipPos int32, opType int32, newEquipID int32, newGuid uint64, oldEquipID int32, oldGuid uint64, oldFelem string, newFelem string,
	retCode int32, codeMask int32, transID uint64) *MazeGameEquipAssembleRecord {
	res := &MazeGameEquipAssembleRecord{
		UserId:     userID,
		EquipPos:   equipPos,
		OpType:     opType,
		NewEquipId: newEquipID,
		NewGuid:    newGuid,
		OldEquipId: oldEquipID,
		OldGuid:    oldGuid,
		OldFElem:   oldFelem,
		NewFElem:   newFelem,
		RetCode:    retCode,
		CodeMask:   codeMask,
		TransID:    transID,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeEquipAssembleRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
