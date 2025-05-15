package dollmazeshoprecord

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/dispatcher"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 人偶迷宫商店流水
type DollMazeShopRecord struct {
	UserId        uint64 `json:"user_id"`         //用户id
	AreaId        int32  `json:"chg_type"`        // 变化类型
	SeqId         int32  `json:"seq_id"`          // 序列id
	BackId        int32  `json:"back_id"`         //备用id
	SlotId        int32  `json:"slot_id"`         //位置id
	CurSeqIndex   int32  `json:"cur_index"`       //当前位置
	BackIndex     int32  `json:"back_index"`      //备用位置
	TotalCount    int32  `json:"total_count"`     //累计购买数量
	AddEquipGuids string `json:"add_equip_guids"` //添加装备guid列表
	OpType        int32  `json:"op_type"`         // 操作来源
	TradeNumber   uint64 `json:"trade_number"`    /// 交易流水号
	RetCode       int32  `json:"ret_code"`        // 0:成功  其他失败
	GroupID       uint32 `json:"group_id"`        // 组id
	CreateTime    int64  `json:"create_time"`     // 操作时间
}

// var mazeShopQueue = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("dollmazeshoprecord", 1001078, mazeShopQueue)
}

var d = dispatcher.NewDispatcher[*DollMazeShopRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *DollMazeShopRecord)) {
	d.Watch(fn)
}

func PushDollMazeShopRecord(agent fklog.FKLogI, data *DollMazeShopRecord) (err error) {
	data.CreateTime = time.Now().UnixNano() / 1000000
	data.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }
	d.Push(agent, data)
	agent.DebugWF("PushDollMazeShopRecord data", zap.Any("userId", data.UserId), zap.Any("detail", data))
	// err = mazeShopQueue.SendWithUserID(data.UserId, cnt)
	// if err != nil {
	// 	agent.ErrorWF("PushDollMazeShopRecord SendWithUserID err", zap.ByteString("cnt", cnt), zap.Error(err))
	// }
	return
}
