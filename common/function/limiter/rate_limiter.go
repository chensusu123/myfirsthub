/*
* @Author: majian
* @Date: 2022-03-10 11:34
 */
package limiter

import (
	"sync"
	"time"

	"gitlab.ifreetalk.com/plate/freetk/fkcore/fkconfig/param"
	"gitlab.ifreetalk.com/plate/freetk/fkcore/fklog"
	"go.uber.org/zap"
)

var qps int32           //频率间隔
var ClearInterval int32 //清理过期时间间隔
var CacheTimeout int32  //缓存过期时间
func init() {
	param.Int32P(&qps, "limter:qps:interval", 100, "最小请求间隔")
	param.Int32P(&ClearInterval, "limiter:clear:interval", 24*3600, "缓存清理间隔")
	param.Int32P(&CacheTimeout, "limiter:cache:timeout", 60, "缓存过期时间")
}

// 限速信息
type LimitInfo struct {
	LastTime int64 `json:"last_time"` //最近一次请求时间  单位 ms
}

// 简单限速器
type Limiter struct {
	CacheList     map[string]*LimitInfo // key:航线ID val:LimitInfo
	LastClearTime int64                 // 上次清理过期时间 单位 秒
	sync.RWMutex
	Name string // 限速器名字
}

func NewLimiter(name string) *Limiter {
	cache := new(Limiter)
	cache.CacheList = make(map[string]*LimitInfo)
	cache.LastClearTime = time.Now().Unix()
	cache.Name = name
	return cache
}

// 是否频率限制
func (m *Limiter) IsRateLimit(key string) bool {
	now := time.Now().UnixNano() / 1000000
	m.RLock()
	defer m.RUnlock()
	if limitInfo, ok := m.CacheList[key]; ok {
		if now-limitInfo.LastTime < int64(qps) {
			return true
		}
	}
	return false
}

// 更新频次
func (m *Limiter) UpdateTime(logger fklog.FKLogI, key string) {
	now := time.Now().UnixNano() / 1000000
	m.Lock()
	if limitInfo, ok := m.CacheList[key]; ok {
		limitInfo.LastTime = now
	} else {
		limitInfo := new(LimitInfo)
		limitInfo.LastTime = now
		m.CacheList[key] = limitInfo
	}
	m.Unlock()
	//被动触发 等有更新的时候检查
	checkExpireCache(logger, m)
}

// 清理过期
func checkExpireCache(logger fklog.FKLogI, limter *Limiter) {
	now := time.Now().Unix()
	span := now - limter.LastClearTime
	if span > int64(ClearInterval) {
		go clearCache(limter, logger, now)
	}
}

func clearCache(limter *Limiter, logger fklog.FKLogI, now int64) {
	limter.Lock()
	defer limter.Unlock()

	for key, info := range limter.CacheList {
		if now-info.LastTime/1000 > int64(CacheTimeout) {
			delete(limter.CacheList, key)
			logger.InfoWF("clear expire lineId cache",
				zap.String("key", key),
				zap.Int32("cacheTimeout", CacheTimeout))
		}
	}
	limter.LastClearTime = now
}
