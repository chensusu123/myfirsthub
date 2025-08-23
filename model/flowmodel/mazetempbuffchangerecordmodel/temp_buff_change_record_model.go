package mazetempbuffchangerecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type AttrChgInfo struct {
	AttrId int32 `json:"attr_id"`
	OldVal int64 `json:"old_val"`
	CurVal int64 `json:"cur_val"`
}
type KafkaCommon = kafkacommonstruct.KafkaCommon

const MazeTempBuffChangeRecordTableName = "maze_temp_buff_change_record"

// 迷宫临时buff变化通知
type MazeTempBuffChangeMsg struct {
	KafkaCommon
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`     // 用户Id
	StageId     int32  `json:"stage_id" gorm:"column:stage_id"`   // 关卡id
	ChgAttrsStr string `json:"chg_attrs" gorm:"column:chg_attrs"` // 变化的属性 ChgAttrs的json格式，数据库存储字段
	ChgType     int32  `json:"chg_type" gorm:"column:chg_type"`   // 变化类型
	ChgDesc     string `json:"chg_desc" gorm:"column:chg_desc"`   // 原因描述
}

func NewMazeTempBuffChangeMsg(userID uint64, stageID int32, chgAttrsStr string, chgType int32, chgDesc string) *MazeTempBuffChangeMsg {
	res := &MazeTempBuffChangeMsg{
		UserId:      userID,
		StageId:     stageID,
		ChgAttrsStr: chgAttrsStr,
		ChgType:     chgType,
		ChgDesc:     chgDesc,
	}

	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeTempBuffChangeRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]

	return res
}
