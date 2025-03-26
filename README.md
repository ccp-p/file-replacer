# File Replacer

一个用于扫描目录并替换文件中特定字符串的工具。

## 功能

- 扫描指定目录下的所有文件
- 支持配置忽略特定目录
- 执行文件内容的字符串查找和替换
- 提供预览模式，不进行实际修改
- 详细的操作日志
- 支持多线程并行处理，加快替换速度

## 使用方法

```bash
# 编译
go build -o file-replacer cmd/file-replacer/main.go

# 使用示例 - 预览模式
./file-replacer -dir ./myproject -search "oldText" -replace "newText" -dry-run

# 使用示例 - 执行实际替换
./file-replacer -dir ./myproject -search "oldText" -replace "newText" -dry-run=false

# 使用示例 - 指定忽略的目录
./file-replacer -dir ./myproject -search "oldText" -replace "newText" -ignore ".git,node_modules,vendor"

# 使用示例 - 开启调试日志
./file-replacer -dir ./myproject -search "oldText" -replace "newText" -debug

# 使用示例 - 指定并发线程数
./file-replacer -dir ./myproject -search "oldText" -replace "newText" -threads 8
```

## 参数说明

- `-dir`: 要扫描的根目录 (默认为 ".")
- `-search`: 要查找的字符串
- `-replace`: 替换成的字符串
- `-dry-run`: 预览模式，不进行实际替换 (默认为 true)
- `-ignore`: 要忽略的目录，用逗号分隔
- `-debug`: 开启调试模式 (默认为 false)
- `-threads`: 指定并发处理的线程数量 (默认为 4)
