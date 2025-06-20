package mazeequipinstancerecord

import (
	"time"

	"maze_game_server/io/dispatcher"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	MazeAddEquip         int32 = 1 // 添加装备
	MazeAddInstanceEquip int32 = 2 // 添加实例化装备
)

// 装备实例化流水
type MazeGameEquipInstanceRecord struct {
	UserId         uint64 `json:"user_id"`          // 用户id
	OpType         int32  `json:"op_type"`          // 业务类型 挂机/锻造/购买
	TradeNum       uint64 `json:"trade_num"`        // 交易单号
	EquipId        int32  `json:"equip_id"`         // 装备配置id
	EquipGuid      int64  `json:"equip_guid"`       // 装备guid
	TotalScore     int32  `json:"total_score"`      // 累计分数
	SentencesLibId int32  `json:"sentences_lib_id"` // 装备词条库ID
	BaseAttrs      string `json:"base_attrs"`       // 基础属性列表
	ChgType        int32  `json:"chg_type"`         // 变化原因 1 添加 2 删除 3 更新 4自动出售 5开始洗练 6 洗练确认
	IsFail         int32  `json:"is_fail"`          // 操作是否失败 0-成功 1-失败
	GroupID        uint32 `json:"group_id"`         // 组id
	CreateTime     int64  `json:"create_time"`      // 操作时间
	StageFactor    int32  `json:"stage_factor"`     // 档位系数
	Conditions     string `json:"conditions"`       // 指定条件
	RuleId         int32  `json:"rule_id"`          // 特色规则id
	SuitId         int32  `json:"suit_id"`          // 套装id
	EquipSubType   int32  `json:"equip_sub_type"`   // 装备子类型
}

// var equipInstanceChgQueue = &fkafka.KafkaProducer{}

func init() {
	// // 1001090 topic-maze-equip-instance-log 迷宫游戏装备实例化流水
	// fkconfig.RegisterNameNode("dollequipinstancerecord", 1001090, equipInstanceChgQueue)
}

var d = dispatcher.NewDispatcher[*MazeGameEquipInstanceRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeGameEquipInstanceRecord)) {
	d.Watch(fn)
}

func PushMazeGameEquipInstanceRecord(agent fklog.FKLogI, data *MazeGameEquipInstanceRecord) (err error) {
	data.CreateTime = time.Now().UnixNano() / 1000000
	// data.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }
	d.Push(agent, data)
	agent.InfoWF("PushMazeGameEquipInstanceRecord data", zap.Any("userId", data.UserId), zap.Any("detail", data))
	// err = equipInstanceChgQueue.SendWithUserID(data.UserId, cnt)
	// if err != nil {
	// 	agent.ErrorWF("PushDollEquipInstanceRecord SendWithUserID err", zap.ByteString("cnt", cnt), zap.Error(err))
	// }
	return
}
