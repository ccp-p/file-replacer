package examples

import (
	"fmt"
	"sync"
	"time"
)

// VisualWaitGroupDemo 通过时间轴展示WaitGroup和监视协程的工作原理
func VisualWaitGroupDemo() {
	var wg sync.WaitGroup
	results := make(chan string, 5) // 使用缓冲通道避免工作协程阻塞

	printWithTime := func(format string, args ...interface{}) {
		fmt.Printf("%s: %s\n",
			time.Now().Format("15:04:05.000"),
			fmt.Sprintf(format, args...))
	}

	printWithTime("程序开始")

	// 创建3个工作协程，模拟不同完成时间
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		printWithTime("启动工作协程 %d，WaitGroup计数+1", i)

		go func(id int) {
			// 模拟不同的工作时间
			workTime := time.Duration(id) * time.Second
			printWithTime("工作协程 %d: 开始工作，预计耗时 %v", id, workTime)
			time.Sleep(workTime)

			printWithTime("工作协程 %d: 工作完成，发送结果", id)
			results <- fmt.Sprintf("结果 %d", id)

			printWithTime("工作协程 %d: 调用WaitGroup.Done()", id)
			wg.Done()
			printWithTime("工作协程 %d: WaitGroup.Done()调用完成", id)
		}(i)
	}

	// 启动监视协程
	printWithTime("启动监视协程")
	go func() {
		printWithTime("监视协程: 开始执行")
		printWithTime("监视协程: 调用WaitGroup.Wait()，将阻塞直到计数为0")
		wg.Wait()
		printWithTime("监视协程: WaitGroup.Wait()返回，所有工作协程已完成")
		printWithTime("监视协程: 关闭结果通道")
		close(results)
		printWithTime("监视协程: 通道已关闭，工作完成")
	}()

	// 接收结果
	printWithTime("主协程: 开始接收结果")
	for result := range results {
		printWithTime("主协程: 收到结果 '%s'", result)
	}

	printWithTime("主协程: 所有结果已接收完毕，程序结束")
}
