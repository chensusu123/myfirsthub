package maze_es_t

import (
	"fmt"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	//"gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testio"
	_ "gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testlogger" // 初始化日志
)

var gTestLogger fklog.FKLogI

func TestMain(m *testing.M) {
	fmt.Println("begin")
	gTestLogger = fklog.AppLogger().Clone("maze_es_t")

	// testio.IOLoad(1, "90032") // 加载N组的io配置

	m.Run()
	fmt.Println("end")

	time.Sleep(time.Second * 2)
}
