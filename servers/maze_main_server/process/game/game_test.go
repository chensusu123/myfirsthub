package game

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/lib/nano/session"
	"maze_game_server/pb/common/MazeCommon"
	"maze_game_server/pb/common/MazeGame"
	"maze_game_server/usecase/business"
	"os"
	"testing"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	ctx := context.Background()
	loadconfigapi.SetLoadConfigFunc(business.GCustomBusiness.LoadCacheConfig)
	err = business.GCustomBusiness.Init(fklog.AppLogger().Clone("loadconfigapi"))
	if err != nil {
		logger.CtxError(ctx, "load file failed", zap.Error(err))
		return
	}
	err = config_manager.Init(context.Background(), logger, nil)
	if err != nil {
		logger.CtxError(ctx, "parse excel failed", zap.Error(err))
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

func TestSaveBarrierData(t *testing.T) {
	game := NewGame()
	session := &session.Session{}
	session.Bind(40000003)
	req := &MazeGame.SaveBarrierDataRQ{}
	req.StageId = proto.Int32(1)
	req.BarrierId = proto.Int32(2)
	req.RescueValue = proto.Int32(100)
	req.BossPower = proto.Int32(1000)
	err := game.OnSaveBarrierDataRQ_10624_10625(session, req)
	if err != nil {
		fmt.Println(err)
		return
	}
}
