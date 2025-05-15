/*
* @Author: majian
* @Date: 2022-07-01 17:56
 */
package copyinterface

import (
	"fmt"
	"sort"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

type CopyParam struct {
	DelOldData bool // 是否删除旧数据
}
type opFunc = func(logger fklog.FKLogI, srcUserId uint64, dstUserId []uint64, param CopyParam) error

type FuncItem struct {
	CallBackF opFunc
	Name      string
	Order     int32
}

var handlers = make(map[string]FuncItem, 20)

func RegistHandler(name string, h opFunc, order int32) {
	_, ok := handlers[name]
	if ok {
		fmt.Println("RegistHandler already register handler", name)
		return
	}
	handlers[name] = FuncItem{CallBackF: h, Name: name, Order: order}
}

// 模板用户拷贝到目标用户列表
func RangeAiCopy(logger fklog.FKLogI, srcUserId uint64, dstUserId []uint64, param CopyParam) error {
	var err error
	var orderFuncList []FuncItem
	for _, f := range handlers {
		orderFuncList = append(orderFuncList, f)
	}
	sort.Slice(orderFuncList, func(i, j int) bool {
		return orderFuncList[i].Order < orderFuncList[j].Order
	})
	totalNow := time.Now()
	for _, item := range orderFuncList {
		moduleStartTime := time.Now()
		err = item.CallBackF(logger, srcUserId, dstUserId, param)
		if err != nil {
			logger.ErrorWF("RangeAiCopy fail",
				zap.Error(err),
				zap.String("name", item.Name),
				zap.Uint64("src", srcUserId),
				zap.Duration("cost", time.Since(moduleStartTime)),
				zap.Int("dstUsersCnt", len(dstUserId)))
		} else {
			logger.WarnWF("RangeAiCopy succ",
				zap.String("name", item.Name),
				zap.Uint64("src", srcUserId),
				zap.Duration("cost", time.Since(moduleStartTime)),
				zap.Int("dstUsersCnt", len(dstUserId)))
		}
	}
	logger.WarnWF("RangeAiCopy result",
		zap.Uint64("src", srcUserId),
		zap.Duration("cost", time.Since(totalNow)),
		zap.Int("dstUsersCnt", len(dstUserId)))

	return nil
}
