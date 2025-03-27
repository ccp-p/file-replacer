package docs

import (
	"fmt"
	"sync"
	"time"
)

// WaitGroupExplanation 详细解释WaitGroup和监视协程的工作原理
func WaitGroupExplanation() {
	var wg sync.WaitGroup
	results := make(chan string)

	fmt.Println("1. 主协程: 开始执行")

	// 创建5个工作协程
	for i := 1; i <= 5; i++ {
		wg.Add(1) // 递增计数器
		fmt.Printf("2. 主协程: 为任务 %d 增加计数器，当前计数: %d\n", i, i)

		go func(id int) {
			fmt.Printf("3. 工作协程 %d: 开始工作\n", id)
			// 模拟工作耗时
			time.Sleep(time.Duration(id*300) * time.Millisecond)

			// 工作完成，发送结果
			fmt.Printf("4. 工作协程 %d: 完成工作，发送结果\n", id)
			results <- fmt.Sprintf("任务 %d 结果", id)

			fmt.Printf("5. 工作协程 %d: 调用 wg.Done() 减少计数\n", id)
			wg.Done() // 递减计数器
		}(i)
	}

	// 创建监视协程
	fmt.Println("6. 主协程: 创建监视协程")
	go func() {
		fmt.Println("7. 监视协程: 开始执行")
		fmt.Println("8. 监视协程: 调用 wg.Wait() 等待所有工作完成")
		wg.Wait() // 在此阻塞，直到计数器变为0
		fmt.Println("9. 监视协程: wg.Wait() 返回，表示所有工作已完成")
		fmt.Println("10. 监视协程: 安全关闭结果通道")
		close(results)
	}()

	// 从通道接收结果
	fmt.Println("11. 主协程: 开始从结果通道接收数据")
	count := 1
	for result := range results {
		fmt.Printf("12. 主协程: 收到第 %d 个结果: %s\n", count, result)
		count++
	}

	fmt.Println("13. 主协程: 通道已关闭，for-range循环结束")
}
