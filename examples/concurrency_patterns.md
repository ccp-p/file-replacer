# Go 并发模式对比

## 1. 工作池模式 (当前代码使用)

```go
fileChan := make(chan string, len(files)) // 有缓冲通道
// 预先放入所有任务
for _, file := range files {
    fileChan <- file
}
close(fileChan)

// 创建固定数量工作协程
for i := 0; i < threads; i++ {
    go func() {
        for file := range fileChan {
            // 处理文件
        }
    }()
}
```

**特点：**
- 控制并发度，避免过多协程
- 通过缓冲通道预先排队所有任务
- 适合大量小任务场景
- 任务分配更均匀

## 2. 无缓冲通道结果收集模式

```go
results := make(chan string) // 无缓冲通道

// 为每个任务创建协程
for _, file := range files {
    go func(f string) {
        // 处理文件
        results <- "处理结果" // 阻塞直到有人接收
    }(file)
}

// 收集所有结果
for i := 0; i < len(files); i++ {
    result := <-results
    // 使用结果
}
```

**特点：**
- 每个任务一个协程，适合任务数量适中的场景
- 实时获取结果
- 结果处理顺序不确定
- 需要额外机制关闭通道
- 如果任务数量太多可能导致资源耗尽
