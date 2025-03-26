package config

// Config 应用程序配置
type Config struct {
	// 要扫描的根目录
	RootDir string
	// 要忽略的目录列表
	IgnoreDirs []string
	// 查找的字符串
	SearchString string
	// 替换的字符串
	ReplaceString string
	// 是否启用调试模式
	Debug bool
	// 是否进行实际替换（false为仅预览）
	DryRun bool
	// 并发线程数
	Threads int
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *Config {
	return &Config{
		RootDir:       "D:\\project\\cx_project\\china_mobile\\gitProject\\qqtgy\\src\\main\\webapp\\res\\wap",
		IgnoreDirs:    []string{"activityPages,.git", "node_modules", "vendor", "build", "dist"},
		SearchString:  "",
		ReplaceString: "",
		Debug:         false,
		DryRun:        false,
		Threads:       12,
	}
}
