package dollequipassmeblekakfa

import (
	"context"
	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazeequipassemblerecordmodel"
	"maze_game_server/services/flowservice"
)

var d = dispatcher.NewDispatcher[*MazeGameEquipAssembleRecord]()

func Watch(fn func(ctx context.Context, msg *MazeGameEquipAssembleRecord)) {
	d.Watch(fn)
}

// var kp = &fkafka.KafkaProducer{}
// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// 1001087 topic-maze-game-equip-assemble-chg-record 迷宫游戏装备装配流水
	// 3548  db_doll_equip_assemble_chg_log 人偶装备装配流水
	// fkconfig.RegisterNameNode("dollequipassmeblekakfa", 1001087, kp)
}

const (
	DollEquipAssembleOpDress    int32 = 1    // 穿戴装备
	DollEquipAssembleOpReplace  int32 = 2    // 更换装备
	DollEquipAssembleOpDown     int32 = 3    // 卸下装备
	DollEquipAssembleOpIdentify int32 = 4    // 鉴定装备
	DollEquipAssembleOpInit     int32 = 5    // 初始穿戴
	DollEquipAssembleOpBag      int32 = 1000 // 背包操作最终类型= DollEquipAssembleOpBag+ 背包ENUM_EQUIP_BAG_OP_TYPE
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

type MazeGameEquipAssembleRecord struct {
	KafkaCommon
	UserId     uint64 `json:"user_id" gorm:"column:user_id"`           // 用户Id
	GroupId    uint32 `json:"group_id" gorm:"column:group_id"`         // 分组ID  当时服务分片所属分组
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
	OpTime     int64  `json:"create_time" gorm:"column:create_time"`   // 流水时间戳
}

// 流水打点使用
func SendMazeGameEquipAssembleRecord(ctx context.Context, record *MazeGameEquipAssembleRecord) error {
	flowData := mazeequipassemblerecordmodel.NewMazeGameEquipAssembleRecord(record.UserId, record.EquipPos, record.OpType, record.NewEquipId, record.NewGuid, record.OldEquipId, record.OldGuid,
		record.OldFElem, record.NewFElem, record.RetCode, record.CodeMask, record.TransID)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	return nil
}
