package main

import (
	"concurrency/config"
	"concurrency/internal/verify"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func maim() {
	cfg := config.LoadConfig()

	service := verify.NewService(cfg.Email, cfg.Password, cfg.Address)
	handler := verify.NewHandler((service))

	http.HandleFunc("/send", handler.SendEmailHandler)
	http.HandleFunc("/verify/", handler.VerifyHandler)

	log.Println("Сервер запушен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
