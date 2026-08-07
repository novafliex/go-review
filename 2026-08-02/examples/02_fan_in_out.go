package main

import (
	"fmt"
	"sync"
	"time"
)

// 扇出：一个 channel 的任务由多个 goroutine 竞争处理。
func fanOutProcessor(input <-chan Order, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range input {
		switch workerID {
		case 1:
			fmt.Printf("处理器%d计算折扣: 订单%d\n", workerID, order.ID)
		case 2:
			fmt.Printf("处理器%d发送通知: 订单%d\n", workerID, order.ID)
		case 3:
			fmt.Printf("处理器%d记录日志: 订单%d\n", workerID, order.ID)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func runFanOut() {
	fmt.Println("\n=== 2. 扇出：一产多消 ===")
	input := make(chan Order, 6)
	var workers sync.WaitGroup

	workers.Add(3)
	for workerID := 1; workerID <= 3; workerID++ {
		go fanOutProcessor(input, workerID, &workers)
	}

	for i := 1; i <= 6; i++ {
		input <- Order{ID: i, Amount: float64(i * 100), Status: "pending"}
	}
	close(input)
	workers.Wait()
}

// 扇入：多个生产者把结果汇总到同一个 channel。
func fanInProducer(producerID int, output chan<- Order, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= 3; i++ {
		output <- Order{
			ID:     producerID*100 + i,
			UserID: producerID,
			Amount: float64(producerID*100 + i),
			Status: "new",
		}
	}
}

func runFanIn() {
	fmt.Println("\n=== 3. 扇入：多产一消 ===")
	output := make(chan Order, 9)
	var producers sync.WaitGroup

	producers.Add(3)
	for producerID := 1; producerID <= 3; producerID++ {
		go fanInProducer(producerID, output, &producers)
	}

	go func() {
		producers.Wait()
		close(output)
	}()

	for order := range output {
		fmt.Printf("汇总订单: %d\n", order.ID)
	}
}
