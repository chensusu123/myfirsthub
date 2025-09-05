package cluster

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/pkg/metricsreport"
)

func safeSend[T any](ch chan T, data T) (ok bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			metricsreport.PanicNum.Incr()
		}
	}()
	ch <- data
	return
}
