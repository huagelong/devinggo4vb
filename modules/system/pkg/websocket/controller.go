// Package websocket
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package websocket

import (
	"context"
	"devinggo/modules/system/pkg/websocket/glob"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
)

// SubscribeController Pusher订阅控制器
func SubscribeController(ctx context.Context, client *Client, req *PusherRequest) {
	// 解析data字段
	var subData SubscribeRequestData
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SubscribeController marshal error:", err)
		client.SendError("Invalid data", CodeNormalClosure)
		return
	}

	err = json.Unmarshal(dataBytes, &subData)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SubscribeController unmarshal error:", err)
		client.SendError("Invalid subscribe data", CodeNormalClosure)
		return
	}

	channel := subData.Channel
	if channel == "" {
		client.SendError("Missing channel", CodeNormalClosure)
		return
	}

	// 获取频道类型
	channelType := GetChannelType(channel)

	// Phase 2: Private Channel认证
	if channelType == ChannelTypePrivate {
		auth := subData.Auth
		if auth == "" {
			glob.WithWsLog().Warning(ctx, "Private channel missing auth:", channel)
			client.SendSubscriptionError(channel, "AuthError", "Auth signature required for private channel", CodeUnauthorized)
			return
		}

		// 验证认证签名
		err := ValidateChannelAuth(client.SocketID, channel, auth, "")
		if err != nil {
			glob.WithWsLog().Warning(ctx, "Private channel auth failed:", err)
			client.SendSubscriptionError(channel, "AuthError", "Invalid auth signature", CodeUnauthorized)
			return
		}

		// 验证成功：加入频道
		if !client.HasChannel(channel) {
			client.AddChannel(channel)
		}

		err = JoinChannel4Redis(ctx, client.SocketID, channel)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "JoinChannel4Redis error:", err)
			client.SendSubscriptionError(channel, "ServerError", err.Error(), 500)
			return
		}

		// 发送订阅成功事件
		client.SendPusherEvent(EventSubscriptionSucceeded, channel, map[string]interface{}{})
		glob.WithWsLog().Debugf(ctx, "Client %s subscribed to private channel %s", client.SocketID, channel)
		return
	}

	// Phase 3: Presence Channel认证
	if channelType == ChannelTypePresence {
		auth := subData.Auth
		channelData := subData.ChannelData

		if auth == "" || channelData == "" {
			glob.WithWsLog().Warning(ctx, "Presence channel missing auth or channel_data:", channel)
			client.SendSubscriptionError(channel, "AuthError", "Auth signature and channel_data required for presence channel", CodeUnauthorized)
			return
		}

		// 验证认证签名（包含channel_data）
		err := ValidateChannelAuth(client.SocketID, channel, auth, channelData)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "Presence channel auth failed:", err)
			client.SendSubscriptionError(channel, "AuthError", "Invalid auth signature", CodeUnauthorized)
			return
		}

		// 解析channel_data获取user_id和user_info
		member, err := ParseChannelData(channelData)
		if err != nil || member.UserID == "" {
			glob.WithWsLog().Warning(ctx, "Invalid channel_data:", err)
			client.SendSubscriptionError(channel, "AuthError", "Invalid channel_data format", CodeUnauthorized)
			return
		}

		// 保存用户信息到Client
		client.UserID = member.UserID
		client.UserInfo = member.UserInfo

		// 添加成员到Redis
		err = AddPresenceMember4Redis(ctx, channel, member.UserID, member.UserInfo)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "AddPresenceMember4Redis error:", err)
			client.SendSubscriptionError(channel, "ServerError", err.Error(), 500)
			return
		}

		// 加入频道
		if !client.HasChannel(channel) {
			client.AddChannel(channel)
		}

		err = JoinChannel4Redis(ctx, client.SocketID, channel)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "JoinChannel4Redis error:", err)
			client.SendSubscriptionError(channel, "ServerError", err.Error(), 500)
			return
		}

		// 获取所有成员列表（使用缓存优化性能）
		members, _, err := GetPresenceCache().GetMembers(ctx, channel)
		if err != nil {
			members = make(map[string]map[string]interface{})
		}

		// 🔍 调试：打印成员列表
		glob.WithWsLog().Debugf(ctx, "🔍 [Presence订阅] channel=%s, members数量=%d", channel, len(members))
		for userID, userInfo := range members {
			glob.WithWsLog().Debugf(ctx, "🔍 [Presence成员] userID=%s, userInfo=%+v", userID, userInfo)
		}

		// 发送订阅成功事件（包含完整成员列表）
		// ⚠️ 传入当前用户ID以构造 me 字段
		// ⚠️ Pusher.js 8.3.0 需要扁平化结构，直接发送 PresenceMemberList，不包装在 presence 字段中
		presenceData := FormatPresenceData(members, member.UserID)

		// 🔍 调试：验证 presenceData.Presence 可以被正确序列化
		testJSON, testErr := json.Marshal(presenceData.Presence) // 使用 .Presence 字段（扁平化）
		if testErr != nil {
			glob.WithWsLog().Errorf(ctx, "❌ [Presence数据] JSON序列化测试失败: %v", testErr)
		} else {
			glob.WithWsLog().Debugf(ctx, "🔍 [Presence数据] JSON序列化成功（扁平化），长度=%d, 内容=%s", len(testJSON), string(testJSON))
		}

		// 直接发送 PresenceMemberList（扁平化），而不是包装在 PresenceData 中
		client.SendPusherEvent(EventSubscriptionSucceeded, channel, presenceData.Presence)

		// 向频道内其他成员广播member_added事件
		memberAddedData := MemberAddedData{
			UserID:   member.UserID,
			UserInfo: member.UserInfo,
		}
		BroadcastToChannel(ctx, channel, EventMemberAdded, memberAddedData, client.SocketID)

		// 使缓存失效（因为新增了成员）
		GetPresenceCache().InvalidateChannel(channel)

		glob.WithWsLog().Debugf(ctx, "Client %s subscribed to presence channel %s as user %s", client.SocketID, channel, member.UserID)
		return
	}

	// Public Channel订阅
	if !client.HasChannel(channel) {
		client.AddChannel(channel)
	}

	err = JoinChannel4Redis(ctx, client.SocketID, channel)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "JoinChannel4Redis error:", err)
		client.SendSubscriptionError(channel, "ServerError", err.Error(), 500)
		return
	}

	// 发送订阅成功事件（⚠️ Public Channel的data为空对象字符串 "{}"）
	client.SendPusherEvent(EventSubscriptionSucceeded, channel, map[string]interface{}{})
	glob.WithWsLog().Debugf(ctx, "Client %s subscribed to channel %s", client.SocketID, channel)
}

// UnsubscribeController Pusher退订控制器
func UnsubscribeController(ctx context.Context, client *Client, req *PusherRequest) {
	// 解析data字段
	var subData SubscribeRequestData
	dataBytes, err := json.Marshal(req.Data)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "UnsubscribeController marshal error:", err)
		return
	}

	err = json.Unmarshal(dataBytes, &subData)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "UnsubscribeController unmarshal error:", err)
		return
	}

	channel := subData.Channel
	if channel == "" {
		return
	}

	if client.HasChannel(channel) {
		client.RemoveChannel(channel)
	}

	LeaveChannel4Redis(ctx, channel, client.SocketID)

	// Phase 3: Presence Channel退订（触发Grace Period）
	channelType := GetChannelType(channel)
	if channelType == ChannelTypePresence && client.UserID != "" {
		// ⚠️ 标记断线，启动Grace Period（30秒）
		MarkPresenceDisconnect4Redis(ctx, client.SocketID)

		// 30秒后移除成员（如果未重连）
		// 注意：实际延迟清理在client.close()中实现
		glob.WithWsLog().Debugf(ctx, "Client %s marked for presence grace period on channel %s", client.SocketID, channel)
	}

	// ⚠️ Pusher规范：unsubscribe不发送响应事件
	glob.WithWsLog().Debugf(ctx, "Client %s unsubscribed from channel %s", client.SocketID, channel)
}

// PingController Pusher心跳控制器
func PingController(ctx context.Context, client *Client, req *PusherRequest) {
	currentTime := int64(gtime.Now().Unix())
	client.Heartbeat(currentTime)
	UpdateSocketIdHeartbeatTime4Redis(ctx, client.SocketID, currentTime)

	// 发送pong响应（⚠️ data为空对象字符串 "{}"）
	client.SendPusherEvent(EventPong, "", map[string]interface{}{})
}

// ClientEventController 客户端事件控制器（Phase 4: Client Events）
func ClientEventController(ctx context.Context, client *Client, req *PusherRequest) {
	// ⚠️ v8.3.0要求：验证事件名必须以client-开头
	if !strings.HasPrefix(req.Event, "client-") {
		glob.WithWsLog().Warning(ctx, "Invalid client event name:", req.Event)
		client.SendError("Client event must start with 'client-'", CodeClientEventForbidden)
		return
	}

	// ⚠️ v8.3.0要求：验证事件名长度（最大50字节）
	if len(req.Event) > 50 {
		glob.WithWsLog().Warning(ctx, "Client event name too long:", req.Event)
		client.SendError("Event name exceeds 50 bytes", CodeClientEventForbidden)
		return
	}

	// 验证channel字段
	if req.Channel == "" {
		glob.WithWsLog().Warning(ctx, "Client event missing channel")
		client.SendError("Channel required for client events", CodeClientEventForbidden)
		return
	}

	// ⚠️ v8.3.0要求：仅允许Private/Presence频道使用Client Events
	channelType := GetChannelType(req.Channel)
	if channelType == ChannelTypePublic {
		glob.WithWsLog().Warning(ctx, "Client events not allowed on public channels:", req.Channel)
		client.SendError("Client events only allowed on private and presence channels", CodeClientEventForbidden)
		return
	}

	// 验证客户端是否已订阅该频道
	if !client.HasChannel(req.Channel) {
		glob.WithWsLog().Warning(ctx, "Client not subscribed to channel:", req.Channel)
		client.SendError("Must subscribe to channel before sending client events", CodeClientEventForbidden)
		return
	}

	// ⚠️ v8.3.0要求：速率限制（10条/秒）
	rateLimiter := GetRateLimiter()
	if !rateLimiter.AllowClientEvent(client.SocketID) {
		glob.WithWsLog().Warning(ctx, "Rate limit exceeded for client:", client.SocketID)
		client.SendError("Rate limit exceeded (max 10 events/sec)", CodeClientEventForbidden)
		return
	}

	// 转发给频道内其他客户端（不包括发送者）
	BroadcastToChannel(ctx, req.Channel, req.Event, req.Data, client.SocketID)

	glob.WithWsLog().Debugf(ctx, "Client event forwarded: socket=%s, channel=%s, event=%s",
		client.SocketID, req.Channel, req.Event)
}

// BroadcastToChannel 向频道内除了指定客户端以外的所有成员广播消息
func BroadcastToChannel(ctx context.Context, channel, event string, data interface{}, excludeSocketID string) {
	// 获取频道内所有socket_id
	socketIds := GetAllSocketIdByChannel4Redis(ctx, channel)

	for _, socketId := range socketIds {
		// 跳过指定的客户端
		if socketId == excludeSocketID {
			continue
		}

		// 获取客户端并发送消息
		targetClient := clientManager.GetClientBySocketID(socketId)
		if targetClient != nil {
			targetClient.SendPusherEvent(event, channel, data)
		}
	}
}

// BroadcastMemberRemoved 广播member_removed事件
func BroadcastMemberRemoved(ctx context.Context, channel, userID string) {
	memberRemovedData := MemberRemovedData{
		UserID: userID,
	}

	// 向频道内所有成员广播
	socketIds := GetAllSocketIdByChannel4Redis(ctx, channel)
	for _, socketId := range socketIds {
		targetClient := clientManager.GetClientBySocketID(socketId)
		if targetClient != nil {
			targetClient.SendPusherEvent(EventMemberRemoved, channel, memberRemovedData)
		}
	}

	// 使Presence缓存失效（成员已移除）
	GetPresenceCache().InvalidateChannel(channel)
}
