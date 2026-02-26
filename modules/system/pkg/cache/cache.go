// Package cache
// @Link  https://github.com/huagelong/devinggo
// @Copyright  Copyright (c) 2024 devinggo
// @Author  Kai <hpuwang@gmail.com>
// @License  https://github.com/huagelong/devinggo/blob/master/LICENSE

package cache

import (
	"bufio"
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

// cache 缓存驱动
var cache *gcache.Cache

// tagCache 全局 TagCache 单例
var tagCache *TagCache

var groupKey = "cache"
var cachePrefixSelectCache = "SelectCache:"

// cache:metadata:
var scanCount = 100

// GetCache 缓存实例
func GetCache() *gcache.Cache {
	if cache == nil {
		panic("cache uninitialized.")
	}
	return cache
}

// SetAdapter 设置缓存适配器，同时初始化 TagCache 单例
func SetAdapter(ctx context.Context) {
	var cacheAdapter gcache.Adapter
	cacheAdapter = NewAdapter()
	g.DB().GetCache().SetAdapter(cacheAdapter)
	// 通用缓存
	cache = gcache.New()
	cache.SetAdapter(cacheAdapter)
	// 初始化 TagCache 单例
	tc, err := newTagCache(ctx, g.Redis(groupKey))
	if err != nil {
		panic("failed to initialize tag cache: " + err.Error())
	}
	tagCache = tc
}

// getTagCacheInstance 获取 TagCache 单例（包内使用）
func getTagCacheInstance() *TagCache {
	if tagCache == nil {
		panic("tag cache uninitialized, call SetAdapter first")
	}
	return tagCache
}

// getRedisClient 获取 Redis 客户端（包内使用）
func getRedisClient() *gredis.Redis {
	return g.Redis(groupKey)
}

// getAdapterRedis 获取 Redis 缓存适配器（包内使用）
func getAdapterRedis() gcache.Adapter {
	return gcache.NewAdapterRedis(g.Redis(groupKey))
}

func ClearByTable(ctx context.Context, table string) (err error) {
	return RemoveByTag(ctx, table)
}

func GetKeys(ctx context.Context) (keys []string, err error) {
	match := "*"
	keys = make([]string, 0)
	iterator := uint64(0)
	var listKeys []string
	for {
		iterator, listKeys, err = getRedisClient().Scan(ctx, iterator, gredis.ScanOption{
			Match: match,
			Count: scanCount,
		})
		if err != nil {
			g.Log().Error(ctx, "Scan error:", err)
			break
		}
		if len(listKeys) > 0 {
			keys = append(keys, listKeys...)
		}
		if iterator == 0 {
			break
		}
	}
	return
}

func GetInfo(ctx context.Context) (map[string]map[string]interface{}, error) {
	info, err := getRedisClient().Do(ctx, "INFO")
	if err != nil {
		return nil, err
	}
	var result = make(map[string][]map[string]interface{})
	scanner := bufio.NewScanner(strings.NewReader(info.String()))
	var key string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			key = strings.TrimSpace(strings.Split(line, "#")[1])
			result[key] = make([]map[string]interface{}, 0)
		} else if len(line) != 0 {
			kv := strings.Split(line, ":")
			m := make(map[string]interface{})
			//判断指标值是否有多个
			if strings.Contains(kv[1], ",") {
				sunValueList := strings.Split(kv[1], ",")
				sunValue := make(map[string]interface{})
				for _, s := range sunValueList {
					skv := strings.Split(s, "=")
					sunValue[skv[0]] = skv[1]
				}
				m[kv[0]] = sunValue
			} else {
				m[kv[0]] = kv[1]
			}
			result[key] = append(result[key], m)
		}
	}

	var res = make(map[string]map[string]interface{})
	for k, vList := range result {
		var value = make(map[string]interface{})
		for _, v := range vList {
			for k1, v1 := range v {
				value[k1] = v1
			}
		}
		res[k] = value
	}

	return res, nil
}
