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
