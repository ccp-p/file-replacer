package main

import (
	"flag"
	"os"
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
	flag.StringVar(&cfg.SearchString, "search", "", "要查找的字符串 (单个替换时使用)")
	flag.StringVar(&cfg.ReplaceString, "replace", "", "替换成的字符串 (单个替换时使用)")
	flag.BoolVar(&cfg.Debug, "debug", false, "开启调试模式")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "预览模式(不进行实际替换)")
	flag.IntVar(&cfg.Threads, "threads", cfg.Threads, "并发线程数")

	ignoreFlag := flag.String("ignore", "", "要忽略的目录，用逗号分隔")
	replacePairsFlag := flag.String("pairs", "", "替换对列表，格式: \"search1:replace1,search2:replace2\"")
	pairsFileFlag := flag.String("pairs-file", "", "包含替换对的文件路径，每行一个替换对，格式: \"search replace\"")

	flag.Parse()

	// 处理忽略目录
	if *ignoreFlag != "" {
		cfg.IgnoreDirs = nil
		for _, dir := range splitCommaList(*ignoreFlag) {
			if dir != "" {
				cfg.IgnoreDirs = append(cfg.IgnoreDirs, dir)
			}
		}
	}

	// 处理替换对列表
	if *replacePairsFlag != "" {
		loadReplacePairsFromString(cfg, *replacePairsFlag)
	}

	// 处理替换对文件
	if *pairsFileFlag != "" {
		err := loadReplacePairsFromFile(cfg, *pairsFileFlag)
		if err != nil {
			logger.Log.Fatalf("加载替换对文件失败: %v", err)
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

// 从字符串加载替换对
func loadReplacePairsFromString(cfg *config.Config, pairsStr string) {
	pairs := splitCommaList(pairsStr)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			cfg.AddReplaceItem(parts[0], parts[1])
		}
	}
}

// 从文件加载替换对
func loadReplacePairsFromFile(cfg *config.Config, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue // 跳过空行和注释行
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			search := parts[0]
			replace := parts[1]
			cfg.AddReplaceItem(search, replace)
		}
	}

	return nil
}
