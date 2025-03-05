package db

import (
	"context"
	"fmt"
	"govue3/config"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Database.Redis,
		Password: config.Cfg.Database.RedisPassword,
		DB:       0,
	})

	if _, err := Redis.Ping(context.Background()).Result(); err != nil {
		panic(fmt.Sprintf("Redis连接失败: %v", err))
	}
	fmt.Println("Redis连接成功")
}
