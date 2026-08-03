package main

import "fmt"

// 今天只做这一个综合练习：并发订单处理。
// 不需要再新建文件夹，直接在本文件中完成。
func runPractice() {
	fmt.Println(`今日练习：并发处理 10 个订单

要求：
1. jobs 使用 chan Order，放入 10 个订单后关闭
2. 启动 3 个 worker，同时处理订单
3. 每个订单处理完成后写入 results
4. 等全部 worker 退出后，只关闭一次 results
5. 主 goroutine 统计 completed 订单数和总金额
6. 使用 go run -race . practice 检查数据竞争

请直接修改 practice.go，完成后再运行。`)

	// TODO: 从这里开始独立实现。
}
