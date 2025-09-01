package mazetempbuffchgmsg

import (
	"context"
	"fmt"
	"time"

	"maze_game_server/io/dispatcher"
	"maze_game_server/io/kafka/kafkacommonstruct"
	"maze_game_server/model/flowmodel/mazetempbuffchangerecordmodel"
	"maze_game_server/services/flowservice"

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
	ChgAttrs    []*AttrChgInfo `json:"chg_attr,omitempty" gorm:"-"`           // 变化的属性
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

// 流水和通知均使用
func PushTempBuffChangeMsg(ctx context.Context, msg *MazeTempBuffChangeMsg) error {
	logger := fklog.ContextAppLogger(ctx)
	if msg.CreateTime == 0 {
		msg.CreateTime = time.Now().UnixNano() / 1000000
	}

	flowData := mazetempbuffchangerecordmodel.NewMazeTempBuffChangeMsg(msg.UserId, msg.StageId, Buff2String(msg.ChgAttrs), msg.ChgType, msg.ChgDesc)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	// cnt, err := json.Marshal(msg)
	// if err != nil {
	// 	logger.CtxError(ctx,"PushTempBuffChangeMsg marshal failed", zap.Uint64("uid", msg.UserId), zap.Error(err))
	// 	return err
	// }

	// err = mazeTempBuffChangeKafka.SendWithUserID(msg.UserId, cnt)
	// if err != nil {
	// 	logger.CtxError(ctx,"PushTempBuffChangeMsg SendWithUserID error", zap.Uint64("uid", msg.UserId),
	// 		zap.Any("msg", msg), zap.Error(err))
	// 	return err
	// }
	d.Push(ctx, msg)
	logger.CtxDebug(ctx, "PushTempBuffChangeMsg end", zap.Any("pushData", msg))
	return nil
}

func Watch(fn func(ctx context.Context, msg *MazeTempBuffChangeMsg)) {
	d.Watch(fn)
}

func Buff2String(buff []*AttrChgInfo) string {
	res := ""
	for _, val := range buff {
		if res == "" {
			res = fmt.Sprintf("%d:%d:%d", val.AttrId, val.OldVal, val.CurVal)
			continue
		}
		res = fmt.Sprintf("%s_%d:%d:%d", res, val.AttrId, val.OldVal, val.CurVal)
	}
	return res

}
