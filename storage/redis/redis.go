package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func ConnectReids() {
	fmt.Println("Testing Golang Redis")

	client := redis.NewClient(&redis.Options{
		Addr:     "118.195.244.48:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(client.Context()).Result()
	fmt.Println(pong, err)
}
