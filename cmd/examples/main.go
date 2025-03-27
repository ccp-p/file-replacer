package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yourusername/file-replacer/examples"
)

func main() {
	demoNumber := "10" // 默认运行新的简化示例
	if len(os.Args) > 1 {
		demoNumber = os.Args[1]
	}

	switch demoNumber {
	case "1":
		examples.WithoutWaitGroupDemo()

	case "2":
		examples.WithWaitGroupDemo()

	case "3":
		examples.DemoRaceConditionRisk()

	case "4":
		examples.GoroutineLeakDemo()
		time.Sleep(2 * time.Second) // 给足够时间观察结果

	case "5":
		examples.VisualWaitGroupDemo()

	case "6":
		examples.ChannelClosingDemo()

	case "7":
		examples.BetterDemoWithTimeout()

	case "8":
		examples.ChannelBlockingVisual()

	case "9":
		examples.SafeChannelDemo()

	case "10":
		examples.SimpleChannelDemo()

	default:
		showHelp()
	}
}

func showHelp() {
	fmt.Println("请选择要运行的示例:")
	fmt.Println("1 - 不使用WaitGroup的通道模式")
	fmt.Println("2 - 使用WaitGroup的通道模式")
	fmt.Println("3 - 演示不使用WaitGroup的风险")
	fmt.Println("4 - 协程泄漏演示")
	fmt.Println("5 - 可视化WaitGroup工作原理")
	fmt.Println("6 - 通道关闭最佳实践")
	fmt.Println("7 - 无缓冲通道阻塞演示")
	fmt.Println("8 - 无缓冲通道阻塞可视化")
	fmt.Println("9 - 安全的通道操作模式")
	fmt.Println("10 - 简化版无缓冲通道阻塞演示")
	fmt.Println("\n使用方法: examples <编号>")
}
