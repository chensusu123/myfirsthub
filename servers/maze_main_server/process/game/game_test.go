package game

import (
	"context"
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/usecase/business"
	"os"
	"testing"
)

var logger = log.Clone("gameTest", 0, 0)

func TestMain(m *testing.M) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_ = os.Chdir("C:/work/maze_game_server/servers/maze_main_server/")
	_, err := fkserver.AppServer.Application.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = globalredis.GCli.Init(fileResolver.New("./conf.d/service.yaml"))
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	err = business.GCustomBusiness.Init(fklog.AppLogger().Clone("loadconfigapi"))
	if err != nil {
		logger.ErrorWF("load file failed", zap.Error(err))
		return
	}
	err = config_manager.Init(context.Background(), logger, nil)
	if err != nil {
		logger.ErrorWF("parse excel failed", zap.Error(err))
		return
	}
	m.Run()
}

func TestPickItem(t *testing.T) {
	game := NewGame()
	session := &session.Session{}
	session.Bind(40000001)
	req := &MazeGame.BarrierPickItemRQ{}
	req.ItemList = append(req.ItemList, &MazeCommon.MazeItem{
		ItemId: proto.Int32(46200002),
	})
	req.BarrierId = proto.Int32(1)

	err := game.OnBarrierPickItemRQ_10527_10528(session, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
