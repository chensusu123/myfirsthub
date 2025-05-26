package mazeequipbagrecord

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/dispatcher"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	MazeAddEquip         int32 = 1 // 添加装备
	MazeDelEquip         int32 = 2 // 出售装备
	MazeInstanceEquip    int32 = 3 // 实例化装备
	MazeDelInstanceEquip int32 = 4 // 删除实例化装备
	MazeDressEquip       int32 = 5 // 装备穿戴
)

// 装备背包流水
type MazeGameEquipBagRecord struct {
	UserId        uint64 `json:"user_id"`         //用户id
	ChgType       int32  `json:"chg_type"`        //变化原因 1 添加 2 删除 3 更新 4 锁定 5 解锁 6 实例化装备 7 删除实例化装备
	TradeNum      uint64 `json:"trade_num"`       //交易单号
	AddEquipGuids string `json:"add_equip_guids"` //新增装备guid列表
	DelEquipGuids string `json:"del_equip_guids"` //删除装备guid列表
	OpType        int32  `json:"op_type"`         // 业务类型 挂机/锻造/购买
	IsFail        int32  `json:"is_fail"`         //操作是否失败 0-成功 1-失败
	GroupID       uint32 `json:"group_id"`        // 组id
	CreateTime    int64  `json:"create_time"`     // 操作时间
}

// var equipBagChgQueue = &fkafka.KafkaProducer{}

func init() {
	// //	1001088 topic-maze-equip-bag-chg-log 迷宫游戏装备背包添加流水
	// fkconfig.RegisterNameNode("dollequipbagrecord", 1001088, equipBagChgQueue)
}

var d = dispatcher.NewDispatcher[*MazeGameEquipBagRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeGameEquipBagRecord)) {
	d.Watch(fn)
}

func PushMazeGameEquipBagRecord(agent fklog.FKLogI, data *MazeGameEquipBagRecord) error {
	// cnt, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }
	d.Push(agent, data)
	agent.InfoWF("PushMazeGameEquipBagRecord data", zap.Any("userId", data.UserId), zap.Any("detail", data))
	// return equipBagChgQueue.SendWithUserID(data.UserId, cnt)
	return nil
}
