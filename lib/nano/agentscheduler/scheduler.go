// Copyright (c) nano Authors. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package agentscheduler

import (
	"sync/atomic"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/metricsreport"
	"go.uber.org/zap"
)

const (
	messageQueueBacklog = 1 << 10
	sessionCloseBacklog = 1 << 8
)

type Task func()

type Hook func()

type AgentScheduler struct {
	agentSession int64
	chDie        chan struct{}
	chExit       chan struct{}
	chTasks      chan Task
	started      int32
	closed       int32
	taskCount    atomic.Int64
}

func NewAgentScheduler(agentSession int64) *AgentScheduler {
	return &AgentScheduler{
		agentSession: agentSession,
		chDie:        make(chan struct{}),
		chExit:       make(chan struct{}),
		chTasks:      make(chan Task, 2048),
	}
}

func (ac *AgentScheduler) SetAgentSession(c int64) {
	ac.agentSession = c
}

// var (
// 	chDie     = make(chan struct{})
// 	chExit    = make(chan struct{})
// 	chTasks   = make(chan Task, 1<<8)
// 	started   int32
// 	closed    int32
// 	taskCount atomic.Int64
// )

func (ac *AgentScheduler) try(f func()) {
	defer func() {
		fkalert.RecoverAlertException()
		c := ac.taskCount.Add(-1)
		_ = c
		// nanometrics.GlobalTaskGauge.Set(float64(c))
	}()
	f()
}

func (ac *AgentScheduler) Sched() {
	if atomic.AddInt32(&ac.started, 1) != 1 {
		return
	}

	defer func() {
		close(ac.chExit)
		fklog.AppLogger().InfoWF("AgentScheduler Sched end 1",
			zap.Int64("agentSession", int64(ac.agentSession)),
		)
		close(ac.chTasks)
		closeTaskCount := 0

		for f := range ac.chTasks {
			ac.try(f)
			closeTaskCount += 1
		}

		fklog.AppLogger().InfoWF("AgentScheduler Sched end 2",
			zap.Int64("closeTaskCount", int64(closeTaskCount)),
			zap.Int64("agentSession", int64(ac.agentSession)),
		)
	}()

	for {
		select {
		case f := <-ac.chTasks:
			ac.try(f)
		case <-ac.chDie:
			return
		}
	}
}

func (ac *AgentScheduler) Close() {
	fklog.AppLogger().InfoWF("AgentScheduler Close 1",
		zap.Int64("agentSession", int64(ac.agentSession)),
	)
	if atomic.AddInt32(&ac.closed, 1) != 1 {
		return
	}
	fklog.AppLogger().InfoWF("AgentScheduler Close 2",

		zap.Int64("agentSession", int64(ac.agentSession)),
	)
	close(ac.chDie)
	<-ac.chExit
	// log.Println("Scheduler stopped")
}

func (ac *AgentScheduler) PushTask(task Task) (int64, bool) {
	c := ac.taskCount.Add(1)
	ok := safeSend(ac.chTasks, task)
	if !ok {
		ac.taskCount.Add(-1)
	}
	return c, ok
}

func (ac *AgentScheduler) TaskCount() int64 {
	return ac.taskCount.Load()
}

func safeSend[T any](ch chan T, data T) (ok bool) {
	ok = true
	defer func() {
		if err := recover(); err != nil {
			ok = false
			metricsreport.PanicNum.Incr()
		}
	}()
	ch <- data
	return
}
