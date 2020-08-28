package main

import (
	"fmt"
	"go.uber.org/ratelimit"
	"time"
)

func main() {
	r := ratelimit.New(1) //per second
	prev := time.Now()
	for i := 0; i < 15; i++ {
		now := r.Take()
		fmt.Println(i, now.Sub(prev))
		prev = now
	}
}
