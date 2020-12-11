/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/12/11
   Description :
-------------------------------------------------
*/

package zlock

import (
	"context"
	"time"
)

var defaultLocker = NewLocker()

// 一直等待获取锁
func Lock() {
	defaultLocker.Lock()
}

// 释放锁
func Unlock() {
	defaultLocker.Unlock()
}

// 尝试获取锁, 如果锁被别的使用者拿到则返回false
func TryLock() bool {
	return defaultLocker.TryLock()
}

// 在规定时间内尝试获取锁, 如果时间结束没有拿到锁则返回false
func TryLockWithTimeout(timeout time.Duration) bool {
	return defaultLocker.TryLockWithTimeout(timeout)
}

// 是否已经有使用者拿到了锁
func IsLocked() bool {
	return defaultLocker.IsLocked()
}

// 一直等待获取锁, 然后执行fn, 结束后释放锁
func LockDo(fn func()) {
	defaultLocker.LockDo(fn)
}

// 尝试获取锁, 如果锁被别的使用者拿到则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func TryLockDo(fn func()) bool {
	return defaultLocker.TryLockDo(fn)
}

// 在规定时间内尝试获取锁, 如果时间结束没有拿到锁则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func TryLockDoWithTimeout(timeout time.Duration, fn func()) bool {
	return defaultLocker.TryLockDoWithTimeout(timeout, fn)
}

// 尝试获取锁, 如果上下文结束则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func TryLockDoWithContext(ctx context.Context, fn func()) bool {
	return defaultLocker.TryLockDoWithContext(ctx, fn)
}
