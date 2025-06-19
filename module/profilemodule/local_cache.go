// @Author pangchenyang 2025/6/9 21:33:00
// @Desc: 
package profilemodule

import (
	"sync"
	lru "github.com/hashicorp/golang-lru"
)

var (
	cache *lru.Cache
	// 初始化本地缓存
	once sync.Once
)

func initUserProfileCache() {
	once.Do(func() {
		var err error
		cache, err = lru.New(10000) // 缓存1万个用户资料
		if err != nil {
			panic(err)
		}
	})
}
