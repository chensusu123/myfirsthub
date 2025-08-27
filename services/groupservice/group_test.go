package groupservice

import (
	"context"
	"fmt"
	"maze_game_server/app"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
)

var logger = log.Clone("p2pservice", 0, 0)

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

func TestSendMessages(t *testing.T) {

	messageID, err := GlobalGroupService.SendMessage(context.Background(), app.Maze, int32(1), uint64(50000001), int32(1), "hello01")
	if err != nil {
		logger.ErrorWF("SendGroupMessageTest error", zap.Error(err))
	}
	fmt.Println("messageID:", messageID)

}

func TestQueryMessages(t *testing.T) {
	messages, err := GlobalGroupService.QueryMessages(context.Background(), app.Maze, int32(1), uint64(50000001), int(10))
	if err != nil {
		logger.ErrorWF("QueryGroupMessagesTest error", zap.Error(err))
	}
	fmt.Println("messages:", messages)

}
