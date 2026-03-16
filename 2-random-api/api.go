package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func randNumber(w http.ResponseWriter, r *http.Request) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := rnd.Intn(6) + 1
	fmt.Fprintf(w, "Number: %d", num)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/rand", randNumber)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	fmt.Println("Сгенерировано рандомное чилсо от 1 до 6")
	server.ListenAndServe()
}
