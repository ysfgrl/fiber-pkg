package cache

import (
	"context"
	"github.com/ysfgrl/gerror"
	"time"
)

type Cache interface {
	Ping(ctx context.Context) *gerror.Error
	SetStr(ctx context.Context, key string, value string) *gerror.Error
	SetStrWithExp(ctx context.Context, key string, value string, expire time.Duration) *gerror.Error
	GetStr(ctx context.Context, key string) (string, *gerror.Error)
	SetStructWithExp(ctx context.Context, key string, value interface{}, expire time.Duration) *gerror.Error
	SetStruct(ctx context.Context, key string, value interface{}) *gerror.Error
	GetStruct(ctx context.Context, key string, value interface{}) *gerror.Error
	DeleteKey(ctx context.Context, key string) (bool, *gerror.Error)
	IsExistKey(ctx context.Context, key string) (int64, *gerror.Error)
}
