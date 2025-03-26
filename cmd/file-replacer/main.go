package main

import (
	"flag"
	"strings"

	"github.com/yourusername/file-replacer/internal/config"
	"github.com/yourusername/file-replacer/internal/replacer"
	"github.com/yourusername/file-replacer/internal/scanner"
	"github.com/yourusername/file-replacer/pkg/logger"
)

func main() {
	// 解析命令行参数
	cfg := config.NewDefaultConfig()

	flag.StringVar(&cfg.RootDir, "dir", cfg.RootDir, "要扫描的根目录")
	flag.StringVar(&cfg.SearchString, "search", "qqt.cmicrwx.cn", "要查找的字符串")
	flag.StringVar(&cfg.ReplaceString, "replace", "qqt.cmicvip.cn", "替换成的字符串")
	flag.BoolVar(&cfg.Debug, "debug", false, "开启调试模式")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "预览模式(不进行实际替换)")
	flag.IntVar(&cfg.Threads, "threads", cfg.Threads, "并发线程数")

	ignoreFlag := flag.String("ignore", "", "要忽略的目录，用逗号分隔")

	flag.Parse()

	// 如果指定了忽略目录，覆盖默认设置
	if *ignoreFlag != "" {
		cfg.IgnoreDirs = nil
		for _, dir := range splitCommaList(*ignoreFlag) {
			if dir != "" {
				cfg.IgnoreDirs = append(cfg.IgnoreDirs, dir)
			}
		}
	}

	// 设置日志级别
	logger.SetDebug(cfg.Debug)

	// 执行扫描
	fileScanner := scanner.NewFileScanner(cfg)
	files, err := fileScanner.Scan()
	if err != nil {
		logger.Log.Fatalf("扫描失败: %v", err)
	}

	// 执行替换
	fileReplacer := replacer.NewReplacer(cfg)
	err = fileReplacer.Replace(files)
	if err != nil {
		logger.Log.Fatalf("替换失败: %v", err)
	}
}

// 辅助函数: 分割逗号分隔列表
func splitCommaList(list string) []string {
	if list == "" {
		return []string{}
	}
	return strings.Split(list, ",")
}
