package mazetempbuffchgmsg

import (
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/21 16:04
 * @Description: 迷宫临时buff变化通知
 */

type AttrChgInfo struct {
	AttrId int32 `json:"attr_id"`
	OldVal int64 `json:"old_val"`
	CurVal int64 `json:"cur_val"`
}
type KafkaCommon = kafkacommonstruct.KafkaCommon

// 迷宫临时buff变化通知
type MazeTempBuffChangeMsg struct {
	KafkaCommon
	UserId      uint64         `json:"user_id" gorm:"column:user_id"`         // 用户Id
	GroupId     uint32         `json:"group_id" gorm:"column:group_id"`       // 分组Id
	StageId     int32          `json:"stage_id" gorm:"column:stage_id"`       // 关卡id
	ChgAttrs    []*AttrChgInfo `json:"chg_attr" gorm:"-"`                     // 变化的属性
	ChgAttrsStr string         `json:"chg_attrs" gorm:"column:chg_attrs"`     // 变化的属性 ChgAttrs的json格式，数据库存储字段
	ChgType     int32          `json:"chg_type" gorm:"column:chg_type"`       // 变化类型
	ChgDesc     string         `json:"chg_desc" gorm:"column:chg_desc"`       // 原因描述
	CreateTime  int64          `json:"create_time" gorm:"column:create_time"` // 时间戳 ms
}

var d = dispatcher.NewDispatcher[*MazeTempBuffChangeMsg]()

// var mazeTempBuffChangeKafka = &fkafka.KafkaProducer{}
// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// _ = fkconfig.RegisterNameNode("mazeTempBuffChangeKafka", 1001104, mazeTempBuffChangeKafka)
}

func PushTempBuffChangeMsg(logger fklog.FKLogI, msg *MazeTempBuffChangeMsg) error {
	// msg.GroupId = fkconfig.EnvVal.GroupID
	if msg.CreateTime == 0 {
		msg.CreateTime = time.Now().UnixNano() / 1000000
	}

	// cnt, err := json.Marshal(msg)
	// if err != nil {
	// 	logger.ErrorWF("PushTempBuffChangeMsg marshal failed", zap.Uint64("uid", msg.UserId), zap.Error(err))
	// 	return err
	// }

	// err = mazeTempBuffChangeKafka.SendWithUserID(msg.UserId, cnt)
	// if err != nil {
	// 	logger.ErrorWF("PushTempBuffChangeMsg SendWithUserID error", zap.Uint64("uid", msg.UserId),
	// 		zap.Any("msg", msg), zap.Error(err))
	// 	return err
	// }
	d.Push(logger, msg)
	logger.DebugWF("PushTempBuffChangeMsg end", zap.Any("pushData", msg))
	return nil
}

func Watch(fn func(logger fklog.FKLogI, msg *MazeTempBuffChangeMsg)) {
	d.Watch(fn)
}
