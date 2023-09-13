package main

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var client *redis.Client

var key string

var mut *redsync.Mutex

var lmut sync.Mutex

var global int

func main() {
	client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	key = "sync"
	pool := goredis.NewPool(client)
	mut = redsync.New(pool).NewMutex(key)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			no := time.Now()
			work()
			fmt.Println(time.Since(no))
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("-------------")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			no := time.Now()
			work1()
			fmt.Println(time.Since(no))
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println("-------------")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			no := time.Now()
			work2()
			fmt.Println(time.Since(no))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func work() {
	lmut.Lock()
	defer lmut.Unlock()
	global++
	fmt.Println("w", global)
}

func work1() {
	err := lock(client, key)
	if err != nil {
		fmt.Println("w1", err)
		return
	}
	defer func() {
		if err == nil {
			unlock(client, key)
		}
	}()
	global++
	fmt.Println("w1", global)
}

func work2() {
	err := mut.Lock()
	if err != nil {
		fmt.Println("w2", err)
		return
	}
	defer func() {
		if err == nil {
			_, _ = mut.Unlock()
		}
	}()
	global++
	fmt.Println("w2", global)
}
