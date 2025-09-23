package sessionservice

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

func TestQueryRecentSessions(t *testing.T) {
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	session, err := Default.QueryRecentSessions(context.Background(), app.Maze, user)
	if err != nil {
		logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	}
	fmt.Println("messages:", session)
}

// func TestCreateNormalSession(t *testing.T) {
// 	userId := uint64(50000001)
// 	user, err := app.WrapUser(userId, "")
// 	if err != nil {
// 		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
// 	}
// 	err = Default.CreateNormalSession(context.Background(), app.Maze, user, 50000002, 1621512345, true)
// 	if err != nil {
// 		logger.ErrorWF("OnSendMessage CreateNormalSession error", zap.Error(err))
// 	}
// }

// func TestCreateGroupSession(t *testing.T) {
// 	userId := uint64(50000001)
// 	user, err := app.WrapUser(userId, "")
// 	if err != nil {
// 		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
// 	}
// 	err = Default.CreateGroupSession(context.Background(), app.Maze, user, 10000001)
// 	if err != nil {
// 		logger.ErrorWF("OnSendMessage CreateGroupSession error", zap.Error(err))
// 	}
// }

func TestRemoveSession(t *testing.T) {
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	err = Default.RemoveSession(context.Background(), app.Maze, user, "123")
	if err != nil {
		logger.ErrorWF("OnSendMessage RemoveSession error", zap.Error(err))
	}
}

func TestGetMessageInfo(t *testing.T) {
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	session, err := Default.QueryRecentSessions(context.Background(), app.Maze, user)
	if err != nil {
		logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	}
	fmt.Println("messages:", session)
	messageInfo, err := Default.GetMessageInfo(context.Background(), app.Maze, user, session)
	if err != nil {
		logger.ErrorWF("OnSendMessage GetMessageInfo error", zap.Error(err))
	}
	fmt.Println("messageInfo:", messageInfo)
}

func TestUpdateNormalSession(t *testing.T) {
	// userId := uint64(50000001)
	// user, err := app.WrapUser(userId, "")
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	// }
	// session, err := Default.QueryRecentSessions(context.Background(), app.Maze, user)
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	// }
	// for _, v := range session {
	// 	if v.PeerInfo.UserID() == 50000002 {
	// 		err = Default.UpdateNormalSession(context.Background(), app.Maze, user, 50000002, v)
	// 		if err != nil {
	// 			logger.ErrorWF("OnSendMessage UpdateNormalSession error", zap.Error(err))
	// 		}
	// 		break
	// 	}
	// }
}

func TestSaveNormalSession(t *testing.T) {
	userId := uint64(50000001)
	user, err := app.WrapUser(userId, "")
	if err != nil {
		logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	}
	err = Default.SaveNormalSession(context.Background(), app.Maze, user, 50000002, 1621512345, true)
	if err != nil {
		logger.ErrorWF("OnSendMessage SaveNormalSession error", zap.Error(err))
	}
}

func TestNotifyNormalSession(t *testing.T) {
	// userId := uint64(50000001)
	// user, err := app.WrapUser(userId, "")
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	// }
	// session, err := Default.QueryRecentSessions(context.Background(), app.Maze, user)
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	// }
	// for _, v := range session {
	// 	if v.PeerInfo.UserID() == 50000002 {
	// 		err = Default.NotifyNormalSession(context.Background(), app.Maze, 50000002, v)
	// 		if err != nil {
	// 			logger.ErrorWF("OnSendMessage NotifyNormalSession error", zap.Error(err))
	// 		}
	// 		break
	// 	}
	// }
}

func TestNotifyRemoveSession(t *testing.T) {
	// userId := uint64(50000001)
	// user, err := app.WrapUser(userId, "")
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage WrapUser error", zap.Error(err))
	// }
	// session, err := Default.QueryRecentSessions(context.Background(), app.Maze, user)
	// if err != nil {
	// 	logger.ErrorWF("OnSendMessage QueryMessages error", zap.Error(err))
	// }
	// for _, v := range session {
	// if v.PeerInfo.UserID() == 50000002 {
	// 	err = Default.NotifyRemoveSession(context.Background(), app.Maze, 50000002, v)
	// 	if err != nil {
	// 		logger.ErrorWF("OnSendMessage NotifyRemoveSession error", zap.Error(err))
	// 	}
	// 	break
	// }
	// 	}
}
