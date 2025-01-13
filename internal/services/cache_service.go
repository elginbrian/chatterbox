package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
    RedisClient *redis.Client
}

func NewCacheService() *CacheService {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", 
        Password: "",               
        DB:       0,                
    })
    return &CacheService{RedisClient: rdb}
}

func (cs *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    return cs.RedisClient.Set(ctx, key, value, expiration).Err()
}

func (cs *CacheService) Get(ctx context.Context, key string) (string, error) {
    return cs.RedisClient.Get(ctx, key).Result()
}

func (cs *CacheService) Delete(ctx context.Context, key string) error {
    return cs.RedisClient.Del(ctx, key).Err()
}