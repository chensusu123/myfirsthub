package tempbuffservice

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
	"maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/model/passareamodel"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/usecase/business"
	"os"
	"testing"
)

var logger = log.Clone("TempBuffTest", 0, 0)

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
func TestTempBuffRedis(t *testing.T) {
	err, model := tempbuffmodel.NewTempBuffInfoModel(logger, 40000001, 31)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(model)
}

func TestPassAreaRedis(t *testing.T) {
	err, model := passareamodel.NewPassAreaModel(logger, 40000001, 31)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(model)
}
func TestGetOptionalTempBuff(t *testing.T) {
	optionalBuffInfo, err := GlobalTempBuffService.GetOptionalTempBuffList(logger, 40000001, 1, 2, 1, 10001, 0, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = optionalBuffInfo
}
