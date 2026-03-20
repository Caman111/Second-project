package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	nums := make(chan int)
	squares := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			num := r.Intn(101)   
			nums <- num
		}
		close(nums)

	}()

	go func() {
		for n := range nums {
			squares <- n * n
		}
		close(squares)
	}()
	for sg := range squares {
		fmt.Println(sg)
	}
}
