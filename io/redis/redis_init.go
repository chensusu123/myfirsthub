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

func (g *GlobalRedis) Set(key string, value []byte) error {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return err
	}
	return db.Set(context.TODO(), string(key), value, 0).Err()
}

func (g *GlobalRedis) Get(key string) ([]byte, error) {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return nil, err
	}

	ret, err := db.Get(context.TODO(), key).Bytes()
	if errors.Is(err, redis.Nil) {
		return []byte{}, nil
	}
	return ret, err
}

func (g *GlobalRedis) Del(key string) error {
	db, err := g.NanoRedis.GetDB()
	if err != nil {
		return err
	}
	return db.Del(context.TODO(), key).Err()
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
