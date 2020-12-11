/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/12/11
   Description :
-------------------------------------------------
*/

package zlock

import (
	"hash/fnv"
	"sync"
	"time"
)

var defaultKeyLocker = func() *keyLocker {
	lockers := make([]map[uint64]*Locker, 8)
	mxs := make([]*sync.RWMutex, 8)
	for i := 0; i < 8; i++ {
		lockers[i] = make(map[uint64]*Locker)
		mxs[i] = new(sync.RWMutex)
	}
	return &keyLocker{
		lockers: lockers,
		mxs:     mxs,
	}
}()

type keyLocker struct {
	lockers []map[uint64]*Locker
	mxs     []*sync.RWMutex
}

// 获取locker
func GetLocker(key string) *Locker {
	c := fnv.New64a()
	_, _ = c.Write([]byte(key))
	id := c.Sum64()
	shard := 7 & c.Sum64()

	lockers, mx := defaultKeyLocker.lockers[shard], defaultKeyLocker.mxs[shard]

	mx.RLock()
	l, ok := lockers[id]
	mx.RUnlock()

	if ok {
		return l
	}

	mx.Lock()
	l, ok = lockers[id]
	if ok {
		mx.Unlock()
		return l
	}

	l = NewLocker()
	lockers[id] = l
	mx.Unlock()
	return l
}

// 一直等待获取锁
func LockKey(key string) {
	GetLocker(key).Lock()
}

// 释放锁
func UnlockKey(key string) {
	GetLocker(key).Unlock()
}

// 尝试获取锁, 如果锁被别的使用者拿到则返回false
func TryLockKey(key string) bool {
	return GetLocker(key).TryLock()
}

// 在规定时间内尝试获取锁, 如果时间结束没有拿到锁则返回false
func TryLockKeyWithTimeout(key string, timeout time.Duration) bool {
	return GetLocker(key).TryLockWithTimeout(timeout)
}

// 是否已经有使用者拿到了锁
func IsLockedKey(key string) bool {
	return GetLocker(key).IsLocked()
}
