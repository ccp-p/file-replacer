package examples

import (
	"fmt"
	"sync"
	"time"
)

// ChannelClosingDemo 演示通道关闭的最佳实践
func ChannelClosingDemo() {
	var wg sync.WaitGroup
	results := make(chan string)

	// 启动 5 个生产者协程
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 模拟处理任务
			time.Sleep(time.Duration(id*100) * time.Millisecond)

			// 发送结果到通道
			results <- fmt.Sprintf("任务 %d 完成", id)
		}(i)
	}

	// 这就是你问的代码 - 启动一个额外协程来关闭通道
	go func() {
		// 等待所有生产者完成
		wg.Wait()
		// 安全地关闭通道
		close(results)
		fmt.Println("所有任务完成，通道已关闭")
	}()

	// 主协程从通道读取所有结果
	for result := range results {
		fmt.Println("收到结果:", result)
	}

	fmt.Println("所有结果已处理完毕")
}
