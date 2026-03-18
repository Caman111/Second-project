package main

import (
	"3-validation-api/config"
	"3-validation-api/internal/verify"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	service := verify.NewService(cfg.Email, cfg.Password, cfg.Address)
	handler := verify.NewHandler((service))

	http.HandleFunc("/send", handler.SendEmailHandler)
	http.HandleFunc("/verify", handler.VerifyHandler)
	log.Println("Сервер запушен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
