package connection

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/Bappy60/BookStore_in_Go/pkg/config"
)

var Client *redis.Client

func RedisConnection() {

	Client = redis.NewClient(&redis.Options{
		Addr:     config.GConfig.REDIS_HOST + ":" + config.GConfig.REDIS_PORT,
		Password: config.GConfig.REDIS_PASS,
		DB:       0,
	})
	red, err := Client.Ping(context.Background()).Result()
	fmt.Println(red)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("redis connection successful...")
}

func Redis() *redis.Client {
	if Client == nil {
		RedisConnection()
	}
	return Client
}
