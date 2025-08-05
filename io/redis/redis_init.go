package globalredis

import (
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/database/nanoredis"
	"gitlab.ifreetalk.com/maze-plate/freetk/fkcore/serverdepend"
)

type GlobalRedis struct {
	*nanoredis.NanoRedis
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
