package bootstrap

import (
	"fmt"
	"game-api/internal/config"
	"sync"

	"github.com/go-redis/redis"
)

var (
	Redis *redis.Client
	redisOnce sync.Once
)

// InitRedis 初始化 Redis（单例）
func InitRedis() error {
    var err error

	redisOnce.Do(func() {

		Redis = redis.NewClient(&redis.Options{
			Addr:     config.Cfg.Redis.Addr,
			Password: config.Cfg.Redis.Password,
			DB:       config.Cfg.Redis.DB,
		})

		_, err = Redis.Ping().Result()
		if err != nil {
			err = fmt.Errorf("redis connect failed: %w", err)
			return
		}

		fmt.Println("Redis connected")
	})

	return err
}

// GetRedis 获取 Redis Client
func GetRedis() *redis.Client {
	return Redis
}

// CloseRedis 关闭 Redis
func CloseRedis() error {

	if Redis == nil {
		return nil
	}

	return Redis.Close()
}