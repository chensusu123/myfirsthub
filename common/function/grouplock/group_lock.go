package grouplock

import "sync"

// 锁管理
type GroupLock struct {
	inner map[uint64]*sync.RWMutex
	size  uint64
	sync.Mutex
}

func NewGroupLock(size uint64) *GroupLock {
	return &GroupLock{inner: make(map[uint64]*sync.RWMutex), size: size}
}

func (g *GroupLock) GetLock(user_id uint64) *sync.RWMutex {
	g.Lock()
	defer g.Unlock()
	sharding := user_id
	if lock, ok := g.inner[sharding]; ok {
		return lock
	}
	lock := &sync.RWMutex{}
	g.inner[sharding] = lock
	return lock
}

func (g *GroupLock) LockOp(user_id uint64) *sync.RWMutex {
	l := g.GetLock(user_id)
	l.Lock()
	return l
}

func (g *GroupLock) RLockOp(user_id uint64) *sync.RWMutex {
	l := g.GetLock(user_id)
	l.RLock()
	return l
}
