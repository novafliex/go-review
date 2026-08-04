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
	Err   error
}

type JobProcessor func(context.Context, Job) (int, error)

// RunWorkerPool 启动固定数量的 worker，并把每个 Job.Value 的平方写入结果。
func RunWorkerPool(ctx context.Context, workerCount int, jobs <-chan Job) (<-chan Result, error) {
	return runWorkerPool(ctx, workerCount, jobs, func(_ context.Context, job Job) (int, error) {
		return job.Value * job.Value, nil
	})
}

func runWorkerPool(
	ctx context.Context,
	workerCount int,
	jobs <-chan Job,
	process JobProcessor,
) (<-chan Result, error) {
	if workerCount <= 0 {
		return nil, errors.New("worker 数量必须大于 0")
	}
	if process == nil {
		return nil, errors.New("任务处理函数不能为空")
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

					value, err := process(ctx, job)
					result := Result{JobID: job.ID, Value: value, Err: err}

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
