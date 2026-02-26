// Package cache
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cache

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
)

// set 设置缓存（包内使用）。全部通过 TagCache 单例存储，tag 可为空。
// duration < 0 或 value == nil 时执行删除操作。
func set(ctx context.Context, key interface{}, value interface{}, duration time.Duration, tag ...interface{}) error {
	if value == nil || duration < 0 {
		_, err := Remove(ctx, key)
		return err
	}
	tc := getTagCacheInstance()
	var tags []string
	for _, t := range gconv.Strings(tag) {
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tc.Set(ctx, gconv.String(key), value, duration, tags)
}

// setIfNotExist 仅在 key 不存在时设置（包内使用）。
func setIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration, tag ...interface{}) (ok bool, err error) {
	// 兼容函数值
	f, isFn := value.(gcache.Func)
	if !isFn {
		f, isFn = value.(func(ctx context.Context) (value interface{}, err error))
	}
	if isFn {
		if value, err = f(ctx); err != nil {
			return false, err
		}
	}
	// 删除情况
	if duration < 0 || value == nil {
		return false, getTagCacheInstance().Delete(ctx, gconv.String(key))
	}
	defaultKey := gconv.String(key)
	ok, err = getRedisClient().SetNX(ctx, defaultKey, value)
	if err != nil || !ok {
		return ok, err
	}
	// SetNX 成功后，通过 TagCache 注册 tag 关联并设置过期
	var tags []string
	for _, t := range gconv.Strings(tag) {
		if t != "" {
			tags = append(tags, t)
		}
	}
	return true, getTagCacheInstance().Set(ctx, defaultKey, value, duration, tags)
}

// setIfNotExistFunc 仅在 key 不存在时执行函数并设置（包内使用）。
func setIfNotExistFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration, tag ...interface{}) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return setIfNotExist(ctx, key, value, duration, tag...)
}

// setIfNotExistFuncLock 带锁的 setIfNotExistFunc（包内使用）。
func setIfNotExistFuncLock(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration, tag ...interface{}) (ok bool, err error) {
	value, err := f(ctx)
	if err != nil {
		return false, err
	}
	return setIfNotExist(ctx, key, value, duration, tag...)
}

// SetIfNotExist 仅在 key 不存在时设置（对外公开，不带 tag）。
func SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	return setIfNotExist(ctx, key, value, duration)
}

// Get 获取缓存值（对外公开）。
func Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	return getTagCacheInstance().Get(ctx, gconv.String(key))
}

// getOrSet 获取缓存，不存在时设置（包内使用）。
func getOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration, tag ...interface{}) (result *gvar.Var, err error) {
	result, err = Get(ctx, key)
	if err != nil {
		return
	}
	if result.IsNil() {
		return gvar.New(value), set(ctx, key, value, duration, tag...)
	}
	return
}

// getOrSetFunc 获取缓存，不存在时执行函数并设置（包内使用）。
func getOrSetFunc(ctx context.Context, key interface{}, f gcache.Func, duration time.Duration, tag ...interface{}) (result *gvar.Var, err error) {
	result, err = Get(ctx, key)
	if err != nil {
		return
	}
	if result.IsNil() {
		value, err := f(ctx)
		if err != nil {
			return nil, err
		}
		if value == nil {
			return nil, nil
		}
		return gvar.New(value), set(ctx, key, value, duration, tag...)
	}
	return
}

// contains 判断 key 是否存在（包内使用）。
// 用 EXISTS 命令而非 GET，避免将 redis.Nil（key 不存在）误当错误传播。
func contains(ctx context.Context, key interface{}) (bool, error) {
	n, err := getRedisClient().Exists(ctx, gconv.String(key))
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// update 更新 key 的值并保留原有过期时间（包内使用）。
func update(ctx context.Context, key interface{}, value interface{}, tag ...interface{}) (oldValue *gvar.Var, exist bool, err error) {
	var (
		v       *gvar.Var
		oldPTTL int64
	)
	defaultKey := gconv.String(key)
	// TTL（毫秒）
	oldPTTL, err = getRedisClient().PTTL(ctx, defaultKey)
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		// key 不存在或已过期
		return
	}
	v, err = Get(ctx, key)
	if err != nil {
		return
	}
	oldValue = v
	if value == nil {
		_, err = Remove(ctx, key)
		return
	}
	if oldPTTL == -1 {
		err = set(ctx, key, value, 0, tag...)
	} else {
		err = set(ctx, key, value, time.Duration(oldPTTL/1000)*time.Second, tag...)
	}
	return oldValue, true, err
}

// updateExpire 更新 key 的过期时间（包内使用）。
func updateExpire(ctx context.Context, key interface{}, duration time.Duration, tag ...interface{}) (oldDuration time.Duration, err error) {
	var oldPTTL int64
	oldPTTL, err = getRedisClient().PTTL(ctx, gconv.String(key))
	if err != nil {
		return
	}
	if oldPTTL == -2 || oldPTTL == 0 {
		return
	}
	oldDuration = time.Duration(oldPTTL) * time.Millisecond
	if duration < 0 {
		_, err = Remove(ctx, key)
		return
	}
	if duration > 0 {
		_, err = getRedisClient().PExpire(ctx, gconv.String(key), duration.Milliseconds())
		return
	}
	// duration == 0：取消过期（持久化），重新写入不带 EX 的 TagCache 条目
	v, err := Get(ctx, key)
	if err != nil {
		return
	}
	var tags []string
	for _, t := range gconv.Strings(tag) {
		if t != "" {
			tags = append(tags, t)
		}
	}
	err = getTagCacheInstance().Set(ctx, gconv.String(key), v.Val(), 0, tags)
	return
}

// getExpire 获取 key 的剩余过期时间（包内使用）。
func getExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	return getAdapterRedis().GetExpire(ctx, key)
}

// Remove 删除缓存 key，支持单个或多个（对外公开）。
func Remove(ctx context.Context, key interface{}) (lastValue *gvar.Var, err error) {
	tc := getTagCacheInstance()

	// 批量 key 处理
	if keys, ok := key.([]interface{}); ok {
		for _, k := range keys {
			strKey := gconv.String(k)
			if lastValue, err = tc.Get(ctx, strKey); err != nil {
				continue
			}
			if err = tc.Delete(ctx, strKey); err != nil {
				return lastValue, err
			}
		}
		return lastValue, nil
	}

	// 单个 key 处理
	strKey := gconv.String(key)
	if lastValue, err = tc.Get(ctx, strKey); err != nil {
		return nil, err
	}
	err = tc.Delete(ctx, strKey)
	return
}

// RemoveByTag 按 tag 批量清除缓存（对外公开）。
func RemoveByTag(ctx context.Context, tags ...interface{}) (err error) {
	g.Log().Debug(ctx, "RemoveByTag:", tags)
	if !g.IsEmpty(tags) && len(tags) > 0 {
		tc := getTagCacheInstance()
		err = tc.InvalidateTags(ctx, gconv.Strings(tags))
		if err != nil {
			g.Log().Debug(ctx, "InvalidateTags err:", err)
		}
	}
	return
}

// ClearCacheAll 清空所有缓存（对外公开）。
func ClearCacheAll(ctx context.Context) error {
	return getAdapterRedis().Clear(ctx)
}
