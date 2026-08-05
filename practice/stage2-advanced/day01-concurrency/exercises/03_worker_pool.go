package exercises

import (
	"context"
	"errors"
	"sync"
)

type Job struct {
	ID    int
	Value int
}

type Result struct {
	JobID int
	Value int
}

// RunWorkerPool 使用固定数量的 worker 计算任务值的平方。
func RunWorkerPool(ctx context.Context, workerCount int, jobs <-chan Job) (<-chan Result, error) {
	if workerCount <= 0 {
		return nil, errors.New("worker 数量必须大于 0")
	}

	results := make(chan Result)
	var workers sync.WaitGroup
	workers.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer workers.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}

					result := Result{JobID: job.ID, Value: job.Value * job.Value}
					select {
					case results <- result:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	go func() {
		workers.Wait()
		close(results)
	}()

	return results, nil
}
