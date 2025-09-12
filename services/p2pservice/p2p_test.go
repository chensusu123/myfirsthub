package p2pservice

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
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}

	messageID, err := GlobalP2PService.SendMessage(context.Background(), app.Maze, user, int64(50000002), int32(1), []byte("hello03"))
	if err != nil {
		logger.ErrorWF("OnSendMessage SendMessage error", zap.Error(err))
	}
	fmt.Println("messageID:", messageID)

}

func TestQueryMessages(t *testing.T) {
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	messages, err := GlobalP2PService.QueryMessages(context.Background(), app.Maze, user, int64(50000002), 621539715899425009, true)
	if err != nil {
		logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	}
	fmt.Println("messages:", messages)

}

// 读消息
func TestReadMessage(t *testing.T) {
	userId := uint64(50000002)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
		return
	}
	err = GlobalP2PService.ReadMessage(context.Background(), app.Maze, user, int64(50000001), 618382372351268509)
	if err != nil {
		logger.ErrorWF("OnSendMessage ReadMessage error", zap.Error(err))
	}
	fmt.Println("read message success")
}

func TestRemoveMessage(t *testing.T) {
	userId := uint64(50000002)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	err = GlobalP2PService.RemoveMessage(context.Background(), app.Maze, user, int64(50000001), 34719007919596739)
	if err != nil {
		logger.ErrorWF("OnSendMessage RemoveMessage error", zap.Error(err))
	}

	fmt.Println("remove message success")
}
