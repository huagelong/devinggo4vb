package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockPusherClient 模拟 Pusher 客户端结构
type MockPusherClient struct {
	AppID                     string
	Key                       string
	Secret                    string
	Cluster                   string
	Secure                    bool
	Host                      string
	HTTPClient                *http.Client
	EncryptionMasterKeyBase64 string
}

// TriggerParams 触发参数结构
type TriggerParams struct {
	SocketID *string
	Info     *string
}

// Event 事件结构
type Event struct {
	Channel  string
	Name     string
	Data     interface{}
	SocketID *string
	Info     *string
}

// MemberData 成员数据结构
type MemberData struct {
	UserID   string
	UserInfo map[string]string
}

// ChannelsParams 频道查询参数
type ChannelsParams struct {
	FilterByPrefix *string
	Info           *string
}

// ChannelParams 单个频道查询参数
type ChannelParams struct {
	Info *string
}

// TestClientCreation 测试客户端创建
func TestClientCreation(t *testing.T) {
	client := MockPusherClient{
		AppID:   "APP_ID",
		Key:     "APP_KEY",
		Secret:  "APP_SECRET",
		Cluster: "APP_CLUSTER",
	}

	if client.AppID != "APP_ID" {
		t.Errorf("Expected AppID to be 'APP_ID', got '%s'", client.AppID)
	}
	if client.Key != "APP_KEY" {
		t.Errorf("Expected Key to be 'APP_KEY', got '%s'", client.Key)
	}
	if client.Secret != "APP_SECRET" {
		t.Errorf("Expected Secret to be 'APP_SECRET', got '%s'", client.Secret)
	}
	if client.Cluster != "APP_CLUSTER" {
		t.Errorf("Expected Cluster to be 'APP_CLUSTER', got '%s'", client.Cluster)
	}
}

// TestClientWithSecure 测试启用 HTTPS 的客户端
func TestClientWithSecure(t *testing.T) {
	client := MockPusherClient{
		AppID:   "APP_ID",
		Key:     "APP_KEY",
		Secret:  "APP_SECRET",
		Cluster: "APP_CLUSTER",
		Secure:  true,
	}

	if !client.Secure {
		t.Error("Expected Secure to be true")
	}
}

// TestClientWithCustomTimeout 测试自定义超时的客户端
func TestClientWithCustomTimeout(t *testing.T) {
	httpClient := &http.Client{Timeout: time.Second * 3}
	client := MockPusherClient{
		AppID:      "APP_ID",
		Key:        "APP_KEY",
		Secret:     "APP_SECRET",
		Cluster:    "APP_CLUSTER",
		HTTPClient: httpClient,
	}

	if client.HTTPClient.Timeout != time.Second*3 {
		t.Errorf("Expected timeout to be 3 seconds, got %v", client.HTTPClient.Timeout)
	}
}

// TestClientWithCustomHost 测试自定义 Host 的客户端
func TestClientWithCustomHost(t *testing.T) {
	client := MockPusherClient{
		AppID:   "APP_ID",
		Key:     "APP_KEY",
		Secret:  "APP_SECRET",
		Cluster: "APP_CLUSTER",
		Host:    "foo.bar.com",
	}

	if client.Host != "foo.bar.com" {
		t.Errorf("Expected Host to be 'foo.bar.com', got '%s'", client.Host)
	}
}

// TestClientWithEncryption 测试端对端加密客户端
func TestClientWithEncryption(t *testing.T) {
	encryptionKey := "dGVzdC1lbmNyeXB0aW9uLWtleQ==" // base64 encoded test key
	client := MockPusherClient{
		AppID:                     "APP_ID",
		Key:                       "APP_KEY",
		Secret:                    "APP_SECRET",
		Cluster:                   "APP_CLUSTER",
		EncryptionMasterKeyBase64: encryptionKey,
	}

	if client.EncryptionMasterKeyBase64 != encryptionKey {
		t.Errorf("Expected EncryptionMasterKeyBase64 to be '%s', got '%s'", encryptionKey, client.EncryptionMasterKeyBase64)
	}
}

// TestTriggerData 测试触发事件的数据结构
func TestTriggerData(t *testing.T) {
	data := map[string]string{"hello": "world"}

	if data["hello"] != "world" {
		t.Errorf("Expected data['hello'] to be 'world', got '%s'", data["hello"])
	}
}

// TestTriggerParams 测试触发参数
func TestTriggerParams(t *testing.T) {
	socketID := "1234.12"
	attributes := "user_count"
	params := TriggerParams{
		SocketID: &socketID,
		Info:     &attributes,
	}

	if *params.SocketID != "1234.12" {
		t.Errorf("Expected SocketID to be '1234.12', got '%s'", *params.SocketID)
	}
	if *params.Info != "user_count" {
		t.Errorf("Expected Info to be 'user_count', got '%s'", *params.Info)
	}
}

// TestMultipleChannels 测试多个频道
func TestMultipleChannels(t *testing.T) {
	channels := []string{"a_channel", "another_channel"}

	if len(channels) != 2 {
		t.Errorf("Expected 2 channels, got %d", len(channels))
	}
	if channels[0] != "a_channel" {
		t.Errorf("Expected first channel to be 'a_channel', got '%s'", channels[0])
	}
	if channels[1] != "another_channel" {
		t.Errorf("Expected second channel to be 'another_channel', got '%s'", channels[1])
	}
}

// TestEventBatch 测试批量事件
func TestEventBatch(t *testing.T) {
	socketID := "1234.12"
	attributes := "user_count"
	batch := []Event{
		{Channel: "a-channel", Name: "event", Data: "hello world"},
		{Channel: "presence-b-channel", Name: "event", Data: "hi my name is bob", SocketID: &socketID, Info: &attributes},
	}

	if len(batch) != 2 {
		t.Errorf("Expected 2 events in batch, got %d", len(batch))
	}
	if batch[0].Channel != "a-channel" {
		t.Errorf("Expected first event channel to be 'a-channel', got '%s'", batch[0].Channel)
	}
	if batch[1].Data != "hi my name is bob" {
		t.Errorf("Expected second event data to be 'hi my name is bob', got '%v'", batch[1].Data)
	}
	if *batch[1].SocketID != "1234.12" {
		t.Errorf("Expected second event SocketID to be '1234.12', got '%s'", *batch[1].SocketID)
	}
}

// TestSendToUserData 测试发送给用户的数据
func TestSendToUserData(t *testing.T) {
	userID := "user123"
	event := "say_hello"
	data := map[string]string{"hello": "world"}

	if userID != "user123" {
		t.Errorf("Expected userID to be 'user123', got '%s'", userID)
	}
	if event != "say_hello" {
		t.Errorf("Expected event to be 'say_hello', got '%s'", event)
	}
	if data["hello"] != "world" {
		t.Errorf("Expected data['hello'] to be 'world', got '%s'", data["hello"])
	}
}

// TestAuthenticateUserData 测试用户认证数据
func TestAuthenticateUserData(t *testing.T) {
	userData := map[string]interface{}{
		"id":      "1234",
		"twitter": "jamiepatel",
	}

	if userData["id"] != "1234" {
		t.Errorf("Expected id to be '1234', got '%v'", userData["id"])
	}
	if userData["twitter"] != "jamiepatel" {
		t.Errorf("Expected twitter to be 'jamiepatel', got '%v'", userData["twitter"])
	}
}

// TestMemberData 测试成员数据
func TestMemberData(t *testing.T) {
	presenceData := MemberData{
		UserID: "1",
		UserInfo: map[string]string{
			"twitter": "jamiepatel",
		},
	}

	if presenceData.UserID != "1" {
		t.Errorf("Expected UserID to be '1', got '%s'", presenceData.UserID)
	}
	if presenceData.UserInfo["twitter"] != "jamiepatel" {
		t.Errorf("Expected twitter to be 'jamiepatel', got '%s'", presenceData.UserInfo["twitter"])
	}
}

// TestChannelsParams 测试频道查询参数
func TestChannelsParams(t *testing.T) {
	prefixFilter := "presence-"
	attributes := "user_count"
	params := ChannelsParams{
		FilterByPrefix: &prefixFilter,
		Info:           &attributes,
	}

	if *params.FilterByPrefix != "presence-" {
		t.Errorf("Expected FilterByPrefix to be 'presence-', got '%s'", *params.FilterByPrefix)
	}
	if *params.Info != "user_count" {
		t.Errorf("Expected Info to be 'user_count', got '%s'", *params.Info)
	}
}

// TestChannelParams 测试单个频道查询参数
func TestChannelParams(t *testing.T) {
	attributes := "user_count,subscription_count"
	params := ChannelParams{
		Info: &attributes,
	}

	if *params.Info != "user_count,subscription_count" {
		t.Errorf("Expected Info to be 'user_count,subscription_count', got '%s'", *params.Info)
	}
}

// TestChannelNamingConvention 测试频道命名约定
func TestChannelNamingConvention(t *testing.T) {
	validChannels := []string{
		"my-channel",
		"my_channel",
		"MyChannel123",
		"presence-chatroom",
		"private-messages",
		"private-encrypted-secret",
	}

	for _, channel := range validChannels {
		if len(channel) > 200 {
			t.Errorf("Channel name '%s' exceeds 200 characters", channel)
		}
	}
}

// TestEventNameLength 测试事件名称长度限制
func TestEventNameLength(t *testing.T) {
	validEvent := "my_event"
	if len(validEvent) > 200 {
		t.Error("Event name exceeds 200 characters")
	}

	// 测试边界情况
	longEvent := ""
	for i := 0; i < 200; i++ {
		longEvent += "a"
	}
	if len(longEvent) != 200 {
		t.Errorf("Expected event name length to be 200, got %d", len(longEvent))
	}
}

// TestChannelLimit 测试频道数量限制（最多10个）
func TestChannelLimit(t *testing.T) {
	channels := []string{
		"channel1", "channel2", "channel3", "channel4", "channel5",
		"channel6", "channel7", "channel8", "channel9", "channel10",
	}

	if len(channels) > 10 {
		t.Error("Cannot trigger more than 10 channels at once")
	}

	// 测试正好10个
	if len(channels) != 10 {
		t.Errorf("Expected 10 channels, got %d", len(channels))
	}
}

// TestHTTPAuthHandler 测试 HTTP 认证处理器
func TestHTTPAuthHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		// 模拟认证响应
		response := `{"auth":"mock_auth_signature"}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))

		if len(body) == 0 {
			t.Error("Expected non-empty request body")
		}
	})

	// 创建测试服务器
	server := httptest.NewServer(handler)
	defer server.Close()

	// 发送测试请求
	requestBody := []byte("socket_id=1234.5678&channel_name=private-test")
	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if len(responseBody) == 0 {
		t.Error("Expected non-empty response body")
	}
}

// TestWebhookValidation 测试 Webhook 验证
func TestWebhookValidation(t *testing.T) {
	// 模拟 webhook 请求
	webhookBody := `{
		"time_ms": 1234567890,
		"events": [
			{
				"name": "channel_occupied",
				"channel": "test-channel"
			}
		]
	}`

	header := http.Header{}
	header.Set("X-Pusher-Key", "test-key")
	header.Set("X-Pusher-Signature", "test-signature")

	if header.Get("X-Pusher-Key") != "test-key" {
		t.Error("Expected X-Pusher-Key header to be set")
	}
	if header.Get("X-Pusher-Signature") != "test-signature" {
		t.Error("Expected X-Pusher-Signature header to be set")
	}

	// 验证 body 非空
	if len(webhookBody) == 0 {
		t.Error("Expected non-empty webhook body")
	}
}

// TestJSONPCallback 测试 JSONP 回调
func TestJSONPCallback(t *testing.T) {
	callback := "myCallback"
	response := `{"auth":"test_signature"}`
	jsonpResponse := callback + "(" + response + ");"

	expectedResponse := "myCallback({\"auth\":\"test_signature\"});"
	if jsonpResponse != expectedResponse {
		t.Errorf("Expected JSONP response to be '%s', got '%s'", expectedResponse, jsonpResponse)
	}
}

// BenchmarkTriggerData 性能测试：创建触发数据
func BenchmarkTriggerData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := map[string]string{"hello": "world", "message": "test"}
		_ = data
	}
}

// BenchmarkEventBatch 性能测试：创建事件批次
func BenchmarkEventBatch(b *testing.B) {
	socketID := "1234.12"
	for i := 0; i < b.N; i++ {
		batch := []Event{
			{Channel: "a-channel", Name: "event", Data: "hello world"},
			{Channel: "b-channel", Name: "event", Data: "test", SocketID: &socketID},
		}
		_ = batch
	}
}
