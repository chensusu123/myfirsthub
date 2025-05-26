/*
 * @Author: majian
 * @Date: 2025-03-11 20:09:57
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-15 14:23:36
 */
package mazeattrmsg

import (
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"maze_game_server/common/structsdef"
	"maze_game_server/io/dispatcher"
)

// var kp = &fkafka.KafkaProducer{}
var d = dispatcher.NewDispatcher[*structsdef.DollAttrChgNotify]()

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// 1001083 topic-maze-attr-chg-notify-msg 迷宫属性变化通知消息
	// fkconfig.RegisterNameNode("mazeattrmsg", 1001083, kp)
}

func SendMazeAttrChgNotify(logger fklog.FKLogI, msg *structsdef.DollAttrChgNotify) error {
	msg.GroupId = fkconfig.EnvVal.GroupID
	if msg.CreateTime == 0 {
		msg.CreateTime = time.Now().UnixNano() / 1000000
	}

	// jbs, e := json.Marshal(msg)
	// if e != nil {
	// 	logger.ErrorWF("SendMazeAttrChgNotify Marshal fail", zap.Error(e), zap.Any("msg", msg))
	// 	return e
	// }
	// e = kp.SendWithUserID(msg.UserId, jbs)
	// if e != nil {
	// 	logger.ErrorWF("SendMazeAttrChgNotify SendWithUserID fail", zap.Error(e), zap.Any("msg", msg))
	// 	return e
	// }
	// logger.InfoWF("SendMazeAttrChgNotify SendWithUserID succ", zap.Any("msg", msg))
	d.Push(logger, msg)
	return nil
}

func Watch(fn func(logger fklog.FKLogI, msg *structsdef.DollAttrChgNotify)) {
	d.Watch(fn)
}
