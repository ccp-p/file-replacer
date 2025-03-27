package examples

import (
	"fmt"
	"time"
)

// SimpleChannelDemo 简化版的无缓冲通道阻塞演示，无需清屏
func SimpleChannelDemo() {
	fmt.Println("=== 简化版无缓冲通道阻塞演示 ===")
	fmt.Println("(按时间顺序观察事件)")
	fmt.Println()

	ch := make(chan string) // 无缓冲通道

	// 将当前时间添加到消息前面
	logWithTime := func(msg string) {
		fmt.Printf("[%s] %s\n", time.Now().Format("15:04:05.000"), msg)
	}

	// 启动发送者协程
	go func() {
		time.Sleep(500 * time.Millisecond) // 小延迟，确保日志按顺序
		logWithTime("发送者: 我已启动，准备发送数据到通道")
		time.Sleep(1 * time.Second)

		logWithTime("发送者: 现在尝试发送数据 -> ch <- \"hello\"")
		logWithTime("发送者: ⚠️ 现在阻塞了! 直到有人接收我的数据...")

		// 这里会阻塞直到有人接收
		ch <- "hello"

		logWithTime("发送者: ✓ 数据已被接收! 我继续执行")
		time.Sleep(500 * time.Millisecond)
		logWithTime("发送者: 任务完成，退出")
	}()

	// 主协程作为接收者
	logWithTime("主协程: 我将在3秒后接收数据，观察发送者状态")
	time.Sleep(3 * time.Second)

	logWithTime("主协程: 现在准备接收数据 -> data := <-ch")
	data := <-ch
	logWithTime(fmt.Sprintf("主协程: ✓ 成功接收数据: '%s'", data))

	// 等待发送者协程完成
	time.Sleep(1 * time.Second)
	logWithTime("主协程: 演示完成")
}
