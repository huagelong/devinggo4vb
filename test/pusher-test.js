// Pusher Protocol v8.3.0 测试脚本
// 使用Node.js的ws库测试WebSocket连接

const WebSocket = require('ws');
const crypto = require('crypto');

// 配置
const CONFIG = {
    wsUrl: 'ws://localhost:8070/system/ws',
    appKey: 'devinggo-app-key',
    appSecret: 'devinggo-app-secret-change-me',
    authEndpoint: 'http://localhost:8070/api/system/pusher/auth'
};

// 颜色输出
const colors = {
    reset: '\x1b[0m',
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
    console.log(`${colors[color]}${message}${colors.reset}`);
}

function logSuccess(message) {
    log(`✅ ${message}`, 'green');
}

function logError(message) {
    log(`❌ ${message}`, 'red');
}

function logInfo(message) {
    log(`ℹ️  ${message}`, 'cyan');
}

function logWarning(message) {
    log(`⚠️  ${message}`, 'yellow');
}

// 生成HMAC-SHA256签名
function generateAuth(socketId, channelName, channelData = null) {
    let stringToSign = `${socketId}:${channelName}`;
    if (channelData) {
        stringToSign += `:${channelData}`;
    }
    
    const hmac = crypto.createHmac('sha256', CONFIG.appSecret);
    hmac.update(stringToSign);
    const signature = hmac.digest('hex');
    
    return `${CONFIG.appKey}:${signature}`;
}

// 测试类
class PusherTest {
    constructor() {
        this.ws = null;
        this.socketId = null;
        this.testResults = {
            passed: 0,
            failed: 0,
            tests: []
        };
    }

    // 连接WebSocket
    connect() {
        return new Promise((resolve, reject) => {
            logInfo('正在连接到WebSocket服务器...');
            
            this.ws = new WebSocket(CONFIG.wsUrl);
            
            this.ws.on('open', () => {
                logSuccess('WebSocket连接已建立');
            });
            
            this.ws.on('message', (data) => {
                try {
                    const message = JSON.parse(data);
                    this.handleMessage(message, resolve);
                } catch (e) {
                    logError(`解析消息失败: ${e.message}`);
                }
            });
            
            this.ws.on('error', (error) => {
                logError(`WebSocket错误: ${error.message}`);
                reject(error);
            });
            
            this.ws.on('close', () => {
                logWarning('WebSocket连接已关闭');
            });
        });
    }

    // 处理接收到的消息
    handleMessage(message, resolveConnect) {
        logInfo(`收到事件: ${message.event}`);
        
        if (message.event === 'pusher:connection_established') {
            const data = JSON.parse(message.data);
            this.socketId = data.socket_id;
            logSuccess(`连接建立成功！Socket ID: ${this.socketId}`);
            this.recordTest('连接建立', true);
            
            if (resolveConnect) {
                resolveConnect();
            }
        } else if (message.event === 'pusher:subscription_succeeded') {
            logSuccess(`频道订阅成功: ${message.channel}`);
            this.recordTest(`订阅频道: ${message.channel}`, true);
            
            if (message.channel.startsWith('presence-')) {
                const data = JSON.parse(message.data);
                logInfo(`  成员数: ${data.presence.count}`);
                logInfo(`  成员列表: ${JSON.stringify(data.presence.ids)}`);
            }
        } else if (message.event === 'pusher:subscription_error') {
            logError(`频道订阅失败: ${message.channel}`);
            const data = JSON.parse(message.data);
            logError(`  错误: ${data.error} (状态: ${data.status})`);
            this.recordTest(`订阅频道: ${message.channel}`, false);
        } else if (message.event === 'pusher:error') {
            const data = JSON.parse(message.data);
            logError(`服务器错误: ${data.message} (代码: ${data.code})`);
        } else if (message.event === 'pusher:pong') {
            logSuccess('收到Pong响应');
            this.recordTest('心跳测试', true);
        } else if (message.event.startsWith('client-')) {
            logSuccess(`收到Client Event: ${message.event}`);
            logInfo(`  数据: ${message.data}`);
            this.recordTest('Client Event接收', true);
        } else if (message.event === 'pusher:member_added') {
            const data = JSON.parse(message.data);
            logSuccess(`新成员加入: ${data.user_id}`);
        } else if (message.event === 'pusher:member_removed') {
            const data = JSON.parse(message.data);
            logSuccess(`成员离开: ${data.user_id}`);
        }
    }

    // 发送消息
    send(message) {
        this.ws.send(JSON.stringify(message));
    }

    // 订阅Public频道
    subscribePublic(channelName) {
        logInfo(`订阅Public频道: ${channelName}`);
        this.send({
            event: 'pusher:subscribe',
            data: {
                channel: channelName
            }
        });
    }

    // 订阅Private频道
    subscribePrivate(channelName) {
        logInfo(`订阅Private频道: ${channelName}`);
        const auth = generateAuth(this.socketId, channelName);
        
        this.send({
            event: 'pusher:subscribe',
            data: {
                channel: channelName,
                auth: auth
            }
        });
    }

    // 订阅Presence频道
    subscribePresence(channelName, userId, userInfo) {
        logInfo(`订阅Presence频道: ${channelName}`);
        const channelData = JSON.stringify({
            user_id: userId,
            user_info: userInfo
        });
        const auth = generateAuth(this.socketId, channelName, channelData);
        
        this.send({
            event: 'pusher:subscribe',
            data: {
                channel: channelName,
                auth: auth,
                channel_data: channelData
            }
        });
    }

    // 发送Client Event
    sendClientEvent(channelName, eventName, data) {
        logInfo(`发送Client Event: ${eventName} 到频道 ${channelName}`);
        this.send({
            event: eventName,
            channel: channelName,
            data: data
        });
    }

    // 发送Ping
    sendPing() {
        logInfo('发送Ping...');
        this.send({
            event: 'pusher:ping',
            data: {}
        });
    }

    // 记录测试结果
    recordTest(testName, passed) {
        this.testResults.tests.push({ name: testName, passed });
        if (passed) {
            this.testResults.passed++;
        } else {
            this.testResults.failed++;
        }
    }

    // 打印测试结果
    printResults() {
        log('\n' + '='.repeat(60), 'cyan');
        log('测试结果汇总', 'cyan');
        log('='.repeat(60), 'cyan');
        
        this.testResults.tests.forEach(test => {
            const status = test.passed ? '✅ 通过' : '❌ 失败';
            const color = test.passed ? 'green' : 'red';
            log(`  ${status}: ${test.name}`, color);
        });
        
        log('\n' + '-'.repeat(60), 'cyan');
        log(`总计: ${this.testResults.tests.length} 个测试`, 'cyan');
        log(`通过: ${this.testResults.passed} 个`, 'green');
        log(`失败: ${this.testResults.failed} 个`, 'red');
        log('='.repeat(60) + '\n', 'cyan');
    }

    // 关闭连接
    close() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// 延迟函数
function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// 主测试函数
async function runTests() {
    const test = new PusherTest();
    
    try {
        log('\n🚀 开始Pusher Protocol v8.3.0测试\n', 'yellow');
        
        // 1. 连接测试
        log('【测试1】连接建立', 'yellow');
        await test.connect();
        await sleep(1000);
        
        // 2. 心跳测试
        log('\n【测试2】心跳机制', 'yellow');
        test.sendPing();
        await sleep(1000);
        
        // 3. Public频道测试
        log('\n【测试3】Public频道订阅', 'yellow');
        test.subscribePublic('chat-room');
        await sleep(1000);
        
        // 4. Private频道测试
        log('\n【测试4】Private频道订阅', 'yellow');
        test.subscribePrivate('private-user-123');
        await sleep(1000);
        
        // 5. Client Event测试（Private频道）
        log('\n【测试5】Client Event（Private频道）', 'yellow');
        test.sendClientEvent('private-user-123', 'client-typing', { user: 'Alice' });
        await sleep(1000);
        
        // 6. Client Event测试（Public频道 - 应该失败）
        log('\n【测试6】Client Event（Public频道 - 预期失败）', 'yellow');
        logWarning('此测试应该返回错误码4301...');
        test.sendClientEvent('chat-room', 'client-test', { message: 'This should fail' });
        await sleep(1000);
        
        // 7. Presence频道测试
        log('\n【测试7】Presence频道订阅', 'yellow');
        test.subscribePresence('presence-lobby', 'user-123', { name: 'Alice', status: 'online' });
        await sleep(2000);
        
        // 8. 速率限制测试
        log('\n【测试8】速率限制测试（快速发送11条消息）', 'yellow');
        logWarning('前10条应该成功，第11条应该失败...');
        for (let i = 0; i < 11; i++) {
            test.sendClientEvent('private-user-123', 'client-ratelimit-test', { seq: i + 1 });
            await sleep(50);
        }
        await sleep(1000);
        
        // 等待所有响应
        log('\n等待所有响应处理完毕...', 'cyan');
        await sleep(3000);
        
        // 打印结果
        test.printResults();
        
        // 关闭连接
        log('关闭连接...', 'cyan');
        test.close();
        
        log('\n✅ 测试完成！', 'green');
        
    } catch (error) {
        logError(`测试失败: ${error.message}`);
        test.close();
        process.exit(1);
    }
}

// 运行测试
if (require.main === module) {
    runTests().catch(err => {
        logError(`测试异常: ${err.message}`);
        process.exit(1);
    });
}

module.exports = { PusherTest, generateAuth };
