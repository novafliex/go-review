package exercises

import (
	"context"
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

	seen := make(map[int]bool)
	for result := range results {
		if seen[result.JobID] {
			t.Fatalf("任务 %d 被重复处理", result.JobID)
		}
		seen[result.JobID] = true
		if result.Value != result.JobID*result.JobID {
			t.Errorf("任务 %d 结果错误: %d", result.JobID, result.Value)
		}
	}
	if len(seen) != 20 {
		t.Fatalf("完成任务数 = %d，期望 20", len(seen))
	}
}

func TestWorkerPoolStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	results, err := RunWorkerPool(ctx, 3, make(chan Job))
	if err != nil {
		t.Fatal(err)
	}
	cancel()

	select {
	case _, ok := <-results:
		if ok {
			t.Fatal("取消后 results 应关闭")
		}
	case <-time.After(time.Second):
		t.Fatal("worker 没有及时退出")
	}
}

func TestWorkerPoolRejectsZeroWorkers(t *testing.T) {
	if _, err := RunWorkerPool(context.Background(), 0, make(chan Job)); err == nil {
		t.Fatal("worker 数量为 0 时应返回错误")
	}
}
