#!/usr/bin/env pwsh
# Encrypted Channels 一键测试脚本
# 功能：启动服务器 → 打开浏览器 → 等待订阅 → 自动推送测试消息

param(
    [switch]$SkipBuild,
    [switch]$OnlyPush
)

$ErrorActionPreference = "Continue"

Write-Host @"

╔════════════════════════════════════════════════╗
║   Encrypted Channels 端到端自动化测试         ║
╚════════════════════════════════════════════════╝

"@ -ForegroundColor Cyan

# 如果只是推送消息
if ($OnlyPush) {
    Write-Host "📤 仅推送模式" -ForegroundColor Yellow
    Write-Host ""
    
    # 检查服务器
    try {
        $null = Invoke-WebRequest -Uri "http://localhost:8070/system/pusher/info" -Method GET -TimeoutSec 2 -ErrorAction Stop
        Write-Host "✅ 服务器正在运行" -ForegroundColor Green
    } catch {
        Write-Host "❌ 服务器未运行" -ForegroundColor Red
        Write-Host "请先启动服务器：.\devinggo.exe" -ForegroundColor Yellow
        exit 1
    }
    
    # 推送消息
    Write-Host ""
    Write-Host "正在推送测试消息..." -ForegroundColor Yellow
    
    $testData = @{
        type = "encrypted-test"
        message = "端到端加密测试消息"
        amount = 99999.99
        timestamp = ([DateTimeOffset]::Now).ToString('yyyy-MM-ddTHH:mm:sszzz')
        test_id = [guid]::NewGuid().ToString()
    }
    
    $body = @{
        name = "encrypted-message"
        channel = "private-encrypted-secure"
        data = ($testData | ConvertTo-Json -Compress)
    } | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8070/system/pusher/events" `
            -Method POST -Body $body -ContentType "application/json"
        
        Write-Host ""
        Write-Host "✅ 推送成功！" -ForegroundColor Green
        Write-Host "📋 消息内容：" -ForegroundColor Cyan
        Write-Host ($testData | ConvertTo-Json) -ForegroundColor White
        Write-Host ""
        Write-Host "🔍 请在浏览器中查看消息" -ForegroundColor Yellow
    } catch {
        Write-Host ""
        Write-Host "❌ 推送失败: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    exit 0
}

# ==== 完整测试流程 ====

# 步骤1：检查并编译
Write-Host "[1/5] 检查可执行文件..." -ForegroundColor Yellow

if (-not (Test-Path ".\devinggo.exe") -or -not $SkipBuild) {
    Write-Host "正在编译服务器..." -ForegroundColor Gray
    $buildOutput = go build -o devinggo.exe .\main.go 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ 编译失败" -ForegroundColor Red
        Write-Host $buildOutput -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ 编译成功" -ForegroundColor Green
} else {
    Write-Host "✅ 可执行文件已存在" -ForegroundColor Green
}

# 步骤2：清理旧进程
Write-Host ""
Write-Host "[2/5] 清理旧进程..." -ForegroundColor Yellow

$oldProcesses = Get-Process -Name "devinggo" -ErrorAction SilentlyContinue
if ($oldProcesses) {
    $oldProcesses | Stop-Process -Force
    Start-Sleep -Seconds 2
    Write-Host "✅ 已清理 $($oldProcesses.Count) 个旧进程" -ForegroundColor Green
} else {
    Write-Host "✅ 无需清理" -ForegroundColor Green
}

# 清理端口占用
$portsToKill = Get-NetTCPConnection -LocalPort 8070 -ErrorAction SilentlyContinue
if ($portsToKill) {
    foreach ($p in $portsToKill) {
        Stop-Process -Id $p.OwningProcess -Force -ErrorAction SilentlyContinue
    }
    Start-Sleep -Seconds 2
    Write-Host "✅ 已释放端口 8070" -ForegroundColor Green
}

# 步骤3：启动服务器
Write-Host ""
Write-Host "[3/5] 启动服务器..." -ForegroundColor Yellow

$serverJob = Start-Job -ScriptBlock {
    Set-Location "E:\code\devinggo-light"
    .\devinggo.exe 2>&1
}

# 等待服务器启动
Write-Host "等待服务器启动..." -ForegroundColor Gray
$retries = 0
$maxRetries = 15
$serverReady = $false

while ($retries -lt $maxRetries) {
    Start-Sleep -Seconds 2
    $retries++
    
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8070/system/pusher/info" `
            -Method GET -TimeoutSec 2 -ErrorAction Stop
        
        if ($response.StatusCode -eq 200) {
            $serverReady = $true
            Write-Host "✅ 服务器启动成功" -ForegroundColor Green
            break
        }
    } catch {
        Write-Host "." -NoNewline -ForegroundColor Gray
    }
}

if (-not $serverReady) {
    Write-Host ""
    Write-Host "❌ 服务器启动失败" -ForegroundColor Red
    Write-Host "查看服务器日志：" -ForegroundColor Yellow
    Receive-Job $serverJob | Select-Object -First 30
    Stop-Job $serverJob
    Remove-Job $serverJob
    exit 1
}

# 步骤4：打开浏览器
Write-Host ""
Write-Host "[4/5] 打开测试页面..." -ForegroundColor Yellow

Start-Process "http://localhost:8070/pusher-test.html"
Write-Host "✅ 浏览器已打开" -ForegroundColor Green

# 步骤5：等待用户订阅并推送
Write-Host ""
Write-Host "[5/5] 准备推送测试消息..." -ForegroundColor Yellow
Write-Host ""
Write-Host "╔════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║           请完成以下操作：                     ║" -ForegroundColor Cyan
Write-Host "╠════════════════════════════════════════════════╣" -ForegroundColor Cyan
Write-Host "║  1. 点击「连接到 Pusher」                     ║" -ForegroundColor White
Write-Host "║  2. 找到「🔐 Encrypted Channel」区域        ║" -ForegroundColor White
Write-Host "║  3. 点击「订阅加密频道」                      ║" -ForegroundColor White
Write-Host "║  4. 看到「✅ 加密频道订阅成功」              ║" -ForegroundColor White
Write-Host "╚════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""
Write-Host "按任意键推送测试消息..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 推送多条测试消息
Write-Host ""
Write-Host "正在推送测试消息..." -ForegroundColor Yellow
Write-Host ""

$testMessages = @(
    @{
        name = "user-notification"
        data = @{
            type = "notification"
            title = "系统通知"
            message = "您有一条新消息"
            timestamp = ([DateTimeOffset]::Now).ToString('yyyy-MM-ddTHH:mm:sszzz')
        }
    },
    @{
        name = "payment-alert"
        data = @{
            type = "payment"
            message = "收到转账"
            amount = 5000.00
            from = "张三"
            card = "6222 **** **** 1234"
            timestamp = ([DateTimeOffset]::Now).ToString('yyyy-MM-ddTHH:mm:sszzz')
        }
    },
    @{
        name = "security-warning"
        data = @{
            type = "security"
            level = "high"
            message = "检测到异地登录"
            ip = "123.45.67.89"
            location = "北京"
            timestamp = ([DateTimeOffset]::Now).ToString('yyyy-MM-ddTHH:mm:sszzz')
        }
    }
)

$successCount = 0
foreach ($msg in $testMessages) {
    $body = @{
        name = $msg.name
        channel = "private-encrypted-secure"
        data = ($msg.data | ConvertTo-Json -Compress)
    } | ConvertTo-Json
    
    try {
        $null = Invoke-RestMethod -Uri "http://localhost:8070/system/pusher/events" `
            -Method POST -Body $body -ContentType "application/json"
        
        Write-Host "✅ [$($msg.name)] 推送成功" -ForegroundColor Green
        $successCount++
        Start-Sleep -Milliseconds 500
    } catch {
        Write-Host "❌ [$($msg.name)] 推送失败: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# 测试总结
Write-Host ""
Write-Host "╔════════════════════════════════════════════════╗" -ForegroundColor Green
Write-Host "║           测试完成                             ║" -ForegroundColor Green
Write-Host "╠════════════════════════════════════════════════╣" -ForegroundColor Green
Write-Host "║  推送成功: $successCount / $($testMessages.Count) 条消息                        ║" -ForegroundColor White
Write-Host "╚════════════════════════════════════════════════╝" -ForegroundColor Green
Write-Host ""
Write-Host "🔍 请在浏览器中查看消息日志" -ForegroundColor Yellow
Write-Host "   应该显示 3 条加密消息（已自动解密）" -ForegroundColor White
Write-Host ""
Write-Host "💡 提示：" -ForegroundColor Cyan
Write-Host "   - 服务器正在后台运行（Job ID: $($serverJob.Id)）" -ForegroundColor Gray
Write-Host "   - 可以继续手动推送：.\test_encrypted_push.ps1" -ForegroundColor Gray
Write-Host "   - 停止服务器：Stop-Job $($serverJob.Id); Remove-Job $($serverJob.Id)" -ForegroundColor Gray
Write-Host ""
Write-Host "按任意键停止服务器并退出..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# 清理
Write-Host ""
Write-Host "正在清理..." -ForegroundColor Yellow
Stop-Job $serverJob -ErrorAction SilentlyContinue
Remove-Job $serverJob -ErrorAction SilentlyContinue

$processes = Get-Process -Name "devinggo" -ErrorAction SilentlyContinue
if ($processes) {
    $processes | Stop-Process -Force
    Write-Host "✅ 服务器已停止" -ForegroundColor Green
}

Write-Host ""
Write-Host "测试完成！" -ForegroundColor Green
