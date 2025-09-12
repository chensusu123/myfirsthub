/*
 * @Author: majian
 * @Date: 2024-12-16 17:08:50
 * @Last Modified by: majian
 * @Last Modified time: 2024-12-16 17:12:24
 */
package asynctask

import (
	"context"

	"gitlab.ifreetalk.com/maze-plate/freetk/common/fkfmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkutil/workergroup"
	"go.uber.org/zap"
)

var wgCacheSize, wgThreadCount uint32
var OpenWq = 1 // 开启异步队列

func init() {
	param.Uint32P(&wgCacheSize, "doll:gm:wg:cache:size", 5000, "channel缓冲大小")
	param.Uint32P(&wgThreadCount, "doll:gm:wg:thread:count", 64, "处理线程数量")
}

type tWorkGroupBusiness struct {
	*workergroup.FkWorkGroup
}

var GWorkGroupBusiness = &tWorkGroupBusiness{}

func (tb *tWorkGroupBusiness) Name() string {
	return "tWorkGroupBusiness"
}

func (tb *tWorkGroupBusiness) OnInit(ctx context.Context, cfg fkconfig.FkConfigerI) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	tb.FkWorkGroup = workergroup.NewFkWrokGroup(int(wgCacheSize), int(wgThreadCount), logger)
	logger.CtxInfo(ctx, "tWorkGroupBusiness OnInit work group",
		zap.Uint32("cacheSize", wgCacheSize),
		zap.Uint32("threadCount", wgThreadCount))
	return
}

func (tb *tWorkGroupBusiness) OnStart(ctx context.Context, cfg fkconfig.FkConfigerI) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "tWorkGroupBusiness OnStart.")
	return nil
}

func (tb *tWorkGroupBusiness) OnStop(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "tWorkGroupBusiness OnStop.")
	tb.FkWorkGroup.Stop()
	return
}

func (tb *tWorkGroupBusiness) OnFinish(ctx context.Context) (err error) {
	logger := fklog.ContextAppLogger(ctx)
	logger.CtxInfo(ctx, "tWorkGroupBusiness OnFinish.")
	return
}

func (tb *tWorkGroupBusiness) SendTask(sharding uint64, fun func()) {
	// fkfmt.Println("SendTask ", sharding)
	if fun == nil {
		fkfmt.Println("tWorkGroupBusiness fun nil")
		return
	}
	if tb.FkWorkGroup == nil {
		fkfmt.Println("tWorkGroupBusiness FkWorkGroup not init")
		return
	}

	if OpenWq == 1 {
		tb.FkWorkGroup.SendTask(sharding, fun)
	} else {
		fun()
	}
}
