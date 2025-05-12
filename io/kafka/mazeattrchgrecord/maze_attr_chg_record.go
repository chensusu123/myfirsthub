/*
 * @Author: majian
 * @Date: 2024-07-11 12:03:32
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 21:14:31
 */
package mazeattrchgrecord

import (
	"time"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/io/dispatcher"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

// var kp = &fkafka.KafkaProducer{}
// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// // 1001086 topic-maze-game-attr-chg-record 迷宫游戏属性变化流水
	// fkconfig.RegisterNameNode("mazeattrchgrecord", 1001086, kp)
}

var d = dispatcher.NewDispatcher[*structsdef.MazeGameAttrChgRecord]()

func Watch(fn func(logger fklog.FKLogI, msg *structsdef.MazeGameAttrChgRecord)) {
	d.Watch(fn)
}

func SendMazeGameAttrChgRecord(logger fklog.FKLogI, record *structsdef.MazeGameAttrChgRecord) error {
	record.GroupId = fkconfig.EnvVal.GroupID
	if record.CreateTime == 0 {
		record.CreateTime = time.Now().UnixNano() / 1000000
	}

	// jbs, e := json.Marshal(record)
	// if e != nil {
	// 	logger.ErrorWF("SendMazeGameAttrChgRecord Marshal fail", zap.Error(e), zap.Any("record", record))
	// 	return e
	// }
	// e = kp.SendWithUserID(record.UserId, jbs)
	// if e != nil {
	// 	logger.ErrorWF("SendMazeGameAttrChgRecord SendWithUserID fail", zap.Error(e), zap.Any("record", record))
	// 	return e
	// }
	logger.InfoWF("SendMazeGameAttrChgRecord SendWithUserID succ", zap.Any("record", record))
	d.Push(logger, record)
	return nil
}
