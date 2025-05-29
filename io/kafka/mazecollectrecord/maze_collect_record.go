package mazecollectrecord

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/dispatcher"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	MazeCollectInit    = 1 // 迷宫挂机初始化
	MazeCollectTimeOut = 2 // 定时收集
	MazeCollectReceive = 3 // 领取
)

// 迷宫挂机变化流水
type MazeCollectChgRecord struct {
	UserId        uint64 `json:"user_id"`        // 用户id
	OpType        int32  `json:"op_type"`        // 1.初始化 2.定时收集 3.领取
	StartTime     int64  `json:"start_time"`     // 开始时间
	LastTime      int64  `json:"last_time"`      // 上次收集结算时间
	NewLastTime   int64  `json:"new_last_time"`  // 本次收集结算时间
	AvailableTime int64  `json:"available_time"` //可领取时间
	EndTime       int64  `json:"end_time"`       // 结束时间
	PeriodTime    int32  `json:"period_time"`    // 产出周期
	CollectTimes  int64  `json:"collect_times"`  //道具产出周期数
	BarrierId     int32  `json:"barrier_id"`     // 关卡id
	TradeNo       uint64 `json:"trade_no"`       // 加物品流水号
	AddItems      string `json:"add_items"`      // 收集的道具/领取的道具
	RemainItems   string `json:"remain_items"`   // 累计产出道具/领取后遗留的道具
	RetCode       int64  `json:"ret_code"`       //0:成功  其他失败
	GroupID       uint32 `json:"group_id"`       // 组id
	CreateTime    int64  `json:"create_time"`    // 操作时间
	ServerId      int32  `json:"server_id"`
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("mazecollectrecord", 1001106, gKafka)
}

var d = dispatcher.NewDispatcher[*MazeCollectChgRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *MazeCollectChgRecord)) {
	d.Watch(fn)
}

func PushMazeCollectChgRecord(agent fklog.FKLogI, record *MazeCollectChgRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	// data, err := json.Marshal(record)
	// if err != nil {
	// return err
	// }
	d.Push(agent, record)
	agent.InfoWF("SendMazeCollectChgRecord data", zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, data)
	return nil
}
