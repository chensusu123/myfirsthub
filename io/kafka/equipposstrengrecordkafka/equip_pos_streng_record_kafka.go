package equipposstrengrecordkafka

import (
	"context"
	"maze_game_server/common/function/itemutil"
	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazeequipposleveluprecordmodel"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/services/flowservice"
)

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 装备位强化流水
type EquipPosLevelUpRecord struct {
	KafkaCommon
	UserId       uint64 `json:"user_id" gorm:"column:user_id"`
	GroupId      uint32 `json:"group_id" gorm:"column:group_id"`
	OpTime       int64  `json:"create_time" gorm:"column:create_time"` // 毫秒时间戳
	PosId        int32  `json:"pos_id" gorm:"column:pos_id"`
	OldPosLv     int32  `json:"old_pos_lv" gorm:"column:old_pos_lv"`           // 强化前装备位等级
	NewPosLv     int32  `json:"new_pos_lv" gorm:"column:new_pos_lv"`           // 强化后装备位等级
	OldPosSuitId int32  `json:"old_pos_suit_id" gorm:"column:old_pos_suit_id"` // 强化前装备位套装id
	NewPosSuitId int32  `json:"new_pos_suit_id" gorm:"column:new_pos_suit_id"` // 强化后装备位套装id
	// OldPkLv      int32  `json:"old_pk_lv"`       // 强化前pk段位
	// NewPkLv      int32  `json:"new_pk_lv"`       // 强化后pk段位
	TradeNo   uint64 `json:"trade_no" gorm:"column:trade_no"`     // 扣物品流水号
	CostItems string `json:"cost_items" gorm:"column:cost_items"` // 扣物品
	Result    int32  `json:"result" gorm:"column:result"`         // 结果 0:成功 1:强化失败 2:存储武力值属性失败 3:存储非武力值属性失败
}

// var logCli = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("equipposstrengrecordkafka", 1001030, logCli)
}

var d = dispatcher.NewDispatcher[*EquipPosLevelUpRecord]()

func Watch(fn func(ctx context.Context, msg *EquipPosLevelUpRecord)) {
	d.Watch(fn)
}

// func pushRecord(ctx context.Context, data *EquipPosLevelUpRecord) error {
// msg, err := json.Marshal(data)
// if err != nil {
// 	return err
// }
// return logCli.SendWithUserID(data.UserId, msg)
// }

// 流水打点使用
func PushEquipPosStrengRecord(ctx context.Context, userId uint64, posId,
	oldPosLv, newPosLv, oldPosSuitId, newPosSuitId int32, tradeNo uint64, items []*MazeCommon.MazeItem, result, mask int32,
) error {
	flowData := mazeequipposleveluprecordmodel.NewEquipPosLevelUpRecord(userId, posId, oldPosLv, newPosLv, oldPosSuitId, newPosSuitId,
		tradeNo, itemutil.CommonItemsToString(items), result)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	return nil
}
