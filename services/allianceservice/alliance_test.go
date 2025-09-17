package allianceservice

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
)

var logger = log.Clone("allianceservice", 0, 0)

func TestMain(m *testing.M) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_ = os.Chdir("D:/src/maze_game_server/servers/maze_main_server/")
	_, err := fkserver.AppServer.Application.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = globalredis.GCli.Init(fileResolver.New("D:/src/maze_game_server/servers/maze_main_server/conf.d/service.yaml"))
	if err != nil {
		fmt.Println("Init redis failed err:", err)
		return
	}
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	m.Run()
}
func TestRedis(t *testing.T) {
	db, err := globalredis.GCli.GetDB()
	if err != nil {
		return
	}
	res, err := db.Get(context.Background(), "123").Result()
	if err != nil {
		if err == redis.Nil {
			return
		}
		return
	}
	fmt.Println(res)
}

func TestCreateAlliance(t *testing.T) {
	_, err := GlobalAllianceService.CreateAlliance(context.Background(), "test002")
	if err != nil {
		logger.ErrorWF("TestCreateAlliance error", zap.Error(err))
	}
	fmt.Println("TestCreateAlliance success")
}

func TestSubscribeAllianceChat(t *testing.T) {
	err := GlobalAllianceService.SubscribeAllianceChat(context.Background(), 50000002, 1, []int64{10000003, 10000005})

	if err != nil {
		logger.ErrorWF("TestSubscribeAllianceChat error", zap.Error(err))
	}
	fmt.Println("TestSubscribeAllianceChat success")
}
