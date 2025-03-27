package replacer

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/yourusername/file-replacer/internal/config"
	"github.com/yourusername/file-replacer/pkg/logger"
)

// 替换结果
type ReplaceResult struct {
	FilePath        string
	Replaced        int
	Error           error
	ContentModified bool
}

// UnbufferedReplacer 使用无缓冲通道的替换器
type UnbufferedReplacer struct {
	config *config.Config
}

// NewUnbufferedReplacer 创建无缓冲通道替换器
func NewUnbufferedReplacer(cfg *config.Config) *UnbufferedReplacer {
	// 同样处理旧版替换项
	if cfg.SearchString != "" && cfg.ReplaceString != "" {
		found := false
		for _, item := range cfg.ReplaceItems {
			if item.SearchString == cfg.SearchString && item.ReplaceString == cfg.ReplaceString {
				found = true
				break
			}
		}
		if !found {
			cfg.AddReplaceItem(cfg.SearchString, cfg.ReplaceString)
		}
	}

	return &UnbufferedReplacer{
		config: cfg,
	}
}

// Replace 实现替换操作
func (r *UnbufferedReplacer) Replace(files []string) error {
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

	// 创建一个无缓冲结果通道
	resultChan := make(chan ReplaceResult)

	// 使用WaitGroup跟踪所有协程完成情况
	var wg sync.WaitGroup

	// 限制同时运行的协程数量
	sem := make(chan struct{}, r.config.Threads)
	logger.Log.Infof("使用无缓冲通道模式，最大并发数: %d", r.config.Threads)

	// 启动协程处理文件
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			// 获取信号量
			sem <- struct{}{}
			defer func() {
				// 释放信号量
				<-sem
				wg.Done()
			}()

			// 处理文件
			result := r.processFile(filePath)
			// 将结果发送到通道
			resultChan <- result
		}(file)
	}

	// 启动协程等待所有处理完成后关闭结果通道
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 处理结果
	var totalFiles, totalReplaced int
	for result := range resultChan {
		if result.Error != nil {
			logger.Log.Warnf("处理文件 %s 时出错: %v", result.FilePath, result.Error)
			continue
		}

		if result.Replaced > 0 {
			totalFiles++
			totalReplaced += result.Replaced
		}
	}

	logger.Log.Infof("替换完成，共处理 %d 个文件，替换 %d 处内容",
		totalFiles, totalReplaced)
	return nil
}

// processFile 处理单个文件
func (r *UnbufferedReplacer) processFile(filePath string) ReplaceResult {
	result := ReplaceResult{
		FilePath: filePath,
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = err
		return result
	}

	contentStr := string(content)
	originalContent := contentStr

	// 对每个替换项进行处理
	for _, item := range r.config.ReplaceItems {
		if !strings.Contains(contentStr, item.SearchString) {
			continue
		}

		// 计算替换数量
		count := strings.Count(contentStr, item.SearchString)
		if count > 0 {
			contentStr = strings.ReplaceAll(contentStr, item.SearchString, item.ReplaceString)
			result.Replaced += count

			logger.Log.Infof("文件 %s: 找到 '%s' %d 处匹配", filePath, item.SearchString, count)
		}
	}

	// 如果有替换
	if result.Replaced > 0 {
		result.ContentModified = contentStr != originalContent

		// 如果不是预览模式且内容有变化，则写入文件
		if !r.config.DryRun && result.ContentModified {
			err = os.WriteFile(filePath, []byte(contentStr), 0644)
			if err != nil {
				result.Error = err
				return result
			}
			logger.Log.Infof("已更新文件 %s，共替换 %d 处内容", filePath, result.Replaced)
		}
	}

	return result
}
