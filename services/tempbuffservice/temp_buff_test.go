package tempbuffservice

import (
	"context"
	"fmt"
	"maze_game_server/io"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/model/passareamodel"
	"maze_game_server/model/tempbuffmodel"
	"maze_game_server/usecase/business"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/fklog"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver/config_manager/loadconfigapi"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	"go.uber.org/zap"
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
	ctx := context.Background()
	if err != nil {
		logger.CtxError(ctx, "load file failed", zap.Error(err))
		return
	}
	err = config_manager.Init(context.Background(), logger, nil)
	if err != nil {
		logger.CtxError(ctx, "parse excel failed", zap.Error(err))
		return
	}
	io.InitBackendCoder(globalredis.GCli, nil)
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
	err, model := tempbuffmodel.NewTempBuffInfoModel(context.TODO(), 40000001, 31)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(model)
}

func TestPassAreaRedis(t *testing.T) {
	err, model := passareamodel.NewPassAreaModel(context.TODO(), 40000001, 31)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(model)
}

func TestGetOptionalTempBuff(t *testing.T) {
	optionalBuffInfo, err := GlobalTempBuffService.GetOptionalTempBuffList(context.TODO(), 40000001, 1, 2, 1, 10001, 0, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = optionalBuffInfo
}

func TestGetTempBuffGroupList(t *testing.T) {
	groupList, err := GlobalTempBuffService.GetTempBuffGroupList(context.TODO(), 50000001, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = groupList
}

func TestCalcLibraryAddWeight(t *testing.T) {
	s := &service{}
	weight1 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 0)
	weight2 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 1)
	weight3 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 2)
	weight4 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 3)
	weight5 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 4)
	weight6 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 5)
	weight7 := s.calcLibraryAddWeight(100, []int32{5000, 6000, 7000, 8000}, 6)
	weight8 := s.calcLibraryAddWeight(100, []int32{}, 1)
	weight9 := s.calcLibraryAddWeight(100, []int32{}, 2)
	fmt.Println(weight1, weight2, weight3, weight4, weight5, weight6, weight7, weight8, weight9)
}

func TestCalcLibrarySubWeight(t *testing.T) {
	s := &service{}
	weight1 := s.calcLibrarySubWeight(100, []int32{}, 0)
	weight2 := s.calcLibrarySubWeight(100, []int32{}, 1)
	weight3 := s.calcLibrarySubWeight(100, []int32{5000}, 1)
	weight4 := s.calcLibrarySubWeight(100, []int32{5000}, 2)
	weight5 := s.calcLibrarySubWeight(100, []int32{5000, 10000}, 1)
	weight6 := s.calcLibrarySubWeight(100, []int32{5000, 10000}, 2)
	weight7 := s.calcLibrarySubWeight(100, []int32{5000, 10000}, 0)
	fmt.Println(weight1, weight2, weight3, weight4, weight5, weight6, weight7)
}
