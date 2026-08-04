package exercises

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolCompletesAllJobs(t *testing.T) {
	jobs := make(chan Job, 20)
	for i := 1; i <= 20; i++ {
		jobs <- Job{ID: i, Value: i}
	}
	close(jobs)

	results, err := RunWorkerPool(context.Background(), 3, jobs)
	if err != nil {
		t.Fatal(err)
	}

	seen := make(map[int]bool, 20)
	for result := range results {
		if result.Err != nil {
			t.Fatalf("任务 %d 失败: %v", result.JobID, result.Err)
		}
		if seen[result.JobID] {
			t.Fatalf("任务 %d 被重复处理", result.JobID)
		}
		seen[result.JobID] = true
		if result.Value != result.JobID*result.JobID {
			t.Errorf("任务 %d 结果 = %d", result.JobID, result.Value)
		}
	}

	if len(seen) != 20 {
		t.Fatalf("完成任务数 = %d，期望 20", len(seen))
	}
}

func TestWorkerPoolCancellationClosesResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan Job)
	results, err := RunWorkerPool(ctx, 3, jobs)
	if err != nil {
		t.Fatal(err)
	}

	cancel()

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("取消后 results 应该关闭")
		}
	case <-time.After(time.Second):
		t.Fatal("取消后 worker 没有及时退出")
	}
}

func TestWorkerPoolLimitsConcurrency(t *testing.T) {
	const workerCount = 3
	jobs := make(chan Job, 6)
	for i := 1; i <= 6; i++ {
		jobs <- Job{ID: i, Value: i}
	}
	close(jobs)

	started := make(chan struct{}, workerCount)
	release := make(chan struct{})
	var active atomic.Int64
	var maximum atomic.Int64

	process := func(_ context.Context, job Job) (int, error) {
		current := active.Add(1)
		for {
			old := maximum.Load()
			if current <= old || maximum.CompareAndSwap(old, current) {
				break
			}
		}
		started <- struct{}{}
		<-release
		active.Add(-1)
		return job.Value * job.Value, nil
	}

	results, err := runWorkerPool(context.Background(), workerCount, jobs, process)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < workerCount; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("没有按预期启动 worker")
		}
	}
	close(release)

	for range results {
	}

	if got := maximum.Load(); got != workerCount {
		t.Fatalf("最大并发数 = %d，期望 %d", got, workerCount)
	}
}

func TestWorkerPoolRejectsInvalidWorkerCount(t *testing.T) {
	jobs := make(chan Job)
	if _, err := RunWorkerPool(context.Background(), 0, jobs); err == nil {
		t.Fatal("worker 数量为 0 时应该返回错误")
	}
}
