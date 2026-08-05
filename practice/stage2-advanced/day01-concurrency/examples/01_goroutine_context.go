package examples

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"
	"time"
)

func fetchData(url string) {
	delay := time.Duration(rand.IntN(500)) * time.Millisecond
	time.Sleep(delay)
	fmt.Printf("完成 %-24s 耗时 %v\n", url, delay)
}

func contextWorker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker", id, "停止:", ctx.Err())
			return
		case <-ticker.C:
			fmt.Println("worker", id, "working")
		}
	}
}

func limitedWorkerPool(tasks []string, max int) {
	semaphore := make(chan struct{}, max)
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(task string) {
			defer wg.Done()
			defer func() { <-semaphore }()
			fetchData(task)
		}(task)
	}
	wg.Wait()
}

func RunGoroutineContext() {
	urls := []string{"google.com", "baidu.com", "yahoo.com", "bing.com", "douyin.com"}

	start := time.Now()
	for _, url := range urls {
		fetchData(url)
	}
	fmt.Println("串行总耗时:", time.Since(start))

	start = time.Now()
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fetchData(url)
		}(url)
	}
	wg.Wait()
	fmt.Println("并发总耗时:", time.Since(start))

	fmt.Println("\n最多同时执行 3 个任务:")
	limitedWorkerPool(urls, 3)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go contextWorker(ctx, i, &wg)
	}
	wg.Wait()

	fmt.Printf("Goroutine: %d，CPU: %d\n", runtime.NumGoroutine(), runtime.NumCPU())
}
