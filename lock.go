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

type Locker struct {
	ch chan struct{}
}

func NewLocker() *Locker {
	return &Locker{make(chan struct{}, 1)}
}

// 一直等待获取锁
func (m *Locker) Lock() {
	m.ch <- struct{}{}
}

// 释放锁
func (m *Locker) Unlock() {
	<-m.ch
}

// 尝试获取锁, 如果锁被别的使用者拿到则返回false
func (m *Locker) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// 在规定时间内尝试获取锁, 如果时间结束没有拿到锁则返回false
func (m *Locker) TryLockWithTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case m.ch <- struct{}{}:
		timer.Stop()
		return true
	case <-timer.C:
		return false
	}
}

// 尝试获取锁, 如果上下文结束则返回false
func (m *Locker) TryLockWithContext(ctx context.Context) bool {
	select {
	case m.ch <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

// 是否已经有使用者拿到了锁
func (m *Locker) IsLocked() bool {
	return len(m.ch) > 0
}

// 一直等待获取锁, 然后执行fn, 结束后释放锁
func (m *Locker) LockDo(fn func()) {
	m.Lock()
	defer m.Unlock() // 防止panic后未解锁
	fn()
}

// 尝试获取锁, 如果锁被别的使用者拿到则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func (m *Locker) TryLockDo(fn func()) bool {
	if m.TryLock() {
		defer m.Unlock() // 防止panic后未解锁
		fn()
		return true
	}
	return false
}

// 在规定时间内尝试获取锁, 如果时间结束没有拿到锁则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func (m *Locker) TryLockDoWithTimeout(timeout time.Duration, fn func()) bool {
	if m.TryLockWithTimeout(timeout) {
		defer m.Unlock() // 防止panic后未解锁
		fn()
		return true
	}
	return false
}

// 尝试获取锁, 如果上下文结束则返回false, 拿到锁后执行fn, 结束后释放锁并返回true
func (m *Locker) TryLockDoWithContext(ctx context.Context, fn func()) bool {
	if m.TryLockWithContext(ctx) {
		defer m.Unlock() // 防止panic后未解锁
		fn()
		return true
	}
	return false
}
