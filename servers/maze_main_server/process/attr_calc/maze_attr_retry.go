/*
 * @Author: majian
 * @Date: 2024-07-17 17:03:00
 * @Last Modified by: majian
 * @Last Modified time: 2025-03-21 20:49:05
 */
package attr_calc

import (
	"context"
	"sync/atomic"
	"time"

	"maze_game_server/common/structsdef"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var (
	RetryTimeInterval  int32 = 0 // 重试时间间隔
	RetryMaxCount      int32 = 0 // 失败重试最大次数
	RetryMaxTime       int32 = 0 // 重试最长时间
	TestSwitch         int32 = 0 // 测试开关
	ConRetryCounter    int64 = 0 // 并发计数
	ConRetryCounterMax int64 = 0 // 并发计数上限
)

func init() {
	param.Int32P(&RetryTimeInterval, "calc:fail:retry:time", 2000, "失败重试时间ms")
	param.Int32P(&RetryMaxCount, "calc:fail:retry:max:times", 3, "失败重试最大次数")
	param.Int32P(&RetryMaxTime, "calc:fail:retry:max:time", 60, "失败重试最长时间s")
	param.Int32P(&TestSwitch, "calc:fail:retry:test", 0, "失败重试测试")
	param.Int64P(&ConRetryCounterMax, "con:retry:cnt:limit", 1000, "同时重试数量上限")
}
func doMazeAttrCalcRetry(ctx context.Context, msg *structsdef.MazeCalcAttrNotifyMsg) {
	logger := fklog.ContextAppLogger(ctx)
	lastTime := msg.Stamp / 1000 // 转成秒
	now := time.Now().Unix()
	if now >= lastTime+int64(RetryMaxTime) {
		logger.CtxWarn(ctx, "doMazeAttrCalcRetry to max time",
			zap.Int64("lastTime", msg.Stamp),
			zap.Int32("maxTime", RetryMaxTime))
		return
	}
	if msg.RetryFlag >= RetryMaxCount {
		logger.CtxWarn(ctx, "doMazeAttrCalcRetry to max retry times",
			zap.Int32("retryCount", msg.RetryFlag),
			zap.Int32("maxTimes", RetryMaxCount))
		return
	}
	msg.RetryFlag++ // 不需要考虑竞争
	atomic.AddInt64(&ConRetryCounter, 1)
	curRetryMax := atomic.LoadInt64(&ConRetryCounter)
	if curRetryMax >= ConRetryCounterMax {
		logger.CtxWarn(ctx, "doMazeAttrCalcRetry Concurrency retry to limit",
			zap.Int64("curConRetryMax", curRetryMax),
			zap.Int64("conretryLimit", ConRetryCounterMax))
		return
	}
	time.AfterFunc(time.Millisecond*time.Duration(RetryTimeInterval), func() {
		// mazeattrcalcnotifyqueue.SendMazeAttrCalcNotify(logger, msg)
		err := OnMazeAttrCalcMsg(ctx, logger, 0, msg)
		if err != nil {
			logger.CtxError(ctx, "doMazeAttrCalcRetry OnMazeAttrCalcMsg failed", zap.Any("msg", msg), zap.Error(err))
		}
		atomic.AddInt64(&ConRetryCounter, -1)
	})
}
