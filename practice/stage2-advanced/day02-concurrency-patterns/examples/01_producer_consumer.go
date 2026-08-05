package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type Order struct {
	ID        int
	UserID    int
	Amount    float64
	Status    string
	CreatedAt time.Time
}

// 生产
func orderProducer(orderChan chan<- Order, number int) {
	defer close(orderChan)
	for i := 1; i <= number; i++ {
		order := Order{
			ID:        i,
			UserID:    rand.IntN(1000) + 1,
			Amount:    rand.Float64() * 1000,
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		fmt.Println("生成订单:", order.ID)
		orderChan <- order
		time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond)
	}
}

// 消费
func orderConsumer(orderChan <-chan Order, wg *sync.WaitGroup) {
	//模拟订单处理
	defer wg.Done()
	for order := range orderChan {
		time.Sleep(200 * time.Millisecond)
		order.Status = "completed"
		fmt.Println("处理订单:", order.ID, "状态:", order.Status)
	}
}

func runProducerConsumer() {
	fmt.Println("\n=== 1. 生产者-消费者 ===")
	orderChan := make(chan Order, 5)
	var consumers sync.WaitGroup

	consumers.Add(2)
	go orderConsumer(orderChan, &consumers)
	go orderConsumer(orderChan, &consumers)
	go orderProducer(orderChan, 6)

	consumers.Wait()
}

func main() {
	runProducerConsumer()
	runFanOut()
	runFanIn()
	runPipeline()
}
