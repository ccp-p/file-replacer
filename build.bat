@echo off
echo 正在初始化Go模块...
go mod tidy
echo 正在构建项目...
go build -o file-replacer.exe cmd\file-replacer\main.go
if %ERRORLEVEL% EQU 0 (
    echo 构建成功！运行示例:
    echo file-replacer.exe -dir . -search "oldText" -replace "newText"
) else (
    echo 构建失败，请检查错误信息
)
