# Pusher Protocol Test Script
# Send events to WebSocket server via HTTP API

$ErrorActionPreference = "Stop"

# Configuration (match with hack/config.yaml)
$appId = "devinggo-app-id"
$appKey = "devinggo-app-key"
$appSecret = "devinggo-app-secret"
$serverHost = "localhost:8070"

Write-Host "Pusher HTTP API Test Script" -ForegroundColor Green
Write-Host "Server: http://$serverHost" -ForegroundColor Cyan
Write-Host "App ID: $appId" -ForegroundColor Cyan
Write-Host "App Key: $appKey" -ForegroundColor Cyan
Write-Host ""

# HMAC-SHA256 signature function
function Get-PusherSignature {
    param(
        [string]$method,
        [string]$path,
        [string]$queryString,
        [string]$body,
        [string]$secret
    )
    
    $stringToSign = "$method`n$path`n$queryString"
    if ($body) {
        $bodyMd5 = [System.Security.Cryptography.MD5]::Create().ComputeHash([System.Text.Encoding]::UTF8.GetBytes($body))
        $bodyMd5Hex = [System.BitConverter]::ToString($bodyMd5).Replace("-", "").ToLower()
        $stringToSign += "`n$bodyMd5Hex"
    }
    
    $hmac = New-Object System.Security.Cryptography.HMACSHA256
    $hmac.Key = [System.Text.Encoding]::UTF8.GetBytes($secret)
    $signature = $hmac.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($stringToSign))
    $signatureHex = [System.BitConverter]::ToString($signature).Replace("-", "").ToLower()
    
    return $signatureHex
}

# Send event function
function Send-PusherEvent {
    param(
        [string]$channel,
        [string]$event,
        [hashtable]$data
    )
    
    $timestamp = [int][double]::Parse((Get-Date -UFormat %s))
    $bodyJson = @{
        name = $event
        channels = @($channel)
        data = ($data | ConvertTo-Json -Compress)
    } | ConvertTo-Json -Compress
    
    # Build query parameters
    $queryParams = @{
        auth_key = $appKey
        auth_timestamp = $timestamp
        auth_version = "1.0"
        body_md5 = [System.Security.Cryptography.MD5]::Create().ComputeHash([System.Text.Encoding]::UTF8.GetBytes($bodyJson)) |
                   ForEach-Object { [System.BitConverter]::ToString($_).Replace("-", "").ToLower() }
    }
    
    $queryString = ($queryParams.GetEnumerator() | Sort-Object Name | ForEach-Object { "$($_.Key)=$($_.Value)" }) -join "&"
    
    # Calculate signature
    $path = "/apps/$appId/events"
    $signature = Get-PusherSignature -method "POST" -path $path -queryString $queryString -body $bodyJson -secret $appSecret
    $queryString += "&auth_signature=$signature"
    
    # Send request
    $url = "http://$serverHost$path`?$queryString"
    
    try {
        $response = Invoke-RestMethod -Uri $url -Method POST -Body $bodyJson -ContentType "application/json" -ErrorAction Stop
        return $true
    } catch {
        Write-Host "Send failed: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $responseBody = $reader.ReadToEnd()
            Write-Host "Response: $responseBody" -ForegroundColor Red
        }
        return $false
    }
}

# Test 1: Send event to Public Channel
Write-Host "Test 1: Send event to Public Channel" -ForegroundColor Yellow
$timestamp = Get-Date -Format "HH:mm:ss"
$data1 = @{
    user = "TestBot"
    message = "Hello from PowerShell Test Script!"
    time = $timestamp
}
Write-Host "  Channel: chat-room" -ForegroundColor Gray
Write-Host "  Event: new-message" -ForegroundColor Gray
Write-Host "  Data: $($data1 | ConvertTo-Json -Compress)" -ForegroundColor Gray

if (Send-PusherEvent -channel "chat-room" -event "new-message" -data $data1) {
    Write-Host "Success!" -ForegroundColor Green
}
Start-Sleep -Seconds 1

# Test 2: Send event to Private Channel
Write-Host "`nTest 2: Send event to Private Channel" -ForegroundColor Yellow
$timestamp = Get-Date -Format "HH:mm:ss"
$data2 = @{
    type = "info"
    title = "System Notification"
    message = "This is a private message from PowerShell script"
    time = $timestamp
}
Write-Host "  Channel: private-user-123" -ForegroundColor Gray
Write-Host "  Event: notification" -ForegroundColor Gray
Write-Host "  Data: $($data2 | ConvertTo-Json -Compress)" -ForegroundColor Gray

if (Send-PusherEvent -channel "private-user-123" -event "notification" -data $data2) {
    Write-Host "Success!" -ForegroundColor Green
}
Start-Sleep -Seconds 1

# Test 3: Send event to Presence Channel
Write-Host "`nTest 3: Send event to Presence Channel" -ForegroundColor Yellow
$timestamp = Get-Date -Format "HH:mm:ss"
$data3 = @{
    type = "broadcast"
    message = "Attention all online users: This is a broadcast from PowerShell"
    time = $timestamp
}
Write-Host "  Channel: presence-lobby" -ForegroundColor Gray
Write-Host "  Event: announcement" -ForegroundColor Gray
Write-Host "  Data: $($data3 | ConvertTo-Json -Compress)" -ForegroundColor Gray

if (Send-PusherEvent -channel "presence-lobby" -event "announcement" -data $data3) {
    Write-Host "Success!" -ForegroundColor Green
}
Start-Sleep -Seconds 1

# Test 4: Send multiple messages
Write-Host "`nTest 4: Send multiple messages" -ForegroundColor Yellow
for ($i = 1; $i -le 3; $i++) {
    $timestamp = Get-Date -Format "HH:mm:ss"
    $data = @{
        seq = $i
        message = "Batch message #$i"
        time = $timestamp
    }
    Write-Host "  Sending message $i..." -ForegroundColor Gray -NoNewline
    if (Send-PusherEvent -channel "chat-room" -event "batch-message" -data $data) {
        Write-Host " OK" -ForegroundColor Green
    }
    Start-Sleep -Milliseconds 500
}

Write-Host "`nTest completed!" -ForegroundColor Green
Write-Host "Tip: Open http://localhost:8070/pusher-test.html in browser to see real-time events" -ForegroundColor Cyan
