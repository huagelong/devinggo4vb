// Package websocket
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package websocket

import (
	"context"
	"devinggo/modules/system/pkg/websocket/glob"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// Redis key 常量定义（Pusher协议）
const (
	KeySocketIdHeartbeatTime = "SocketId2HeartbeatTime" // 客户端心跳时间
	KeyClearExpireLock       = "ClearExpire4Redis"      // 清理过期数据的锁
	KeySocketId2ServerName   = "SocketId2ServerName:"   // 客户端对应的服务器名称
	KeyServerNames           = "ServerNames"            // 所有服务器名称集合
	KeyChannels              = "Channels"               // 所有频道集合
	KeyChannel2SocketId      = "Channel2SocketId:"      // 频道对应的客户端集合
	KeySocketId2Channel      = "SocketId2Channel:"      // 客户端对应的频道集合
	KeyChannel2ServerName    = "Channel2ServerName:"    // 频道对应的服务器名称集合
)

func getRedisClient() *gredis.Redis {
	return g.Redis("websocket")
}

// 删除心跳数据
func RemoveSocketIdHeartbeatTime4Redis(ctx context.Context, socketId string) (err error) {
	if g.IsEmpty(socketId) {
		return
	}
	_, err = getRedisClient().Do(ctx, "HDEL", KeySocketIdHeartbeatTime, socketId)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2HeartbeatTime HDEL error:", err)
		return
	}
	return
}

// 更新心跳数据
func UpdateSocketIdHeartbeatTime4Redis(ctx context.Context, socketId string, currentTime int64) (err error) {
	if g.IsEmpty(socketId) {
		return
	}
	_, err = getRedisClient().HSet(ctx, KeySocketIdHeartbeatTime, g.Map{socketId: currentTime})
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2HeartbeatTime HSET error:", err)
		return
	}
	return
}

// 清理心跳过期数据,清除所有客户端数据
func ClearExpire4Redis(ctx context.Context) (err error) {
	rs, err := getRedisClient().SetNX(ctx, KeyClearExpireLock, 1)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "ClearExpire4Redis SetNX error:", err)
		return
	}
	if !rs {
		return
	}
	getRedisClient().Expire(ctx, KeyClearExpireLock, 3600)
	value, err := getRedisClient().HGetAll(ctx, KeySocketIdHeartbeatTime)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2HeartbeatTime HGETALL error:", err)
		getRedisClient().Del(ctx, KeyClearExpireLock)
		return
	}
	for socketId, currentTime := range value.Map() {
		now := int(gtime.Now().Unix())
		currentTimeInt := gconv.Int(currentTime)
		glob.WithWsLog().Debug(ctx, "ClearExpire4Redis:", socketId)
		if heartbeatExpirationTime+currentTimeInt <= now {
			ClearSocketId4Redis(ctx, socketId)
		}
	}
	getRedisClient().Del(ctx, KeyClearExpireLock)
	return
}

// 清除所有客户端数据，包含心跳数据，订阅数据，全局数据
func ClearSocketId4Redis(ctx context.Context, socketId string) (err error) {
	err = RemoveSocketIdHeartbeatTime4Redis(ctx, socketId)
	for _, channel := range GetAllChannelBySocketId(ctx, socketId) {
		LeaveChannel4Redis(ctx, channel, socketId)
	}
	err = DeleteServerNameBySocketId4Redis(ctx, socketId)
	return
}

// 删除客户端订阅数据
func DeleteServerNameBySocketId4Redis(ctx context.Context, socketId string) (err error) {
	key := KeySocketId2ServerName + socketId
	_, err = getRedisClient().Del(ctx, key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "DeleteServerNameBySocketId error:", err)
	}
	return
}

// 获取客户端订阅数据
func GetServerNameBySocketId4Redis(ctx context.Context, socketId string) string {
	key := KeySocketId2ServerName + socketId
	serverName, err := getRedisClient().Get(ctx, key)

	if err != nil {
		glob.WithWsLog().Warning(ctx, "GetServerNameBySocketId4Redis error:", err)
		return ""
	}
	return gconv.String(serverName)
}

// 添加客户端订阅数据,并确认在那个服务器上
func AddServerNameSocketId4Redis(ctx context.Context, socketId string, serverName string) (err error) {
	key := KeySocketId2ServerName + socketId
	getRedisClient().Set(ctx, key, serverName)
	_, err = getRedisClient().Do(ctx, "SADD", KeyServerNames, serverName)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "ServerNames SADD error:", err)
	}
	return
}

// 获取所有服务器名称
func GetAllServerNames(ctx context.Context) []string {
	ls, err := getRedisClient().Do(ctx, "SMEMBERS", KeyServerNames)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "ServerNames error:", err)
		return nil
	}
	return gconv.Strings(ls)
}

// 加入频道
func JoinChannel4Redis(ctx context.Context, socketId string, channel string) (err error) {
	if g.IsEmpty(channel) {
		return
	}
	getRedisClient().Do(ctx, "MULTI")
	_, err = getRedisClient().Do(ctx, "SADD", KeyChannels, channel)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channels SADD error:", err)
		getRedisClient().Do(ctx, "DISCARD")
		return
	}

	key := KeyChannel2SocketId + channel
	_, err = getRedisClient().Do(ctx, "SADD", key, socketId)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channel2SocketId SADD error:", err)
		getRedisClient().Do(ctx, "DISCARD")
		return
	}

	keySocket2Channel := KeySocketId2Channel + socketId
	_, err = getRedisClient().Do(ctx, "SADD", keySocket2Channel, channel)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2Channel SADD error:", err)
		getRedisClient().Do(ctx, "DISCARD")
		return
	}

	getRedisClient().Do(ctx, "EXEC")

	keyServername := KeyChannel2ServerName + channel
	serverName := GetServerNameBySocketId4Redis(ctx, socketId)

	if !g.IsEmpty(serverName) {
		_, err = getRedisClient().Do(ctx, "SADD", keyServername, serverName)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "Channel2ServerName SADD error:", err)
			return
		}
	}
	return
}

// 离开频道
func LeaveChannel4Redis(ctx context.Context, channel string, socketId string) (err error) {
	if g.IsEmpty(channel) {
		return
	}
	key := KeyChannel2SocketId + channel
	_, err = getRedisClient().Do(ctx, "SREM", key, socketId)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channel2SocketId SREM error:", err)
		return
	}

	keyServername := KeyChannel2ServerName + channel
	serverName := GetServerNameBySocketId4Redis(ctx, socketId)
	if !g.IsEmpty(serverName) {
		_, err = getRedisClient().Do(ctx, "SREM", keyServername, serverName)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "Channel2ServerName SREM error:", err)
			return
		}
	}

	keySocket2Channel := KeySocketId2Channel + socketId
	_, err = getRedisClient().Do(ctx, "SREM", keySocket2Channel, channel)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2Channel SREM error:", err)
		return
	}

	keyChannel2SocketId := KeyChannel2SocketId + channel
	count, err := getRedisClient().Do(ctx, "SCARD", keyChannel2SocketId)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channel2SocketId SCARD error:", err)
		return
	}

	if gconv.Int(count) == 0 {
		_, err = getRedisClient().Do(ctx, "SREM", KeyChannels, channel)
		if err != nil {
			glob.WithWsLog().Warning(ctx, "Channels SREM error:", err)
			return
		}
	}
	return
}

// GetAllSocketIdByChannel4Redis 获取频道内所有socket_id
func GetAllSocketIdByChannel4Redis(ctx context.Context, channel string) []string {
	if g.IsEmpty(channel) {
		return nil
	}

	key := KeyChannel2SocketId + channel
	ls, err := getRedisClient().Do(ctx, "SMEMBERS", key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channel2SocketId SMEMBERS error:", err)
		return nil
	}
	return gconv.Strings(ls)
}

// 获取频道的服务器名称
func GetAllServerNameByChannel(ctx context.Context, channel string) []string {
	if g.IsEmpty(channel) {
		return nil
	}

	keyServername := KeyChannel2ServerName + channel
	ls, err := getRedisClient().Do(ctx, "SMEMBERS", keyServername)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channel2ServerName error:", err)
		return nil
	}
	return gconv.Strings(ls)
}

// 获取客户端订阅的所有频道
func GetAllChannelBySocketId(ctx context.Context, socketId string) []string {
	if g.IsEmpty(socketId) {
		return nil
	}

	key := KeySocketId2Channel + socketId
	ls, err := getRedisClient().Do(ctx, "SMEMBERS", key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "SocketId2Channel SMEMBERS error:", err)
		return nil
	}
	return gconv.Strings(ls)
}

// 获取所有频道
func GetAllChannels(ctx context.Context) []string {
	ls, err := getRedisClient().Do(ctx, "SMEMBERS", KeyChannels)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channels SMEMBERS error:", err)
		return nil
	}
	return gconv.Strings(ls)
}

// 判断频道是否存在
func isChannelExist(ctx context.Context, channel string) bool {
	if g.IsEmpty(channel) {
		return false
	}
	ls, err := getRedisClient().Do(ctx, "SISMEMBER", KeyChannels, channel)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "Channels SISMEMBER error:", err)
		return false
	}
	return gconv.Int(ls) == 1
}

// ========== Presence Channel Redis操作 ==========

const (
	KeyPresenceChannel    = "PresenceChannel:"    // Presence频道成员信息 Hash
	KeyPresenceDisconnect = "PresenceDisconnect:" // Presence断线标记（用于Grace Period）
	PresenceGracePeriod   = 30                    // Grace Period 30秒
)

// AddPresenceMember4Redis 添加Presence成员
func AddPresenceMember4Redis(ctx context.Context, channel, userID string, userInfo map[string]interface{}) error {
	if g.IsEmpty(channel) || g.IsEmpty(userID) {
		return nil
	}

	key := KeyPresenceChannel + channel
	userInfoJSON := gconv.String(userInfo)

	// 🔍 调试：打印添加的成员信息
	glob.WithWsLog().Debugf(ctx, "🔍 [AddPresenceMember] channel=%s, userID=%s, userInfo=%+v", channel, userID, userInfo)
	glob.WithWsLog().Debugf(ctx, "🔍 [AddPresenceMember] userInfoJSON=%s", userInfoJSON)

	_, err := getRedisClient().HSet(ctx, key, g.Map{userID: userInfoJSON})
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceChannel HSET error:", err)
		return err
	}

	glob.WithWsLog().Debugf(ctx, "AddPresenceMember: channel=%s, userID=%s", channel, userID)
	return nil
}

// RemovePresenceMember4Redis 移除Presence成员
func RemovePresenceMember4Redis(ctx context.Context, channel, userID string) error {
	if g.IsEmpty(channel) || g.IsEmpty(userID) {
		return nil
	}

	key := KeyPresenceChannel + channel
	_, err := getRedisClient().HDel(ctx, key, userID)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceChannel HDEL error:", err)
		return err
	}

	glob.WithWsLog().Debugf(ctx, "RemovePresenceMember: channel=%s, userID=%s", channel, userID)
	return nil
}

// GetPresenceMembers4Redis 获取Presence频道所有成员
// 返回 map[userID]userInfo
func GetPresenceMembers4Redis(ctx context.Context, channel string) (map[string]map[string]interface{}, error) {
	if g.IsEmpty(channel) {
		return make(map[string]map[string]interface{}), nil
	}

	key := KeyPresenceChannel + channel
	value, err := getRedisClient().HGetAll(ctx, key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceChannel HGETALL error:", err)
		return nil, err
	}

	// 🔍 调试：打印从Redis获取的原始数据
	glob.WithWsLog().Debugf(ctx, "🔍 [GetPresenceMembers] channel=%s, Redis返回数据: %+v", channel, value.Map())

	members := make(map[string]map[string]interface{})
	for userID, userInfoJSON := range value.Map() {
		glob.WithWsLog().Debugf(ctx, "🔍 [GetPresenceMembers] 解析: userID=%s, userInfoJSON=%v", userID, userInfoJSON)

		var userInfo map[string]interface{}
		if err := gconv.Struct(userInfoJSON, &userInfo); err == nil {
			members[userID] = userInfo
			glob.WithWsLog().Debugf(ctx, "🔍 [GetPresenceMembers] 解析成功: userInfo=%+v", userInfo)
		} else {
			glob.WithWsLog().Warning(ctx, "🔍 [GetPresenceMembers] 解析失败: userID=%s, error=%v", userID, err)
		}
	}

	glob.WithWsLog().Debugf(ctx, "🔍 [GetPresenceMembers] channel=%s, 最终members数量=%d", channel, len(members))

	return members, nil
}

// GetPresenceCount4Redis 获取Presence频道成员数量
func GetPresenceCount4Redis(ctx context.Context, channel string) int {
	if g.IsEmpty(channel) {
		return 0
	}

	key := KeyPresenceChannel + channel
	count, err := getRedisClient().HLen(ctx, key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceChannel HLEN error:", err)
		return 0
	}

	return int(count)
}

// MarkPresenceDisconnect4Redis 标记Presence断线（Grace Period开始）
func MarkPresenceDisconnect4Redis(ctx context.Context, socketId string) error {
	if g.IsEmpty(socketId) {
		return nil
	}

	key := KeyPresenceDisconnect + socketId
	timestamp := gtime.Now().Unix()

	err := getRedisClient().SetEX(ctx, key, timestamp, PresenceGracePeriod)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceDisconnect SETEX error:", err)
		return err
	}

	return nil
}

// ClearPresenceDisconnect4Redis 清除Presence断线标记（重连时）
func ClearPresenceDisconnect4Redis(ctx context.Context, socketId string) error {
	if g.IsEmpty(socketId) {
		return nil
	}

	key := KeyPresenceDisconnect + socketId
	_, err := getRedisClient().Del(ctx, key)
	if err != nil {
		glob.WithWsLog().Warning(ctx, "PresenceDisconnect DEL error:", err)
		return err
	}

	return nil
}

// IsPresenceDisconnected4Redis 检查是否在Grace Period内断线
func IsPresenceDisconnected4Redis(ctx context.Context, socketId string) bool {
	if g.IsEmpty(socketId) {
		return false
	}

	key := KeyPresenceDisconnect + socketId
	exists, err := getRedisClient().Exists(ctx, key)
	if err != nil {
		return false
	}

	return exists > 0
}
