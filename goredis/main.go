package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	fmt.Println(rdb)
	ctx := context.Background()
	val, err := rdb.Get(ctx, "hello").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}
