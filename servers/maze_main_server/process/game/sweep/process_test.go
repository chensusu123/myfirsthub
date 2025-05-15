// @Author pangchenyang 2025/3/21 19:28:00
// @Desc:
package sweep

import (
	"context"
	"fmt"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testio"
	_ "gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testlogger"
)

var (
	logger = fklog.AppLogger().Clone("process_t")
	tcpCtx = fknet.TCPContext{Context: context.Background(), FkTransportBase: nil, FKLogI: logger}
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = testio.IOLoad(9)
	m.Run()
	fmt.Println("end")
}

// func TestOnStartMazeSweepRQ(t *testing.T) {
// 	logger.SetLogId(time.Now().UnixNano())
// 	defer func() {
// 		time.Sleep(time.Second * 2)
// 	}()
// 	startSweepRq := &MazeGame.StartMazeSweepRQ{
// 		BarrierId: proto.Int32(1),
// 	}
// 	startSweepRs := &MazeGame.StartMazeSweepRS{}
// 	_ = OnStartMazeSweepRQ(tcpCtx, 1, startSweepRq, startSweepRs)
// }
