/*
* @Author: majian
* @Date: 2021-03-04 22:34
 */
package monitor

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkmonitor"
)

type MyMonitor struct {
	fkmonitor.TimeMonitor
	endFun func()
}

var GMonitorManager = map[string]fkmonitor.TimeMonitor{}

/*注册rpc监控*/
func RegRpcMonitor(name string) {
	GMonitorManager[name] = fkmonitor.DefaultTimeMonitor(fmt.Sprintf("RPC.%s", name))
}

/*注册redis监控*/
func RegRedisMonitor(name string, cgkType int32) {
	GMonitorManager[name] = fkmonitor.DefaultTimeMonitor(fmt.Sprintf("Redis.%d%s", cgkType, name))
}

/*获取监控*/
func GetMonitor(name string) MyMonitor {
	st := MyMonitor{}
	if monitor, ok := GMonitorManager[name]; ok {
		st.TimeMonitor = monitor
	}
	return st
}

func (m *MyMonitor) StartV1() {
	if m.TimeMonitor != nil {
		m.endFun = m.TimeMonitor.Start()
	}
}

func (m *MyMonitor) EndV1() {
	if m.endFun != nil {
		m.endFun()
	}
}
