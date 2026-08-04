package exercises

import (
	"os"
	"sync"
	"testing"
)

func TestUnsafeCounter(t *testing.T) {
	if os.Getenv("RUN_UNSAFE_RACE") != "1" {
		t.Skip("设置 RUN_UNSAFE_RACE=1 后用 -race 观察故意保留的数据竞争")
	}

	counter := &UnsafeCounter{}
	runConcurrentIncrements(counter)
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
	runConcurrentIncrements(counter)

	const want int64 = 100000
	if got := counter.Value(); got != want {
		t.Fatalf("计数结果 = %d，期望 %d", got, want)
	}
}

func runConcurrentIncrements(counter Counter) {
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
