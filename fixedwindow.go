package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:       "localhost:6379",
		Password:   "",
		DB:         0,
		PoolSize:   3,
		MaxRetries: 3,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	//test1()
	test2()
	var block = make(chan bool)
	<-block
}

//test fixed window
func test1() {
	for i := 0; i < 10; i++ {
		go func() {
			rs := fixedWindowRateLimit("test1", 1*time.Second, 5)
			fmt.Println("result is:", rs)
		}()
	}
}

//weak point for fixed window
func test2() {
	fillInteval := 1 * time.Minute
	var limitNum int64 = 5
	waitTime := 30
	fmt.Printf("time range from 0 to %d\n", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)
	for i := 0; i < 10; i++ {
		fmt.Printf("time range from %d to %d\n", i*10+waitTime, (i+1)*10+waitTime)
		rs := fixedWindowRateLimit("test2", fillInteval, limitNum)
		fmt.Println("result is:", rs)
		time.Sleep(10 * time.Second)
	}
}

//@param key string object for rate limit such as uid/ip+url
//@param fillInterval time.Duration such as 1*time.Second
//@param limitNum max int64 allowed number per fillInterval
//@return whether reach rate limit, false means reach
func fixedWindowRateLimit(key string, fillInterval time.Duration, limitNum int64) bool {
	//current tick time window
	tick := int64(time.Now().Unix() / int64(fillInterval.Seconds()))
	currentKey := fmt.Sprintf("%s_%d_%d_%d", key, fillInterval, limitNum, tick)
	fmt.Println(currentKey)

	startCount := 0
	_, err := client.SetNX(currentKey, startCount, fillInterval).Result()
	if err != nil {
		panic(err)
	}
	//number in current time window
	quantum, err := client.Incr(currentKey).Result()
	if err != nil {
		panic(err)
	}
	if quantum > limitNum {
		return false
	}
	return true
}
