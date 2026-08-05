package exercises

import (
	"sync"
	"sync/atomic"
)

type Counter interface {
	Inc()
	Value() int64
}

// UnsafeCounter 故意不做并发保护，只用于观察数据竞争。
type UnsafeCounter struct {
	value int64
}

func (c *UnsafeCounter) Inc() {
	c.value++
}

func (c *UnsafeCounter) Value() int64 {
	return c.value
}

type MutexCounter struct {
	mu    sync.Mutex
	value int64
}

func (c *MutexCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *MutexCounter) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

type AtomicCounter struct {
	value atomic.Int64
}

func (c *AtomicCounter) Inc() {
	c.value.Add(1)
}

func (c *AtomicCounter) Value() int64 {
	return c.value.Load()
}
