package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
func main() {
	test()
}

func test() {
	fillInteval := 1 * time.Minute
	var limitNum int64 = 5
	var segmentNum int64 = 6
	waitTime := 30
	fmt.Printf("time range from 0 to %d\n", waitTime)
	time.Sleep(time.Duration(waitTime) * time.Second)
	for i := 0; i < 10; i++ {
		fmt.Printf("time range from %d to %d\n", i*10+waitTime, (i+1)*10+waitTime)
		for j := 0; j < 8; j++ {
			rs := slidingWindowRatelimit("test", fillInteval, segmentNum, limitNum)
			fmt.Println("result is:", rs)
		}
		time.Sleep(10 * time.Second)
	}
}

//segmentNum split inteval time into smaller segments
func slidingWindowRatelimit(key string, fillInteval time.Duration, segmentNum int64, limitNum int64) bool {
	segmentInteval := fillInteval.Seconds() / float64(segmentNum)
	tick := float64(time.Now().Unix()) / segmentInteval
	currentKey := fmt.Sprintf("%s_%d_%d_%d_%f", key, fillInteval, segmentNum, limitNum, tick)
	//fmt.Println(currentKey)

	startCount := 0
	_, err := client.SetNX(currentKey, startCount, fillInteval).Result()
	if err != nil {
		panic(err)
	}
	quantum, err := client.Incr(currentKey).Result()
	if err != nil {
		panic(err)
	}
	//add in the number of the previous time
	for tickStart := segmentInteval; tickStart < fillInteval.Seconds(); tickStart += segmentInteval {
		tick = tick - 1
		preKey := fmt.Sprintf("%s_%d_%d_%d_%f", key, fillInteval, segmentNum, limitNum, tick)
		val, err := client.Get(preKey).Result()
		if err != nil {
			val = "0"
		}
		num, err := strconv.ParseInt(val, 0, 64)
		quantum = quantum + num
		if quantum > limitNum {
			client.Decr(currentKey).Result()
			return false
		}
	}
	return true

}
