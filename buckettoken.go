package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"math"
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
	key := "111"
	var capacity int64 = 2
	for i := 0; i < 10; i++ {
		rs := bucketTokenRateLimit(key, 1*time.Second, 1, capacity)
		fmt.Println(rs)
	}
}

//rate increment number per second
//capacity total number in the bucket
func bucketTokenRateLimit(key string, fillInterval time.Duration, limitNum int64, capacity int64) bool {
	currentKey := fmt.Sprintf("%s_%d_%d_%d", key, fillInterval, limitNum, capacity)
	numKey := "num"
	lastTimeKey := "lasttime"
	currentTime := time.Now().Unix()
	//only init once fill fully
	client.HSetNX(currentKey, numKey, capacity).Result()
	client.HSetNX(currentKey, lastTimeKey, currentTime).Result()

	//compute current available number
	result, _ := client.HMGet(currentKey, numKey, lastTimeKey).Result()
	lastNum, _ := strconv.ParseInt(result[0].(string), 0, 64)
	lastTime, _ := strconv.ParseInt(result[1].(string), 0, 64)
	rate := float64(limitNum) / float64(fillInterval.Seconds())
	incrNum := int64(math.Ceil(float64(currentTime-lastTime) * rate)) //increment number from lasttime to currenttime
	currentNum := min(lastNum+incrNum, capacity)

	//can access
	if currentNum > 0 {
		var fields = map[string]interface{}{lastTimeKey: currentTime, numKey: currentNum - 1}
		client.HMSet(currentKey, fields)
		return true
	}
	return false

}

func min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
