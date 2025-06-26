package tempbuffservice

import (
	"fmt"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
	tempbuffredis "maze_game_server/io/redis/tempbuff"
	"maze_game_server/lib/log"
	"maze_game_server/model/passareamodel"
	"maze_game_server/model/tempbuffmodel"
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
	err = tempbuffredis.GCli.Init(fileResolver.New("./conf.d/service.yaml"))
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	m.Run()
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
