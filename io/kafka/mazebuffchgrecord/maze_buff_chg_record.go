/*
 * @Author: majian
 * @Date: 2025-03-15 10:46:52
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 10:48:23
 */
package mazebuffchgrecord

import (
	"context"
	"time"

	"maze_game_server/common/structsdef"
	"maze_game_server/io/dispatcher"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var kp = &fkafka.KafkaProducer{}
// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// // 1001091 topic-maze-buff-attr-chg-log 迷宫游戏buff属性变化流水
	// fkconfig.RegisterNameNode("mazebuffchgrecord", 1001091, kp)
}

var d = dispatcher.NewDispatcher[*structsdef.MazeGameBuffAttrChgRecord]()

func Watch(fn func(ctx context.Context, msg *structsdef.MazeGameBuffAttrChgRecord)) {
	d.Watch(fn)
}

func SendMazeBuffAttrRecord(ctx context.Context, record *structsdef.MazeGameBuffAttrChgRecord) error {
	logger := fklog.ContextAppLogger(ctx)
	if record.CreateTime == 0 {
		record.CreateTime = time.Now().UnixNano() / 1000000
	}

	d.Push(ctx, record)
	logger.InfoWF("SendMazeBuffAttrRecord SendWithUserID succ", zap.Any("record", record))
	return nil
}
