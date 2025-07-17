package usersection

import (
	"context"
	"fmt"

	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
)

type UserSection struct {
	*nanoredis.NanoRedis
}

func New(serviceName string, name string) *UserSection {
	rt := &UserSection{}
	rt.NanoRedis = nanoredis.NewNanoRedis(serviceName, name)
	return rt
}

func (g *UserSection) Set(ctx context.Context, userID uint64, section string) error {
	key := fmt.Sprintf("user:section:%d", userID)
	db, err := g.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), key, section, 0).Err()
}

var gCli *UserSection

func init() {
	gCli = New("user.section.redis", "user.section.redis.maze_main")
	serverdepend.RegisterDepend(gCli)
}

func Set(ctx context.Context, userID uint64, section string) error {
	return gCli.Set(ctx, userID, section)
}
