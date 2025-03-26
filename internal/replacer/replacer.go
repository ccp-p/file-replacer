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
	return &Replacer{
		config:   cfg,
		replaced: 0,
		files:    0,
	}
}

// Replace 对指定文件列表执行替换操作
func (r *Replacer) Replace(files []string) error {
	if r.config.SearchString == "" {
		return fmt.Errorf("搜索字符串不能为空")
	}

	logger.Log.Infof("开始替换操作，搜索 '%s' 替换为 '%s'",
		r.config.SearchString, r.config.ReplaceString)
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

	// 将内容转换为字符串并检查是否包含搜索字符串
	contentStr := string(content)
	if !strings.Contains(contentStr, r.config.SearchString) {
		logger.Log.Debugf("文件 %s 不包含搜索字符串，跳过", filePath)
		return nil
	}

	// 计算替换数量
	count := strings.Count(contentStr, r.config.SearchString)
	atomic.AddInt64(&r.replaced, int64(count))
	atomic.AddInt64(&r.files, 1)

	// 执行替换
	newContent := strings.ReplaceAll(contentStr, r.config.SearchString, r.config.ReplaceString)

	logger.Log.Infof("文件 %s: 找到 %d 处匹配", filePath, count)

	// 如果不是预览模式，则写入文件
	if !r.config.DryRun {
		err = os.WriteFile(filePath, []byte(newContent), 0644)
		if err != nil {
			return err
		}
		logger.Log.Infof("已更新文件 %s", filePath)
	}

	return nil
}
