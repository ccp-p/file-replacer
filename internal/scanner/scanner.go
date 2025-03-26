package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/file-replacer/internal/config"
	"github.com/yourusername/file-replacer/pkg/logger"
)

// FileScanner 文件扫描器
type FileScanner struct {
	config *config.Config
	files  []string
}

// NewFileScanner 创建新的文件扫描器
func NewFileScanner(cfg *config.Config) *FileScanner {
	return &FileScanner{
		config: cfg,
		files:  make([]string, 0),
	}
}

// Scan 扫描目录下的所有文件
func (s *FileScanner) Scan() ([]string, error) {
	logger.Log.Infof("开始扫描目录: %s", s.config.RootDir)

	err := filepath.Walk(s.config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Log.Errorf("访问路径 %s 时出错: %v", path, err)
			return err
		}

		// 检查是否为目录
		if info.IsDir() {
			// 检查是否应该忽略该目录
			if s.shouldIgnoreDir(path) {
				logger.Log.Debugf("忽略目录: %s", path)
				return filepath.SkipDir
			}
			return nil
		}

		// 将文件添加到扫描列表
		s.files = append(s.files, path)
		logger.Log.Debugf("找到文件: %s", path)
		return nil
	})

	if err != nil {
		logger.Log.Errorf("扫描过程中出错: %v", err)
		return nil, err
	}

	logger.Log.Infof("扫描完成，共发现 %d 个文件", len(s.files))
	return s.files, nil
}

// shouldIgnoreDir 检查是否应该忽略该目录
func (s *FileScanner) shouldIgnoreDir(path string) bool {
	dir := filepath.Base(path)
	for _, ignoreDir := range s.config.IgnoreDirs {
		if strings.EqualFold(dir, ignoreDir) {
			return true
		}
	}
	return false
}
