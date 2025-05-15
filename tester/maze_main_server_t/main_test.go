/*
* @Author: majian
* @Date: 2022-07-28 20:19
 */
package maze_main_server_t

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testio"
)

var gTestLogger fklog.FKLogI

func TestMain(m *testing.M) {
	fmt.Println("begin")
	initLog()
	gTestLogger = fklog.AppLogger().Clone("maze_main_server_t")

	testio.IOLoad(12) //加载N组的io配置
	rand.Seed(time.Now().UnixNano())
	m.Run()
	fmt.Println("end")
}

func initLog() {
	logConfig := fklog.LogConfig{
		LogDir:     ".",
		LogLevel:   "debug",
		LogName:    "new_test",
		LogType:    "zap",
		WithCaller: true,
	}
	fklog.InitAppFkLog(&logConfig)
}
