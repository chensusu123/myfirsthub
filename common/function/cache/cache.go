//@Desc : \\todo
//@Author : zhangdengyuan 2022/3/16 11:05 下午
//@Update: xxx 2022/3/16 11:05 下午

package cache

import (
	"sync"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/common/fkfmt"
)

type cacheData struct {
	ts      int64
	content interface{}
}
type mapCache struct {
	sync.RWMutex
	cache    map[uint64]*cacheData // cache
	ctime    int64                 // 缓存时长
	size     int                   // 缓存数量
	cleanGap int                   // 定时清理过期缓存,定时间隔
	ticker   *time.Ticker
}

// cacheTime 缓存时间
// cacheSize 缓存大小
// cleanGap 定期检查时间
func NewCache(cacheTime int64, cacheSize, cleanGap int) *mapCache {
	cache := &mapCache{ctime: cacheTime, size: cacheSize, cleanGap: cleanGap,
		cache: make(map[uint64]*cacheData, cacheSize)}
	if cache.cleanGap > 0 {
		cache.ticker = time.NewTicker(time.Second * time.Duration(cache.cleanGap))
		go cache.loopClean()
	}
	return cache
}

func (r *mapCache) SetCache(key uint64, content interface{}) {
	r.Lock()
	defer r.Unlock()

	cell := &cacheData{}
	cell.ts = time.Now().Unix()
	cell.content = content
	r.cache[key] = cell
}

func (r *mapCache) DelCache(key uint64) {
	r.Lock()
	defer r.Unlock()

	delete(r.cache, key)
}

func (r *mapCache) GetCache(key uint64) (interface{}, bool) {
	r.RLock()
	defer r.RUnlock()

	cell, ok := r.cache[key]
	if !ok {
		return 0, false
	}
	return cell.content, true
}

func (r *mapCache) loopClean() {
	for {

		select {
		case <-r.ticker.C:

			start := time.Now()
			r.Lock()
			oldLen := len(r.cache)
			for k, v := range r.cache {
				if start.Unix()-v.ts > r.ctime {
					delete(r.cache, k)
				}
			}
			newLen := len(r.cache)
			r.Unlock()

			duration := time.Since(start)
			if duration > 1*time.Second {
				fkfmt.Println("cache loopClean time:", duration, "oldLen:", oldLen, "newLen:", newLen)
			}
		}
	}
}
