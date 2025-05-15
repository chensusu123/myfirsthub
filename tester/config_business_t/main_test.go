package config_business_t

import (
	"fmt"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	_ "gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testlogger" // 初始化日志
)

var gTestLogger fklog.FKLogI

func TestMain(m *testing.M) {
	fmt.Println("begin")
	gTestLogger = fklog.AppLogger().Clone("config_business_t")

	m.Run()
	fmt.Println("end")

	time.Sleep(time.Second * 2)
}
