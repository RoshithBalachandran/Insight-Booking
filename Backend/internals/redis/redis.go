package redis

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"), // localhost:6379
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Redis connected successfully")
}
