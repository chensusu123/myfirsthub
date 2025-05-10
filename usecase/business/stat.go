package business

import (
	"fmt"
	"sync/atomic"
)

var StatIns = &Stat{}

type Stat struct {
	loadNum, loadAllNum, noLoadAllNum, bufferFullNum uint64 //配置垃取次数，加载全部次数，非加载全部次数，写入kafka的通道满了之后丢弃的消息数
	sheetMd5CalNum, sheetMd5CalCost                  uint64 //sheetDataMd5计算次数、sheetDataMd5计算消耗
}

func (s *Stat) AddLoadNum() {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.loadNum, 1)
}

func (s *Stat) AddBufferFullNum() {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.bufferFullNum, 1)
}

func (s *Stat) AddLoadAllNum() {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.loadAllNum, 1)
}

func (s *Stat) AddNoLoadAllNum() {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.noLoadAllNum, 1)
}

func (s *Stat) AddSheetMd5CalNum() {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.sheetMd5CalNum, 1)
}

func (s *Stat) AddSheetMd5CalCost(cost uint64) {
	if statOpen == 0 {
		return
	}
	atomic.AddUint64(&s.sheetMd5CalCost, cost)
}

func (s *Stat) Show() string {
	if statOpen == 0 {
		return "未打开程序统计"
	}
	return fmt.Sprintf(
		"配置垃取次数:%v\t加载全部次数:%v\t非加载全部次数:%v\t写入kafka的通道满了之后丢弃的消息数:%v\t计算sheetDataMd5次数:%v\t计算sheetDataMd5消耗:%vms\n",
		atomic.LoadUint64(&s.loadNum), atomic.LoadUint64(&s.loadAllNum), atomic.LoadUint64(&s.noLoadAllNum),
		atomic.LoadUint64(&s.bufferFullNum), atomic.LoadUint64(&s.sheetMd5CalNum), atomic.LoadUint64(&s.sheetMd5CalCost))
}

func (s *Stat) Reset() {
	StatIns = &Stat{}
}
