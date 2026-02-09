package config

import (
	"context"
	"fmt"
	"log"
	"os"
	//"time"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // 默认Redis地址
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // 你的Redis密码，如果没有则留空
		DB:       0,    // 默认DB
	})

	// 尝试Ping一下Redis，检查连接是否成功
	pong, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Redis connected:", pong)
}