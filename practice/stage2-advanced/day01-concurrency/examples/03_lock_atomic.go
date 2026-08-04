package examples

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type DemoInventory struct {
	stock int
	mu    sync.RWMutex
}

func (i *DemoInventory) Stock() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.stock
}

func (i *DemoInventory) Deduct(quantity int) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if quantity > i.stock {
		return false
	}
	i.stock -= quantity
	return true
}

type DemoAtomicCounter struct {
	value int64
}

func (c *DemoAtomicCounter) Inc()       { atomic.AddInt64(&c.value, 1) }
func (c *DemoAtomicCounter) Dec()       { atomic.AddInt64(&c.value, -1) }
func (c *DemoAtomicCounter) Get() int64 { return atomic.LoadInt64(&c.value) }

func RunLockAtomic() {
	var wg sync.WaitGroup
	inventory := &DemoInventory{stock: 100}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inventory.Deduct(1)
		}()
	}
	wg.Wait()
	fmt.Println("最终库存:", inventory.Stock())

	var likes DemoAtomicCounter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			likes.Inc()
		}()
	}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			likes.Dec()
		}()
	}
	wg.Wait()
	fmt.Println("最终点赞数:", likes.Get())
}
