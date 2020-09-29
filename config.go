package main

import (
	"github.com/go-redis/redis"
	"os"
)

func jwtConfig() {
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
}

func redisInit() {
	dsn := os.Getenv("REDIS_DSN")
	redisClient = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}
