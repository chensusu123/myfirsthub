package equipposstrengrecordkafka

import (
	"encoding/json"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkafka"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/protodef/MazeCommon"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/function/itemutil"
)

// 装备位强化流水
type EquipPosLevelUpRecord struct {
	UserId       uint64 `json:"user_id"`
	GroupId      uint32 `json:"group_id"`
	OpTime       int64  `json:"op_time"` // 毫秒时间戳
	PosId        int32  `json:"pos_id"`
	OldPosLv     int32  `json:"old_pos_lv"`      // 强化前装备位等级
	NewPosLv     int32  `json:"new_pos_lv"`      // 强化后装备位等级
	OldPosSuitId int32  `json:"old_pos_suit_id"` // 强化前装备位套装id
	NewPosSuitId int32  `json:"new_pos_suit_id"` // 强化后装备位套装id
	// OldPkLv      int32  `json:"old_pk_lv"`       // 强化前pk段位
	// NewPkLv      int32  `json:"new_pk_lv"`       // 强化后pk段位
	TradeNo   uint64 `json:"trade_no"`   // 扣物品流水号
	CostItems string `json:"cost_items"` // 扣物品
	Result    int32  `json:"result"`     // 结果 0:成功 1:强化失败 2:存储武力值属性失败 3:存储非武力值属性失败
}

var logCli = &fkafka.KafkaProducer{}

func init() {
	fkconfig.RegisterNameNode("equipposstrengrecordkafka", 1001030, logCli)
}

func pushRecord(logger fklog.FKLogI, data *EquipPosLevelUpRecord) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return logCli.SendWithUserID(data.UserId, msg)
}

func PushEquipPosStrengRecord(logger fklog.FKLogI, userId uint64, posId,
	oldPosLv, newPosLv, oldPosSuitId, newPosSuitId int32, tradeNo uint64, items []*MazeCommon.MazeItem, result, mask int32) error {
	data := &EquipPosLevelUpRecord{
		UserId:       userId,
		GroupId:      fkconfig.EnvVal.GroupID,
		OpTime:       time.Now().UnixMilli(),
		PosId:        posId,
		OldPosLv:     oldPosLv,
		NewPosLv:     newPosLv,
		OldPosSuitId: oldPosSuitId,
		NewPosSuitId: newPosSuitId,
		// OldPkLv:      oldPkLv,
		// NewPkLv:      newPkLv,
		TradeNo:   tradeNo,
		CostItems: itemutil.CommonItemsToString(items),
		Result:    result,
	}
	return pushRecord(logger, data)
}
