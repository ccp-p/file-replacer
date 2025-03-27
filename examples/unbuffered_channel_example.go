package examples

import (
	"fmt"
	"sync"
)

// UnbufferedChannelDemo 演示使用无缓冲通道实现并发处理
func UnbufferedChannelDemo() {
	files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}
	results := make(chan string) // 无缓冲通道用于接收结果

	var wg sync.WaitGroup

	// 启动处理协程
	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			// 模拟文件处理
			result := fmt.Sprintf("处理文件 %s 完成", f)
			// 将结果发送到通道
			results <- result // 这里会阻塞直到有人接收
		}(file)
	}

	// 启动单独的协程来关闭结果通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 从通道接收并处理所有结果
	for result := range results {
		fmt.Println(result)
	}
}
