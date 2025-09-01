package cluster

import "gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkalert"

func safeSend[T any](ch chan T, data T) {
	defer fkalert.RecoverAlertException()
	ch <- data
}
