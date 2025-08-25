package mailservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	globalredis "maze_game_server/io/redis"
	"maze_game_server/lib/log"
	"maze_game_server/model/mailmodel"

	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkserver"
	fileResolver "gitlab.ifreetalk.com/maze-plate/freetk/registry/fileresolver"
)

var logger = log.Clone("EquipDropTest", 0, 0)

func TestMain(m *testing.M) {
	originalStdout := os.Stdout
	originalStderr := os.Stderr
	_ = os.Chdir("D:/work/maze_game_server/servers/maze_main_server/")
	_, err := fkserver.AppServer.Application.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = globalredis.GCli.Init(fileResolver.New("D:/work/maze_game_server/servers/maze_main_server/conf.d/service.yaml"))
	if err != nil {
		fmt.Println(err)
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

func TestGetMailList(t *testing.T) {
	mailList, err := GlobalMailService.GetMailListByLabel(logger, 40000005, 0, 0, 30)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mailList)
}

func TestMail(t *testing.T) {
	attachments := []*mailmodel.Attachment{
		{
			ItemID: 46700001,
			Count:  10,
			Extra:  "equip",
		},
		{
			ItemID: 46200001,
			Count:  10,
			Extra:  "item",
		},
	}

	err := GlobalMailService.SendMail(logger, "测试邮件", "邮件内容", "发送者", 0, 40000005, attachments, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	mailList, err := GlobalMailService.ReadAllMail(logger, 40000005, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mailList)

	mailList, attachs, err := GlobalMailService.GetAllMailAttachment(logger, 40000005, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mailList)
	fmt.Println(attachs)

	err = GlobalMailService.GetMailAttachmentAfter(logger, 40000005, mailList)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = GlobalMailService.DelAllMail(logger, 40000005, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	mailMap, err := GlobalMailService.GetAllMailList(logger, 40000005)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mailMap)
}
