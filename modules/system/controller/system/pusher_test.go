// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"devinggo/modules/system/pkg/websocket"
)

// TestSplitInfoParams 测试 info 参数分割功能
func TestSplitInfoParams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "单个参数",
			input:    "user_count",
			expected: []string{"user_count"},
		},
		{
			name:     "多个参数",
			input:    "user_count,subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "带空格的参数",
			input:    " user_count , subscription_count ",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
		{
			name:     "只有逗号",
			input:    ",,,",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitInfoParams(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitInfoParams(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("splitInfoParams(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestSplitInfo 测试频道查询的 info 参数分割
func TestSplitInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "单个参数",
			input:    "user_count",
			expected: []string{"user_count"},
		},
		{
			name:     "多个参数",
			input:    "user_count,subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "带空格",
			input:    "user_count, subscription_count",
			expected: []string{"user_count", "subscription_count"},
		},
		{
			name:     "空字符串",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitInfo(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("splitInfo(%q) length = %d, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("splitInfo(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestAbs 测试绝对值函数
func TestAbs(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{input: 10, expected: 10},
		{input: -10, expected: 10},
		{input: 0, expected: 0},
		{input: 100, expected: 100},
		{input: -100, expected: 100},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

// TestVerifyWebhookSignature 测试 Webhook 签名验证
func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"time_ms":123456,"events":[]}`)

	// 生成正确的签名
	signature := generateTestSignature(body, secret)

	// 测试正确的签名
	if !verifyWebhookSignature(body, signature, secret) {
		t.Error("verifyWebhookSignature should return true for valid signature")
	}

	// 测试错误的签名
	if verifyWebhookSignature(body, "invalid-signature", secret) {
		t.Error("verifyWebhookSignature should return false for invalid signature")
	}

	// 测试空签名
	if verifyWebhookSignature(body, "", secret) {
		t.Error("verifyWebhookSignature should return false for empty signature")
	}
}

// generateTestSignature 生成测试签名
func generateTestSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// TestChannelNameValidation 测试频道名称验证（Pusher 规范）
func TestChannelNameValidation(t *testing.T) {
	tests := []struct {
		name        string
		channelName string
		shouldPass  bool
	}{
		{"有效的公开频道", "my-channel", true},
		{"有效的私有频道", "private-my-channel", true},
		{"有效的在线状态频道", "presence-chatroom", true},
		{"有效的加密频道", "private-encrypted-secret", true},
		{"包含数字和下划线", "channel_123", true},
		{"超长频道名（201字符）", strings.Repeat("a", 201), false},
		{"包含空格", "my channel", false},
		{"包含特殊字符", "my@channel", false},
		{"空字符串", "", false},
		{"只有前缀", "private-", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := websocket.ValidateChannelName(tt.channelName)
			if tt.shouldPass && err != nil {
				t.Errorf("ValidateChannelName(%q) 应该通过但返回错误: %v", tt.channelName, err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("ValidateChannelName(%q) 应该失败但未返回错误", tt.channelName)
			}
		})
	}
}

// TestEventNameValidation 测试事件名称验证（Pusher 规范）
func TestEventNameValidation(t *testing.T) {
	tests := []struct {
		name       string
		eventName  string
		shouldPass bool
	}{
		{"普通事件名", "my-event", true},
		{"客户端事件", "client-message", true},
		{"包含点号", "user.login", true},
		{"最大长度（200字符）", strings.Repeat("a", 200), true},
		{"超长事件名（201字符）", strings.Repeat("a", 201), false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := websocket.ValidateEventName(tt.eventName)
			if tt.shouldPass && err != nil {
				t.Errorf("ValidateEventName(%q) 应该通过但返回错误: %v", tt.eventName, err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("ValidateEventName(%q) 应该失败但未返回错误", tt.eventName)
			}
		})
	}
}

// TestChannelsListValidation 测试频道列表验证（Pusher 规范：最多10个频道）
func TestChannelsListValidation(t *testing.T) {
	tests := []struct {
		name       string
		channels   []string
		shouldPass bool
	}{
		{"单个频道", []string{"channel-1"}, true},
		{"10个频道（最大限制）", []string{"ch1", "ch2", "ch3", "ch4", "ch5", "ch6", "ch7", "ch8", "ch9", "ch10"}, true},
		{"11个频道（超过限制）", []string{"ch1", "ch2", "ch3", "ch4", "ch5", "ch6", "ch7", "ch8", "ch9", "ch10", "ch11"}, false},
		{"空列表", []string{}, false},
		{"包含无效频道名", []string{"valid-channel", "invalid channel"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := websocket.ValidateChannels(tt.channels)
			if tt.shouldPass && err != nil {
				t.Errorf("ValidateChannels(%v) 应该通过但返回错误: %v", tt.channels, err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("ValidateChannels(%v) 应该失败但未返回错误", tt.channels)
			}
		})
	}
}

// TestPusherConstants 测试 Pusher 常量定义
func TestPusherConstants(t *testing.T) {
	if websocket.MaxChannelNameLength != 200 {
		t.Errorf("MaxChannelNameLength = %d, want 200", websocket.MaxChannelNameLength)
	}

	if websocket.MaxEventNameLength != 200 {
		t.Errorf("MaxEventNameLength = %d, want 200", websocket.MaxEventNameLength)
	}

	if websocket.MaxChannelsPerTrigger != 10 {
		t.Errorf("MaxChannelsPerTrigger = %d, want 10", websocket.MaxChannelsPerTrigger)
	}
}
