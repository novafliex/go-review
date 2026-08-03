package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	UserID string
	Amount float64
	Status string
}

func produceOrders(out chan<- Order, number int) {
	defer close(out)
	for i := 1; i <= number; i++ {
		order := Order{
			ID: i, UserID: fmt.Sprintf("user-%d", rand.Intn(100)),
			Amount: rand.Float64() * 1000, Status: "pending",
		}
		out <- order
		fmt.Println("生成订单:", order.ID)
	}
}

func processOrders(id int, in <-chan Order, out chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range in {
		time.Sleep(100 * time.Millisecond)
		order.Status = "completed"
		fmt.Printf("worker %d 处理订单 %d\n", id, order.ID)
		out <- order
	}
}

func selectTimeoutDemo() {
	ch := make(chan string)
	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- "请求成功"
	}()

	select {
	case message := <-ch:
		fmt.Println(message)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("请求超时")
	}
}

func runLesson2() {
	orders := make(chan Order, 5)
	results := make(chan Order, 5)
	var workers sync.WaitGroup

	go produceOrders(orders, 10)
	for i := 1; i <= 3; i++ {
		workers.Add(1)
		go processOrders(i, orders, results, &workers)
	}

	// 只有协调者关闭 results，worker 不能各自关闭它。
	go func() {
		workers.Wait()
		close(results)
	}()

	for order := range results {
		fmt.Printf("订单 %d 状态: %s，金额: %.2f\n", order.ID, order.Status, order.Amount)
	}

	selectTimeoutDemo()
}
