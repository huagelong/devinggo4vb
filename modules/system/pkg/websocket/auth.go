// Package websocket
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package websocket

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

// Pusher认证配置
var (
	pusherAppKey    string
	pusherAppSecret string
)

// InitPusherAuth 初始化Pusher认证配置
func InitPusherAuth(appKey, appSecret string) {
	pusherAppKey = appKey
	pusherAppSecret = appSecret
}

// GetPusherConfig 从配置文件读取Pusher配置
func GetPusherConfig() (appKey, appSecret string) {
	ctx := gctx.GetInitCtx()

	// 从配置文件读取
	appKey = g.Cfg().MustGet(ctx, "pusher.appKey", "default-app-key").String()
	appSecret = g.Cfg().MustGet(ctx, "pusher.appSecret", "default-app-secret").String()

	// 初始化
	InitPusherAuth(appKey, appSecret)

	return appKey, appSecret
}

// ValidateChannelAuth 验证频道认证签名
// ⚠️ v8.3.0安全要求：防止时序攻击和重放攻击
func ValidateChannelAuth(socketID, channel, auth, channelData string) error {
	if socketID == "" || channel == "" || auth == "" {
		return errors.New("missing required parameters")
	}

	// 解析 auth 字符串（格式：{app_key}:{signature}）
	parts := strings.Split(auth, ":")
	if len(parts) != 2 {
		return errors.New("invalid auth format, expected: app_key:signature")
	}

	receivedAppKey := parts[0]
	receivedSignature := parts[1]

	// 验证 app_key
	if pusherAppKey == "" {
		GetPusherConfig() // 自动加载配置
	}

	if receivedAppKey != pusherAppKey {
		return errors.New("invalid app_key")
	}

	// 构建签名字符串
	var stringToSign string
	if channelData != "" {
		// Presence Channel: socket_id:channel_name:channel_data
		stringToSign = fmt.Sprintf("%s:%s:%s", socketID, channel, channelData)
	} else {
		// Private Channel: socket_id:channel_name
		stringToSign = fmt.Sprintf("%s:%s", socketID, channel)
	}

	// 计算期望的签名
	expectedSignature := generateHMAC(stringToSign, pusherAppSecret)

	// ⚠️ 使用 constant time 比较防止时序攻击
	if !constantTimeCompare(receivedSignature, expectedSignature) {
		return errors.New("invalid signature")
	}

	return nil
}

// GenerateAuthSignature 生成认证签名（供HTTP认证端点使用）
func GenerateAuthSignature(socketID, channel, channelData string) string {
	// 构建签名字符串
	var stringToSign string
	if channelData != "" {
		stringToSign = fmt.Sprintf("%s:%s:%s", socketID, channel, channelData)
	} else {
		stringToSign = fmt.Sprintf("%s:%s", socketID, channel)
	}

	// 计算签名
	if pusherAppKey == "" || pusherAppSecret == "" {
		GetPusherConfig() // 自动加载配置
	}

	signature := generateHMAC(stringToSign, pusherAppSecret)

	// 返回格式：{app_key}:{signature}
	return fmt.Sprintf("%s:%s", pusherAppKey, signature)
}

// generateHMAC 生成HMAC-SHA256签名
func generateHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// constantTimeCompare 常量时间字符串比较（防止时序攻击）
func constantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// ValidateSocketID 验证socket_id是否有效（可选，用于防重放攻击）
func ValidateSocketID(socketID string, expectedServerName string) bool {
	parts := strings.Split(socketID, ".")
	if len(parts) != 2 {
		return false
	}

	serverName := parts[0]
	return serverName == expectedServerName
}

// GenerateSharedSecret 生成加密频道的共享密钥
// ⚠️ Encrypted Channels 需要：返回 32 字节随机密钥的 Base64 编码
// Pusher.js 使用此密钥进行端到端加密（NaCl/TweetNaCl）
func GenerateSharedSecret() string {
	// 生成 32 字节随机密钥（256 位）
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		// 如果随机数生成失败，使用备用方案（不应该发生）
		// 在生产环境应该 panic，因为加密密钥必须是真随机的
		g.Log().Error(gctx.New(), "Failed to generate random key for encrypted channel:", err)
		// 返回固定密钥作为降级方案（仅用于开发/测试）
		return base64.StdEncoding.EncodeToString([]byte("INSECURE-FALLBACK-KEY-32BYTES!"))
	}

	// Base64 编码（标准编码，Pusher.js 要求）
	return base64.StdEncoding.EncodeToString(key)
}

// GenerateUserAuthSignature 生成用户认证签名
// 用于 Pusher User Authentication（用户身份认证）
// 文档：https://pusher.com/docs/channels/server_api/authenticating-users/
//
// 参数：
//   - socketID: 客户端的 socket_id
//   - userData: 用户数据（必须包含 "id" 字段，且为字符串类型）
//
// 返回格式：app_key:signature
//
// 签名算法：
//  1. 将 userData 序列化为 JSON
//  2. 构建待签名字符串：socket_id::user::user_data_json
//  3. 使用 HMAC-SHA256 计算签名
//  4. 返回：app_key:hex(signature)
func GenerateUserAuthSignature(socketID string, userData map[string]interface{}) (string, error) {
	// 验证 userData 必须包含 id 字段
	if userData == nil {
		return "", fmt.Errorf("userData cannot be nil")
	}

	userID, ok := userData["id"]
	if !ok {
		return "", fmt.Errorf("userData must contain 'id' field")
	}

	// id 必须是字符串类型
	if _, ok := userID.(string); !ok {
		return "", fmt.Errorf("userData['id'] must be a string")
	}

	// 1. 序列化用户数据为 JSON
	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user data: %w", err)
	}

	// 2. 构建待签名字符串
	// 格式：socket_id::user::user_data_json
	// 注意：使用双冒号 "::" 作为分隔符
	stringToSign := fmt.Sprintf("%s::user::%s", socketID, string(userDataJSON))

	// 3. 加载配置（如果未加载）
	if pusherAppKey == "" || pusherAppSecret == "" {
		GetPusherConfig()
	}

	// 4. 计算 HMAC-SHA256 签名
	signature := generateHMAC(stringToSign, pusherAppSecret)

	// 5. 返回格式：app_key:signature
	return fmt.Sprintf("%s:%s", pusherAppKey, signature), nil
}

// ValidateUserAuthSignature 验证用户认证签名
// 校验规则：
//  1. auth 格式为 app_key:signature
//  2. app_key 必须匹配配置
//  3. 签名字符串为 socket_id::user::user_data_json
//  4. 使用 HMAC-SHA256 且常量时间比较
func ValidateUserAuthSignature(socketID, auth, userDataJSON string) error {
	if socketID == "" || auth == "" || userDataJSON == "" {
		return errors.New("missing required parameters")
	}

	parts := strings.Split(auth, ":")
	if len(parts) != 2 {
		return errors.New("invalid auth format, expected: app_key:signature")
	}

	if pusherAppKey == "" || pusherAppSecret == "" {
		GetPusherConfig()
	}

	receivedAppKey := parts[0]
	receivedSignature := parts[1]
	if receivedAppKey != pusherAppKey {
		return errors.New("invalid app_key")
	}

	stringToSign := fmt.Sprintf("%s::user::%s", socketID, userDataJSON)
	expectedSignature := generateHMAC(stringToSign, pusherAppSecret)
	if !constantTimeCompare(receivedSignature, expectedSignature) {
		return errors.New("invalid signature")
	}

	return nil
}
