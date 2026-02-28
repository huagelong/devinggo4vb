// Package system
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package system

import (
	"devinggo/modules/system/model"

	"github.com/gogf/gf/v2/frame/g"
)

// PusherAuthReq Pusher频道认证请求
type PusherAuthReq struct {
	g.Meta `path:"/pusher/auth" method:"post" tags:"WebSocket" summary:"Pusher频道认证" x-exceptAuth:"true"`
	model.AuthorHeader
	SocketId    string `json:"socket_id" dc:"WebSocket连接的socket_id" v:"required"`
	ChannelName string `json:"channel_name" dc:"要订阅的频道名称" v:"required"`
}

// PusherAuthRes Pusher频道认证响应
type PusherAuthRes struct {
	g.Meta      `mime:"application/json"`
	Auth        string `json:"auth" dc:"认证签名，格式：{app_key}:{signature}"`
	ChannelData string `json:"channel_data,omitempty" dc:"Presence频道的用户数据（可选），格式：{user_id, user_info}"`
}
