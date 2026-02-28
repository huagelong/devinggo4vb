// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"devinggo/modules/system/api/system"
	"devinggo/modules/system/pkg/websocket"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	PusherEvents = cPusherEvents{}
)

type cPusherEvents struct{}

// Events Pusher HTTP Events API - 推送事件到频道
func (c *cPusherEvents) Events(ctx context.Context, req *system.PusherEventsReq) (res *system.PusherEventsRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 1. 验证应用配置
	config, err := getAppConfig(ctx)
	if err != nil {
		return nil, err
	}

	if req.AppId != config.AppID {
		return nil, invalidAppIdError(r)
	}

	if req.AuthKey != config.Key {
		return nil, invalidAppKeyError(r)
	}

	// 2. 验证时间戳（防止重放攻击，允许±600秒误差）
	now := time.Now().Unix()
	if abs(now-req.AuthTimestamp) > 600 {
		return nil, timestampExpiredError(r)
	}

	// 3. 验证签名
	bodyBytes := r.GetBody()
	if !verifySignature(req.AuthKey, req.AuthTimestamp, req.AuthVersion, req.BodyMd5, req.AuthSignature, config.Secret, "POST", fmt.Sprintf("/apps/%s/events", req.AppId), bodyBytes) {
		return nil, signatureInvalidError(r)
	}

	// 4. 推送事件到各个频道
	g.Log().Debugf(ctx, "HTTP Events API: event=%s, channels=%v, data=%s", req.Name, req.Channels, req.Data)

	for _, channel := range req.Channels {
		// 构建Pusher响应消息
		pusherResponse := &websocket.PusherResponse{
			Event:   req.Name,
			Channel: channel,
			Data:    req.Data, // 已经是JSON字符串
		}

		g.Log().Debugf(ctx, "Sending to channel: %s, event: %s", channel, req.Name)

		// 1) 先发送给本地服务器的客户端
		websocket.SendToChannel(channel, pusherResponse)

		// 2) 再发布到其他服务器（通过Redis PubSub）
		topicMsg := &websocket.TopicWResponse{
			Topic:          channel,
			PusherResponse: pusherResponse,
		}

		// 排除特定socket_id（如果指定）
		if req.SocketId != "" {
			// TODO: 实现socket_id排除逻辑
			g.Log().Debug(ctx, "Exclude socket_id:", req.SocketId)
		}

		err = websocket.PublishChannelMessage(ctx, channel, topicMsg)
		if err != nil {
			g.Log().Warning(ctx, "Failed to publish message to channel:", channel, err)
		}
	}

	// 5. 返回成功响应
	res = &system.PusherEventsRes{}
	return
}

// BatchEvents Pusher HTTP Batch Events API - 批量推送事件
func (c *cPusherEvents) BatchEvents(ctx context.Context, req *system.PusherBatchEventsReq) (res *system.PusherBatchEventsRes, err error) {
	r := g.RequestFromCtx(ctx)

	// 1. 验证应用配置
	config, err := getAppConfig(ctx)
	if err != nil {
		return nil, err
	}

	if req.AppId != config.AppID {
		return nil, invalidAppIdError(r)
	}

	if req.AuthKey != config.Key {
		return nil, invalidAppKeyError(r)
	}

	// 2. 验证时间戳
	now := time.Now().Unix()
	if abs(now-req.AuthTimestamp) > 600 {
		return nil, timestampExpiredError(r)
	}

	// 3. 验证签名
	bodyBytes := r.GetBody()
	if !verifySignature(req.AuthKey, req.AuthTimestamp, req.AuthVersion, req.BodyMd5, req.AuthSignature, config.Secret, "POST", fmt.Sprintf("/apps/%s/batch_events", req.AppId), bodyBytes) {
		return nil, signatureInvalidError(r)
	}

	// 4. 批量推送事件
	for _, event := range req.Batch {
		pusherResponse := &websocket.PusherResponse{
			Event:   event.Name,
			Channel: event.Channel,
			Data:    event.Data,
		}

		// 1) 先发送给本地服务器的客户端
		websocket.SendToChannel(event.Channel, pusherResponse)

		// 2) 再发布到其他服务器（通过Redis PubSub）
		topicMsg := &websocket.TopicWResponse{
			Topic:          event.Channel,
			PusherResponse: pusherResponse,
		}

		if event.SocketId != "" {
			g.Log().Debug(ctx, "Exclude socket_id:", event.SocketId)
		}

		err = websocket.PublishChannelMessage(ctx, event.Channel, topicMsg)
		if err != nil {
			g.Log().Warning(ctx, "Failed to publish message to channel:", event.Channel, err)
		}
	}

	// 5. 返回成功响应
	res = &system.PusherBatchEventsRes{}
	return
}

// verifySignature 验证Pusher HTTP API签名
func verifySignature(authKey string, authTimestamp int64, authVersion string, bodyMd5Provided string, authSignature string, appSecret string, method string, path string, bodyBytes []byte) bool {
	// 1. 验证body_md5
	bodyMd5Computed := fmt.Sprintf("%x", md5.Sum(bodyBytes))

	g.Log().Debugf(context.Background(), "Signature Verification:")
	g.Log().Debugf(context.Background(), "  Body MD5 (provided): %s", bodyMd5Provided)
	g.Log().Debugf(context.Background(), "  Body MD5 (computed): %s", bodyMd5Computed)
	g.Log().Debugf(context.Background(), "  Body length: %d bytes", len(bodyBytes))

	if bodyMd5Provided != bodyMd5Computed {
		g.Log().Warning(context.Background(), "Body MD5 mismatch!")
		return false
	}

	// 2. 构建查询字符串（按字母顺序排序）
	queryParams := map[string]string{
		"auth_key":       authKey,
		"auth_timestamp": gconv.String(authTimestamp),
		"auth_version":   authVersion,
		"body_md5":       bodyMd5Provided,
	}

	keys := make([]string, 0, len(queryParams))
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	queryParts := make([]string, 0, len(keys))
	for _, k := range keys {
		queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, queryParams[k]))
	}
	queryString := strings.Join(queryParts, "&")

	// 3. 构建待签名字符串
	stringToSign := fmt.Sprintf("%s\n%s\n%s", method, path, queryString)

	g.Log().Debugf(context.Background(), "  Query string: %s", queryString)
	g.Log().Debugf(context.Background(), "  String to sign: %s", stringToSign)

	// 4. 计算HMAC-SHA256签名
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(stringToSign))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	g.Log().Debugf(context.Background(), "  Expected signature: %s", expectedSignature)
	g.Log().Debugf(context.Background(), "  Provided signature: %s", authSignature)

	// 5. 比对签名
	match := authSignature == expectedSignature
	if !match {
		g.Log().Warning(context.Background(), "Signature mismatch!")
	}
	return match
}

// getAppConfig 获取应用配置（复用pusher_auth.go中的逻辑）
func getAppConfig(ctx context.Context) (*AppConfig, error) {
	config := g.Cfg()
	appID := config.MustGet(ctx, "pusher.appId", "").String()
	appKey := config.MustGet(ctx, "pusher.appKey", "").String()
	appSecret := config.MustGet(ctx, "pusher.appSecret", "").String()

	if appID == "" || appKey == "" || appSecret == "" {
		return nil, fmt.Errorf("WebSocket Pusher configuration not found in config file")
	}

	return &AppConfig{
		AppID:  appID,
		Key:    appKey,
		Secret: appSecret,
	}, nil
}

// AppConfig 应用配置
type AppConfig struct {
	AppID  string
	Key    string
	Secret string
}

// 错误响应辅助函数
func invalidAppIdError(r *ghttp.Request) error {
	r.Response.Status = 400
	r.Response.WriteJson(g.Map{
		"error": "Invalid app_id",
	})
	r.ExitAll()
	return nil
}

func invalidAppKeyError(r *ghttp.Request) error {
	r.Response.Status = 401
	r.Response.WriteJson(g.Map{
		"error": "Invalid app_key",
	})
	r.ExitAll()
	return nil
}

func timestampExpiredError(r *ghttp.Request) error {
	r.Response.Status = 401
	r.Response.WriteJson(g.Map{
		"error": "Timestamp expired",
	})
	r.ExitAll()
	return nil
}

func signatureInvalidError(r *ghttp.Request) error {
	r.Response.Status = 401
	r.Response.WriteJson(g.Map{
		"error": "Invalid signature",
	})
	r.ExitAll()
	return nil
}

// abs 计算绝对值
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
