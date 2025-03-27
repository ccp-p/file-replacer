package examples

import (
	"fmt"
	"strings"
	"time"
)

// ChannelBlockingVisual 可视化展示无缓冲通道的阻塞行为
func ChannelBlockingVisual() {
	ch := make(chan string)

	// 格式化时间戳
	timeNow := func() string {
		return time.Now().Format("15:04:05.000")
	}

	// 绘制分隔线
	drawLine := func() {
		fmt.Println(strings.Repeat("=", 60))
	}

	// 绘制状态
	drawStatus := func(step int, desc string, senderStatus, channelStatus, receiverStatus string) {
		drawLine()
		fmt.Printf("步骤 %d: %s [时间: %s]\n\n", step, desc, timeNow())

		// 创建状态框
		fmt.Println("┌────────────发送者────────────┬──────通道──────┬────────────接收者────────────┐")
		fmt.Printf("│ %-28s│ %-14s│ %-32s│\n", senderStatus, channelStatus, receiverStatus)
		fmt.Println("└────────────────────────────┴────────────────┴────────────────────────────────┘")

		// 添加可视化箭头
		var arrow string
		switch channelStatus {
		case "发送中...":
			arrow = "  发送者 =====> [ 通道 ] -----> 接收者  "
		case "接收中...":
			arrow = "  发送者 -----> [ 通道 ] <===== 接收者  "
		case "数据传输中!":
			arrow = "  发送者 =====> [ 通道 ] <===== 接收者  "
		default:
			arrow = "  发送者 -----> [ 通道 ] -----> 接收者  "
		}
		fmt.Println("\n" + arrow + "\n")

		time.Sleep(1500 * time.Millisecond) // 让每一步清晰可见
	}

	// 打印标题
	fmt.Println("\n\n")
	fmt.Println(strings.Repeat("*", 80))
	fmt.Println("                       无缓冲通道阻塞行为可视化演示                      ")
	fmt.Println(strings.Repeat("*", 80))
	fmt.Println("\n这个演示将展示无缓冲通道如何导致发送和接收操作阻塞\n")
	time.Sleep(2 * time.Second)

	// 启动发送者协程
	go func() {
		drawStatus(1, "发送者协程启动", "准备发送数据", "空闲", "未启动")
		time.Sleep(1 * time.Second)

		drawStatus(2, "发送者尝试发送数据", "执行: ch <- \"数据\"", "空闲", "未启动")

		drawStatus(3, "发送被阻塞!", "⚠️ 阻塞在发送操作", "发送中...", "未启动")

		// 这里会阻塞，直到有人接收
		ch <- "数据"

		drawStatus(6, "发送成功完成", "✓ 发送完成，继续执行", "数据传输中!", "✓ 已接收数据")
		time.Sleep(1 * time.Second)

		drawStatus(7, "发送者完成任务", "任务完成，退出协程", "空闲", "处理数据中")
	}()

	// 主协程模拟接收者，但故意延迟
	time.Sleep(4 * time.Second) // 给发送者足够时间展示阻塞状态

	drawStatus(4, "接收者准备就绪", "⚠️ 仍在阻塞", "发送中...", "启动，准备接收")
	time.Sleep(1 * time.Second)

	drawStatus(5, "接收者执行接收", "⚠️ 仍在阻塞", "接收中...", "执行: data := <-ch")

	// 接收数据
	data := <-ch

	drawStatus(6, "数据成功传递", "✓ 发送完成，继续执行", "数据传输中!", fmt.Sprintf("✓ 已接收数据: \"%s\"", data))

	drawStatus(8, "演示完成", "已退出", "空闲", "演示结束")
}

// EnhancedChannelBlockingDemo 增强版的通道阻塞演示，显示多个阶段
func EnhancedChannelBlockingDemo() {
	fmt.Println("\n=== 增强版无缓冲通道阻塞演示 ===\n")

	// 创建一个简单的进度条
	progressBar := func(percent int) string {
		width := 30
		complete := width * percent / 100
		bar := "["
		for i := 0; i < width; i++ {
			if i < complete {
				bar += "="
			} else {
				bar += " "
			}
		}
		return bar + fmt.Sprintf("] %d%%", percent)
	}

	// 记录事件
	events := []struct {
		time     string
		sender   string
		channel  string
		receiver string
	}{
		{"00:00.000", "启动", "空闲", "启动"},
		{"00:01.000", "准备发送", "空闲", "等待中"},
		{"00:02.000", "⚠️ 阻塞在ch<-", "等待接收者", "未准备接收"},
		{"00:03.000", "⚠️ 阻塞在ch<-", "等待接收者", "未准备接收"},
		{"00:04.000", "⚠️ 阻塞在ch<-", "等待接收者", "准备接收"},
		{"00:05.000", "✓ 发送完成", "数据传递", "✓ 接收完成"},
		{"00:06.000", "继续执行", "空闲", "处理数据"},
		{"00:07.000", "退出", "空闲", "完成"},
	}

	// 显示整个过程的时间轴
	fmt.Println("       时间轴       |       发送者      |      通道      |      接收者       ")
	fmt.Println("---------------------|-------------------|----------------|------------------")

	for _, e := range events {
		fmt.Printf(" %-18s | %-17s | %-14s | %-18s\n",
			e.time, e.sender, e.channel, e.receiver)
		time.Sleep(1 * time.Second)
	}

	// 显示阻塞期间的模拟进度
	fmt.Println("\n\n=== 阻塞期间的详细视图 ===\n")
	fmt.Println("发送者执行: ch <- \"数据\"")

	// 模拟阻塞等待
	for i := 0; i <= 100; i += 10 {
		if i < 80 {
			fmt.Printf("\r发送者: 等待接收者 %s", progressBar(i))
		} else {
			fmt.Printf("\r发送者: 正在发送   %s", progressBar(i))
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\n\n发送者: ✓ 发送完成!")
	fmt.Println("接收者: ✓ 接收完成!")

	fmt.Println("\n演示结束")
}

// SafeChannelDemo 展示如何安全地处理可能阻塞的通道操作
func SafeChannelDemo() {
	resultChan := make(chan string)
	done := make(chan struct{}) // 使用空结构体通道作为信号

	// 启动5个工作协程
	for i := 1; i <= 5; i++ {
		go func(id int) {
			time.Sleep(time.Duration(id*300) * time.Millisecond)

			// 使用select避免永久阻塞
			select {
			case resultChan <- fmt.Sprintf("结果 %d", id):
				fmt.Printf("协程 %d: 成功发送结果\n", id)
			case <-done:
				fmt.Printf("协程 %d: 收到取消信号，放弃发送\n", id)
			case <-time.After(2 * time.Second): // 添加超时
				fmt.Printf("协程 %d: 发送超时，放弃发送\n", id)
			}
		}(i)
	}

	// 只接收3个结果后，发送取消信号
	fmt.Println("接收3个结果:")
	for i := 0; i < 3; i++ {
		fmt.Printf("收到: %s\n", <-resultChan)
	}

	// 通知剩余协程取消操作
	close(done)
	fmt.Println("\n已发送取消信号，剩余的工作协程将不会永久阻塞")

	// 等待一段时间让所有协程有机会响应取消信号
	time.Sleep(1 * time.Second)
}
