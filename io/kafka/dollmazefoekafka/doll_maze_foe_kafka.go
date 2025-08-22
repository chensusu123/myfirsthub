package dollmazefoekafka

import (
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

type KafkaCommon = kafkacommonstruct.KafkaCommon

// 用户打怪变化流水
type DollMazeFoeRecord struct {
	KafkaCommon
	UserId      uint64 `json:"user_id" gorm:"column:user_id"`           // 用户id
	Barrier     int32  `json:"barrier" gorm:"column:barrier"`           // 关卡id
	Area        int32  `json:"area" gorm:"column:area"`                 // 区域id
	Level       int32  `json:"level" gorm:"column:level"`               // 用户等级
	FoeList     string `json:"foe_list,omitempty" gorm:"-"`             // 上报的打败的怪物列表
	AwardList   string `json:"award_list" gorm:"column:award_list"`     // 打怪获得的奖励 1银子 2装备积分 3经验
	Equips      string `json:"equips" gorm:"column:equips"`             // 打怪获得的装备
	EquipPoints int32  `json:"equip_points" gorm:"column:equip_points"` // 本次打怪后的当前装备积分
	GroupID     uint32 `json:"group_id" gorm:"column:group_id"`         // 组id
	CreateTime  int64  `json:"create_time" gorm:"column:create_time"`   // 操作时间
	MasterId    int64  `json:"master_id" gorm:"column:master_id"`       // 怪物id
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
	record.FoeList = ""
	// record.GroupID = fkconfig.EnvVal.GroupID
	// cnt, err := json.Marshal(record)
	// if err != nil {
	// return err
	// }
	d.Push(agent, record)
	agent.InfoWF("PushDollMazeFoeRecord data", zap.Any("userId", record.UserId), zap.Any("record", record))
	// return gKafka.SendWithUserID(record.UserId, cnt)
	return nil
}
