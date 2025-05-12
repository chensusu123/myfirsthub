/*
 * @Author: majian
 * @Date: 2024-07-10 15:23:44
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-14 22:19:14
 */
package mazeattrcalcnotifyqueue

import (
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkredis"

	"context"

	"gitlab.ifreetalk.com/maze/maze_game_server/common/structsdef"
	"gitlab.ifreetalk.com/maze/maze_game_server/servers/maze_main_server/process/attr_calc"
	"go.uber.org/zap"
)

var (
	gRedis = &fkredis.QueueRedis{}
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	// 21645 maze:attr:calc:notify:que 迷宫游戏buff变化通知队列
	fkconfig.RegisterNameNode("dollattrcalcnotifyqueue", 21645, gRedis)
}

type MazeCalcAttrNotifyMsg struct {
	UserId     uint64 `json:"user_id"`     // 用户Id
	FromServer string `json:"from_server"` // 服务来源(服务类型+名字)
	BuffSrc    int32  `json:"buff_src"`    // buff来源
	ChgType    int32  `json:"chg_type"`    // 变化类型
	ChgDesc    string `json:"chg_desc"`    // 变化原因描述
	Session    string `json:"session"`     // session
	Stamp      int64  `json:"stamp"`       // 消息时间戳 ms
	RetryFlag  int32  `json:"retry_flag"`  // 失败重试用,内部用不用设置
}

func SendMazeAttrCalcNotify(agent fklog.FKLogI, msg *structsdef.MazeCalcAttrNotifyMsg) error {
	agent.InfoWF("SendMazeAttrCalcNotify start", zap.Any("msg", msg))
	return attr_calc.OnMazeAttrCalcMsg(context.TODO(), agent, 0, msg)
	if msg.Stamp == 0 {
		msg.Stamp = time.Now().UnixNano() / 1000000
	}

	key := fmt.Sprintf("maze:attr:calc:notify:que")

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
