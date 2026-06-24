package main

import "fmt"

func main() {
	// 匿名函数：定义后立即执行
	func() {
		fmt.Println("hi")
	}()

	// 匿名函数：先保存，再调用
	Nihat := func(a, b int) int {
		return a * b
	}
	fmt.Println(Nihat(2, 3))

	sayHello()
	sayName("Alice")

	fmt.Println(getNumber())
	fmt.Println(add(1, 2))

	sumResult, subResult := calculate(20, 10)
	fmt.Println(sumResult, subResult)

	fmt.Println(sum(1, 2))
	fmt.Println(sum(1, 2, 3, 4))

	// 函数作为参数
	fmt.Println(calculate1(1, 2, add))

	a := getAddFunction()
	fmt.Println(a(10, 20))

	counter := createCounter()

	fmt.Println(counter())
	fmt.Println(counter())
	fmt.Println(counter())
}

// 无参数、无返回值
func sayHello() {
	fmt.Println("hello")
}

// 有参数、无返回值
func sayName(name string) {
	fmt.Println(name)
}

// 无参数、有返回值
func getNumber() int {
	return 10
}

// 有参数、有返回值
func add(a, b int) int {
	return a + b
}

// 多个返回值
func calculate(a, b int) (int, int) {
	return a + b, a - b
}

// 可变参数函数
func sum(nums ...int) int {
	total := 0

	for _, num := range nums {
		total += num
	}

	return total
}

// 函数作为参数
func calculate1(a, b int, c func(int, int) int) int {
	return c(a, b)
}

// 函数作为返回值
func getAddFunction() func(int, int) int {
	return func(a, b int) int {
		return a + b
	}
}

// 闭包
func createCounter() func() int {
	count := 0

	return func() int {
		count++
		return count
	}
}
