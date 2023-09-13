package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func lock(cli *redis.Client, key string) error {
	for {
		result, err := cli.SetNX(context.Background(), key, 1, time.Second).Result()
		if err != nil {
			return err
		}
		if result {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
}

func unlock(cli *redis.Client, key string) {
	cli.Del(context.Background(), key)
}
