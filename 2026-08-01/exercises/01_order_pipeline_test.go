package exercises

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

type Order struct {
	ID     int
	Amount float64
	Status string
}

func produceOrders(jobs chan<- Order) {
	//循环创建10个订单 把订单发送到jobs 全部发送完毕关闭jobs
	defer close(jobs)
	for i := 1; i <= 10; i++ {
		order := Order{
			ID:     i,
			Amount: float64(rand.Intn(100)),
			Status: "pending",
		}
		jobs <- order
		fmt.Println("生成订单: ", order.ID)
	}
}

func worker(id int, jobs <-chan Order, results chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done()
	// 读取订单 改状态 发送到result
	for order := range jobs {
		fmt.Printf("worker %d 正在处理订单 %d\n", id, order.ID)
		order.Status = "completed"
		results <- order
	}
}
func TestOrderPipeline(t *testing.T) {
	jobs := make(chan Order)
	results := make(chan Order)

	var wg sync.WaitGroup

	go produceOrders(jobs)

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	completed := 0
	totalAmount := 0.0
	for order := range results {
		completed++
		totalAmount += order.Amount
		if order.Status != "completed" {
			t.Errorf("订单 %d 状态错误: %s", order.ID, order.Status)
		}
	}

	if completed != 10 {
		t.Fatalf("完成订单数 = %d，期望 10", completed)
	}
	t.Log("完成订单数:", completed, "订单总金额:", totalAmount)
}
