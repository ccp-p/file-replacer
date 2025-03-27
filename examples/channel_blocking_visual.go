package examples

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// ChannelBlockingVisual 可视化展示无缓冲通道的阻塞行为
func ChannelBlockingVisual() {
	ch := make(chan string)

	// 定义跨平台清屏函数
	clear := func() {
		// 针对不同操作系统使用不同清屏命令
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/c", "cls")
		} else {
			cmd = exec.Command("clear")
		}
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}

	// 替代清屏方法 - 如果上面的方法不起作用，使用多行分隔
	alternateClear := func() {
		// 打印多行分隔来模拟清屏效果
		fmt.Println(strings.Repeat("\n", 50))
		fmt.Println("============================================================")
	}

	// 绘制可视化界面
	draw := func(senderState, receiverState, channelState string) {
		// 尝试两种方式，至少一种会生效
		clear()
		alternateClear()

		fmt.Println("=== 无缓冲通道阻塞可视化演示 ===\n")

		// 发送者状态框
		fmt.Println("┌────────────发送者─────────────┐")
		fmt.Printf("│ 状态: %-23s │\n", senderState)
		fmt.Println("└──────────────────────────────┘")

		// 通道状态
		arrowLeft := " "
		arrowRight := " "
		if channelState == "sending" {
			arrowLeft = ">"
		}
		if channelState == "receiving" {
			arrowRight = "<"
		}
		if channelState == "transferring" {
			arrowLeft = ">"
			arrowRight = "<"
		}

		fmt.Println()
		fmt.Printf("   %s──────────────%s   \n", arrowLeft, arrowRight)
		fmt.Printf("   │  [通道]   │   \n")
		fmt.Printf("   %s──────────────%s   \n", arrowLeft, arrowRight)
		fmt.Println()

		// 接收者状态框
		fmt.Println("┌────────────接收者─────────────┐")
		fmt.Printf("│ 状态: %-23s │\n", receiverState)
		fmt.Println("└──────────────────────────────┘")

		fmt.Println("\n按 Ctrl+C 退出演示")
		fmt.Println("\n当前时间:", time.Now().Format("15:04:05.000"))

		// 添加一个小延迟，让动画效果更明显
		time.Sleep(800 * time.Millisecond)
	}

	// 启动发送者协程
	go func() {
		draw("初始化中", "未启动", " ")
		time.Sleep(1 * time.Second)

		draw("准备发送数据", "未启动", " ")
		time.Sleep(1 * time.Second)

		draw("尝试发送 -> 阻塞!", "未启动", "sending")

		// 这里会阻塞，直到有人接收
		ch <- "数据"

		draw("发送成功，继续执行", "已接收数据", "transferring")
		time.Sleep(1 * time.Second)

		draw("任务完成", "处理数据中", " ")
	}()

	// 主协程模拟接收者，但故意延迟
	time.Sleep(3 * time.Second)
	draw("等待中...(已阻塞)", "启动中", " ")
	time.Sleep(1 * time.Second)

	draw("等待中...(已阻塞)", "准备接收", "receiving")
	time.Sleep(1 * time.Second)

	// 接收数据
	data := <-ch

	draw("发送成功，继续执行", "已接收数据: "+data, "transferring")
	time.Sleep(2 * time.Second)
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
