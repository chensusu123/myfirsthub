/*
 * @Author: majian
 * @Date: 2024-07-10 15:23:44
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 22:19:14
 */
package mazeattrcalcnotifyqueue

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"
	"go.uber.org/zap"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/common/vardef"
)

var (
	gRedis = &fkredis.QueueRedis{}
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	// 21645 maze:attr:calc:notify:que 迷宫游戏buff变化通知队列
	fkconfig.RegisterNameNode("dollattrcalcnotifyqueue", 21645, gRedis)
}

func SendMazeAttrCalcNotify(agent fklog.FKLogI, msg *structsdef.MazeCalcAttrNotifyMsg) error {
	if msg.Stamp == 0 {
		msg.Stamp = time.Now().UnixNano() / 1000000
	}

	if msg.ChgDesc == "" {
		msg.ChgDesc = vardef.MazeBuffChgTypeDesc[msg.ChgType]
	}
	key := gRedis.GetKey()

	data, err := json.Marshal(msg)
	if err != nil {
		agent.ErrorWF("SendMazeAttrCalcNotify marshal error",
			zap.Any("msg", msg), zap.String("key", key), zap.Error(err))
		return err
	}

	err = gRedis.PushMsgWithKey(key, data)
	if err != nil {
		agent.ErrorWF("SendMazeAttrCalcNotify send error", zap.Error(err),
			zap.Any("msg", msg),
			zap.String("key", key))
		return err
	}

	agent.InfoWF("SendMazeAttrCalcNotify send succ", zap.Any("msg", msg), zap.String("key", key))
	return nil
}
