# WebSocket改造支持Pusher协议 - 实施计划（v8.3.0兼容版）

## 概述

将现有自定义WebSocket协议完全改造为兼容Pusher.com的标准协议，支持Public/Private/Presence三种Channel类型和Client Events。采用标准Pusher认证机制（HMAC-SHA256签名），利用现有的Redis集群架构实现分布式支持。

**关键决策**：
- ✅ 完全切换到Pusher协议（不保留旧协议）
- ✅ 实现完整功能集（Public/Private/Presence + Client Events）
- ✅ 使用标准Pusher认证端点（`/pusher/auth`）
- ✅ 复用现有Redis集群架构
- ✅ 完全兼容pusher-js v8.3.0客户端库

**预计工期**: 14-17天（分4个阶段，已针对v8.3.0优化）

**v8.3.0兼容性声明**：本计划已根据pusher-js v8.3.0的协议要求进行全面审查和修正。

---

## 实施步骤

### Phase 1: 核心协议改造（P0，必须完成）- 5天

#### 1. 重构消息数据模型
**文件**: [model.go](model.go)

**改造内容**:
- 删除现有 `Request` 和 `WResponse` 结构
- 新增 `PusherRequest` 结构：
  - `event` string - 事件名（如 "pusher:subscribe"）
  - `data` interface{} - JSON对象或字符串
  - `channel` string - 频道名（可选）
- 新增 `PusherResponse` 结构：
  - `event` string - 事件名
  - `channel` string - 频道名（可选）
  - `data` string - JSON字符串（⚠️ v8.3.0严格要求：必须是字符串）
- 新增 `ConnectionEstablishedData` 结构：
  - `socket_id` string - 连接唯一标识
  - `activity_timeout` int - 活动超时（秒，v8.3.0推荐120）
- 新增 `ErrorMessage` 结构：
  - `event` string - "pusher:error"
  - `message` string - 错误消息
  - `code` int - 错误码
- 新增 `SubscriptionErrorData` 结构（⚠️ v8.3.0要求）：
  - `type` string - 错误类型（"AuthError" | "ChannelLimitReached"）
  - `error` string - 错误描述
  - `status` int - HTTP状态码 (401/403/500)

**代码示例**:
```go
// ⚠️ v8.3.0关键要求：服务器→客户端的data字段必须是JSON字符串
type PusherResponse struct {
    Event   string `json:"event"`
    Channel string `json:"channel,omitempty"`
    Data    string `json:"data"` // 必须是JSON字符串，不能是对象
}

// 新增：subscription_error专用结构（v8.3.0要求）
type SubscriptionErrorData struct {
    Type   string `json:"type"`   // "AuthError" | "ChannelLimitReached"
    Error  string `json:"error"`  // 错误描述
    Status int    `json:"status"` // HTTP状态码
}

// 错误码常量（v8.3.0完整版）
const (
    CodeNormalClosure        = 4000 // 正常关闭
    CodeGoingAway            = 4001 // 服务器下线
    CodeMaxConnections       = 4004 // 连接数超限
    CodePathNotFound         = 4005 // 路径错误
    CodeOverCapacity         = 4008 // 服务过载
    CodeUnauthorized         = 4009 // 未授权
    CodeAppDisabled          = 4100 // 应用已禁用
    CodePongTimeout          = 4200 // 心跳超时
    CodeClosedByClient       = 4201 // 客户端主动关闭
    CodeClientEventForbidden = 4301 // 客户端事件错误
)
```

#### 2. 修改客户端结构
**文件**: [client.go](client.go)

**改造内容**:
- 将 `ID` 字段重命名/增加 `SocketID` 字段（格式：`{serverName}.{uniqueId}`）
- 将 `topics` 字段改为 `channels`（字符串数组）
- 新增字段（v8.3.0兼容）：
  - `UserID` string - Presence Channel用户ID
  - `UserInfo` map[string]interface{} - Presence Channel用户信息
  - `LastEventTime` time.Time - 用于速率限制
  - `EventCount` int - 1秒内事件计数
- 修改 `SendMsg` 方法：
  - 支持Pusher消息格式
  - ⚠️ 确保 `data` 字段序列化为JSON字符串
- 删除 `ResponseSuccess`/`ResponseFail` 方法
- 新增 `SendPusherEvent` 方法：发送标准Pusher事件（自动序列化data）
- 新增 `SendError` 方法：发送错误消息
- 新增 `SendSubscriptionError` 方法：发送订阅错误（v8.3.0格式）

**代码示例**:
```go
type Client struct {
    // 现有字段...
    SocketID       string                 // v8.3.0格式：{serverName}.{uniqueId}
    Channels       []string               // 改名自topics
    UserID         string                 // Presence Channel用户ID
    UserInfo       map[string]interface{} // Presence Channel用户信息
    LastEventTime  time.Time              // 用于速率限制
    EventCount     int                    // 1秒内事件计数
    IsDisconnected bool                   // 断线标记（Grace Period用）
}

// 新增：发送Pusher事件（⚠️ 确保data序列化为JSON字符串）
func (c *Client) SendPusherEvent(event, channel string, data interface{}) error {
    var dataStr string
    if str, ok := data.(string); ok {
        dataStr = str
    } else {
        jsonBytes, err := json.Marshal(data)
        if err != nil {
            return err
        }
        dataStr = string(jsonBytes)
    }
    
    response := &PusherResponse{
        Event:   event,
        Channel: channel,
        Data:    dataStr,
    }
    return c.SendMsg(response)
}

// 新增：发送subscription_error（v8.3.0格式）
func (c *Client) SendSubscriptionError(channel, errorType, message string, status int) {
    errData := SubscriptionErrorData{
        Type:   errorType,
        Error:  message,
        Status: status,
    }
    c.SendPusherEvent("pusher:subscription_error", channel, errData)
}
```

#### 3. 改造事件路由器
**文件**: [router.go](router.go)

**改造内容**:
- 修改 `Reg` 函数，注册Pusher系统事件：
  - `pusher:subscribe` → SubscribeController
  - `pusher:unsubscribe` → UnsubscribeController
  - `pusher:ping` → PingController
  - 以 `client-` 开头的事件 → ClientEventController（新增）
- 修改 `ProcessData` 函数：
  - 解析 `PusherRequest` 格式
  - 移除 `BindEvent` 机制（Pusher不需要）
  - 优先检查 `client-` 事件前缀
- 删除旧事件注册（如 `subscribe`、`broadcastMsg` 等）

#### 4. 重写事件控制器
**文件**: [controller.go](controller.go)

**改造内容**:

**SubscribeController**: 
- 解析 `data.channel` 字段
- 暂时只支持Public Channel（`private-`/`presence-`前缀返回错误）
- 调用 `JoinChannel4Redis`（重命名自 `JoinTopic4Redis`）
- 发送 `pusher:subscription_succeeded` 事件（⚠️ data为空对象字符串 `"{}"`）

**UnsubscribeController**:
- 解析 `data.channel`
- 调用 `LeaveChannel4Redis`
- ⚠️ 不发送响应（Pusher规范无响应事件）

**PingController**:
- 发送 `pusher:pong` 事件（⚠️ data为空对象字符串 `"{}"`）

**删除控制器**:
- `BindEventController`
- `BroadcastMessageController`
- `IdMessageController`
- `PublishController`

#### 5. 修改连接初始化
**文件**: [init.go](init.go)

**改造内容**:
- ⚠️ v8.3.0新增：验证查询参数 `?protocol=7`（必须是7，否则返回4005错误）
- 修改 `WsPage` 函数，连接建立后：
  - 生成 `socket_id`（格式：`{serverName}.{timestamp}{random}`）
  - 发送 `pusher:connection_established` 事件
  - ⚠️ 包含 `socket_id` 和 `activity_timeout: 120`（v8.3.0推荐值）
  - ⚠️ data必须是JSON字符串格式
- 移除旧的 `Connected` 事件逻辑
- ⚠️ 心跳超时调整为150秒（120+30，v8.3.0推荐）

**代码示例**:
```go
func WsPage(r *ghttp.Request) {
    ctx := r.GetCtx()
    
    // ⚠️ v8.3.0要求：验证协议版本
    protocol := r.GetQuery("protocol").String()
    if protocol != "" && protocol != "7" {
        r.Response.WriteStatus(400)
        r.Response.WriteJson(g.Map{"error": "Unsupported protocol version"})
        return
    }
    
    conn, err := GetConnection(r)
    if err != nil {
        return
    }
    
    // 生成socket_id（v8.3.0格式）
    socketID := fmt.Sprintf("%s.%d%d", 
        clientManager.GetServerName(), 
        time.Now().Unix(), 
        rand.Intn(100000))
    
    client := NewClient(socketID, conn, r)
    
    // 发送connection_established（⚠️ activity_timeout改为120秒）
    establishedData := map[string]interface{}{
        "socket_id":        socketID,
        "activity_timeout": 120, // v8.3.0推荐值
    }
    client.SendPusherEvent("pusher:connection_established", "", establishedData)
    
    // 心跳超时调整为150秒（120+30）
    go client.heartbeatMonitor(150 * time.Second)
}
```

#### 6. 调整Redis操作
**文件**: [redis.go](redis.go)

**改造内容**:
- 将所有 `Topic` 相关命名改为 `Channel`：
  - `Topics` → `Channels`
  - `Topic2ClientId:{topic}` → `Channel2SocketId:{channel}`
  - `ClientId2Topic:{id}` → `SocketId2Channel:{id}`
  - `Topic2ServerName:{topic}` → `Channel2ServerName:{channel}`
- 重命名函数：
  - `JoinTopic4Redis` → `JoinChannel4Redis`
  - `LeaveTopic4Redis` → `LeaveChannel4Redis`
- 更新 `ClientId2HeartbeatTime` 为 `SocketId2HeartbeatTime`

#### 7. 更新Pub/Sub消息格式
**文件**: [pubsub.go](pubsub.go)

**改造内容**:
- 修改 `SubscribeRedis` 中的消息解析，适配Pusher格式
- 修改 `PublishToRedis` 相关函数，发送Pusher格式消息
- 确保跨服务器消息传递时保持 `event` 和 `channel` 字段
- ⚠️ 确保data字段为JSON字符串

#### 8. 调整ClientManager
**文件**: [client_manager.go](client_manager.go)

**改造内容**:
- 将 `TopicBroadcast` 改为 `ChannelBroadcast`
- 修改 `SendToTopic` 为 `SendToChannel`
- 新增 `GetClientBySocketID` 方法（用于认证验证）
- 更新广播逻辑，使用Pusher消息格式

**验证步骤**:
- 使用 `websocat` 工具连接：`websocat "ws://localhost:8000/ws?protocol=7"`
- 验证收到 `pusher:connection_established` 事件（data为JSON字符串）
- 发送 `pusher:subscribe` 订阅Public Channel（如 `my-channel`）
- 验证收到 `pusher:subscription_succeeded`（data为 `"{}"`）
- 发送 `pusher:ping`，验证收到 `pusher:pong`
- 通过HTTP API向 `my-channel` 发送消息，验证接收

---

### Phase 2: Private Channel + 认证（P1）- 4天

#### 9. 新增认证模块
**文件**: [auth.go](auth.go) - **新建**

**实现内容**:
- `ValidateChannelAuth` 函数：
  - 解析 `auth` 字符串（格式：`{app_key}:{signature}`）
  - 使用HMAC-SHA256验证签名
  - 签名内容：`{socket_id}:{channel_name}`（Private）
  - 或 `{socket_id}:{channel_name}:{channel_data}`（Presence）
  - ⚠️ 使用 `hmac.Equal` 或 `subtle.ConstantTimeCompare` 防止时序攻击
  - ⚠️ 可选参数：验证socket_id来自当前连接（防重放攻击）
- `GenerateAuthSignature` 函数（供测试/HTTP端点使用）
  - ⚠️ v8.3.0安全增强：添加时间戳（5分钟有效期）
- 从配置文件读取 `app_key` 和 `app_secret`

**代码示例**:
```go
// ⚠️ v8.3.0安全要求：验证socket_id防止重放攻击
func ValidateChannelAuth(socketID, channel, auth, channelData string, validateSocketID bool) error {
    parts := strings.Split(auth, ":")
    if len(parts) != 2 {
        return errors.New("invalid auth format")
    }
    
    appKey := parts[0]
    signature := parts[1]
    
    // 可选：验证socket_id来自当前连接
    if validateSocketID {
        if !isValidSocketID(socketID) {
            return errors.New("invalid socket_id")
        }
    }
    
    // 构建签名字符串
    var stringToSign string
    if channelData != "" {
        stringToSign = fmt.Sprintf("%s:%s:%s", socketID, channel, channelData)
    } else {
        stringToSign = fmt.Sprintf("%s:%s", socketID, channel)
    }
    
    // HMAC-SHA256验证（⚠️ 使用constant time比较）
    appSecret := getAppSecret()
    expectedSig := generateHMAC(stringToSign, appSecret)
    
    if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
        return errors.New("invalid signature")
    }
    
    return nil
}
```

#### 10. 新增Channel类型管理
**文件**: [channel.go](channel.go) - **新建**

**实现内容**:
- 定义 `ChannelType` 常量：
  - `ChannelTypePublic`
  - `ChannelTypePrivate`
  - `ChannelTypePresence`
- `GetChannelType` 函数（根据前缀返回类型）
- `IsPrivateChannel` 辅助函数
- `IsPresenceChannel` 辅助函数

#### 11. 更新SubscribeController（支持Private）
**文件**: [controller.go](controller.go)

**改造内容**:
- 检查channel类型，如果是 `private-*`：
  - 要求 `data.auth` 字段存在
  - 调用 `ValidateChannelAuth` 验证签名
  - ⚠️ 验证失败发送 `pusher:subscription_error`（v8.3.0格式：含type/status）
- 验证成功后正常订阅
- `presence-*` 暂时返回错误（Phase 3实现）

#### 12. 新增HTTP认证端点
**位置**: 控制器层（需在路由层添加）

**实现内容**:
- 路径：`POST /api/system/pusher/auth`
- 接收参数：`socket_id`，`channel_name`（支持Form和JSON）
- ⚠️ v8.3.0安全要求：验证socket_id来自当前用户的连接
- 验证用户身份（使用现有session/token机制）
- 生成并返回：`{"auth": "{app_key}:{signature}"}`
- 将控制器添加到路由配置

**代码示例**:
```go
// 控制器：POST /api/system/pusher/auth
func PusherAuthController(r *ghttp.Request) {
    ctx := r.GetCtx()
    
    socketID := r.Get("socket_id").String()
    channelName := r.Get("channel_name").String()
    
    if socketID == "" || channelName == "" {
        r.Response.WriteStatus(400)
        r.Response.WriteJson(g.Map{"error": "Missing parameters"})
        return
    }
    
    // ⚠️ v8.3.0安全要求：验证socket_id来自当前用户
    client := clientManager.GetClientBySocketID(socketID)
    if client == nil || client.SessionID != getCurrentSessionID(r) {
        r.Response.WriteStatus(403)
        r.Response.WriteJson(g.Map{"error": "Invalid socket_id"})
        return
    }
    
    // 验证用户身份
    user := getCurrentUser(r)
    if user == nil {
        r.Response.WriteStatus(401)
        r.Response.WriteJson(g.Map{"error": "Unauthorized"})
        return
    }
    
    // 生成认证签名
    auth, _ := GenerateAuthSignature(socketID, channelName, "")
    r.Response.WriteJson(g.Map{"auth": auth})
}
```

#### 13. 配置文件更新
**文件**: `manifest/config/config.yaml`

**新增配置节**:
```yaml
pusher:
  appKey: "your-app-key"
  appSecret: "your-app-secret"
  activityTimeout: 120        # ⚠️ v8.3.0推荐值（秒）
  heartbeatCheckInterval: 150 # ⚠️ 修正为150秒（120+30）
  presenceGracePeriod: 30     # ⚠️ 新增：成员断线grace period（秒）
  clientEventRateLimit: 10    # ⚠️ 新增：客户端事件速率（条/秒）
  maxChannelsPerConnection: 100 # ⚠️ 新增：单连接最大订阅数
  enableUserAuth: false       # ⚠️ 新增：暂不支持v8.0 User Authentication
```

**验证步骤**:
- 订阅 `private-test-channel`，不带 `auth` → 收到 `pusher:subscription_error`（v8.3.0格式）
- 请求 `/api/system/pusher/auth`，获取签名
- 带签名订阅 → 收到 `pusher:subscription_succeeded`
- 验证只有已认证的客户端能接收私有频道消息
- 测试socket_id防重放验证

---

### Phase 3: Presence Channel + 成员管理（P1/P2）- 4天

#### 14. 新增Presence管理模块
**文件**: [presence.go](presence.go) - **新建**

**实现内容**:
- 定义 `PresenceMember` 结构：
  - `user_id` string
  - `user_info` map[string]interface{}
- `AddPresenceMember` 函数（添加成员到Redis Hash）
- `RemovePresenceMember` 函数（从Redis Hash删除）
- `GetPresenceMembers` 函数（获取全部成员）
- `FormatPresenceData` 函数（生成Pusher格式的成员列表）：
  - ⚠️ v8.3.0格式要求：hash字段只包含user_info，不含user_id

**代码示例**:
```go
// ⚠️ v8.3.0格式要求：hash字段只包含user_info
func FormatPresenceData(members map[string]interface{}) string {
    ids := make([]string, 0, len(members))
    hash := make(map[string]interface{})
    
    for userID, userInfo := range members {
        ids = append(ids, userID)
        hash[userID] = userInfo // ⚠️ 只存储user_info，不包含user_id
    }
    
    data := map[string]interface{}{
        "presence": map[string]interface{}{
            "count": len(members),
            "ids":   ids,
            "hash":  hash,
        },
    }
    
    jsonData, _ := json.Marshal(data)
    return string(jsonData)
}
```

#### 15. 扩展Redis存储
**文件**: [redis.go](redis.go)

**新增内容**:
- Redis Key：`PresenceChannel:{channel}` Hash存储（`user_id` → `user_info JSON`）
- Redis Key：`PresenceDisconnect:{socket_id}` String存储（断线时间戳，用于Grace Period）
- 新增函数：
  - `AddPresenceMember4Redis(ctx, channel, userId, userInfo)`
  - `RemovePresenceMember4Redis(ctx, channel, userId)`
  - `GetPresenceMembers4Redis(ctx, channel)` 返回 `map[string]interface{}`
  - `GetPresenceCount4Redis(ctx, channel)` 返回成员数量
  - `MarkPresenceDisconnect4Redis(ctx, socketId)` 标记断线
  - `ClearPresenceDisconnect4Redis(ctx, socketId)` 清除断线标记

#### 16. 更新SubscribeController（支持Presence）
**文件**: [controller.go](controller.go)

**改造内容**:
- 检测 `presence-*` channel：
  - 验证 `auth` 签名（包含 `channel_data`）
  - 解析 `data.channel_data`（JSON字符串，包含 `user_id` 和 `user_info`）
  - 调用 `AddPresenceMember4Redis` 添加成员
  - 获取所有成员列表
  - 发送 `pusher:subscription_succeeded`，data包含完整成员列表（⚠️ v8.3.0格式）
  - 向channel内其他成员广播 `pusher:member_added` 事件

#### 17. 更新UnsubscribeController（Presence支持）
**文件**: [controller.go](controller.go)

**改造内容**:
- 如果是 `presence-*` channel：
  - 从 `channel_data` 或客户端记录中获取 `user_id`
  - ⚠️ 标记断线但延迟删除（Grace Period 30秒）
  - 30秒后如未重连，调用 `RemovePresenceMember4Redis`
  - 向channel内其他成员广播 `pusher:member_removed` 事件

#### 18. 客户端断开时清理
**文件**: [init.go](init.go) 或 [client.go](client.go)

**改造内容**:
- 在客户端 `Close` 方法中：
  - 标记为断开 `IsDisconnected = true`
  - 立即清理非Presence订阅
  - ⚠️ Presence Channel延迟30秒清理（Grace Period）
  - 30秒内重连则取消清理

**代码示例**:
```go
func (c *Client) Close() {
    c.IsDisconnected = true
    
    // 立即清理非Presence订阅
    for _, channel := range c.Channels {
        if !strings.HasPrefix(channel, "presence-") {
            LeaveChannel4Redis(ctx, channel, c.SocketID)
        }
    }
    
    // ⚠️ Presence Channel延迟清理（v8.3.0 Grace Period）
    time.AfterFunc(30*time.Second, func() {
        if c.IsDisconnected { // 30秒内未重连
            for _, channel := range c.Channels {
                if strings.HasPrefix(channel, "presence-") {
                    RemovePresenceMember4Redis(ctx, channel, c.UserID)
                    broadcastMemberRemoved(channel, c.UserID)
                }
            }
        }
    })
}
```

#### 19. 跨服务器成员同步
**文件**: [pubsub.go](pubsub.go)

**改造内容**:
- Redis Pub/Sub消息新增类型：`presence_join` 和 `presence_leave`
- 当本地服务器成员加入/离开时，发布消息到Redis
- 其他服务器接收到消息后，向本地订阅该channel的客户端广播

#### 20. 更新HTTP认证端点
**改造内容**:
- 对于 `presence-*` channel：
  - 返回 `auth` 和 `channel_data`
  - `channel_data` 格式：`{"user_id":"123", "user_info":{"name":"Alice"}}`
  - 签名时包含 `channel_data`

**验证步骤**:
- 两个客户端订阅 `presence-lobby`
- 验证第二个客户端收到的 `pusher:subscription_succeeded` 包含第一个成员信息（v8.3.0格式）
- 验证第一个客户端收到 `pusher:member_added` 事件
- 断开一个客户端，验证30秒后另一个收到 `pusher:member_removed`
- 测试30秒内重连，验证成员不被移除
- 多服务器环境测试成员同步

---

### Phase 4: Client Events + 优化（P2/P3）- 3天

#### 21. 新增Client Event处理
**文件**: [client_events.go](client_events.go) - **新建** 或在 [controller.go](controller.go) 中添加

**实现内容**:
- `ClientEventController` 函数：
  - 检查 `event` 是否以 `client-` 开头
  - ⚠️ v8.3.0要求：验证事件名长度（最大50字节）
  - 检查 `channel` 是否为 `private-*` 或 `presence-*`（Public不允许）
  - ⚠️ v8.3.0要求：速率限制（10条/秒）
  - 原样转发消息到channel内其他客户端（不发送给自己）
  - 错误情况发送 `pusher:error`（code 4301）

**代码示例**:
```go
func ClientEventController(ctx context.Context, client *Client, req *PusherRequest) {
    // 验证事件名
    if !strings.HasPrefix(req.Event, "client-") {
        client.SendError("Invalid client event name", CodeClientEventForbidden)
        return
    }
    
    // ⚠️ v8.3.0要求：验证事件名长度（最大50字节）
    if len(req.Event) > 50 {
        client.SendError("Event name too long", CodeClientEventForbidden)
        return
    }
    
    // 验证Channel类型
    channelType := GetChannelType(req.Channel)
    if channelType == ChannelTypePublic {
        client.SendError("Client events not allowed on public channels", CodeClientEventForbidden)
        return
    }
    
    // ⚠️ v8.3.0要求：速率限制（10条/秒）
    if !rateLimiter.AllowClientEvent(client.SocketID) {
        client.SendError("Rate limit exceeded", CodeClientEventForbidden)
        return
    }
    
    // 转发给其他客户端
    BroadcastToChannelExcept(req.Channel, client.SocketID, &PusherResponse{
        Event:   req.Event,
        Channel: req.Channel,
        Data:    req.Data,
    })
}
```

#### 21.1 新增速率限制模块
**文件**: [rate_limit.go](rate_limit.go) - **新建**

**实现内容**:
- `RateLimiter` 结构体（Token Bucket算法）
- `AllowClientEvent` 方法：检查是否允许发送（10条/秒）
- 自动清理过期bucket

**代码示例**:
```go
// ⚠️ v8.3.0要求：Client Events速率限制（10条/秒）
type RateLimiter struct {
    mu      sync.Mutex
    buckets map[string]*tokenBucket
}

type tokenBucket struct {
    tokens     int
    lastRefill time.Time
}

func (rl *RateLimiter) AllowClientEvent(socketID string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    
    bucket, exists := rl.buckets[socketID]
    if !exists {
        bucket = &tokenBucket{tokens: 10, lastRefill: time.Now()}
        rl.buckets[socketID] = bucket
    }
    
    // 每秒补充令牌
    now := time.Now()
    elapsed := now.Sub(bucket.lastRefill).Seconds()
    if elapsed >= 1 {
        bucket.tokens = 10
        bucket.lastRefill = now
    }
    
    if bucket.tokens > 0 {
        bucket.tokens--
        return true
    }
    
    return false
}
```

#### 22. 路由器注册Client Events
**文件**: [router.go](router.go)

**改造内容**:
- 在 `ProcessData` 中优先检查事件前缀
- 如果以 `client-` 开头，直接路由到 `ClientEventController`

#### 23. 错误处理完善
**涉及**: 所有相关文件

**改造内容**:
- 统一错误响应格式：
  ```json
  {
    "event": "pusher:error",
    "data": "{\"message\":\"...\",\"code\":4xxx}"
  }
  ```
- 订阅错误使用专用格式（v8.3.0）：
  ```json
  {
    "event": "pusher:subscription_error",
    "channel": "private-test",
    "data": "{\"type\":\"AuthError\",\"error\":\"...\",\"status\":401}"
  }
  ```
- 定义完整错误码常量（见步骤1）

#### 24. 性能优化
**改造内容**:
- [model.go](model.go)：使用 `sync.Pool` 复用消息对象
- [client.go](client.go)：批量发送消息（减少写锁开销）
- [presence.go](presence.go)：成员列表缓存（TTL 5秒）
- [pubsub.go](pubsub.go)：批量处理成员事件（100ms窗口）

#### 25. 日志增强
**文件**: [log.go](log.go)

**改造内容**:
- 记录关键事件：
  - 连接建立（含socket_id）
  - 订阅/退订（含channel）
  - 认证失败
  - 速率限制触发
- 使用结构化日志（JSON格式）
- 增加性能指标：消息延迟、连接数、channel订阅数

#### 26. 单元测试（可选但推荐）
**新建文件**:
- `auth_test.go` - 签名生成/验证测试
- `channel_test.go` - channel类型识别测试
- `presence_test.go` - 成员管理逻辑测试
- `rate_limit_test.go` - 速率限制测试

**测试内容**:
- 覆盖核心逻辑单元
- 集成测试：使用 `gorilla/websocket` 客户端模拟完整流程

**验证步骤**:
- 在 `private-chat` channel中，客户端A发送 `client-message` 事件
- 验证客户端B收到，但A不收到（不回显）
- 在Public Channel中尝试发送Client Event → 收到错误（code 4301）
- 测试速率限制：连续发送11条Client Event，第11条应返回错误
- 压力测试：1000并发连接，每秒100条消息，监控内存和CPU
- 测试Presence Channel在高频加入/离开时的性能

---

## 关键文件清单

### 需修改（按优先级）
1. 🔴 [model.go](model.go) - 消息结构定义（v8.3.0格式）
2. 🔴 [router.go](router.go) - 事件路由
3. 🔴 [controller.go](controller.go) - 控制器逻辑
4. 🔴 [init.go](init.go) - 连接初始化（协议验证）
5. 🟡 [client.go](client.go) - 客户端结构（新增字段）
6. 🟡 [redis.go](redis.go) - Redis操作（Grace Period）
7. 🟡 [pubsub.go](pubsub.go) - 消息订阅
8. 🟢 [client_manager.go](client_manager.go) - 管理器

### 需新增
1. 🟡 [auth.go](auth.go) - 认证签名（防重放攻击）
2. 🟡 [channel.go](channel.go) - Channel类型
3. 🟢 [presence.go](presence.go) - Presence管理（v8.3.0格式）
4. 🟢 [client_events.go](client_events.go) - 客户端事件
5. 🟢 [rate_limit.go](rate_limit.go) - 速率限制（v8.3.0要求）

### 配置和路由
- `../../../manifest/config/config.yaml` - 添加Pusher配置（v8.3.0参数）
- `../../router/router.go` - 注册 `/pusher/auth` 端点

---

## 验证方案

### 开发阶段测试
1. **命令行测试**（使用 `websocat`）
   ```bash
   websocat "ws://localhost:8000/ws?protocol=7"
   ```
2. **单元测试**（Go test覆盖核心逻辑）
3. **集成测试**（编写脚本模拟多客户端场景）

### 最终验证

#### 1. 使用Pusher官方JS SDK测试（v8.3.0）
```javascript
// ⚠️ 关键：v8.3.0必须的配置
const pusher = new Pusher('your-app-key', {
  wsHost: 'localhost',
  wsPort: 8000,
  forceTLS: false,           // ⚠️ 开发环境必须设置
  enabledTransports: ['ws'],
  disableStats: true,        // ⚠️ 关闭统计请求
  
  channelAuthorization: {    // ⚠️ v8.x新格式
    endpoint: '/api/system/pusher/auth',
    transport: 'ajax'
  },
  
  activityTimeout: 120000,   // 120秒（毫秒）
  pongTimeout: 30000         // 30秒
});

// 测试Public Channel
const publicChannel = pusher.subscribe('my-channel');
publicChannel.bind('test-event', data => {
  console.log('Public:', data);
});

// 测试Private Channel
const privateChannel = pusher.subscribe('private-chat');
privateChannel.bind('message', data => {
  console.log('Private:', data);
});

// 测试Client Event
privateChannel.trigger('client-typing', {user: 'Alice'});

// 测试Presence Channel
const presenceChannel = pusher.subscribe('presence-lobby');
presenceChannel.bind('pusher:subscription_succeeded', members => {
  console.log('Members:', members);
  console.log('Count:', members.count);
  console.log('IDs:', members.ids);
  console.log('Hash:', members.hash);
});

presenceChannel.bind('pusher:member_added', member => {
  console.log('Member added:', member);
});

presenceChannel.bind('pusher:member_removed', member => {
  console.log('Member removed:', member);
});
```

#### 2. 功能检查清单（v8.3.0兼容性）
- ✅ 连接时带`?protocol=7`参数
- ✅ 连接建立收到 `socket_id`（data为JSON字符串）
- ✅ Public Channel订阅和消息接收
- ✅ Private Channel认证流程（socket_id验证）
- ✅ Presence Channel成员列表（v8.3.0格式：hash只含user_info）
- ✅ Presence成员事件（member_added/removed）
- ✅ Client Event转发（仅私有频道，不回显）
- ✅ Client Event速率限制（11条/秒时报错）
- ✅ subscription_error包含type/status字段
- ✅ 心跳120秒内有响应
- ✅ Presence成员30秒Grace Period
- ✅ 多服务器环境消息同步
- ✅ 异常断开时成员清理

#### 3. 性能指标
- 1000并发连接稳定运行
- 消息延迟 < 50ms（P99）
- Presence成员列表查询 < 10ms
- 内存使用合理（< 100MB/1000连接）

---

## 技术注意事项

### v8.3.0关键要求
1. **⚠️ Pusher data字段序列化**：
   - 服务器→客户端：`data` 必须是JSON字符串
   - 客户端→服务器：`data` 可以是对象或字符串
   - 错误示例：`{"event":"...", "data": {"key":"value"}}`  ❌
   - 正确示例：`{"event":"...", "data": "{\"key\":\"value\"}"}`  ✅

2. **⚠️ socket_id格式**：
   - 格式：`{serverName}.{timestamp}{random}`
   - 确保全局唯一
   - 示例：`server1.17093456781234`

3. **⚠️ 签名验证时序攻击**：
   - 使用 `hmac.Equal` 或 `subtle.ConstantTimeCompare` 比较签名
   - 禁止使用 `==` 直接比较

4. **⚠️ Presence hash结构**：
   - `hash` 字段只存储 `user_info`，不包含 `user_id`
   - 正确：`{"u1": {"name":"Alice"}}`
   - 错误：`{"u1": {"user_id":"u1", "name":"Alice"}}`

5. **⚠️ 心跳超时**：
   - `activity_timeout`: 120秒（客户端110秒发ping）
   - 服务器超时：150秒（120+30）
   - 客户端pong超时：30秒

6. **⚠️ Client Event安全**：
   - 仅Private/Presence Channel允许
   - 事件名最大50字节
   - 速率限制10条/秒
   - 不回显给发送者

7. **⚠️ Grace Period**：
   - Presence成员断线后30秒再移除
   - 防止短暂网络波动导致频繁member_added/removed

8. **⚠️ 协议版本验证**：
   - 连接时必须带 `?protocol=7`
   - 其他版本返回4005错误

---

## Pusher协议关键点总结

### 消息格式对比

#### 现有格式（废弃）
```json
{
  "be": "bindEvent",
  "e": "subscribe", 
  "d": {"topic": "test"},
  "r": "req123"
}
```

#### Pusher格式（新标准，v8.3.0）
```json
{
  "event": "pusher:subscribe",
  "data": {
    "channel": "test"
  }
}
```

### 系统事件列表
| 事件名称 | 方向 | 描述 | v8.3.0注意事项 |
|---------|------|-----|---------------|
| `pusher:connection_established` | S→C | 连接建立 | data必须是JSON字符串 |
| `pusher:ping` | C→S | 客户端心跳请求 | - |
| `pusher:pong` | S→C | 服务器心跳响应 | data为 `"{}"` |
| `pusher:subscribe` | C→S | 订阅频道 | 支持auth/channel_data |
| `pusher:subscription_succeeded` | S→C | 订阅成功 | Public: data为`"{}"`<br>Presence: 含成员列表 |
| `pusher:subscription_error` | S→C | 订阅失败 | 必须含type/status字段 |
| `pusher:unsubscribe` | C→S | 取消订阅 | 无响应事件 |
| `pusher:error` | S→C | 错误消息 | 含message/code字段 |
| `pusher:member_added` | S→C | 成员加入（Presence） | 含user_id/user_info |
| `pusher:member_removed` | S→C | 成员离开（Presence） | 只含user_id |

### Channel类型
| Channel类型 | 格式 | 用途 | 认证要求 | v8.3.0特性 |
|------------|------|------|---------|-----------|
| **Public** | `my-channel` | 公开频道 | 无 | 不支持Client Events |
| **Private** | `private-*` | 私有频道 | HMAC签名 | 支持Client Events |
| **Presence** | `presence-*` | 在线状态频道 | HMAC签名 + 成员信息 | 30秒Grace Period |

---

## v8.3.0兼容性说明

### ✅ 完全支持
- Pusher Protocol v7
- pusher-js v8.3.0客户端库
- Public/Private/Presence Channels
- Client Events（含速率限制）
- HMAC-SHA256认证（含防重放）
- 标准错误格式
- Grace Period机制

### ❌ 暂不支持（v8.0+新特性）
- **User Authentication**（v8.0新增）
  - 独立的用户认证（`pusher:signin`）
  - 建议使用Presence Channel代替
- **Webhook签名验证**
  - 服务器→应用的事件推送
- **Channel查询API**
  - `GET /channels/{channel}` 返回成员列表
- **统计端点**
  - `POST /pusher/stats`

### 📝 使用建议
1. 开发环境必须设置 `forceTLS: false`
2. 建议设置 `disableStats: true` 避免额外请求
3. 使用 `channelAuthorization` 配置（v8.x新格式）
4. 需要用户认证时，使用Presence Channel实现

---

## 参考资源

- **Pusher协议文档**: https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol/
- **pusher-js v8 Changelog**: https://github.com/pusher/pusher-js/blob/master/CHANGELOG.md
- **Channel认证**: https://pusher.com/docs/channels/server_api/authenticating-users/
- **Presence Channels**: https://pusher.com/docs/channels/using_channels/presence-channels/
- **错误码参考**: https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol/#error-codes
- **Pusher JS SDK**: https://github.com/pusher/pusher-js

---

## 进度追踪

### Phase 1: 核心协议改造（5天）⏳
- [ ] 1. 重构消息数据模型（v8.3.0格式）
- [ ] 2. 修改客户端结构（新增字段）
- [ ] 3. 改造事件路由器
- [ ] 4. 重写事件控制器
- [ ] 5. 修改连接初始化（协议验证）
- [ ] 6. 调整Redis操作
- [ ] 7. 更新Pub/Sub消息格式
- [ ] 8. 调整ClientManager

### Phase 2: Private Channel + 认证（4天）⏳
- [ ] 9. 新增认证模块（防重放）
- [ ] 10. 新增Channel类型管理
- [ ] 11. 更新SubscribeController（支持Private）
- [ ] 12. 新增HTTP认证端点（socket_id验证）
- [ ] 13. 配置文件更新（v8.3.0参数）

### Phase 3: Presence Channel + 成员管理（4天）⏳
- [ ] 14. 新增Presence管理模块（v8.3.0格式）
- [ ] 15. 扩展Redis存储（Grace Period）
- [ ] 16. 更新SubscribeController（支持Presence）
- [ ] 17. 更新UnsubscribeController（Presence支持）
- [ ] 18. 客户端断开时清理（30秒延迟）
- [ ] 19. 跨服务器成员同步
- [ ] 20. 更新HTTP认证端点（channel_data）

### Phase 4: Client Events + 优化（3天）⏳
- [ ] 21. 新增Client Event处理（速率限制）
- [ ] 21.1 新增速率限制模块
- [ ] 22. 路由器注册Client Events
- [ ] 23. 错误处理完善（v8.3.0格式）
- [ ] 24. 性能优化
- [ ] 25. 日志增强
- [ ] 26. 单元测试

---

## 修订历史

**v2.0 (2026-02-27)**
- 基于pusher-js v8.3.0完整兼容性审查
- 修正data字段序列化规则
- 修正Presence hash结构
- 修正subscription_error格式
- 新增协议版本验证
- 新增socket_id防重放验证
- 新增Client Events速率限制
- 新增Presence Grace Period（30秒）
- 调整心跳超时为150秒
- 工期调整为14-17天

**v1.0 (2026-02-26)**
- 初始版本
