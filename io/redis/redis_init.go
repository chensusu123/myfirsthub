package globalredis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
)

type GlobalRedis struct {
	*nanoredis.NanoRedis
}

func (g *GlobalRedis) Set(ctx context.Context, key string, value []byte) error {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return err
	}
	return db.Set(ctx, key, value, 0).Err()
}

func (g *GlobalRedis) Get(ctx context.Context, key string) ([]byte, error) {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return nil, err
	}

	ret, err := db.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return []byte{}, nil
	}
	return ret, err
}

func (g *GlobalRedis) Del(ctx context.Context, key string) error {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return err
	}
	return db.Del(ctx, key).Err()
}

func New(serviceName string, name string) *GlobalRedis {
	rt := &GlobalRedis{}
	rt.NanoRedis = nanoredis.NewNanoRedis(serviceName, name)
	return rt
}

// 全局redis实例
var GCli *GlobalRedis

func init() {
	GCli = New("maze-game.redis", "maze-game.redis.main")
	serverdepend.RegisterDepend(GCli)
}
