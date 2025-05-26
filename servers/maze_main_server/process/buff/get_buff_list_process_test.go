package buff

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"gitlab.ifreetalk.com/maze-plate/extra/protobuf/proto"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fkconfig"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fknet"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/frontcache_service"
	"gitlab.ifreetalk.com/maze-plate/freetk/fktestutil/testio"
	"maze_game_server/pb/common/MazeTempBuff"
)

/**
 * @Author: liushuhang
 * @Date: 2025/3/22 16:13
 * @Description:
 */

var gTestLogger fklog.FKLogI

func TestMain(m *testing.M) {
	fmt.Println("begin")
	// initLog()
	rand.Seed(time.Now().UnixNano())
	gTestLogger = fklog.AppLogger().Clone("test_maze_temp_buff")

	testio.IOLoad(12, "90038") // 加载N组的io配置

	frontcache_service.PlugFCService()
	fkconfig.EnvVal.AppName = "maze-temp-buff-server"
	fkconfig.EnvVal.GroupID = 12
	fkconfig.EnvVal.ServerID = 20005
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

func TestGetMazeTempBuffListRQ(t *testing.T) {
	type args struct {
		logger     fknet.TCPContext
		shardingID uint64
		request    *MazeTempBuff.GetMazeTempBuffListRQ
		response   *MazeTempBuff.GetMazeTempBuffListRS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "查询buff列表",
			args: args{
				logger: fknet.TCPContext{
					Context: context.Background(),
					FKLogI:  gTestLogger,
				},
				shardingID: 9003200130019765,
				request: &MazeTempBuff.GetMazeTempBuffListRQ{
					Header:  nil,
					StageId: proto.Int32(1),
				},
				response: &MazeTempBuff.GetMazeTempBuffListRS{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetMazeTempBuffListRQ(tt.args.logger, tt.args.shardingID, tt.args.request, tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("GetMazeTempBuffListRQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
