package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/ysfgrl/gerror"
	"time"
)

type BaseRedis struct {
	Client     *redis.Client
	Expiration time.Duration
	Prefix     string
}

func (cache *BaseRedis) SetStrWithExp(ctx context.Context, key string, value string, expire time.Duration) *gerror.Error {
	_, err := cache.Client.Set(ctx, cache.Prefix+key, value, expire).Result()
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}

func (cache *BaseRedis) SetStr(ctx context.Context, key string, value string) *gerror.Error {
	return cache.SetStrWithExp(ctx, key, value, cache.Expiration)
}

func (cache *BaseRedis) GetStr(ctx context.Context, key string) (string, *gerror.Error) {
	value, err := cache.Client.Get(ctx, cache.Prefix+key).Result()
	if err != nil {
		return "", gerror.GetError(err)
	}
	return value, nil
}

func (cache *BaseRedis) SetStructWithExp(ctx context.Context, key string, value interface{}, expire time.Duration) *gerror.Error {
	_, err := cache.Client.HSet(ctx, cache.Prefix+key, value).Result()
	if err != nil {
		return gerror.GetError(err)
	}
	cache.Client.Expire(ctx, cache.Prefix+key, expire)
	return nil
}

func (cache *BaseRedis) SetStruct(ctx context.Context, key string, value interface{}) *gerror.Error {
	return cache.SetStructWithExp(ctx, key, value, cache.Expiration)
}

func (cache *BaseRedis) GetStruct(ctx context.Context, key string, value interface{}) *gerror.Error {
	err := cache.Client.HGetAll(ctx, cache.Prefix+key).Scan(value)
	if err != nil {
		return gerror.GetError(err)
	}
	return nil
}

func (cache *BaseRedis) DeleteKey(ctx context.Context, key string) (bool, *gerror.Error) {
	boolCmd := cache.Client.Expire(ctx, cache.Prefix+key, 1)
	if err := boolCmd.Err(); err != nil {
		return false, gerror.GetError(err)
	}
	return boolCmd.Val(), nil
}

func (cache *BaseRedis) IsExistKey(ctx context.Context, key string) (int64, *gerror.Error) {
	intCmd := cache.Client.Exists(ctx, key)
	if err := intCmd.Err(); err != nil {
		return 0, gerror.GetError(err)
	}
	return intCmd.Val(), nil
}

func (cache *BaseRedis) Ping(ctx context.Context) *gerror.Error {
	if _, err := cache.Client.Ping(ctx).Result(); err != nil {
		return gerror.GetError(err)
	}
	return nil
}
