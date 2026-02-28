# Encrypted Channels 端到端测试指南

## 问题说明

**`go test` 和真实推送的区别**：

- ❌ **`go test`**：只是单元测试，测试代码逻辑，不会启动 WebSocket 服务器，不会推送消息到浏览器
- ✅ **端到端测试**：启动服务器 → 浏览器订阅频道 → 服务器推送消息 → 浏览器接收并解密

---

## 端到端测试步骤

### 第一步：启动服务器

```powershell
# 编译并启动服务器
go build -o devinggo.exe .\main.go
.\devinggo.exe
```

**或者直接运行**：
```powershell
gf run main.go
```

等待看到日志：
```
服务器启动成功
WebSocket 服务已启动：ws://localhost:8070/system/pusher/ws
```

---

### 第二步：打开浏览器测试页面

在浏览器中打开：
```
http://localhost:8070/pusher-test.html
```

---

### 第三步：连接并订阅加密频道

1. **点击「连接到 Pusher」按钮**
   - 等待看到绿色状态：`已连接 (socket_id: xxxxx)`

2. **找到「🔐 Encrypted Channel」区域**
   - 默认频道名：`private-encrypted-secure`
   - 点击「订阅加密频道」按钮

3. **确认订阅成功**
   - 日志中应显示：  
     ```
     ✅ 加密频道订阅成功: private-encrypted-secure
     ⚠️ 注意：Encrypted Channels 不支持 Client Events（仅接收服务器推送）
     ```

---

### 第四步：推送测试消息

**打开新的 PowerShell 窗口**，运行测试脚本：

```powershell
cd e:\code\devinggo-light
.\test_encrypted_push.ps1
```

**或者手动推送**：

```powershell
# 构造消息
$testData = @{
    type = "payment-notification"
    message = "您的账户收到一笔转账"
    amount = 5000.00
    card_number = "6222 **** **** 1234"
    timestamp = ([DateTimeOffset]::Now).ToString("yyyy-MM-ddTHH:mm:sszzz")
} | ConvertTo-Json

# 推送请求
$body = @{
    name = "encrypted-message"
    channel = "private-encrypted-secure"
    data = $testData
} | ConvertTo-Json

# 发送推送
Invoke-RestMethod -Uri "http://localhost:8070/system/pusher/events" `
    -Method POST `
    -Body $body `
    -ContentType "application/json; charset=utf-8"
```

---

### 第五步：验证接收

**在浏览器日志中查看**，应该看到：

```
📨 收到加密频道事件: encrypted-message => {...}
```

消息内容会自动解密显示，例如：
```json
{
  "type": "payment-notification",
  "message": "您的账户收到一笔转账",
  "amount": 5000.00,
  "card_number": "6222 **** **** 1234",
  "timestamp": "2026-02-28T18:48:54+08:00"
}
```

---

## 测试验证清单

- [ ] 服务器正常启动（端口 8070）
- [ ] 浏览器成功连接到 Pusher
- [ ] 加密频道订阅成功（需要 HTTP 认证）
- [ ] 服务器推送消息成功返回
- [ ] 浏览器收到消息并自动解密
- [ ] 日志中显示完整的消息内容

---

## 常见问题

### 1. 浏览器没收到消息

**检查项**：
- [ ] 服务器是否在运行（`http://localhost:8070` 能访问）
- [ ] 浏览器是否已连接（状态显示「已连接」）
- [ ] 是否已订阅加密频道（日志有「订阅成功」）
- [ ] 频道名称是否匹配（推送和订阅必须一致）

### 2. 订阅失败：No shared_secret

**原因**：服务器未返回 `shared_secret`

**解决**：
1. 确认 `pusher_auth.go` 中有 Encrypted Channel 处理：
```go
if websocket.IsEncryptedChannel(req.ChannelName) {
    sharedSecret := websocket.GenerateSharedSecret()
    auth := websocket.GenerateAuthSignature(req.SocketId, req.ChannelName, "")
    
    r.Response.WriteJson(g.Map{
        "auth":          auth,
        "shared_secret": sharedSecret,
    })
}
```

2. 重启服务器

### 3. 推送失败 401/403

**原因**：认证失败或无权限

**解决**：
- 使用内部 API（开发环境）：`http://localhost:8070/system/pusher/events`
- 生产环境需要在 Header 中添加 JWT Token：
  ```powershell
  -Headers @{ "Authorization" = "Bearer YOUR_JWT_TOKEN" }
  ```

### 4. NaCl 库未加载

**表现**：订阅时提示「NaCl 加密库未加载」

**解决**：
1. 检查 pusher-test.html 中是否引入 TweetNaCl：
```html
<script src="https://cdn.jsdelivr.net/npm/tweetnacl@1.0.3/nacl-fast.min.js"></script>
```

2. 刷新浏览器页面
3. 检查浏览器控制台是否有加载错误

---

## 测试脚本说明

### test_encrypted_push.ps1

**功能**：
- 自动检查服务器状态
- 提示用户完成浏览器订阅
- 构造并推送测试消息
- 显示详细的推送结果

**使用方法**：
```powershell
# 1. 启动服务器（终端1）
.\devinggo.exe

# 2. 浏览器中订阅加密频道
# 打开: http://localhost:8070/pusher-test.html
# 连接 → 订阅 private-encrypted-secure

# 3. 运行推送脚本（终端2）
.\test_encrypted_push.ps1
```

---

## 高级测试

### 多消息推送测试

```powershell
# 连续推送 10 条消息
for ($i = 1; $i -le 10; $i++) {
    $data = @{
        message_id = $i
        content = "测试消息 #$i"
        timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
    } | ConvertTo-Json
    
    $body = @{
        name = "test-event"
        channel = "private-encrypted-secure"
        data = $data
    } | ConvertTo-Json
    
    Invoke-RestMethod -Uri "http://localhost:8070/system/pusher/events" `
        -Method POST -Body $body -ContentType "application/json"
    
    Write-Host "✅ 消息 $i 已发送"
    Start-Sleep -Milliseconds 500
}
```

### 性能压力测试

```powershell
# 并发推送 100 条消息
$jobs = 1..100 | ForEach-Object {
    Start-Job -ScriptBlock {
        param($id)
        $data = @{
            test_id = $id
            timestamp = [DateTimeOffset]::Now.ToUnixTimeSeconds()
        } | ConvertTo-Json
        
        $body = @{
            name = "stress-test"
            channel = "private-encrypted-secure"
            data = $data
        } | ConvertTo-Json
        
        Invoke-RestMethod -Uri "http://localhost:8070/system/pusher/events" `
            -Method POST -Body $body -ContentType "application/json"
    } -ArgumentList $_
}

$jobs | Wait-Job | Receive-Job
$jobs | Remove-Job

Write-Host "✅ 100 条消息推送完成"
```

---

## 总结

### 单元测试 vs 端到端测试

| 类型 | 命令 | 作用 | 是否推送到浏览器 |
|------|------|------|------------------|
| **单元测试** | `go test ./modules/system/pkg/websocket` | 测试代码逻辑 | ❌ 否 |
| **端到端测试** | 启动服务器 + 浏览器 + 推送脚本 | 测试完整功能 | ✅ 是 |

### 完整测试流程

1. ✅ **单元测试**（代码验证）：`go test -v ./modules/system/pkg/websocket`
2. ✅ **性能测试**（性能验证）：`go test -bench=Benchmark -benchmem`
3. ✅ **端到端测试**（功能验证）：启动服务器 + 浏览器订阅 + 推送消息
4. ✅ **压力测试**（稳定性验证）：并发推送多条消息

---

**注意**：Encrypted Channels 的加密/解密由客户端（TweetNaCl）完成，服务器只负责路由消息，不会解密内容（端到端加密）。

