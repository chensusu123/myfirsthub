package dollmazefoekafka

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
	"maze_game_server/io/dispatcher"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 用户打怪变化流水
type DollMazeFoeRecord struct {
	UserId      uint64 `json:"user_id"`      // 用户id
	Barrier     int32  `json:"barrier"`      // 关卡id
	Area        int32  `json:"area"`         // 区域id
	Level       int32  `json:"level"`        // 用户等级
	FoeList     string `json:"foe_list"`     // 上报的打败的怪物列表
	AwardList   string `json:"award_list"`   // 打怪获得的奖励 1银子 2装备积分 3经验
	Equips      string `json:"equips"`       // 打怪获得的装备
	EquipPoints int32  `json:"equip_points"` // 本次打怪后的当前装备积分
	GroupID     uint32 `json:"group_id"`     // 组id
	CreateTime  int64  `json:"create_time"`  // 操作时间
}

// var gKafka = &fkafka.KafkaProducer{}

func init() {
	// fkconfig.RegisterNameNode("dollmazefoekafka", 1001081, gKafka)
}

var d = dispatcher.NewDispatcher[*DollMazeFoeRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *DollMazeFoeRecord)) {
	d.Watch(fn)
}

func PushDollMazeFoeRecord(agent fklog.FKLogI, record *DollMazeFoeRecord) error {
	record.CreateTime = time.Now().UnixNano() / 1e6
	record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// return err
	// }
	d.Push(agent, record)
	agent.InfoWF("PushDollMazeFoeRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	return nil
}
