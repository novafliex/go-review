package exercises

import (
	"sync"
)

// FanOutSquares 启动固定数量的 worker，多个 worker 竞争处理同一个输入。
func FanOutSquares(input <-chan int, workerCount int) <-chan int {
	output := make(chan int)
	var workers sync.WaitGroup

	if workerCount <= 0 {
		close(output)
		return output
	}

	workers.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer workers.Done()
			for value := range input {
				output <- value * value
			}
		}()
	}

	go func() {
		workers.Wait()
		close(output)
	}()
	return output
}

// FanIn 把多个输入 channel 合并成一个输出 channel。
func FanIn(inputs ...<-chan int) <-chan int {
	output := make(chan int)
	var producers sync.WaitGroup
	producers.Add(len(inputs))

	for _, input := range inputs {
		go func(input <-chan int) {
			defer producers.Done()
			for value := range input {
				output <- value
			}
		}(input)
	}

	go func() {
		producers.Wait()
		close(output)
	}()
	return output
}

func Generate(values ...int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for _, value := range values {
			output <- value
		}
	}()
	return output
}

func Positive(input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for value := range input {
			if value <= 0 {
				continue
			}
			output <- value
		}
	}()
	return output
}

func Square(input <-chan int) <-chan int {
	output := make(chan int)
	go func() {
		defer close(output)
		for value := range input {
			output <- value * value
		}
	}()
	return output
}
