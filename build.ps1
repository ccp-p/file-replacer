Write-Host "正在初始化Go模块..." -ForegroundColor Cyan
go mod tidy
chcp 65001 > $null
Write-Host "正在构建项目..." -ForegroundColor Cyan
go build -o file-replacer.exe cmd/file-replacer/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "构建成功！" -ForegroundColor Green
    Write-Host "运行示例:" -ForegroundColor Yellow
    Write-Host ".\file-replacer.exe -dir . -search `"oldText`" -replace `"newText`"" -ForegroundColor Yellow
} else {
    Write-Host "构建失败，请检查错误信息" -ForegroundColor Red
}
