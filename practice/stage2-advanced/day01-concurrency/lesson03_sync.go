package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Inventory struct {
	stock int
	mu    sync.RWMutex
}

func (i *Inventory) Stock() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.stock
}

func (i *Inventory) Deduct(quantity int) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	if quantity > i.stock {
		return false
	}
	i.stock -= quantity
	return true
}

type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Inc()       { atomic.AddInt64(&c.value, 1) }
func (c *AtomicCounter) Dec()       { atomic.AddInt64(&c.value, -1) }
func (c *AtomicCounter) Get() int64 { return atomic.LoadInt64(&c.value) }

func runLesson3() {
	var wg sync.WaitGroup
	inventory := &Inventory{stock: 100}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inventory.Deduct(1)
		}()
	}
	wg.Wait()
	fmt.Println("最终库存:", inventory.Stock())

	var likes AtomicCounter
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
