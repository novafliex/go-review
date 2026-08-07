package main

import (
	"fmt"
	"time"
)

func validationStage(input <-chan Order) <-chan Order {
	output := make(chan Order, 10)
	go func() {
		defer close(output)
		for order := range input {
			time.Sleep(50 * time.Millisecond)
			if order.Amount <= 0 {
				fmt.Printf("验证失败: 订单%d\n", order.ID)
				continue
			}
			order.Status = "validated"
			output <- order
		}
	}()
	return output
}

func paymentStage(input <-chan Order) <-chan Order {
	output := make(chan Order, 10)
	go func() {
		defer close(output)
		for order := range input {
			time.Sleep(100 * time.Millisecond)
			order.Status = "paid"
			output <- order
		}
	}()
	return output
}

func shippingStage(input <-chan Order) <-chan Order {
	output := make(chan Order, 10)
	go func() {
		defer close(output)
		for order := range input {
			time.Sleep(150 * time.Millisecond)
			order.Status = "shipped"
			output <- order
		}
	}()
	return output
}

func runPipeline() {
	fmt.Println("\n=== 4. Pipeline：验证 → 支付 → 发货 ===")
	input := make(chan Order, 3)
	results := shippingStage(paymentStage(validationStage(input)))

	go func() {
		defer close(input)
		for i := 1; i <= 3; i++ {
			input <- Order{ID: i, Amount: float64(i * 50), Status: "new"}
		}
	}()

	for order := range results {
		fmt.Printf("完成订单: %d，状态: %s\n", order.ID, order.Status)
	}
}
