package exercises

import (
	"os"
	"sync"
	"testing"
)

func TestUnsafeCounter(t *testing.T) {
	if os.Getenv("RUN_UNSAFE_RACE") != "1" {
		t.Skip("设置 RUN_UNSAFE_RACE=1 后，用 -race 单独观察数据竞争")
	}

	counter := &UnsafeCounter{}
	runIncrements(counter)
	t.Log("理论计数: 100000，实际计数:", counter.Value())
}

func TestMutexCounter(t *testing.T) {
	testSafeCounter(t, &MutexCounter{})
}

func TestAtomicCounter(t *testing.T) {
	testSafeCounter(t, &AtomicCounter{})
}

func testSafeCounter(t *testing.T, counter Counter) {
	t.Helper()
	runIncrements(counter)
	if got := counter.Value(); got != 100000 {
		t.Fatalf("计数 = %d，期望 100000", got)
	}
}

func runIncrements(counter Counter) {
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Inc()
			}
		}()
	}

	wg.Wait()
}
