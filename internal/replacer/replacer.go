package replacer

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/yourusername/file-replacer/internal/config"
	"github.com/yourusername/file-replacer/pkg/logger"
)

// Replacer 文件内容替换器
type Replacer struct {
	config   *config.Config
	replaced int64
	files    int64
}

// NewReplacer 创建新的替换器
func NewReplacer(cfg *config.Config) *Replacer {
	// 如果有旧版的单个替换项，添加到替换项列表中
	if cfg.SearchString != "" && cfg.ReplaceString != "" {
		found := false
		// 检查是否已存在相同的替换项
		for _, item := range cfg.ReplaceItems {
			if item.SearchString == cfg.SearchString && item.ReplaceString == cfg.ReplaceString {
				found = true
				break
			}
		}
		// 如果不存在，则添加
		if !found {
			cfg.AddReplaceItem(cfg.SearchString, cfg.ReplaceString)
		}
	}

	return &Replacer{
		config:   cfg,
		replaced: 0,
		files:    0,
	}
}

// Replace 对指定文件列表执行替换操作
func (r *Replacer) Replace(files []string) error {
	if len(r.config.ReplaceItems) == 0 {
		return fmt.Errorf("没有指定替换项")
	}

	logger.Log.Infof("开始替换操作，共有 %d 个替换项", len(r.config.ReplaceItems))
	for i, item := range r.config.ReplaceItems {
		logger.Log.Infof("替换项 #%d: 搜索 '%s' 替换为 '%s'",
			i+1, item.SearchString, item.ReplaceString)
	}

	if r.config.DryRun {
		logger.Log.Info("当前为预览模式，不会进行实际替换")
	}

	logger.Log.Infof("使用 %d 个线程进行并行处理", r.config.Threads)

	// 使用有界的并发模型
	fileChan := make(chan string, len(files))
	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup
	for i := 0; i < r.config.Threads; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for file := range fileChan {
				err := r.replaceInFile(file)
				if err != nil {
					logger.Log.Warnf("[线程 %d] 处理文件 %s 时出错: %v", workerId, file, err)
				}
			}
		}(i)
	}

	// 等待所有替换任务完成
	wg.Wait()

	logger.Log.Infof("替换完成，共处理 %d 个文件，替换 %d 处内容",
		r.files, r.replaced)
	return nil
}

// replaceInFile 在单个文件中执行替换
func (r *Replacer) replaceInFile(filePath string) error {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentStr := string(content)
	originalContent := contentStr
	fileProcessed := false
	totalReplacements := 0

	// 对每个替换项进行处理
	for _, item := range r.config.ReplaceItems {
		if !strings.Contains(contentStr, item.SearchString) {
			continue
		}

		// 计算替换数量
		count := strings.Count(contentStr, item.SearchString)
		if count > 0 {
			contentStr = strings.ReplaceAll(contentStr, item.SearchString, item.ReplaceString)
			totalReplacements += count
			fileProcessed = true

			logger.Log.Infof("文件 %s: 找到 '%s' %d 处匹配", filePath, item.SearchString, count)
		}
	}

	// 如果文件被处理了
	if fileProcessed {
		atomic.AddInt64(&r.files, 1)
		atomic.AddInt64(&r.replaced, int64(totalReplacements))

		// 如果不是预览模式且内容有变化，则写入文件
		if !r.config.DryRun && contentStr != originalContent {
			err = os.WriteFile(filePath, []byte(contentStr), 0644)
			if err != nil {
				return err
			}
			logger.Log.Infof("已更新文件 %s，共替换 %d 处内容", filePath, totalReplacements)
		}
	}

	return nil
}
