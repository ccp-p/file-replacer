package examples

import (
	"fmt"
	"runtime"
	"time"
)

// GoroutineLeakDemo 演示不当关闭通道导致的协程泄漏
func GoroutineLeakDemo() {
	fmt.Println("==== 协程泄漏演示 ====")

	// 打印初始协程数量
	fmt.Printf("初始协程数量: %d\n", runtime.NumGoroutine())

	// 模拟一个可能导致泄漏的函数
	demonstrateLeakRisk()

	// 等待一段时间，观察是否有协程未退出
	time.Sleep(time.Second)
	fmt.Printf("程序结束前协程数量: %d (如果数量没有恢复到初始值，说明有协程泄漏)\n", runtime.NumGoroutine())
}

func demonstrateLeakRisk() {
	resultChan := make(chan string) // 无缓冲通道

	// 创建10个工作协程
	for i := 1; i <= 10; i++ {
		go func(id int) {
			// 模拟处理时间不一致的任务
			time.Sleep(time.Duration(id*100) * time.Millisecond)

			// 尝试发送结果，这里可能会阻塞
			fmt.Printf("协程 %d: 准备发送结果 (阻塞前)\n", id)

			// 这里发生阻塞，直到有人从通道接收
			resultChan <- fmt.Sprintf("结果 %d", id)

			// 只有当结果被接收后，下面这行才会执行
			fmt.Printf("协程 %d: 已发送结果 (阻塞后)\n", id)
		}(i)
	}

	// 只接收部分结果
	fmt.Println("\n=== 只接收5个结果(共10个) ===")
	for i := 0; i < 5; i++ {
		result := <-resultChan
		fmt.Printf("主协程: 接收到 %s\n", result)
	}

	fmt.Println("\n=== 函数退出，但仍有5个协程被阻塞在发送操作上! ===")
	fmt.Printf("当前协程数量: %d (其中至少5个将永远阻塞)\n", runtime.NumGoroutine())
}

// BetterDemoWithTimeout 通过超时机制展示无缓冲通道阻塞的更清晰例子
func BetterDemoWithTimeout() {
	fmt.Println("==== 无缓冲通道阻塞演示 ====")

	ch := make(chan string) // 无缓冲通道

	// 启动发送者协程
	go func() {
		fmt.Println("发送者: 准备发送数据 (阻塞前)")
		fmt.Println("发送者: 尝试发送 -> 现在阻塞了，等待接收者...")
		ch <- "hello"
		fmt.Println("发送者: 数据已被接收! (阻塞结束)")
	}()

	// 等待一会儿，让我们观察发送者的状态
	fmt.Println("主协程: 故意等待2秒，看发送者是否阻塞")
	time.Sleep(2 * time.Second)

	// 接收数据
	fmt.Println("主协程: 现在准备接收数据")
	msg := <-ch
	fmt.Printf("主协程: 接收到: '%s'\n", msg)

	// 等待一会儿，让发送者打印完成
	time.Sleep(10 * time.Millisecond)
}
