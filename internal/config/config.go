package config

// ReplaceItem 表示一个替换项
type ReplaceItem struct {
	// 查找的字符串
	SearchString string
	// 替换的字符串
	ReplaceString string
}

// Config 应用程序配置
type Config struct {
	// 要扫描的根目录
	RootDir string
	// 要忽略的目录列表
	IgnoreDirs []string
	// 替换项列表
	ReplaceItems []ReplaceItem
	// 是否启用调试模式
	Debug bool
	// 是否进行实际替换（false为仅预览）
	DryRun bool
	// 并发线程数
	Threads int
	// 兼容旧版的单个替换项
	SearchString  string
	ReplaceString string
}

// NewDefaultConfig 返回默认配置
func NewDefaultConfig() *Config {
	return &Config{
		RootDir:    "D:\\project\\cx_project\\china_mobile\\gitProject\\qqtgy\\src\\main\\webapp\\res\\wap",
		IgnoreDirs: []string{"activityPages,.git", "node_modules", "vendor", "build", "dist"},
		ReplaceItems: []ReplaceItem{
			{
				SearchString:  "qqt.cmicrwx.cn",
				ReplaceString: "qqt.cmicvip.cn",
			},
			{
				SearchString:  "qqt-res.cmicrwx.cn",
				ReplaceString: "qqt.cmicvip.cn",
			},
		},
		Debug:   false,
		DryRun:  false,
		Threads: 12,
	}
}

// AddReplaceItem 添加一个替换项
func (c *Config) AddReplaceItem(search, replace string) {
	c.ReplaceItems = append(c.ReplaceItems, ReplaceItem{
		SearchString:  search,
		ReplaceString: replace,
	})
}
