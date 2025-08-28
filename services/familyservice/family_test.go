package familyservice

import (
	"context"
	"fmt"
	globalredis "maze_game_server/io/redis"
	"maze_game_server/model/familymodel"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
)

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

func TestCreateFamily(t *testing.T) {
	ctx := context.Background()
	family, err := GlobalFamilyService.CreateFamily(ctx, uint64(50000001), 1, "family_name1111", 1, familymodel.FamilyMember{
		UserID:         uint64(50000001),
		NickName:       "nick_name1111",
		Avatar:         "avatar1111",
		Sex:            1,
		Level:          1,
		PrivilegeLevel: 1,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("family----------:", family)
}
