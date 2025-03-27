package examples

import (
	"fmt"
	"sync"
	"time"
)

// WithoutWaitGroupDemo 展示不使用WaitGroup的通道模式
func WithoutWaitGroupDemo() {
	fmt.Println("=== 不使用 WaitGroup 的模式 ===")

	tasks := []int{1, 2, 3, 4, 5}
	resultChan := make(chan string)

	// 启动工作协程
	for _, task := range tasks {
		go func(t int) {
			// 模拟处理任务
			time.Sleep(time.Duration(t*100) * time.Millisecond)
			result := fmt.Sprintf("任务 %d 完成", t)
			resultChan <- result // 将结果发送到通道
		}(task)
	}

	// 收集所有结果
	for i := 0; i < len(tasks); i++ {
		result := <-resultChan
		fmt.Println("收到结果:", result)
	}

	fmt.Println("所有结果收集完毕")
	// 注意: 此时通道未关闭
}

// WithWaitGroupDemo 展示使用WaitGroup的通道模式
func WithWaitGroupDemo() {
	fmt.Println("=== 使用 WaitGroup 的模式 ===")

	var wg sync.WaitGroup
	tasks := []int{1, 2, 3, 4, 5}
	resultChan := make(chan string)

	// 启动工作协程
	for _, task := range tasks {
		wg.Add(1)
		go func(t int) {
			defer wg.Done()
			// 模拟处理任务
			time.Sleep(time.Duration(t*100) * time.Millisecond)
			result := fmt.Sprintf("任务 %d 完成", t)
			resultChan <- result
		}(task)
	}

	// 启动监视协程，等待所有工作完成后关闭通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集所有结果
	for result := range resultChan {
		fmt.Println("收到结果:", result)
	}

	fmt.Println("所有结果收集完毕")
}

// DemoRaceConditionRisk 演示不使用WaitGroup可能的竞态条件风险
func DemoRaceConditionRisk() {
	fmt.Println("=== 演示不使用WaitGroup的风险 ===")

	// 这个函数故意设计了一种可能失败的场景
	tasks := []int{1, 2, 3}
	resultChan := make(chan string)

	// 启动工作协程，注意有些任务耗时更长
	for i, task := range tasks {
		delay := task * 100
		if i == 2 {
			delay = 1000 // 最后一个任务明显更慢
		}

		go func(t int, d int) {
			time.Sleep(time.Duration(d) * time.Millisecond)
			result := fmt.Sprintf("任务 %d 完成 (延迟: %dms)", t, d)
			resultChan <- result
		}(task, delay)
	}

	// 在收集结果之前等待一段时间，模拟主程序执行其他逻辑
	time.Sleep(200 * time.Millisecond)

	// 在只收集部分结果后再进行一些处理
	for i := 0; i < 2; i++ { // 故意只收集2个结果（共3个任务）
		result := <-resultChan
		fmt.Println("收到结果:", result)
	}

	fmt.Println("部分结果收集完毕，程序继续执行其他逻辑...")

	// 模拟程序结束
	fmt.Println("程序即将结束...最后一个任务的结果无人接收！")
	time.Sleep(100 * time.Millisecond)
	// 在实际程序中，如果此时退出，最后一个协程会被挂起，因为无人接收它的结果
}
