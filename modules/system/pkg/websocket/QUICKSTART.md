# 🚀 Pusher WebSocket 快速启动指南

## 一、环境准备

### 1. 确认配置文件
编辑 `hack/config.yaml` 或 `manifest/config/config.yaml`：

```yaml
pusher:
  appKey: "devinggo-app-key"                  # 客户端连接用
  appSecret: "devinggo-app-secret-change-me"  # ⚠️ 生产环境必须修改！
  activityTimeout: 120
  heartbeatCheckInterval: 60
  presenceGracePeriod: 30
  maxChannelsPerConnection: 100
  clientEventRateLimit: 10
```

### 2. 确保Redis运行
```bash
# 检查Redis是否启动
redis-cli ping
# 应返回: PONG
```

---

## 二、编译与启动

### 编译项目
```powershell
cd e:\code\devinggo-light
go build -o devinggo.exe .\main.go
```

### 启动服务
```powershell
.\devinggo.exe
```

**成功标志**：
```
2026-02-28 10:00:00.000 [INFO] Server started on :8070
2026-02-28 10:00:00.000 [INFO] WebSocket endpoint: /system/ws
```

---

## 三、快速测试

### 方法1：使用浏览器测试客户端（推荐）

1. 启动服务后，浏览器打开：
   ```
   http://localhost:8070/pusher-test.html
   ```

2. 点击"连接"按钮

3. 测试功能：
   - ✅ Public频道订阅
   - ✅ Private频道订阅 + Client Event
   - ✅ Presence频道成员管理
   - ✅ 速率限制测试

---

### 方法2：使用JavaScript客户端

#### 安装pusher-js
```bash
npm install pusher-js@8.3.0
```

#### 基础连接
```javascript
const Pusher = require('pusher-js');

const pusher = new Pusher('devinggo-app-key', {
  wsHost: 'localhost',
  wsPort: 8070,
  forceTLS: false,
  enabledTransports: ['ws'],
  cluster: 'mt1',  // 自托管服务器必需（虚拟cluster名称）
  disableStats: true,  // 禁用统计信息
  authEndpoint: 'http://localhost:8070/system/pusher/auth',
  auth: {
    headers: {
      'Authorization': 'Bearer YOUR_JWT_TOKEN'  // 如需认证
    }
  }
});

pusher.connection.bind('connected', function() {
  console.log('✅ 连接成功！Socket ID:', pusher.connection.socket_id);
});
```

#### Public频道订阅
```javascript
const publicChannel = pusher.subscribe('chat-room');

publicChannel.bind('pusher:subscription_succeeded', function() {
  console.log('✅ Public频道订阅成功');
});

publicChannel.bind('new-message', function(data) {
  console.log('📨 收到消息:', data);
});
```

#### Private频道 + Client Event
```javascript
const privateChannel = pusher.subscribe('private-user-123');

privateChannel.bind('pusher:subscription_succeeded', function() {
  console.log('✅ Private频道订阅成功');
  
  // 发送Client Event
  privateChannel.trigger('client-typing', {user: 'Alice'});
});

privateChannel.bind('client-typing', function(data) {
  console.log('👤 用户正在输入:', data.user);
});
```

#### Presence频道
```javascript
const presenceChannel = pusher.subscribe('presence-lobby');

presenceChannel.bind('pusher:subscription_succeeded', function(members) {
  console.log('✅ Presence频道订阅成功');
  console.log('👥 当前在线:', members.count, '人');
  console.log('👥 成员列表:', members.members);
});

presenceChannel.bind('pusher:member_added', function(member) {
  console.log('👤 新成员加入:', member.id, member.info);
});

presenceChannel.bind('pusher:member_removed', function(member) {
  console.log('👤 成员离开:', member.id);
});
```

---

### 方法3：使用curl测试HTTP认证端点

```bash
# Private频道认证
curl -X POST http://localhost:8070/system/pusher/auth \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "socket_id": "123456.789",
    "channel_name": "private-test"
  }'

# 期望返回:
# {"auth":"devinggo-app-key:HMAC_SHA256_SIGNATURE"}
```

```bash
# Presence频道认证
curl -X POST http://localhost:8070/system/pusher/auth \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "socket_id": "123456.789",
    "channel_name": "presence-lobby"
  }'

# 期望返回:
# {
#   "auth": "devinggo-app-key:HMAC_SHA256_SIGNATURE",
#   "channel_data": "{\"user_id\":\"123\",\"user_info\":{\"name\":\"User\"}}"
# }
```

---

## 四、常见问题

### 1. 连接失败：WebSocket connection to 'ws://...' failed

**原因**: 服务未启动或端口错误

**解决**:
```bash
# 检查服务是否运行
netstat -ano | findstr :8070

# 确认WebSocket端点
curl http://localhost:8070/system/ws
# 应返回: 400 Bad Request (Upgrade required)
```

---

### 2. Private频道订阅失败：pusher:subscription_error

**原因**: HTTP认证端点返回错误或签名验证失败

**检查**:
1. 认证端点是否可访问：
   ```bash
   curl -X POST http://localhost:8070/system/pusher/auth
   ```

2. JWT Token是否有效（如需认证）

3. Redis是否运行正常

---

### 3. Client Event发送失败：错误码4301

**可能原因**:
- ✅ Public频道禁止Client Event（正常行为）
- ❌ 事件名不是 `client-` 前缀
- ❌ 事件名超过50字节
- ❌ 未订阅该频道

**验证**:
```javascript
// ❌ 错误：Public频道
publicChannel.trigger('client-test', {});  // 返回4301

// ✅ 正确：Private频道
privateChannel.trigger('client-test', {});  // 成功
```

---

### 4. 超过速率限制

**现象**: 快速发送Client Event时，部分消息失败

**原因**: 默认限制10 events/sec

**解决**:
1. 调整配置（config.yaml）：
   ```yaml
   pusher:
     clientEventRateLimit: 20  # 提高到20/秒
   ```

2. 使用节流/防抖：
   ```javascript
   // 防抖：500ms内只发送1次
   const debouncedTrigger = _.debounce(() => {
     channel.trigger('client-typing', {});
   }, 500);
   ```

---

### 5. Presence成员列表不更新

**检查**:
1. Grace Period生效（30秒）：
   - 断线后30秒内重连不会触发离开事件

2. 缓存问题：
   - Presence缓存TTL为5秒
   - 成员加入/离开会自动失效缓存

3. Redis连接：
   ```bash
   redis-cli
   > KEYS ws:presence:*
   > HGETALL ws:presence:channel:presence-lobby
   ```

---

## 五、生产环境部署

### 1. 安全配置

#### ⚠️ 必须修改appSecret
```yaml
pusher:
  appSecret: "$(openssl rand -hex 32)"  # 生成64字符随机密钥
```

#### 启用TLS/SSL
```yaml
server:
  address: ":443"
  httpsAddr: ":443"
  httpsCertPath: "/path/to/cert.pem"
  httpsKeyPath: "/path/to/key.pem"
```

#### Redis密码
```yaml
redis:
  default:
    address: "127.0.0.1:6379"
    pass: "your-redis-password"
```

---

### 2. 性能优化

#### 调整连接限制
```yaml
pusher:
  maxChannelsPerConnection: 200      # 根据业务调整
  heartbeatCheckInterval: 30         # 减少检查间隔（提高响应）
  presenceGracePeriod: 60            # 增加容错时间（降低误判）
```

#### Redis连接池
```yaml
redis:
  default:
    maxIdle: 50
    maxActive: 100
    idleTimeout: 300
```

---

### 3. 监控告警

#### 关键指标
- 当前连接数
- 频道订阅总数
- 消息吞吐量（/秒）
- Client Event速率限制触发次数
- Redis连接池使用率

#### Prometheus Metrics（建议）
```go
// 在代码中添加metrics
var (
    wsConnections = prometheus.NewGauge(...)
    wsMessages = prometheus.NewCounter(...)
    wsLatency = prometheus.NewHistogram(...)
)
```

---

### 4. 水平扩展

#### 多实例部署
```bash
# 实例1
.\devinggo.exe --port 8001

# 实例2
.\devinggo.exe --port 8002
```

#### Nginx负载均衡
```nginx
upstream websocket_backend {
    least_conn;
    server 127.0.0.1:8001;
    server 127.0.0.1:8002;
}

server {
    listen 8070;
    
    location /api/system/ws {
        proxy_pass http://websocket_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        
        # Sticky session (IP hash)
        ip_hash;
    }
}
```

---

## 六、更多资源

### 文档
- 📖 [完整实现总结](IMPLEMENTATION_SUMMARY.md)
- 📋 [测试验证清单](TEST_VALIDATION.md)
- 📘 [Pusher Protocol官方文档](https://pusher.com/docs/channels/library_auth_reference/pusher-websockets-protocol/)

### 示例代码
- 🌐 [浏览器测试客户端](../../../resource/public/pusher-test.html)
- 📦 [pusher-js官方示例](https://github.com/pusher/pusher-js/tree/master/example)

### 工具推荐
- **调试**: wscat, websocat, DevTools Network面板
- **压测**: k6, artillery, vegeta
- **监控**: Prometheus + Grafana

---

## 🎉 完成！

现在你可以：
1. ✅ 打开 http://localhost:8070/pusher-test.html 测试
2. ✅ 集成pusher-js到你的前端应用
3. ✅ 使用Laravel Echo等框架（配置认证端点）
4. ✅ 部署到生产环境（记得改appSecret！）

**祝开发顺利！** 🚀

