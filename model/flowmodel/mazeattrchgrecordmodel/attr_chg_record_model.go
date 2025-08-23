package mazeattrchgrecordmodel

import (
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/io/mysql"
	"strings"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

type MazeGameAttrChgRecordFlow struct {
	KafkaCommon
	UserId     uint64 `json:"user_id"`      //用户Id
	AttrId     int32  `json:"attr_id"`      //属性ID
	AttrType   int32  `json:"attr_type"`    //属性类型
	NewVal     int64  `json:"new_val"`      //新值
	OldVal     int64  `json:"old_val"`      //旧值
	ChgType    int32  `json:"chg_type"`     //变化类型
	ChgSubType int32  `json:"chg_sub_type"` //变化子类型
	ChgDesc    string `json:"chg_desc"`     //原因描述
	Extra      string `json:"extra"`        // 扩展信息
}

const MazeAttrChgRecordTableName = "maze_attr_chg_record"

func NewMazeGameAttrChgRecordFlow(useID uint64, attrID int32, attrType int32, newVal int64, oldVal int64,
	chgType int32, chgSubType int32, chgDesc string, extra string) *MazeGameAttrChgRecordFlow {
	res := &MazeGameAttrChgRecordFlow{
		UserId:     useID,
		AttrId:     attrID,
		AttrType:   attrType,
		NewVal:     newVal,
		OldVal:     oldVal,
		ChgType:    chgType,
		ChgSubType: chgSubType,
		ChgDesc:    chgDesc,
		Extra:      extra,
	}
	nowDbTable := strings.Split(mysql.GetFullyQualifiedTableName(MazeAttrChgRecordTableName), ".")

	res.DataBase = nowDbTable[0]
	res.Table = nowDbTable[1]
	return res
}
