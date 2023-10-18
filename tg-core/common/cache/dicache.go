/**
*	function:	local cache based on FreeCache
*	Author	:	dayunzhangyunfeng@didiglobal.com
*	Since	:	2019-05-31 17:25:26
 */

package cache

import (
	"context"
	"fmt"
	"git.xiaojukeji.com/nuwa/golibs/redis"
	"github.com/coocood/freecache"
)

const (
	LOCALCACHE_TTL = 1800 //unit:s
)

type DiCache struct {
	RedisCache *redis.Manager
	Cache      *freecache.Cache
	Size       int
	TTL        int
}

func NewDiCache(redisCache *redis.Manager, size int, ttl int) *DiCache {
	var freeCache = freecache.NewCache(size)
	return &DiCache{
		RedisCache: redisCache,
		Cache:      freeCache,
		Size:       size,
		TTL:        ttl}
}

func (this *DiCache) Get(ctx context.Context, key string) (string, error) {
	vBytes, err := this.Cache.Get([]byte(key))
	if err == nil {
		return string(vBytes), nil
	}

	val, err := this.RedisCache.Get(ctx, key)
	if err == nil {
		this.Cache.Set([]byte(key), []byte(val), this.TTL)
	}

	return val, err
}

func (this *DiCache) GetMust(ctx context.Context, key, defaultValue string) string {
	vBytes, err := this.Cache.Get([]byte(key))
	if err == nil {
		return string(vBytes)
	}

	val, err := this.RedisCache.Get(ctx, key)
	if err == nil {
		this.Cache.Set([]byte(key), []byte(val), this.TTL)
		return val
	}

	return defaultValue
}

func (this *DiCache) MGet(ctx context.Context, keys []string) (map[string]string, error) {
	valMap := make(map[string]string)
	if keys == nil || len(keys) == 0 {
		return valMap, nil
	}

	notHitKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		vBytes, err := this.Cache.Get([]byte(key))
		if err == nil {
			valMap[key] = string(vBytes)
		} else {
			notHitKeys = append(notHitKeys, key)
		}
	}

	notHitCount := len(notHitKeys)
	if notHitCount == 0 {
		return valMap, nil
	}

	//不知道err不空时tempMap还是否可能有值，兜一下
	tempMap, err := this.RedisCache.MGet(ctx, notHitKeys)
	if tempMap == nil || len(tempMap) == 0 {
		return valMap, fmt.Errorf("redis keys not found all keys:%v, err:%v", notHitKeys, err)
	}

	tempKeys := make([]string, 0, notHitCount)
	for _, notHitKey := range notHitKeys {
		if notHitVal, ok := tempMap[notHitKey]; ok {
			valMap[notHitKey] = notHitVal
			this.Cache.Set([]byte(notHitKey), []byte(notHitVal), this.TTL)
		} else {
			tempKeys = append(tempKeys, notHitKey)
		}
	}

	if len(tempKeys) > 0 {
		err = fmt.Errorf("redis keys not found keys:%v", tempKeys)
	}

	return valMap, err
}
