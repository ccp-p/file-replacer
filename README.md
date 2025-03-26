# File Replacer

一个用于扫描目录并替换文件中特定字符串的工具。

## 功能

- 扫描指定目录下的所有文件
- 支持配置忽略特定目录
- 执行文件内容的字符串查找和替换
- **支持同时执行多组替换操作**
- 提供预览模式，不进行实际修改
- 详细的操作日志
- 支持多线程并行处理，加快替换速度

## 使用方法

```bash
# 编译
go build -o file-replacer cmd/file-replacer/main.go

# 使用示例 - 单个替换
./file-replacer -dir ./myproject -search "oldText" -replace "newText"

# 使用示例 - 多组替换 (使用命令行参数)
./file-replacer -dir ./myproject -pairs "oldText1:newText1,oldText2:newText2"

# 使用示例 - 从文件读取替换对
./file-replacer -dir ./myproject -pairs-file replace_config.txt

# 使用示例 - 预览模式
./file-replacer -dir ./myproject -pairs "oldText1:newText1,oldText2:newText2" -dry-run

# 使用示例 - 指定忽略的目录
./file-replacer -dir ./myproject -pairs "oldText1:newText1" -ignore ".git,node_modules,vendor"

# 使用示例 - 开启调试日志
./file-replacer -dir ./myproject -pairs "oldText1:newText1" -debug

# 使用示例 - 指定并发线程数
./file-replacer -dir ./myproject -pairs "oldText1:newText1" -threads 8
```

## 参数说明

- `-dir`: 要扫描的根目录
- `-search`, `-replace`: 单个替换项的搜索和替换字符串
- `-pairs`: 多个替换项，格式为 "搜索1:替换1,搜索2:替换2,..."
- `-pairs-file`: 包含替换对的文件路径，每行一个替换对，格式为 "搜索 替换"
- `-dry-run`: 预览模式，不进行实际替换 (默认为 false)
- `-ignore`: 要忽略的目录，用逗号分隔
- `-debug`: 开启调试模式 (默认为 false)
- `-threads`: 指定并发处理的线程数量 (默认为 12)

## 替换对文件格式示例

```
# 这是注释行
qqt.cmicrwx.cn qqt.cmicvip.cn
qqt-res.cmicrwx.cn qqt-res.cmicvip.cn
```
