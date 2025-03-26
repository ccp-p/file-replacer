# 常见问题解决方案

## 构建错误: 缺少 go.sum 条目

如果遇到类似以下错误:

```
missing go.sum entry for module providing package github.com/sirupsen/logrus
```

### 解决方法:

1. 运行 `go mod tidy` 命令生成正确的 go.sum 文件:

```bash
cd d:\project\my_go_project\file-replacer
go mod tidy
```

2. 然后重新构建项目:

```bash
go build -o file-replacer.exe cmd\file-replacer\main.go
```

## 如果使用了 ioutil 导致的警告

新版本 Go 中，`ioutil` 包已被弃用，应改用 `os` 包中对应的函数:

- 将 `ioutil.ReadFile()` 替换为 `os.ReadFile()`
- 将 `ioutil.WriteFile()` 替换为 `os.WriteFile()`

## Windows 下运行脚本

### 如果批处理脚本 (build.bat) 无法运行:

可能是文件格式或换行符问题。尝试以下方法:

1. 使用 PowerShell 脚本代替:
```
.\build.ps1
```

2. 或直接运行命令:
```
go mod tidy
go build -o file-replacer.exe cmd\file-replacer\main.go
```

3. 如果仍有问题，可以手动创建 build.bat 文件:
   - 打开记事本
   - 粘贴以下内容:
     ```
     @echo off
     echo 正在初始化Go模块...
     go mod tidy
     echo 正在构建项目...
     go build -o file-replacer.exe cmd\file-replacer\main.go
     if %ERRORLEVEL% EQU 0 (
         echo 构建成功！
     ) else (
         echo 构建失败，请检查错误信息
     )
     pause
     ```
   - 保存为 "build.bat" (注意使用引号确保文件扩展名正确)
