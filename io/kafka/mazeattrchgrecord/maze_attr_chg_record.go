/*
 * @Author: majian
 * @Date: 2024-07-11 12:03:32
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 21:14:31
 */
package mazeattrchgrecord

import (
	"context"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/dispatcher"
	"maze_game_server/model/flowmodel/mazeattrchgrecordmodel"
	"maze_game_server/services/flowservice"
)

// var kp = &fkafka.KafkaProducer{}
// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// // 1001086 topic-maze-game-attr-chg-record 迷宫游戏属性变化流水
	// fkconfig.RegisterNameNode("mazeattrchgrecord", 1001086, kp)
}

var d = dispatcher.NewDispatcher[*structsdef.MazeGameAttrChgRecord]()

func Watch(fn func(ctx context.Context, msg *structsdef.MazeGameAttrChgRecord)) {
	d.Watch(fn)
}

// 流水打点使用
func SendMazeGameAttrChgRecord(ctx context.Context, record *structsdef.MazeGameAttrChgRecord) error {
	flowData := mazeattrchgrecordmodel.NewMazeGameAttrChgRecordFlow(record.UserId, record.AttrId, record.AttrType, record.NewVal, record.OldVal,
		record.ChgType, record.ChgSubType, record.ChgDesc, record.Extra)
	flowservice.GflowService.SendFlowData(ctx, flowData)
	return nil
}
