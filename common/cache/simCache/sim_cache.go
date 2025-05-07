// @Author: ZhaoXiming 2022/7/19 15:32
// @Desc: 简单缓存 需要自己clear
package simCache

import (
	"sync"
)

type cacheData struct {
	//ts      int64
	content interface{}
}
type mapCache struct {
	sync.RWMutex
	cache map[interface{}]*cacheData // cache
}

func NewCache() *mapCache {
	cache := &mapCache{cache: make(map[interface{}]*cacheData)}
	return cache
}

// 已存在的 不能set, exist=true
func (r *mapCache) SetCache(key interface{}, content interface{}) (exist bool) {
	r.Lock()
	defer r.Unlock()

	_, ok := r.cache[key]
	if ok {
		return true
	}

	cell := &cacheData{}
	//cell.ts = time.Now().Unix()
	cell.content = content
	r.cache[key] = cell
	return false
}

func (r *mapCache) GetCache(key interface{}) (interface{}, bool) {
	r.RLock()
	defer r.RUnlock()

	cell, ok := r.cache[key]
	if !ok {
		return nil, false
	}
	return cell.content, true
}

func (r *mapCache) DelCache(key interface{}) {
	r.Lock()
	defer r.Unlock()

	delete(r.cache, key)
}
